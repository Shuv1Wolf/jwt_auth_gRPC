package jwtapp

import (
	"jwt_auth_gRPC/sso/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["epx"] = time.Now().Add(duration).Unix()
	claims["aid"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
