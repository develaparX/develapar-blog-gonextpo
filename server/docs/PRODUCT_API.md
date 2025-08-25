# Product API Documentation

## Overview

API untuk mengelola produk afiliasi dengan kategori dan link afiliasi yang terhubung.

## Features

- ✅ **DTO Pattern**: Request dan Response menggunakan DTO
- ✅ **Affiliate Links**: Setiap product include affiliate links saat GET
- ✅ **UUIDv7**: Menggunakan UUIDv7 untuk semua ID
- ✅ **Validation**: Comprehensive validation untuk semua input
- ✅ **Error Handling**: Proper error handling dengan context timeout

## Endpoints

### Product Categories

#### 1. Create Product Category

```http
POST /api/v1/product-categories
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "Laptop Gaming",
  "slug": "laptop-gaming", // optional, auto-generated from name
  "description": "Kategori untuk laptop gaming terbaik"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Product Category created successfully",
  "data": {
    "product_category": {
      "id": "01234567-89ab-cdef-0123-456789abcdef",
      "name": "Laptop Gaming",
      "slug": "laptop-gaming",
      "description": "Kategori untuk laptop gaming terbaik",
      "created_at": "2025-01-20T10:00:00Z",
      "updated_at": "2025-01-20T10:00:00Z"
    }
  }
}
```

#### 2. Get All Product Categories

```http
GET /api/v1/product-categories
```

#### 3. Get Product Category by ID

```http
GET /api/v1/product-categories/{id}
```

#### 4. Get Product Category by Slug

```http
GET /api/v1/product-categories/s/{slug}
```

#### 5. Update Product Category

```http
PUT /api/v1/product-categories/{id}
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "Gaming Laptops Updated",
  "description": "Updated description"
}
```

#### 6. Delete Product Category

```http
DELETE /api/v1/product-categories/{id}
Authorization: Bearer <admin_token>
```

### Products

#### 1. Create Product

```http
POST /api/v1/products
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "product_category_id": "01234567-89ab-cdef-0123-456789abcdef",
  "name": "ASUS ROG Strix G15",
  "description": "Gaming laptop dengan performa tinggi",
  "image_url": "https://example.com/image.jpg",
  "is_active": true
}
```

**Response:**

```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "product": {
      "id": "01234567-89ab-cdef-0123-456789abcdef",
      "product_category_id": "01234567-89ab-cdef-0123-456789abcdef",
      "product_category": {
        "id": "01234567-89ab-cdef-0123-456789abcdef",
        "name": "Laptop Gaming",
        "slug": "laptop-gaming",
        "description": "Kategori untuk laptop gaming terbaik",
        "created_at": "2025-01-20T10:00:00Z",
        "updated_at": "2025-01-20T10:00:00Z"
      },
      "name": "ASUS ROG Strix G15",
      "description": "Gaming laptop dengan performa tinggi",
      "image_url": "https://example.com/image.jpg",
      "is_active": true,
      "affiliate_links": [],
      "created_at": "2025-01-20T10:00:00Z",
      "updated_at": "2025-01-20T10:00:00Z"
    }
  }
}
```

#### 2. Get All Products (List View - No Affiliate Links)

```http
GET /api/v1/products
```

**Response:**

```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": {
    "products": [
      {
        "id": "01234567-89ab-cdef-0123-456789abcdef",
        "product_category_id": "01234567-89ab-cdef-0123-456789abcdef",
        "product_category": {
          "id": "01234567-89ab-cdef-0123-456789abcdef",
          "name": "Laptop Gaming",
          "slug": "laptop-gaming",
          "description": "Kategori untuk laptop gaming terbaik",
          "created_at": "2025-01-20T10:00:00Z",
          "updated_at": "2025-01-20T10:00:00Z"
        },
        "name": "ASUS ROG Strix G15",
        "description": "Gaming laptop dengan performa tinggi",
        "image_url": "https://example.com/image.jpg",
        "is_active": true,
        "created_at": "2025-01-20T10:00:00Z",
        "updated_at": "2025-01-20T10:00:00Z"
      }
    ]
  }
}
```

#### 3. Get Product by ID (Detail View - With Affiliate Links)

```http
GET /api/v1/products/{id}
```

**Response:**

```json
{
  "success": true,
  "message": "Product retrieved successfully",
  "data": {
    "product": {
      "id": "01234567-89ab-cdef-0123-456789abcdef",
      "product_category_id": "01234567-89ab-cdef-0123-456789abcdef",
      "product_category": {
        "id": "01234567-89ab-cdef-0123-456789abcdef",
        "name": "Laptop Gaming",
        "slug": "laptop-gaming",
        "description": "Kategori untuk laptop gaming terbaik",
        "created_at": "2025-01-20T10:00:00Z",
        "updated_at": "2025-01-20T10:00:00Z"
      },
      "name": "ASUS ROG Strix G15",
      "description": "Gaming laptop dengan performa tinggi",
      "image_url": "https://example.com/image.jpg",
      "is_active": true,
      "affiliate_links": [
        {
          "id": "01234567-89ab-cdef-0123-456789abcdef",
          "platform_name": "Tokopedia",
          "url": "https://tokopedia.com/product/asus-rog-strix-g15",
          "created_at": "2025-01-20T10:00:00Z",
          "updated_at": "2025-01-20T10:00:00Z"
        },
        {
          "id": "01234567-89ab-cdef-0123-456789abcdef",
          "platform_name": "Shopee",
          "url": "https://shopee.co.id/product/asus-rog-strix-g15",
          "created_at": "2025-01-20T10:00:00Z",
          "updated_at": "2025-01-20T10:00:00Z"
        }
      ],
      "created_at": "2025-01-20T10:00:00Z",
      "updated_at": "2025-01-20T10:00:00Z"
    }
  }
}
```

