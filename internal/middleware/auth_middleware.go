package middleware

import (
	"net/http"
	"strings"

	"inventory-api/internal/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// JWTMiddleware crea un middleware para validar tokens JWT
func JWTMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	authService := services.NewAuthService(db)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Obtener el header Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Authorization header required",
				})
			}

			// Verificar formato Bearer
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid authorization header format",
				})
			}

			// Extraer token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Token is required",
				})
			}

			// Validar token
			claims, err := authService.ValidateJWT(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid or expired token",
				})
			}

			// Almacenar información del usuario en el contexto
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)

			// Continuar con el siguiente handler
			return next(c)
		}
	}
}

// GetUserID obtiene el ID del usuario desde el contexto
func GetUserID(c echo.Context) (uint, bool) {
	userID, ok := c.Get("user_id").(uint)
	return userID, ok
}

// GetUserEmail obtiene el email del usuario desde el contexto
func GetUserEmail(c echo.Context) (string, bool) {
	email, ok := c.Get("user_email").(string)
	return email, ok
}

// RequireAuth es un alias más semántico para JWTMiddleware
func RequireAuth(db *gorm.DB) echo.MiddlewareFunc {
	return JWTMiddleware(db)
}
