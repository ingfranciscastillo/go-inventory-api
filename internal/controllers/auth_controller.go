package controllers

import (
	"net/http"

	"inventory-api/internal/models"
	"inventory-api/internal/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// AuthController maneja los endpoints de autenticación
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController crea una nueva instancia del controlador de autenticación
func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{
		authService: services.NewAuthService(db),
	}
}

// Register maneja el registro de nuevos usuarios
// @Summary Registrar un nuevo usuario
// @Description Crea una cuenta de usuario nueva
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "Datos del usuario"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /auth/register [post]
func (ac *AuthController) Register(c echo.Context) error {
	var req models.UserRequest

	// Bind JSON request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
	}

	// Validar campos requeridos
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Email is required",
		})
	}

	if req.Password == "" || len(req.Password) < 6 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Password must be at least 6 characters long",
		})
	}

	// Registrar usuario
	user, err := ac.authService.RegisterUser(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user already exists" {
			statusCode = http.StatusConflict
		}

		return c.JSON(statusCode, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	})
}

// Login maneja la autenticación de usuarios
// @Summary Iniciar sesión
// @Description Autentica un usuario y retorna un token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserRequest true "Credenciales"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (ac *AuthController) Login(c echo.Context) error {
	var req models.UserRequest

	// Bind JSON request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
	}

	// Validar campos requeridos
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Email and password are required",
		})
	}

	// Autenticar usuario
	token, user, err := ac.authService.LoginUser(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid credentials" {
			statusCode = http.StatusUnauthorized
		}

		return c.JSON(statusCode, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}

// Profile obtiene el perfil del usuario autenticado
// @Summary Obtener perfil de usuario
// @Description Retorna la información del usuario autenticado
// @Tags auth
// @Produce json
// @Security Bearer
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Router /auth/profile [get]
func (ac *AuthController) Profile(c echo.Context) error {
	// Obtener ID del usuario del contexto (añadido por el middleware JWT)
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Invalid token",
		})
	}

	// Obtener usuario de la base de datos
	user, err := ac.authService.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "User not found",
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user.ToResponse(),
	})
}

// RefreshToken genera un nuevo token para el usuario autenticado
// @Summary Refrescar token JWT
// @Description Genera un nuevo token JWT para el usuario autenticado
// @Tags auth
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (ac *AuthController) RefreshToken(c echo.Context) error {
	// Obtener ID del usuario del contexto
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Invalid token",
		})
	}

	// Generar nuevo token
	token, err := ac.authService.RefreshToken(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to refresh token",
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Token refreshed successfully",
		"token":   token,
	})
}
