package app

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppContext struct {
	DBpool        *pgxpool.Pool
	JwtSigningKey []byte
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}