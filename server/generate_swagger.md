# Cara Generate Swagger Documentation

## 1. Install Swag CLI (jika belum ada)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## 2. Generate Swagger Documentation

Jalankan command ini di folder `server/`:

```bash
# Generate dari root folder server
swag init

# Atau dengan parameter lebih spesifik
swag init -g main.go -o ./docs --parseDependency --parseInternal
```

## 3. Parameter Penjelasan

- `-g main.go`: File utama yang berisi info Swagger
- `-o ./docs`: Output directory untuk generated files
- `--parseDependency`: Parse dependencies untuk struct definitions
- `--parseInternal`: Parse internal packages

## 4. Files yang di-generate

Setelah run command, akan terbuat:

- `docs/docs.go` - Go code untuk embed docs
- `docs/swagger.json` - JSON specification
- `docs/swagger.yaml` - YAML specification

## 5. Akses Swagger UI

Setelah server running, buka:

```
http://localhost:4300/swagger/index.html
```

## 6. Auto-regenerate saat development

Buat script untuk auto-generate:

```bash
#!/bin/bash
# save as generate-docs.sh
cd server
swag init -g main.go -o ./docs --parseDependency --parseInternal
echo "Swagger documentation generated successfully!"
```

## 7. Troubleshooting

Jika ada error:

- Pastikan semua import statements benar
- Check syntax annotations Swagger
- Pastikan struct DTO sudah di-export (huruf kapital)
- Run `go mod tidy` sebelum generate
