package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"inventory-api/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// AuthService maneja la lógica de autenticación
type AuthService struct {
	db *gorm.DB
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// JWTClaims define las claims personalizadas del JWT
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// RegisterUser registra un nuevo usuario
func (as *AuthService) RegisterUser(req models.UserRequest) (*models.UserResponse, error) {
	// Verificar si el usuario ya existe
	var existingUser models.User
	if err := as.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Crear nuevo usuario
	user := models.User{
		Email:    req.Email,
		Password: req.Password, // Se hasheará automáticamente en el hook BeforeCreate
	}

	if err := as.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// LoginUser autentica un usuario y genera un JWT
func (as *AuthService) LoginUser(req models.UserRequest) (string, *models.UserResponse, error) {
	// Buscar usuario por email
	var user models.User
	if err := as.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid credentials")
		}
		return "", nil, fmt.Errorf("database error: %w", err)
	}

	// Verificar contraseña
	if !user.CheckPassword(req.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	// Generar JWT token
	token, err := as.GenerateJWT(&user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	response := user.ToResponse()
	return token, &response, nil
}

// GenerateJWT genera un token JWT para el usuario
func (as *AuthService) GenerateJWT(user *models.User) (string, error) {
	// Obtener la clave secreta
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not configured")
	}

	// Crear claims
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "inventory-api",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// Crear token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT valida un token JWT y retorna las claims
func (as *AuthService) ValidateJWT(tokenString string) (*JWTClaims, error) {
	// Obtener la clave secreta
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not configured")
	}

	// Parsear token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validar método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extraer claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// GetUserByID obtiene un usuario por su ID
func (as *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := as.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

// RefreshToken genera un nuevo token para un usuario autenticado
func (as *AuthService) RefreshToken(userID uint) (string, error) {
	user, err := as.GetUserByID(userID)
	if err != nil {
		return "", err
	}

	return as.GenerateJWT(user)
}
