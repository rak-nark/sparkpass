package middleware

import (
	
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rak-nark/sparkpass/models"
	"github.com/rak-nark/sparkpass/utils"
	"gorm.io/gorm"
)

// En tu middleware/auth.go
func AuthMiddleware(db *gorm.DB) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // 1. Extraer token
            authHeader := c.Request().Header.Get("Authorization")
            if authHeader == "" {
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Token no proporcionado"})
            }

            tokenString := strings.TrimPrefix(authHeader, "Bearer ")

            // 2. Parsear token
            token, err := utils.ParseJWT(tokenString)
            if err != nil || !token.Valid {
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Token inválido"})
            }

            // 3. Extraer claims
            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Claims inválidos"})
            }

            // 4. Obtener userID de los claims
            claimUserID, ok := claims["user_id"]
            if !ok {
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Claim user_id no encontrado"})
            }

            // Convertir userID a uint
            var userID uint
            switch v := claimUserID.(type) {
            case float64:
                userID = uint(v)
            case string:
                parsedID, err := strconv.ParseUint(v, 10, 64)
                if err != nil {
                    return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Formato de user_id inválido"})
                }
                userID = uint(parsedID)
            default:
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Tipo de user_id no soportado"})
            }

            // 5. Buscar usuario en la base de datos
            var user models.User
            if err := db.First(&user, userID).Error; err != nil {
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Usuario no encontrado"})
            }

            // 6. Guardar en contexto
            c.Set("user", &user)
            c.Set("userID", user.ID)
            c.Set("db", db)

            return next(c)
        }
    }
}