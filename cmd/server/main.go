package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/fdanctl/piggytron/config"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/application/appbudget"
	"github.com/fdanctl/piggytron/internal/application/appexpensecategory"
	"github.com/fdanctl/piggytron/internal/application/appincomecategory"
	"github.com/fdanctl/piggytron/internal/application/appledger"
	"github.com/fdanctl/piggytron/internal/application/appuser"
	"github.com/fdanctl/piggytron/internal/auth"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/fdanctl/piggytron/internal/interface/http/handlers"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/interface/http/shared"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("load config failed: ", err.Error())
		return
	}

	var logger *slog.Logger

	if cfg.IsDev {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		logger.Error("failed to open db", "error", err.Error())
		os.Exit(1)
		return
	}
	if err := db.Ping(); err != nil {
		logger.Error("failed to connect to db", "error", err.Error())
		os.Exit(1)
		return
	}
	defer db.Close()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("failed to connect to redis", "error", err.Error())
		os.Exit(1)
		return
	}
	defer client.Close()

	hasher := appuser.NewPasswordHasher(
		cfg.HashConfig.Time,
		cfg.HashConfig.Memory,
		cfg.HashConfig.Threads,
		cfg.HashConfig.KeyLen,
		cfg.HashConfig.SaltLen,
	)

	sessionStore := rdb.NewSessionStore(client)
	sessionVersionStore := rdb.NewSessionVersionStore(client)
	sessionManager := auth.NewSessionManager(sessionStore, sessionVersionStore)

	sessionCM := shared.NewCookieMaker(http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   !cfg.IsDev,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((time.Hour * 24).Seconds()),
	})

	// repositories
	accountRepo := postgres.NewAccountRepository(db)
	ledgerRepo := postgres.NewLedgerRepository(db)
	expenseCatRepo := postgres.NewExpenseCategoryRepository(db)
	incomeCatRepo := postgres.NewIncomeCategoryRepository(db)
	userRepo := postgres.NewUserRepository(db)
	budgetRepo := postgres.NewBudgetRepository(db)

	// query services
	catQueryService := postgres.NewCategoryQueryService(db)
	ledgerQueryService := postgres.NewLedgerQueryService(db)
	accountQueryService := postgres.NewAccountQueryService(db)

	// services
	accountService := appaccount.NewService(accountRepo, db)
	ledgerService := appledger.NewService(ledgerRepo, db)
	expenseCatService := appexpensecategory.NewService(expenseCatRepo)
	incomeCatService := appincomecategory.NewService(incomeCatRepo)
	userService := appuser.NewService(userRepo, hasher, sessionManager)
	budgetService := appbudget.NewService(budgetRepo)

	// web mux - returns full HTML page (or, in most cases, just the main element if Hx-Request)
	webMux := http.NewServeMux()
	webMux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("web/static")),
		),
	)
	webMux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/assets/favicon.ico")
	})

	dashboardHandler := handlers.NewDashboardHandler(
		ledgerQueryService,
		accountQueryService,
		catQueryService,
	)
	webMux.Handle("/", middleware.AuthProtectedRoute(dashboardHandler))

	bh := handlers.NewBudgetPageHandler(catQueryService, ledgerQueryService)
	webMux.Handle("/budget", middleware.AuthProtectedRoute(bh))

	goalsHandler := handlers.NewGoalsHandler(
		accountService,
		ledgerQueryService,
		accountQueryService,
	)
	webMux.Handle("/goals", middleware.AuthProtectedRoute(goalsHandler))
	webMux.Handle("/goals/{id}", middleware.AuthProtectedRoute(goalsHandler))

	banksHandler := handlers.NewBanksHandler(
		accountService,
		ledgerQueryService,
		accountQueryService,
	)
	webMux.Handle("/banks", middleware.AuthProtectedRoute(banksHandler))
	webMux.Handle("/banks/{id}", middleware.AuthProtectedRoute(banksHandler))

	ledgerPageHandler := handlers.NewLedgerPageHandler(
		ledgerQueryService,
	)
	webMux.Handle("/ledger", middleware.AuthProtectedRoute(ledgerPageHandler))

	eh := handlers.NewExpensesHandler(ledgerQueryService)
	webMux.Handle("/reports/expenses", middleware.AuthProtectedRoute(eh))

	ih := handlers.NewIncomeHandler(ledgerQueryService)
	webMux.Handle("/reports/income", middleware.AuthProtectedRoute(ih))

	categoriesHandler := handlers.NewCategoriesHandler(
		expenseCatService,
		incomeCatService,
		ledgerQueryService,
	)

	webMux.Handle("/categories", middleware.AuthProtectedRoute(categoriesHandler))
	webMux.Handle("/categories/{id}", middleware.AuthProtectedRoute(categoriesHandler))

	lh := handlers.NewLoginHandler(cfg.IsDev)
	webMux.Handle("/login", middleware.AuthenticatedRedirect(lh))

	sh := handlers.SignupHandler{}
	webMux.Handle("/signup", middleware.AuthenticatedRedirect(&sh))

	ah := handlers.NewAccountHandler(userService)
	webMux.Handle("/account", middleware.AuthProtectedRoute(ah))

	ph := handlers.NewPreferencesHandler(userService)
	webMux.Handle("/preferences", middleware.AuthProtectedRoute(ph))

	if cfg.IsDev {
		th := handlers.TestHandler{}
		webMux.Handle("/test", &th)
	}

	// partials mux - returns HTMX fragments
	partialsMux := http.NewServeMux()

	authHandler := handlers.NewAuthHandler(userService, sessionCM)
	partialsMux.Handle("/partials/auth/{action}", authHandler)

	userHandler := handlers.NewUserHandler(userService)
	partialsMux.Handle("/partials/user/change-name", userHandler)

	budgetHandler := handlers.NewBudgetHandler(
		budgetService,
		catQueryService,
	)
	partialsMux.Handle("/partials/budget", budgetHandler)

	importBudgetHandler := handlers.NewImportBudgetHandler(budgetService)
	partialsMux.Handle("/partials/budget/import", importBudgetHandler)

	incomeCatHandler := handlers.NewIncomeCategoriesHandler(incomeCatService)
	partialsMux.Handle("/partials/income-category", incomeCatHandler)

	expenseCatHandler := handlers.NewExpenseCategoriesHandler(expenseCatService)
	partialsMux.Handle("/partials/expense-category", expenseCatHandler)

	filteredLedgerHandler := handlers.NewFilteredLedgerHandler(ledgerQueryService)
	partialsMux.Handle("/partials/ledger", filteredLedgerHandler)

	ledgerHandler := handlers.NewLedgerHandler(
		ledgerService,
		catQueryService,
		accountService,
	)
	partialsMux.Handle("/partials/ledger/entry", ledgerHandler)

	ledgerEntryHandler := handlers.NewLedgerEntryHandler(
		ledgerService,
		catQueryService,
		accountService,
	)
	partialsMux.Handle("/partials/ledger/entry/{id}", ledgerEntryHandler)

	goalContributeHandler := handlers.NewGoalContributeHandler(
		ledgerService,
		catQueryService,
		accountService,
	)
	partialsMux.Handle("/partials/goal-contribute/{id}", goalContributeHandler)

	transactionDetails := handlers.NewTransactionDetailsHandler(ledgerQueryService)
	partialsMux.Handle("/partials/ledger/entry/details/{id}", transactionDetails)

	goalContributions := handlers.NewGoalContributionsHandler(ledgerQueryService)
	partialsMux.Handle("/partials/contributions", goalContributions)

	catHistChartHandler := handlers.NewCategoryChartHandler(catQueryService)
	partialsMux.Handle("/partials/charts/cat-hist/{id}", catHistChartHandler)

	accountChartHandler := handlers.NewAccountChartHandler(
		accountQueryService,
	)
	partialsMux.Handle("/partials/charts/account-hist/{id}", accountChartHandler)

	accountsHistChartHandler := handlers.NewAccountsHistoryChartHandler(
		accountQueryService,
	)
	partialsMux.Handle("/partials/charts/accounts-history", accountsHistChartHandler)

	dashboardBudgetCharts := handlers.NewDashboardBudgetCharts(
		catQueryService,
	)
	partialsMux.Handle("/partials/charts/dashboard-budget-spent", dashboardBudgetCharts)

	bankChartHandler := handlers.NewBankChartHandler(
		accountQueryService,
	)
	partialsMux.Handle("/partials/charts/bank-hist/{id}", bankChartHandler)

	banksChartsHandler := handlers.NewBanksChartsHandler(
		accountQueryService,
	)
	partialsMux.Handle("/partials/charts/banks", banksChartsHandler)

	budgetChartHandler := handlers.NewBudgetChartHandler(
		catQueryService,
	)
	partialsMux.Handle("/partials/charts/budget-chart/{month}", budgetChartHandler)

	ledgerFiltersHandler := handlers.NewFilterDialogHandler(
		catQueryService,
		accountService,
		ledgerQueryService,
		accountQueryService,
	)
	partialsMux.Handle("/partials/ledger-filters", ledgerFiltersHandler)

	bankHandler := handlers.NewBankHandler(accountService)
	partialsMux.Handle("/partials/bank", bankHandler)

	goalHandler := handlers.NewGoalHandler(
		accountService, accountQueryService, catQueryService,
	)
	partialsMux.Handle("/partials/goal", goalHandler)
	partialsMux.Handle("/partials/goal/{id}", goalHandler)

	logger.Info("server starting", "addr", ":8080")

	if cfg.IsDev {
		// Try to find local IP
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			for _, addr := range addrs {
				// check the address type and ignore loopback
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil { // IPv4
						logger.Info(
							fmt.Sprintf(
								"Accessible on your LAN at: http://%s:%s",
								ipnet.IP.String(),
								cfg.ServerPort,
							),
						)
					}
				}
			}
		}
	}

	rootMux := http.NewServeMux()
	rootMux.Handle("/", webMux)
	rootMux.Handle("/partials/", middleware.RequireHTMX(partialsMux))

	http.ListenAndServe(
		":"+cfg.ServerPort,
		middleware.Chain(
			rootMux,
			middleware.RecoveryMiddleware,
			middleware.RequestIDMiddleware,
			middleware.LoggingMiddleware(logger),
			middleware.AuthMiddleware(sessionManager),
		),
	)
}
