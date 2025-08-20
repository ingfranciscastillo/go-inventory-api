package controllers

import (
	"net/http"
	"strconv"

	"inventory-api/internal/models"
	"inventory-api/internal/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ProductController maneja los endpoints de productos
type ProductController struct {
	productService *services.ProductService
}

// NewProductController crea una nueva instancia del controlador de productos
func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{
		productService: services.NewProductService(db),
	}
}

// CreateProduct maneja la creación de nuevos productos
// @Summary Crear un nuevo producto
// @Description Crea un producto en el inventario
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param product body models.ProductRequest true "Datos del producto"
// @Success 201 {object} models.ProductResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /products [post]
func (pc *ProductController) CreateProduct(c echo.Context) error {
	var req models.ProductRequest

	// Bind JSON request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
	}

	// Validar campos requeridos
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Product name is required",
		})
	}

	if req.Price < 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Price cannot be negative",
		})
	}

	if req.Quantity < 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Quantity cannot be negative",
		})
	}

	if req.Category == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Category is required",
		})
	}

	// Crear producto
	product, err := pc.productService.CreateProduct(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to create product",
			"details": err.Error(),
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Product created successfully",
		"product": product,
	})
}

// GetAllProducts maneja la obtención de todos los productos
// @Summary Listar todos los productos
// @Description Obtiene una lista de todos los productos del inventario
// @Tags products
// @Produce json
// @Success 200 {array} models.ProductResponse
// @Failure 500 {object} map[string]interface{}
// @Router /products [get]
func (pc *ProductController) GetAllProducts(c echo.Context) error {
	// Verificar si hay parámetro de búsqueda
	search := c.QueryParam("search")
	if search != "" {
		products, err := pc.productService.SearchProducts(search)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   "Failed to search products",
				"details": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"products": products,
			"total":    len(products),
		})
	}

	// Verificar si hay filtro por categoría
	category := c.QueryParam("category")
	if category != "" {
		products, err := pc.productService.GetProductsByCategory(category)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   "Failed to get products by category",
				"details": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"products": products,
			"total":    len(products),
		})
	}

	// Obtener todos los productos
	products, err := pc.productService.GetAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to fetch products",
			"details": err.Error(),
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
		"total":    len(products),
	})
}

// GetProductByID maneja la obtención de un producto por ID
// @Summary Obtener producto por ID
// @Description Obtiene los detalles de un producto específico
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.ProductResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id} [get]
func (pc *ProductController) GetProductByID(c echo.Context) error {
	// Obtener ID del parámetro URL
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid product ID",
		})
	}

	// Obtener producto
	product, err := pc.productService.GetProductByID(uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Product not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to fetch product",
			"details": err.Error(),
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"product": product,
	})
}

// UpdateProduct maneja la actualización de productos
// @Summary Actualizar producto
// @Description Actualiza los datos de un producto existente
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Product ID"
// @Param product body models.ProductRequest true "Datos actualizados del producto"
// @Success 200 {object} models.ProductResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id} [put]
func (pc *ProductController) UpdateProduct(c echo.Context) error {
	// Obtener ID del parámetro URL
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid product ID",
		})
	}

	var req models.ProductRequest

	// Bind JSON request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
	}

	// Validar campos requeridos
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Product name is required",
		})
	}

	if req.Price < 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Price cannot be negative",
		})
	}

	if req.Quantity < 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Quantity cannot be negative",
		})
	}

	if req.Category == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Category is required",
		})
	}

	// Actualizar producto
	product, err := pc.productService.UpdateProduct(uint(id), req)
	if err != nil {
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Product not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to update product",
			"details": err.Error(),
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product updated successfully",
		"product": product,
	})
}

// DeleteProduct maneja la eliminación de productos
// @Summary Eliminar producto
// @Description Elimina un producto del inventario
// @Tags products
// @Security Bearer
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id} [delete]
func (pc *ProductController) DeleteProduct(c echo.Context) error {
	// Obtener ID del parámetro URL
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid product ID",
		})
	}

	// Eliminar producto
	err = pc.productService.DeleteProduct(uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Product not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to delete product",
			"details": err.Error(),
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product deleted successfully",
	})
}

// GetLowStockProducts maneja la obtención de productos con stock bajo
// @Summary Productos con stock bajo
// @Description Obtiene productos con cantidad menor al umbral especificado
// @Tags products
// @Produce json
// @Param threshold query int false "Umbral de stock bajo (default: 5)"
// @Success 200 {array} models.ProductResponse
// @Failure 500 {object} map[string]interface{}
// @Router /products/low-stock [get]
func (pc *ProductController) GetLowStockProducts(c echo.Context) error {
	// Obtener umbral de los query parameters
	threshold := 5 // Valor por defecto
	if thresholdParam := c.QueryParam("threshold"); thresholdParam != "" {
		if t, err := strconv.Atoi(thresholdParam); err == nil && t > 0 {
			threshold = t
		}
	}

	// Obtener productos con stock bajo
	products, err := pc.productService.GetLowStockProducts(threshold)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to fetch low stock products",
			"details": err.Error(),
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"products":  products,
		"total":     len(products),
		"threshold": threshold,
	})
}

// GenerateAlerts maneja la generación de alertas usando concurrencia
// @Summary Generar alertas de stock
// @Description Genera alertas de productos con stock bajo usando goroutines
// @Tags products
// @Produce json
// @Security Bearer
// @Param threshold query int false "Umbral para alertas (default: 5)"
// @Success 200 {array} models.ProductAlert
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/alerts [get]
func (pc *ProductController) GenerateAlerts(c echo.Context) error {
	// Obtener umbral de los query parameters
	threshold := 5 // Valor por defecto
	if thresholdParam := c.QueryParam("threshold"); thresholdParam != "" {
		if t, err := strconv.Atoi(thresholdParam); err == nil && t > 0 {
			threshold = t
		}
	}

	// Generar alertas con concurrencia
	alerts, err := pc.productService.GenerateAlertsWithConcurrency(threshold)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to generate alerts",
			"details": err.Error(),
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"alerts":    alerts,
		"total":     len(alerts),
		"threshold": threshold,
		"message":   "Alerts generated using concurrent processing",
	})
}

// GetInventoryStats maneja la obtención de estadísticas del inventario
// @Summary Estadísticas del inventario
// @Description Obtiene estadísticas generales del inventario
// @Tags products
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/stats [get]
func (pc *ProductController) GetInventoryStats(c echo.Context) error {
	// Obtener estadísticas
	stats, err := pc.productService.GetInventoryStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to fetch inventory statistics",
			"details": err.Error(),
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"stats": stats,
	})
}

// UpdateStock maneja la actualización solo del stock de un producto
// @Summary Actualizar stock
// @Description Actualiza únicamente la cantidad en stock de un producto
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Product ID"
// @Param stock body map[string]int true "Nueva cantidad"
// @Success 200 {object} models.ProductResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id}/stock [put]
func (pc *ProductController) UpdateStock(c echo.Context) error {
	// Obtener ID del parámetro URL
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid product ID",
		})
	}

	var req struct {
		Quantity int `json:"quantity" validate:"required,min=0"`
	}

	// Bind JSON request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
	}

	// Validar cantidad
	if req.Quantity < 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Quantity cannot be negative",
		})
	}

	// Actualizar stock
	product, err := pc.productService.UpdateStock(uint(id), req.Quantity)
	if err != nil {
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Product not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to update stock",
			"details": err.Error(),
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Stock updated successfully",
		"product": product,
	})
}
