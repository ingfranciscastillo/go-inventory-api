package main

import (
	"fmt"
	"log"

	"inventory-api/internal/db"
	"inventory-api/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🌱 Starting database seeding...")

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Conectar a la base de datos
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.CloseDB(database)

	// Verificar si ya hay datos
	var userCount int64
	database.Model(&models.User{}).Count(&userCount)

	var productCount int64
	database.Model(&models.Product{}).Count(&productCount)

	if userCount > 0 || productCount > 0 {
		fmt.Println("⚠️  Database already contains data. Skipping seed...")
		fmt.Printf("   Users: %d, Products: %d\n", userCount, productCount)
		return
	}

	// Crear usuarios de ejemplo
	users := []models.User{
		{
			Email:    "admin@inventory.com",
			Password: "admin123", // Se hasheará automáticamente
		},
		{
			Email:    "manager@inventory.com",
			Password: "manager123",
		},
		{
			Email:    "user@inventory.com",
			Password: "user123",
		},
	}

	fmt.Println("👤 Creating example users...")
	for _, user := range users {
		if err := database.Create(&user).Error; err != nil {
			log.Printf("❌ Failed to create user %s: %v", user.Email, err)
		} else {
			fmt.Printf("   ✅ Created user: %s\n", user.Email)
		}
	}

	// Crear productos de ejemplo
	products := []models.Product{
		{
			Name:        "Laptop Dell XPS 13",
			Description: "Laptop ultradelgada de 13 pulgadas con procesador Intel Core i7",
			Quantity:    15,
			Price:       1299.99,
			Category:    "Electronics",
		},
		{
			Name:        "iPhone 14 Pro",
			Description: "Smartphone Apple con cámara profesional de 48MP",
			Quantity:    8,
			Price:       1099.99,
			Category:    "Electronics",
		},
		{
			Name:        "Escritorio de Oficina",
			Description: "Escritorio ergonómico de madera con cajones",
			Quantity:    25,
			Price:       299.99,
			Category:    "Furniture",
		},
		{
			Name:        "Silla Ejecutiva",
			Description: "Silla ergonómica con soporte lumbar y reposabrazos",
			Quantity:    12,
			Price:       199.99,
			Category:    "Furniture",
		},
		{
			Name:        "Monitor 4K Samsung",
			Description: "Monitor de 27 pulgadas con resolución 4K UHD",
			Quantity:    20,
			Price:       399.99,
			Category:    "Electronics",
		},
		{
			Name:        "Teclado Mecánico",
			Description: "Teclado mecánico RGB para gaming con switches Cherry MX",
			Quantity:    30,
			Price:       129.99,
			Category:    "Electronics",
		},
		{
			Name:        "Mouse Inalámbrico",
			Description: "Mouse ergonómico inalámbrico con sensor óptico",
			Quantity:    45,
			Price:       59.99,
			Category:    "Electronics",
		},
		{
			Name:        "Lámpara LED",
			Description: "Lámpara de escritorio LED con control táctil",
			Quantity:    18,
			Price:       79.99,
			Category:    "Lighting",
		},
		{
			Name:        "Cafetera Automática",
			Description: "Cafetera programable con molinillo integrado",
			Quantity:    10,
			Price:       249.99,
			Category:    "Appliances",
		},
		{
			Name:        "Auriculares Bluetooth",
			Description: "Auriculares inalámbricos con cancelación de ruido",
			Quantity:    35,
			Price:       199.99,
			Category:    "Electronics",
		},
		// Productos con stock bajo para testing
		{
			Name:        "Tablet iPad Pro",
			Description: "Tablet profesional con pantalla Liquid Retina",
			Quantity:    3, // Stock bajo
			Price:       799.99,
			Category:    "Electronics",
		},
		{
			Name:        "Impresora Láser",
			Description: "Impresora láser multifunción para oficina",
			Quantity:    2, // Stock crítico
			Price:       349.99,
			Category:    "Office Equipment",
		},
		{
			Name:        "Webcam HD",
			Description: "Cámara web Full HD para videoconferencias",
			Quantity:    4, // Stock bajo
			Price:       89.99,
			Category:    "Electronics",
		},
		{
			Name:        "Disco Duro SSD",
			Description: "Disco sólido de 1TB con interfaz SATA III",
			Quantity:    1, // Stock crítico
			Price:       149.99,
			Category:    "Electronics",
		},
		{
			Name:        "Router WiFi 6",
			Description: "Router inalámbrico de alta velocidad WiFi 6",
			Quantity:    0, // Sin stock
			Price:       179.99,
			Category:    "Electronics",
		},
	}

	fmt.Println("📦 Creating example products...")
	for _, product := range products {
		if err := database.Create(&product).Error; err != nil {
			log.Printf("❌ Failed to create product %s: %v", product.Name, err)
		} else {
			status := "normal"
			if product.Quantity == 0 {
				status = "out of stock"
			} else if product.Quantity <= 2 {
				status = "critical"
			} else if product.Quantity <= 5 {
				status = "low"
			}
			fmt.Printf("   ✅ Created product: %s (Stock: %d - %s)\n", product.Name, product.Quantity, status)
		}
	}

	// Mostrar resumen
	database.Model(&models.User{}).Count(&userCount)
	database.Model(&models.Product{}).Count(&productCount)

	fmt.Println("\n📊 Seeding Summary:")
	fmt.Printf("   👤 Users created: %d\n", userCount)
	fmt.Printf("   📦 Products created: %d\n", productCount)

	// Mostrar estadísticas de stock
	var lowStockCount int64
	database.Model(&models.Product{}).Where("quantity < ?", 5).Count(&lowStockCount)

	var outOfStockCount int64
	database.Model(&models.Product{}).Where("quantity = 0").Count(&outOfStockCount)

	fmt.Println("\n⚠️  Stock Status:")
	fmt.Printf("   📉 Low stock products (< 5): %d\n", lowStockCount)
	fmt.Printf("   🚫 Out of stock products: %d\n", outOfStockCount)

	fmt.Println("\n🔑 Test Credentials:")
	fmt.Println("   Email: admin@inventory.com | Password: admin123")
	fmt.Println("   Email: manager@inventory.com | Password: manager123")
	fmt.Println("   Email: user@inventory.com | Password: user123")

	fmt.Println("\n✅ Database seeding completed successfully!")
}
