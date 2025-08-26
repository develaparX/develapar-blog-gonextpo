# Swagger Documentation Workflow

## 📋 Overview

Proyek ini menggunakan **Swaggo** untuk auto-generate API documentation dari komentar di kode Go.

## 🏗️ Arsitektur Swagger

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Annotations   │───▶│   swag init      │───▶│  Generated Docs │
│   (Comments)    │    │   (CLI Tool)     │    │  (JSON/YAML)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                        │                       │
         ▼                        ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│ Controller      │    │ Parse & Analyze  │    │ Swagger UI      │
│ Functions       │    │ Go Code          │    │ /swagger/       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## 🔧 Setup yang Sudah Ada

### 1. Dependencies (go.mod)

```go
github.com/swaggo/swag v1.16.4
github.com/swaggo/gin-swagger v1.6.0
github.com/swaggo/files v1.0.1
```

### 2. Main Configuration (main.go)

```go
// @title Develapar API
// @version 1.0
// @description REST API untuk aplikasi blog Develapar
// @host localhost:4300
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### 3. Server Setup (server.go)

```go
import (
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "develapar-server/docs" // Import generated docs
)

// Swagger UI endpoint
s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

## 📝 Cara Menulis Annotations

### Format Dasar

```go
// FunctionName godoc
// @Summary Short description
// @Description Detailed description
// @Tags Tag Name
// @Accept json
// @Produce json
// @Security BearerAuth (jika perlu auth)
// @Param name type dataType required "description"
// @Success 200 {object} ResponseType
// @Failure 400 {object} ErrorType
// @Router /endpoint [method]
func FunctionName(ctx *gin.Context) {
    // Implementation
}
```

### Contoh Lengkap

```go
// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with category, name, description, and image
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dto.CreateProductRequest true "Product creation details"
// @Success 201 {object} dto.APIResponse{data=object{message=string,product=dto.ProductResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Router /products [post]
func (c *ProductController) CreateProduct(ginCtx *gin.Context) {
    // Implementation
}
```

## 🚀 Generate Documentation

### 1. Manual Generate

```bash
cd server
swag init -g main.go -o ./docs --parseDependency --parseInternal
```

### 2. Using Script

```bash
cd server
./generate-docs.sh
```

### 3. Parameters Explanation

- `-g main.go`: Entry point file dengan Swagger info
- `-o ./docs`: Output directory
- `--parseDependency`: Parse external dependencies
- `--parseInternal`: Parse internal packages

## 📁 Generated Files

Setelah generate, akan terbuat:

```
docs/
├── docs.go          # Go code untuk embed documentation
├── swagger.json     # JSON specification (OpenAPI 3.0)
└── swagger.yaml     # YAML specification
```

## 🌐 Akses Documentation

### Development

```
http://localhost:4300/swagger/index.html
```

### Production

```
https://your-domain.com/swagger/index.html
```

## 🔄 Workflow Development

### 1. Tambah/Update Endpoint

```go
// Tambah annotations di controller function
// @Summary New endpoint
// @Router /new-endpoint [post]
func NewEndpoint(ctx *gin.Context) {}
```

### 2. Generate Documentation

```bash
./generate-docs.sh
```

### 3. Test Documentation

- Start server: `go run main.go`
- Open: `http://localhost:4300/swagger/index.html`
- Test endpoints langsung dari UI

### 4. Commit Changes

```bash
git add docs/
git commit -m "docs: update swagger documentation"
```

## 📋 Best Practices

### 1. Consistent Naming

```go
// ✅ Good
// @Tags Products
// @Tags Product Categories

// ❌ Bad
// @Tags product
// @Tags ProductCat
```

### 2. Detailed Descriptions

```go
// ✅ Good
// @Summary Create a new product
// @Description Create a new product with category, name, description, and image URL. Requires admin authentication.

// ❌ Bad
// @Summary Create product
// @Description Creates product
```

### 3. Proper Response Types

```go
// ✅ Good - Specific response structure
// @Success 200 {object} dto.APIResponse{data=object{products=[]dto.ProductResponse}}

// ❌ Bad - Generic response
// @Success 200 {object} interface{}
```

### 4. Complete Error Handling

```go
// ✅ Good - All possible errors
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Not found"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
```

## 🐛 Troubleshooting

### 1. Generate Errors

```bash
# Check syntax
swag init --parseDependency --parseInternal

# Common issues:
# - Missing struct tags
# - Incorrect import paths
# - Circular dependencies
```

### 2. Missing Endpoints

```bash
# Pastikan:
# - Annotations format benar
# - Function di-export (huruf kapital)
# - Router terdaftar di server.go
```

### 3. Struct Not Found

```bash
# Pastikan:
# - DTO struct di-export
# - Import path benar
# - Run dengan --parseDependency
```

## 📚 Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Gin Swagger Integration](https://github.com/swaggo/gin-swagger)

## 🔄 Auto-Generation (Optional)

Untuk auto-generate saat file berubah:

```bash
# Install air untuk hot reload
go install github.com/cosmtrek/air@latest

# Buat .air.toml dengan post-command
# post_cmd = ["./generate-docs.sh"]
```
