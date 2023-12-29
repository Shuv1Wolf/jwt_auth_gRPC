package jwtapp

import (
	"jwt_auth_gRPC/sso/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["aid"] = app.ID
	secret := []byte(app.Secret) // Преобразуем app.Secret в тип []byte
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
