package handlers

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"josephwest2.com/go-list/app"
	"josephwest2.com/go-list/components"
	"josephwest2.com/go-list/sqlc"
)

func TodoListPageHandler(appContext app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
		}
		tokenString := authCookie.Value
		token, err := jwt.ParseWithClaims(tokenString, &app.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return appContext.JwtSigningKey, nil
		})
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
		}
		username := token.Claims.(*app.Claims).Username
		q := sqlc.New(appContext.DBpool)
		RenderPage(components.TodoListPage(), w)
	}
}
