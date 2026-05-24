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
	"github.com/fdanctl/piggytron/internal/application/account"
	"github.com/fdanctl/piggytron/internal/application/budget"
	"github.com/fdanctl/piggytron/internal/application/charts"
	expensecategory "github.com/fdanctl/piggytron/internal/application/expense_category"
	incomecategory "github.com/fdanctl/piggytron/internal/application/income_category"
	"github.com/fdanctl/piggytron/internal/application/user"
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

	hasher := user.NewPasswordHasher(
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
	// transactionRepo := postgres.NewTransactionRepository(db)
	expenseCatRepo := postgres.NewExpenseCategoryRepository(db)
	incomeCatRepo := postgres.NewIncomeCategoryRepository(db)
	userRepo := postgres.NewUserRepository(db)
	budgetRepo := postgres.NewBudgetRepository(db)

	// query services
	var catQueryService query.CategoryQueryService = postgres.NewCategoryQueryService(db)
	var transactionQueryService query.TransactionQueryService = postgres.NewTransactionQueryService(db)
	var accountQueryService query.AccountQueryService = postgres.NewAccountQueryService(db)

	// services
	accountService := account.NewService(accountRepo)
	// transactionService := transaction.NewService(transactionRepo)
	expenseCatService := expensecategory.NewService(expenseCatRepo)
	incomeCatService := incomecategory.NewService(incomeCatRepo)
	userService := user.NewService(userRepo, hasher, sessionStore)
	budgetService := budget.NewService(budgetRepo)
	chartsService := charts.NewService()

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

	bh := handlers.NewBudgetPageHandler(catQueryService, transactionQueryService)
	webMux.Handle("/budget", middleware.AuthProtectedRoute(bh))

	goalsHandler := handlers.NewGoalsHandler(
		accountService,
		transactionQueryService,
		accountQueryService,
	)
	webMux.Handle("/goals", middleware.AuthProtectedRoute(goalsHandler))
	webMux.Handle("/goals/{id}", middleware.AuthProtectedRoute(goalsHandler))

	banksHandler := handlers.NewBanksHandler(
		accountService,
		transactionQueryService,
		accountQueryService,
	)
	webMux.Handle("/banks", middleware.AuthProtectedRoute(banksHandler))
	webMux.Handle("/banks/{id}", middleware.AuthProtectedRoute(banksHandler))

	allTransactionsHandler := handlers.NewAllTransactionsHandler(
		transactionQueryService,
	)
	webMux.Handle("/transactions/all", middleware.AuthProtectedRoute(allTransactionsHandler))

	eh := handlers.ExpensesHandler{}
	webMux.Handle("/transactions/expenses", middleware.AuthProtectedRoute(&eh))

	categoriesHandler := handlers.NewCategoriesHandler(
		expenseCatService,
		incomeCatService,
		transactionQueryService,
	)

	webMux.Handle("/categories", middleware.AuthProtectedRoute(categoriesHandler))
	webMux.Handle("/categories/{id}", middleware.AuthProtectedRoute(categoriesHandler))

	lh := handlers.LoginHandler{}
	webMux.Handle("/login", middleware.AuthenticatedRedirect(&lh))

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

	filteredTransaction := handlers.NewFilteredTransactionsHandler(transactionQueryService)
	partialsMux.Handle("/partials/transactions", filteredTransaction)

	goalContributions := handlers.NewGoalContributionsHandler(transactionQueryService)
	partialsMux.Handle("/partials/contributions", goalContributions)

	catHistChartHandler := handlers.NewCatHistChartHandler()
	partialsMux.Handle("/partials/charts/cat-hist/{id}", catHistChartHandler)

	accountChartHandler := handlers.NewAccountChartHandler(
		chartsService,
		accountQueryService,
	)
	partialsMux.Handle("/partials/charts/account-hist/{id}", accountChartHandler)

	banksChartsHandler := handlers.NewBanksChartsHandler(
		chartsService,
		transactionQueryService,
		accountQueryService,
	)
	partialsMux.Handle("/partials/charts/banks", banksChartsHandler)

	budgetChartHandler := handlers.NewBudgetChartHandler(
		chartsService,
		catQueryService,
	)
	partialsMux.Handle("/partials/charts/budget-chart/{month}", budgetChartHandler)

	transactionFiltersHandler := handlers.NewFilterDialogHandler(
		catQueryService,
		accountService,
		transactionQueryService,
		accountQueryService,
	)
	partialsMux.Handle("/partials/transaction-filters", transactionFiltersHandler)

	bankHandler := handlers.NewBankHandler(accountService)
	partialsMux.Handle("/partials/bank", bankHandler)

	goalHandler := handlers.NewGoalHandler(accountService, expenseCatService)
	partialsMux.Handle("/partials/goal", goalHandler)

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
