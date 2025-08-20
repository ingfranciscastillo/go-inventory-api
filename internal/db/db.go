package db

import (
	"fmt"
	"log"
	"os"

	"inventory-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB inicializa la conexión a la base de datos Neon
func InitDB() (*gorm.DB, error) {
	// Obtener variables de entorno
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// Valores por defecto
	if port == "" {
		port = "5432"
	}
	if sslmode == "" {
		sslmode = "require"
	}

	// Validar variables requeridas
	if host == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	// Construir DSN para Neon DB
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode,
	)

	log.Println("🔗 Connecting to Neon DB...")

	// Configuración de GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Establecer conexión
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verificar conexión
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configurar pool de conexiones
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("✅ Successfully connected to Neon DB")
	return db, nil
}

// AutoMigrate ejecuta las migraciones automáticamente
func AutoMigrate(db *gorm.DB) error {
	log.Println("📦 Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Migrations completed successfully")
	return nil
}

// CreateIndexes crea índices para optimizar consultas
func CreateIndexes(db *gorm.DB) error {
	log.Println("🔍 Creating database indexes...")

	// Índice para búsquedas por categoría
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_products_category ON products(category)").Error; err != nil {
		log.Printf("Warning: Failed to create category index: %v", err)
	}

	// Índice para búsquedas por stock bajo
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_products_quantity ON products(quantity)").Error; err != nil {
		log.Printf("Warning: Failed to create quantity index: %v", err)
	}

	// Índice para emails de usuarios (único)
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
		log.Printf("Warning: Failed to create email index: %v", err)
	}

	log.Println("✅ Indexes created successfully")
	return nil
}

// CloseDB cierra la conexión a la base de datos
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
