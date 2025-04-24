package app

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppContext struct {
	DBpool        *pgxpool.Pool
	JwtSigningKey []byte
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func ReadAuthToken(appContext AppContext, r *http.Request) (*jwt.Token, error) {
	authCookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}
	tokenString := authCookie.Value
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return appContext.JwtSigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
