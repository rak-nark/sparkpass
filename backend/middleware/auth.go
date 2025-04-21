package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/utils"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token no proporcionado"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.ParseJWT(tokenString)
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token inválido"})
		}

		// Guardar userID en el contexto
		claims := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["user_id"])

		return next(c)
	}
}