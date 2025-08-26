#!/bin/bash

# Script untuk generate Swagger documentation
# Usage: ./generate-docs.sh

echo "🚀 Generating Swagger documentation..."

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "❌ Swag CLI not found. Installing..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate documentation
echo "📝 Running swag init..."
swag init -g main.go -o ./docs --parseDependency --parseInternal

if [ $? -eq 0 ]; then
    echo "✅ Swagger documentation generated successfully!"
    echo "📖 Access documentation at: http://localhost:4300/swagger/index.html"
    echo ""
    echo "Generated files:"
    echo "  - docs/docs.go"
    echo "  - docs/swagger.json" 
    echo "  - docs/swagger.yaml"
else
    echo "❌ Failed to generate Swagger documentation"
    exit 1
fi