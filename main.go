package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"josephwest2.com/go-list/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env")
	}
	postgresConnectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	if postgresConnectionString == "" {
		log.Fatal("POSTGRES_CONNECTION_STRING is not set in .env")
	}
	dbpool, err := pgxpool.New(context.Background(), postgresConnectionString)
	if err != nil {
		log.Fatal("Failed to create database pool:", err)
	}
	defer dbpool.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.IndexPageHandler)
	mux.HandleFunc("GET /todo-list", handlers.TodoListPageHandler)
	mux.HandleFunc("GET /login", handlers.LoginPageHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler(dbpool))
	mux.HandleFunc("GET /register", handlers.RegisterPageHandler)
	mux.HandleFunc("POST /register", handlers.RegisterHandler(dbpool))
	mux.HandleFunc("GET /logout", handlers.LogoutPageHandler)
	mux.Handle("GET /web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.ListenAndServe(":3000", mux)
}
