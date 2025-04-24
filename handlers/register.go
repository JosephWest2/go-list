package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"josephwest2.com/go-list/app"
	"josephwest2.com/go-list/components"
	"josephwest2.com/go-list/sqlc"
)

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	RenderPage(components.RegisterPage(), w)
}

func RegisterHandler(app app.AppContext) http.HandlerFunc {
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
		err := ValidateUsername(usernameInput, app.DBpool)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = ValidatePassword(passwordInput, app.DBpool)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		q := sqlc.New(app.DBpool)
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(passwordInput), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("failed to hash password: ", err.Error())
		}

		_, err = q.CreateUser(context.Background(), sqlc.CreateUserParams{
			Username:     usernameInput,
			PasswordHash: string(passwordHash),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println("failed to create user: ", err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully registered"))
	}
}

func ValidateUsername(username string, dbpool *pgxpool.Pool) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(username) > 50 {
		return errors.New("username must 50 characters long or less")
	}
	q := sqlc.New(dbpool)
	_, err := q.GetUser(context.Background(), username)
	if err == nil {
		return errors.New("username already taken")
	}
	return nil
}

func ValidatePassword(password string, dbpool *pgxpool.Pool) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 50 {
		return errors.New("password must 50 characters long or less")
	}
	if !strings.ContainsAny(password, "0123456789") {
		return errors.New("password must contain at least one number")
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return errors.New("password must contain at least one uppercase letter")
	}
	return nil
}
