package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/config"
	"github.com/fdanctl/piggytron/internal/application/user"
	"github.com/fdanctl/piggytron/internal/interface/http/handlers"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"

	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	templates "github.com/fdanctl/piggytron/web/templates/layout"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	db, err := sql.Open("postgres", cfg.DBURL)
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	userService := user.NewService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	webMux := http.NewServeMux() // returns full HTML page
	webMux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("web/static")),
		),
	)
	webMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		templates.Base("title").Render(r.Context(), w)
	})

	partialsMux := http.NewServeMux() // returns HTMX fragment
	partialsMux.Handle("/partials/auth/{action}", userHandler)

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
	rootMux.Handle("/", webMux)
	rootMux.Handle(
		"/partials/",
		middleware.RequireHTMX(partialsMux),
	)

	http.ListenAndServe(":"+cfg.ServerPort, middleware.LoggingMiddleware(rootMux))
}
