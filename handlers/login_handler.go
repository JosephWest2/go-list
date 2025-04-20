package handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"josephwest2.com/go-list/sqlc"
)

func LoginHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usernameInput := r.FormValue("username")
		if usernameInput == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		passwordInput := r.FormValue("password")
		if passwordInput == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		q := sqlc.New(dbpool)

		user, err := q.GetUser(context.Background(), usernameInput)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordInput))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully logged in as " + usernameInput))

	}
}
