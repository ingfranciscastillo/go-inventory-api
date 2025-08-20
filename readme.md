# Inventory Management API - Go + Echo + Neon DB

Un proyecto completo de API REST para gestión de inventario desarrollado en Go con Echo framework y Neon DB (PostgreSQL serverless).

## 📁 Estructura del Proyecto

```
inventory-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── controllers/
│   │   ├── auth_controller.go
│   │   └── product_controller.go
│   ├── middleware/
│   │   └── auth_middleware.go
│   ├── models/
│   │   ├── product.go
│   │   └── user.go
│   ├── routes/
│   │   └── routes.go
│   ├── services/
│   │   ├── auth_service.go
│   │   └── product_service.go
│   └── db/
│       ├── db.go
│       └── migrations.go
├── scripts/
│   ├── migrate.go
│   └── seed.go
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Características

- ✅ **Framework**: Echo v4
- ✅ **Base de datos**: Neon DB (PostgreSQL serverless) con GORM
- ✅ **Autenticación**: JWT para endpoints protegidos
- ✅ **Concurrencia**: Goroutines y channels para alertas de stock
- ✅ **Arquitectura modular**: Separación clara de responsabilidades
- ✅ **Docker**: Contenedores para desarrollo y producción
- ✅ **Migraciones y seeds**: Scripts para inicializar la BD
- ✅ **Documentación**: README completo con instrucciones

## 📋 Prerrequisitos

- Go
- Cuenta en [Neon DB](https://neon.tech)
- Docker (opcional)

## 🛠️ Configuración

### 1. Clonar el repositorio

```bash
git clone https://github.com/ingfranciscastillo/go-inventory-api
cd inventory-api
```

### 2. Configurar variables de entorno

Copia el archivo `.env.example` a `.env` y configura las variables:

```bash
cp .env.example .env
```

Edita `.env` con tus datos de Neon DB:

```env
# Database
DB_HOST=your-neon-host.neon.tech
DB_USER=your-username
DB_PASSWORD=your-password
DB_NAME=your-database
DB_PORT=5432
DB_SSLMODE=require

# JWT
JWT_SECRET=your-super-secret-jwt-key-here

# Server
PORT=8080
```

### 3. Instalar dependencias

```bash
go mod tidy
```

### 4. Ejecutar migraciones

```bash
go run scripts/migrate/main.go
```

### 5. Poblar con datos de ejemplo (opcional)

```bash
go run scripts/seed/main.go
```

### 6. Ejecutar la aplicación

```bash
go run cmd/server/main.go
```

La API estará disponible en `http://localhost:8080`

## 🐳 Docker

### Desarrollo con Docker Compose

```bash
# Construir e iniciar los contenedores
docker-compose up --build

# Solo iniciar (si ya están construidos)
docker-compose up

# Ejecutar en background
docker-compose up -d

# Ver logs
docker-compose logs -f api

# Parar los contenedores
docker-compose down
```

### Solo contenedor de la API

```bash
# Construir imagen
docker build -t inventory-api .

# Ejecutar contenedor
docker run -p 8080:8080 --env-file .env inventory-api
```

## 🔌 Endpoints de la API

### Autenticación

| Método | Endpoint         | Descripción       | Auth |
| ------ | ---------------- | ----------------- | ---- |
| POST   | `/auth/register` | Registrar usuario | No   |
| POST   | `/auth/login`    | Iniciar sesión    | No   |

### Productos

| Método | Endpoint              | Descripción          | Auth |
| ------ | --------------------- | -------------------- | ---- |
| GET    | `/products`           | Listar productos     | No   |
| GET    | `/products/:id`       | Obtener producto     | No   |
| POST   | `/products`           | Crear producto       | JWT  |
| PUT    | `/products/:id`       | Actualizar producto  | JWT  |
| DELETE | `/products/:id`       | Eliminar producto    | JWT  |
| GET    | `/products/low-stock` | Stock bajo           | No   |
| GET    | `/products/alerts`    | Alertas concurrentes | JWT  |

## 📝 Ejemplos de uso

### 1. Registrar usuario

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

### 2. Iniciar sesión

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

### 3. Crear producto

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Laptop Dell XPS 13",
    "description": "Laptop ultradelgada para profesionales",
    "quantity": 10,
    "price": 1299.99,
    "category": "Electronics"
  }'
```

### 4. Listar productos

```bash
curl http://localhost:8080/products
```

### 5. Productos con stock bajo

```bash
curl http://localhost:8080/products/low-stock
```

### 6. Alertas concurrentes

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/products/alerts
```

## 🏗️ Arquitectura

### Capas de la aplicación

1. **Controllers**: Manejan las peticiones HTTP
2. **Services**: Lógica de negocio
3. **Models**: Estructuras de datos y validación
4. **Middleware**: Autenticación y validaciones
5. **DB**: Conexión y configuración de base de datos

### Patrones utilizados

- **Repository Pattern**: Para acceso a datos
- **Dependency Injection**: Para desacoplamiento
- **Middleware Pattern**: Para funcionalidades transversales
- **Service Layer**: Para lógica de negocio

## 🔧 Desarrollo

### Agregar nuevas funcionalidades

1. **Modelo**: Crear en `internal/models/`
2. **Servicio**: Lógica en `internal/services/`
3. **Controlador**: Handlers en `internal/controllers/`
4. **Rutas**: Registrar en `internal/routes/`

### Migraciones

```bash
# Crear nueva migración
go run scripts/migrate/main.go

## 📊 Monitoreo y Logs

La aplicación incluye:

- Logging estructurado con Echo middleware
- Manejo de errores HTTP estandarizado
- Métricas básicas de requests

## 🛡️ Seguridad

- Autenticación JWT
- Validación de entrada
- Rate limiting (configurable)
- CORS habilitado
- Headers de seguridad

## 🤝 Contribuir

1. Fork el proyecto
2. Crear branch (`git checkout -b feature/nueva-caracteristica`)
3. Commit (`git commit -am 'Agregar nueva característica'`)
4. Push (`git push origin feature/nueva-caracteristica`)
5. Crear Pull Request

## 📄 Licencia

MIT License - ver archivo [LICENSE](LICENSE) para detalles.

**¡Feliz coding! 🚀**
```