#### 4. Get Products by Category (With Affiliate Links)

```http
GET /api/v1/products/category/{category_id}
```

#### 5. Update Product

```http
PUT /api/v1/products/{id}
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "ASUS ROG Strix G15 Updated",
  "description": "Updated description",
  "is_active": false
}
```

#### 6. Delete Product

```http
DELETE /api/v1/products/{id}
Authorization: Bearer <admin_token>
```

### Product Affiliate Links

#### 1. Create Affiliate Link

```http
POST /api/v1/products/{product_id}/affiliate
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "platform_name": "Tokopedia",
  "url": "https://tokopedia.com/product/asus-rog-strix-g15"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Product affiliate link created successfully",
  "data": {
    "product_affiliate_link": {
      "id": "01234567-89ab-cdef-0123-456789abcdef",
      "platform_name": "Tokopedia",
      "url": "https://tokopedia.com/product/asus-rog-strix-g15",
      "created_at": "2025-01-20T10:00:00Z",
      "updated_at": "2025-01-20T10:00:00Z"
    }
  }
}
```

#### 2. Get Affiliate Links by Product ID

```http
GET /api/v1/products/{product_id}/affiliate
Authorization: Bearer <admin_token>
```

#### 3. Update Affiliate Link

```http
PUT /api/v1/products/{product_id}/affiliate/{affiliate_id}
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "platform_name": "Tokopedia Official",
  "url": "https://tokopedia.com/official/asus-rog-strix-g15"
}
```

#### 4. Delete Affiliate Link

```http
DELETE /api/v1/products/{product_id}/affiliate/{affiliate_id}
Authorization: Bearer <admin_token>
```

### Article-Product Relations

#### 1. Add Product to Article

```http
POST /api/v1/products/{product_id}/article/{article_id}
Authorization: Bearer <admin_token>
```

#### 2. Remove Product from Article

```http
DELETE /api/v1/products/{product_id}/article/{article_id}
Authorization: Bearer <admin_token>
```

#### 3. Get Products by Article ID (With Affiliate Links)

```http
GET /api/v1/products/article/{article_id}
```

## Key Features

### 1. DTO Pattern

- **Request DTOs**: Validation dan type safety untuk input
- **Response DTOs**: Consistent response format
- **Separation**: Request dan response terpisah untuk flexibility

### 2. Affiliate Links Integration

- **Get Product by ID**: Include affiliate links
- **Get Products by Category**: Include affiliate links untuk semua products
- **Get Products by Article**: Include affiliate links untuk products dalam artikel
- **List View**: Tidak include affiliate links untuk performance

### 3. Performance Optimization

- **List endpoints**: Tidak load affiliate links untuk performance
- **Detail endpoints**: Load affiliate links untuk informasi lengkap
- **Batch loading**: Efficient loading untuk multiple products

### 4. Validation

- **URL validation**: Affiliate links harus valid URL
- **Slug validation**: Auto-generate dan validate slug format
- **Required fields**: Proper validation untuk semua required fields

### 5. Error Handling

- **Context timeout**: 15 detik timeout untuk semua operations
- **Proper error codes**: HTTP status codes yang sesuai
- **Error messages**: Descriptive error messages

## Usage Examples

### Frontend Integration

```javascript
// Get product with affiliate links
const product = await fetch(
  "/api/v1/products/01234567-89ab-cdef-0123-456789abcdef"
).then((res) => res.json());

// Display affiliate links
product.data.product.affiliate_links.forEach((link) => {
  console.log(`${link.platform_name}: ${link.url}`);
});

// Get products for article (with affiliate links)
const articleProducts = await fetch(
  "/api/v1/products/article/01234567-89ab-cdef-0123-456789abcdef"
).then((res) => res.json());
```

### Admin Operations

```javascript
// Create product with category
const newProduct = await fetch("/api/v1/products", {
  method: "POST",
  headers: {
    Authorization: "Bearer " + adminToken,
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    product_category_id: categoryId,
    name: "New Product",
    description: "Product description",
    is_active: true,
  }),
});

// Add affiliate link
const affiliateLink = await fetch(`/api/v1/products/${productId}/affiliate`, {
  method: "POST",
  headers: {
    Authorization: "Bearer " + adminToken,
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    platform_name: "Tokopedia",
    url: "https://tokopedia.com/product-link",
  }),
});
```
