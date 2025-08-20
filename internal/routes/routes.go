package routes

import (
	"inventory-api/internal/controllers"
	"inventory-api/internal/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// Inicializar controladores
	authController := controllers.NewAuthController(db)
	productController := controllers.NewProductController(db)

	// Grupo de rutas de autenticación (públicas)
	authGroup := e.Group("/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)

		// Rutas protegidas de auth
		authProtected := authGroup.Group("", middleware.RequireAuth(db))
		authProtected.GET("/profile", authController.Profile)
		authProtected.POST("/refresh", authController.RefreshToken)
	}

	// Grupo de rutas de productos
	productsGroup := e.Group("/products")
	{
		// Rutas públicas de productos
		productsGroup.GET("", productController.GetAllProducts)                // GET /products
		productsGroup.GET("/:id", productController.GetProductByID)            // GET /products/:id
		productsGroup.GET("/low-stock", productController.GetLowStockProducts) // GET /products/low-stock
		productsGroup.GET("/stats", productController.GetInventoryStats)       // GET /products/stats

		// Rutas protegidas de productos (requieren autenticación)
		protectedProducts := productsGroup.Group("", middleware.RequireAuth(db))
		protectedProducts.POST("", productController.CreateProduct)        // POST /products
		protectedProducts.PUT("/:id", productController.UpdateProduct)     // PUT /products/:id
		protectedProducts.DELETE("/:id", productController.DeleteProduct)  // DELETE /products/:id
		protectedProducts.PUT("/:id/stock", productController.UpdateStock) // PUT /products/:id/stock
		protectedProducts.GET("/alerts", productController.GenerateAlerts) // GET /products/alerts
	}

	// Rutas adicionales de API
	apiGroup := e.Group("/api/v1")
	{
		// Rutas de autenticación con versionado
		apiAuthGroup := apiGroup.Group("/auth")
		{
			apiAuthGroup.POST("/register", authController.Register)
			apiAuthGroup.POST("/login", authController.Login)

			apiAuthProtected := apiAuthGroup.Group("", middleware.RequireAuth(db))
			apiAuthProtected.GET("/profile", authController.Profile)
			apiAuthProtected.POST("/refresh", authController.RefreshToken)
		}

		// Rutas de productos con versionado
		apiProductsGroup := apiGroup.Group("/products")
		{
			// Públicas
			apiProductsGroup.GET("", productController.GetAllProducts)
			apiProductsGroup.GET("/:id", productController.GetProductByID)
			apiProductsGroup.GET("/low-stock", productController.GetLowStockProducts)
			apiProductsGroup.GET("/stats", productController.GetInventoryStats)

			// Protegidas
			apiProtectedProducts := apiProductsGroup.Group("", middleware.RequireAuth(db))
			apiProtectedProducts.POST("", productController.CreateProduct)
			apiProtectedProducts.PUT("/:id", productController.UpdateProduct)
			apiProtectedProducts.DELETE("/:id", productController.DeleteProduct)
			apiProtectedProducts.PUT("/:id/stock", productController.UpdateStock)
			apiProtectedProducts.GET("/alerts", productController.GenerateAlerts)
		}
	}
}
