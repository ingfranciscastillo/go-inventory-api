package main

import (
	"fmt"
	"log"

	"inventory-api/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🔄 Starting database migration...")

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Conectar a la base de datos
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Ejecutar migraciones
	if err := db.AutoMigrate(database); err != nil {
		log.Fatal("❌ Failed to run migrations:", err)
	}

	// Crear índices
	if err := db.CreateIndexes(database); err != nil {
		log.Printf("⚠️  Warning: Failed to create some indexes: %v", err)
	}

	// Cerrar conexión
	if err := db.CloseDB(database); err != nil {
		log.Printf("⚠️  Warning: Failed to close database connection: %v", err)
	}

	fmt.Println("✅ Migration completed successfully!")
	fmt.Println("📊 Tables created:")
	fmt.Println("   - users")
	fmt.Println("   - products")
	fmt.Println("🔍 Indexes created:")
	fmt.Println("   - idx_users_email (unique)")
	fmt.Println("   - idx_products_category")
	fmt.Println("   - idx_products_quantity")
}
