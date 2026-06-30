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
	"github.com/fdanctl/piggytron/internal/application/appcharts"
	"github.com/fdanctl/piggytron/internal/application/appexpensecategory"
	"github.com/fdanctl/piggytron/internal/application/appincomecategory"
	"github.com/fdanctl/piggytron/internal/application/appledger"
	"github.com/fdanctl/piggytron/internal/application/appuser"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/fdanctl/piggytron/internal/interface/http/handlers"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/interface/http/shared"
	"github.com/fdanctl/piggytron/internal/query"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("load config failed: ", err.Error())
		return
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalln("failed to open db", err.Error())
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
		log.Fatalln("failed to connect to redis", err.Error())
		return
	}
	defer client.Close()

	var logger *slog.Logger

	if cfg.IsDev {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	hasher := appuser.NewPasswordHasher(
		cfg.HashConfig.Time,
		cfg.HashConfig.Memory,
		cfg.HashConfig.Threads,
		cfg.HashConfig.KeyLen,
		cfg.HashConfig.SaltLen,
	)
	sessionStore := rdb.NewSessionStore(client)
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
	var catQueryService query.CategoryQueryService = postgres.NewCategoryQueryService(db)
	var ledgerQueryService query.LedgerQueryService = postgres.NewLedgerQueryService(db)
	var accountQueryService query.AccountQueryService = postgres.NewAccountQueryService(db)

	// services
	accountService := appaccount.NewService(accountRepo, db)
	ledgerService := appledger.NewService(ledgerRepo, db)
	expenseCatService := appexpensecategory.NewService(expenseCatRepo)
	incomeCatService := appincomecategory.NewService(incomeCatRepo)
	userService := appuser.NewService(userRepo, hasher, sessionStore)
	budgetService := appbudget.NewService(budgetRepo)
	chartsService := appcharts.NewService()

	webMux := http.NewServeMux() // returns full HTML page
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

	hh := handlers.HomeHandler{}
	webMux.Handle("/", middleware.AuthProtectedRoute(&hh))

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

	eh := handlers.ExpensesHandler{}
	webMux.Handle("/reports/expenses", middleware.AuthProtectedRoute(&eh))

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

	partialsMux := http.NewServeMux() // returns HTMX fragment

	userHandler := handlers.NewUserHandler(userService, sessionCM)
	partialsMux.Handle("/partials/auth/{action}", userHandler)

	budgetHandler := handlers.NewBudgetHandler(
		budgetService,
		chartsService,
		catQueryService,
	)
	partialsMux.Handle("/partials/budget", budgetHandler)

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

	catHistChartHandler := handlers.NewCategoryChartHandler(chartsService, catQueryService)
	partialsMux.Handle("/partials/charts/cat-hist/{id}", catHistChartHandler)

	accountChartHandler := handlers.NewAccountChartHandler(
		chartsService,
		accountQueryService,
	)
	partialsMux.Handle("/partials/charts/account-hist/{id}", accountChartHandler)

	bankChartHandler := handlers.NewBankChartHandler(
		chartsService,
		accountQueryService,
	)
	partialsMux.Handle("/partials/charts/bank-hist/{id}", bankChartHandler)

	banksChartsHandler := handlers.NewBanksChartsHandler(
		chartsService,
		ledgerQueryService,
		accountQueryService,
	)
	partialsMux.Handle("/partials/charts/banks", banksChartsHandler)

	budgetChartHandler := handlers.NewBudgetChartHandler(
		chartsService,
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

	// TODO remove
	partialsMux.HandleFunc("/partials/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprint(w, time.Now().Format(time.TimeOnly))
	})

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
			middleware.AuthMiddleware(sessionStore),
		),
	)
}
