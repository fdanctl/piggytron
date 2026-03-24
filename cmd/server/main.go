package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/config"
	"github.com/fdanctl/piggytron/internal/application/account"
	expensecategory "github.com/fdanctl/piggytron/internal/application/expense_category"
	incomecategory "github.com/fdanctl/piggytron/internal/application/income_category"
	"github.com/fdanctl/piggytron/internal/application/transaction"
	"github.com/fdanctl/piggytron/internal/application/user"
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

	bh := handlers.BudgetHandler{}
	webMux.Handle("/budget", middleware.AuthProtectedRoute(&bh))

	accountRepository := postgres.NewAccountRepository(db)
	bankService := account.NewService(accountRepository)
	banksHandler := handlers.NewBanksHandler(bankService)
	webMux.Handle("/banks", middleware.AuthProtectedRoute(banksHandler))
	webMux.Handle("/banks/{id}", middleware.AuthProtectedRoute(banksHandler))

	transactionRepo := postgres.NewTransactionRepository(db)
	transactionService := transaction.NewService(transactionRepo)
	allTransactionsHandler := handlers.NewAllTransactionsHandler(transactionService)
	webMux.Handle("/transactions/all", middleware.AuthProtectedRoute(allTransactionsHandler))

	eh := handlers.ExpensesHandler{}
	webMux.Handle("/transactions/expenses", middleware.AuthProtectedRoute(&eh))

	expenseCatRepo := postgres.NewExpenseCategoryRepository(db)
	expenseCatService := expensecategory.NewService(expenseCatRepo)

	incomeCatRepo := postgres.NewIncomeCategoryRepository(db)
	incomeCatService := incomecategory.NewService(incomeCatRepo)

	categoriesHandler := handlers.NewCategoriesHandler(
		expenseCatService,
		incomeCatService,
		transactionService,
	)

	webMux.Handle("/categories", middleware.AuthProtectedRoute(categoriesHandler))
	webMux.Handle("/categories/{id}", middleware.AuthProtectedRoute(categoriesHandler))

	lh := handlers.LoginHandler{}
	webMux.Handle("/login", middleware.AuthenticatedRedirect(&lh))

	sh := handlers.SignupHandler{}
	webMux.Handle("/signup", middleware.AuthenticatedRedirect(&sh))

	partialsMux := http.NewServeMux() // returns HTMX fragment

	userRepo := postgres.NewUserRepository(db)
	userService := user.NewService(userRepo, hasher, sessionStore)
	userHandler := handlers.NewUserHandler(userService, sessionCM)
	partialsMux.Handle("/partials/auth/{action}", userHandler)

	incomeCatHandler := handlers.NewIncomeCategoriesHandler(incomeCatService)
	partialsMux.Handle("/partials/income-category", incomeCatHandler)

	expenseCatHandler := handlers.NewExpenseCategoriesHandler(expenseCatService)
	partialsMux.Handle("/partials/expense-category", expenseCatHandler)

	filteredTransaction := handlers.NewFilteredTransactionsHandler(transactionService)
	partialsMux.Handle("/partials/transactions", filteredTransaction)

	catHistChartHandler := handlers.NewCatHistChartHandler()
	partialsMux.Handle("/partials/charts/cat-hist/{id}", catHistChartHandler)

	dialogHandler := handlers.NewDialogHandler(
		expenseCatService,
		incomeCatService,
		transactionService,
		bankService,
	)
	partialsMux.Handle("/partials/dialog/{dialog}", dialogHandler)

	// TODO remove
	partialsMux.HandleFunc("/partials/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprint(w, time.Now().Format(time.TimeOnly))
	})

	fmt.Println(time.Now())
	fmt.Println("Server running at http://localhost:" + cfg.ServerPort)
	// Try to find local IP
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			// check the address type and ignore loopback
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil { // IPv4
					fmt.Printf(
						"Accessible on your LAN at: http://%s:%s\n",
						ipnet.IP.String(),
						cfg.ServerPort,
					)
				}
			}
		}
	}

	rootMux := http.NewServeMux()
	rootMux.Handle("/", middleware.AuthMiddleware(sessionStore)(webMux))
	rootMux.Handle(
		"/partials/",
		middleware.Chain(
			partialsMux,
			middleware.RequireHTMX,
			middleware.AuthMiddleware(sessionStore),
		),
	)

	http.ListenAndServe(":"+cfg.ServerPort, middleware.LoggingMiddleware(rootMux))
}
