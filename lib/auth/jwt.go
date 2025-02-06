package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func GenerateJWT(userID uuid.UUID, role string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен истекает через 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func ValidateJWT(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid || claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
