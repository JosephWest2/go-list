package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"josephwest2.com/go-list/app"
	"josephwest2.com/go-list/handlers"
)

func main() {
	isDev := os.Getenv("GO_LIST_IS_DEV") == "true"
	if isDev {
		println("Running in dev mode")
		err := godotenv.Load(".dev.env")
		if err != nil {
			log.Fatal("Failed to load .dev.env")
		}
	} else {
		println("Running in production mode")
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Failed to load .env")
		}
	}
	postgresConnectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	if postgresConnectionString == "" {
		log.Fatal("POSTGRES_CONNECTION_STRING is not set in .env")
	}
	jwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	if jwtSigningKey == "" {
		log.Fatal("JWT_SIGNING_KEY is not set in .env")
	}
	dbpool, err := pgxpool.New(context.Background(), postgresConnectionString)
	if err != nil {
		log.Fatal("Failed to create database pool:", err)
	}
	defer dbpool.Close()
	app := app.AppContext{
		DBpool:        dbpool,
		JwtSigningKey: []byte(jwtSigningKey),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.IndexPageHandler)
	mux.HandleFunc("GET /todo-list", handlers.TodoListPageHandler(app))
	mux.HandleFunc("GET /login", handlers.LoginPageHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler(app))
	mux.HandleFunc("GET /register", handlers.RegisterPageHandler)
	mux.HandleFunc("POST /register", handlers.RegisterHandler(app))
	mux.HandleFunc("GET /logout", handlers.LogoutPageHandler)
	mux.Handle("GET /web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.ListenAndServe(":3000", mux)
}
