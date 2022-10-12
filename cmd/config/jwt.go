package config

import (
	"backend-2/api/cmd/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Username  string `json:"username"`
	IsAdmin bool   `json:"is_admin"`
	jwt.StandardClaims
}

func GenerateJwt(username string, isAdmin bool) (string, error) {
	claims := &jwtCustomClaims{
		username,
		isAdmin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(GetJwtSecret()))

	return t, err
}

func GetJwtSecret() string {
	return utils.GetEnv("JWT_SECRET")
}

func GetJwtMiddleware() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(GetJwtSecret()),
	}

	return middleware.JWTWithConfig(config)
}
