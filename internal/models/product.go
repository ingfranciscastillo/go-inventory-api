package models

import (
	"time"

	"gorm.io/gorm"
)

// Product representa un producto en el inventario
type Product struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"not null;index" json:"name" validate:"required,min=2,max=100"`
	Description string          `gorm:"type:text" json:"description" validate:"max=500"`
	Quantity    int             `gorm:"not null;index" json:"quantity" validate:"required,min=0"`
	Price       float64         `gorm:"not null;type:decimal(10,2)" json:"price" validate:"required,min=0"`
	Category    string          `gorm:"not null;index" json:"category" validate:"required,min=2,max=50"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// ProductRequest representa la estructura para crear/actualizar productos
type ProductRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description string  `json:"description" validate:"max=500"`
	Quantity    int     `json:"quantity" validate:"required,min=0"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Category    string  `json:"category" validate:"required,min=2,max=50"`
}

// ProductResponse representa la respuesta con información completa del producto
type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	StockStatus string    `json:"stock_status"`
}

// ProductSummary representa un resumen del producto para listas
type ProductSummary struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// ProductAlert representa una alerta de stock bajo
type ProductAlert struct {
	ProductID   uint      `json:"product_id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Quantity    int       `json:"quantity"`
	Threshold   int       `json:"threshold"`
	Severity    string    `json:"severity"`
	Message     string    `json:"message"`
	GeneratedAt time.Time `json:"generated_at"`
}

// BeforeCreate hook que se ejecuta antes de crear un producto
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	// Validaciones adicionales o transformaciones si son necesarias
	return nil
}

// BeforeUpdate hook que se ejecuta antes de actualizar un producto
func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	// Lógica adicional antes de actualizar
	return nil
}

// IsLowStock verifica si el producto tiene stock bajo
func (p *Product) IsLowStock(threshold int) bool {
	return p.Quantity < threshold
}

// GetStockStatus retorna el estado del stock
func (p *Product) GetStockStatus(lowThreshold, criticalThreshold int) string {
	switch {
	case p.Quantity == 0:
		return "out_of_stock"
	case p.Quantity <= criticalThreshold:
		return "critical"
	case p.Quantity <= lowThreshold:
		return "low"
	default:
		return "normal"
	}
}

// ToResponse convierte Product a ProductResponse
func (p *Product) ToResponse() ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Quantity:    p.Quantity,
		Price:       p.Price,
		Category:    p.Category,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		StockStatus: p.GetStockStatus(5, 2), // Umbral bajo: 5, crítico: 2
	}
}

// ToSummary convierte Product a ProductSummary
func (p *Product) ToSummary() ProductSummary {
	return ProductSummary{
		ID:       p.ID,
		Name:     p.Name,
		Quantity: p.Quantity,
		Price:    p.Price,
		Category: p.Category,
	}
}

// GenerateAlert crea una alerta para el producto si es necesario
func (p *Product) GenerateAlert(threshold int) *ProductAlert {
	if !p.IsLowStock(threshold) {
		return nil
	}

	severity := "low"
	message := "Stock below threshold"

	if p.Quantity == 0 {
		severity = "critical"
		message = "Product out of stock"
	} else if p.Quantity <= 2 {
		severity = "high"
		message = "Critical stock level"
	}

	return &ProductAlert{
		ProductID:   p.ID,
		Name:        p.Name,
		Category:    p.Category,
		Quantity:    p.Quantity,
		Threshold:   threshold,
		Severity:    severity,
		Message:     message,
		GeneratedAt: time.Now(),
	}
}

// TableName especifica el nombre de la tabla
func (Product) TableName() string {
	return "products"
}
