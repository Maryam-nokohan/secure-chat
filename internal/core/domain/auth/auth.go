package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthResult struct {
	Token     string
	UserID    string
	Username  string
	PublicKey string
}