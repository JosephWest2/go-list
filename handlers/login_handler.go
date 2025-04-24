package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"josephwest2.com/go-list/app"
	"josephwest2.com/go-list/sqlc"
)

func LoginHandler(appContext app.AppContext) http.HandlerFunc {
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

		q := sqlc.New(appContext.DBpool)

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
		claims := app.Claims{
			Username: usernameInput,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			},
		}
		authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := authToken.SignedString(appContext.JwtSigningKey)
		if err != nil {
			log.Fatal("failed to sign auth token: ", err.Error())
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(2 * time.Hour),
			SameSite: http.SameSiteStrictMode,
		})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully logged in as " + usernameInput))

	}
}
