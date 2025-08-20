package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"inventory-api/internal/models"

	"gorm.io/gorm"
)

// ProductService maneja la lógica de negocio de productos
type ProductService struct {
	db *gorm.DB
}

// NewProductService crea una nueva instancia del servicio de productos
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// CreateProduct crea un nuevo producto
func (ps *ProductService) CreateProduct(req models.ProductRequest) (*models.ProductResponse, error) {
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		Price:       req.Price,
		Category:    req.Category,
	}

	if err := ps.db.Create(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	response := product.ToResponse()
	return &response, nil
}

// GetAllProducts obtiene todos los productos
func (ps *ProductService) GetAllProducts() ([]models.ProductResponse, error) {
	var products []models.Product
	if err := ps.db.Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	var responses []models.ProductResponse
	for _, product := range products {
		responses = append(responses, product.ToResponse())
	}

	return responses, nil
}

// GetProductByID obtiene un producto por su ID
func (ps *ProductService) GetProductByID(id uint) (*models.ProductResponse, error) {
	var product models.Product
	if err := ps.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	response := product.ToResponse()
	return &response, nil
}

// UpdateProduct actualiza un producto existente
func (ps *ProductService) UpdateProduct(id uint, req models.ProductRequest) (*models.ProductResponse, error) {
	var product models.Product
	if err := ps.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	// Actualizar campos
	product.Name = req.Name
	product.Description = req.Description
	product.Quantity = req.Quantity
	product.Price = req.Price
	product.Category = req.Category

	if err := ps.db.Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	response := product.ToResponse()
	return &response, nil
}

// DeleteProduct elimina un producto (soft delete)
func (ps *ProductService) DeleteProduct(id uint) error {
	var product models.Product
	if err := ps.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return fmt.Errorf("failed to fetch product: %w", err)
	}

	if err := ps.db.Delete(&product).Error; err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// GetLowStockProducts obtiene productos con stock bajo
func (ps *ProductService) GetLowStockProducts(threshold int) ([]models.ProductResponse, error) {
	if threshold <= 0 {
		threshold = 5 // Valor por defecto
	}

	var products []models.Product
	if err := ps.db.Where("quantity < ?", threshold).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch low stock products: %w", err)
	}

	var responses []models.ProductResponse
	for _, product := range products {
		responses = append(responses, product.ToResponse())
	}

	return responses, nil
}

// GetProductsByCategory obtiene productos por categoría
func (ps *ProductService) GetProductsByCategory(category string) ([]models.ProductResponse, error) {
	var products []models.Product
	if err := ps.db.Where("category = ?", category).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products by category: %w", err)
	}

	var responses []models.ProductResponse
	for _, product := range products {
		responses = append(responses, product.ToResponse())
	}

	return responses, nil
}

// GenerateAlertsWithConcurrency genera alertas de stock bajo usando concurrencia
func (ps *ProductService) GenerateAlertsWithConcurrency(threshold int) ([]models.ProductAlert, error) {
	if threshold <= 0 {
		threshold = 5
	}

	// Obtener todos los productos
	var products []models.Product
	if err := ps.db.Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	// Canal para recibir alertas
	alertsChan := make(chan *models.ProductAlert, len(products))
	var wg sync.WaitGroup

	// Procesar productos en paralelo
	for _, product := range products {
		wg.Add(1)
		go func(p models.Product) {
			defer wg.Done()
			
			// Simular procesamiento más complejo
			time.Sleep(10 * time.Millisecond)
			
			// Generar alerta si es necesario
			if alert := p.GenerateAlert(threshold); alert != nil {
				alertsChan <- alert
			}
		}(product)
	}

	// Goroutine para cerrar el canal cuando terminen todos los workers
	go func() {
		wg.Wait()
		close(alertsChan)
	}()

	// Recolectar alertas
	var alerts []models.ProductAlert
	for alert := range alertsChan {
		alerts = append(alerts, *alert)
	}

	return alerts, nil
}

// GetInventoryStats obtiene estadísticas del inventario
func (ps *ProductService) GetInventoryStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total de productos
	var totalProducts int64
	if err := ps.db.Model(&models.Product{}).Count(&totalProducts).Error; err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}
	stats["total_products"] = totalProducts

	// Valor total del inventario
	var totalValue float64
	if err := ps.db.Model(&models.Product{}).Select("SUM(price * quantity)").Scan(&totalValue).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total value: %w", err)
	}
	stats["total_value"] = totalValue

	// Productos con stock bajo (< 5)
	var lowStockCount int64
	if err := ps.db.Model(&models.Product{}).Where("quantity < ?", 5).Count(&lowStockCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count low stock products: %w", err)
	}
	stats["low_stock_count"] = lowStockCount

	// Productos sin stock
	var outOfStockCount int64
	if err := ps.db.Model(&models.Product{}).Where("quantity = 0").Count(&outOfStockCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count out of stock products: %w", err)
	}
	stats["out_of_stock_count"] = outOfStockCount

	// Categorías disponibles
	var categories []string
	if err := ps.db.Model(&models.Product{}).Distinct("category").Pluck("category", &categories).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	stats["categories"] = categories
	stats["categories_count"] = len(categories)

	return stats, nil
}

// SearchProducts busca productos por nombre o descripción
func (ps *ProductService) SearchProducts(query string) ([]models.ProductResponse, error) {
	var products []models.Product
	searchPattern := "%" + query + "%"
	
	if err := ps.db.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	var responses []models.ProductResponse
	for _, product := range products {
		responses = append(responses, product.ToResponse())
	}

	return responses, nil
}

// UpdateStock actualiza solo el stock de un producto
func (ps *ProductService) UpdateStock(id uint, newQuantity int) (*models.ProductResponse, error) {
	var product models.Product
	if err := ps.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	product.Quantity = newQuantity
	if err := ps.db.Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	response := product.ToResponse()
	return &response, nil
}