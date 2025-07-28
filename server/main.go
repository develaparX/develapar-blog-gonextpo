package main

// @title Develapar API
// @version 1.0
// @description REST API untuk aplikasi blog Develapar dengan fitur lengkap untuk manajemen artikel, komentar, kategori, tag, bookmark, dan like. API menggunakan standard response format dengan metadata, request tracking, rate limiting, dan comprehensive error handling.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:4300
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	NewServer().Start()

}