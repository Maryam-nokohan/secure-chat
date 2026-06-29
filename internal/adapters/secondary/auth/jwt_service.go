package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/auth"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type JWTService struct {
	secret string
}

func NewJWTService(secret string) ports.TokenService {
	pkg.LogInfo("Initializing JWTService...")
	return &JWTService{
		secret: secret,
	}
}

func (j *JWTService) Generate(userID, username string) (string, error) {

	now := time.Now()

	claims := auth.Claims{
		UserID:   userID,
		Username: username,

		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    "secure-chat",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(j.secret))
}

func (j *JWTService) Validate(tokenStr string) (*auth.Claims, error) {

	token , err := jwt.ParseWithClaims(
		tokenStr,
		&auth.Claims{},
		func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					t.Header["alg"],
				)
			}

			return []byte(j.secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*auth.Claims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}