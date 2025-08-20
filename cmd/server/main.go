package main

import (
	"fmt"
	"log"
	"os"

	"inventory-api/internal/db"
	"inventory-api/internal/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Inicializar conexión a base de datos
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Ejecutar migraciones automáticamente
	if err := db.AutoMigrate(database); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Crear instancia de Echo
	e := echo.New()

	// Middleware globales
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Configurar middleware de rate limiting
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Headers de seguridad
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	// Configurar rutas
	routes.SetupRoutes(e, database)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":  "healthy",
			"service": "inventory-api",
			"version": "1.0.0",
		})
	})

	// Puerto del servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Mensaje de inicio
	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Println("📚 Endpoints available:")
	fmt.Println("   GET  /health")
	fmt.Println("   POST /auth/register")
	fmt.Println("   POST /auth/login")
	fmt.Println("   GET  /products")
	fmt.Println("   POST /products (Auth required)")
	fmt.Println("   GET  /products/:id")
	fmt.Println("   PUT  /products/:id (Auth required)")
	fmt.Println("   DELETE /products/:id (Auth required)")
	fmt.Println("   GET  /products/low-stock")
	fmt.Println("   GET  /products/alerts (Auth required)")

	// Iniciar servidor
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
