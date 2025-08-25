# Product API Documentation (With Pagination)

## Overview

API untuk mengelola produk afiliasi dengan kategori dan link afiliasi yang terhubung. **Semua endpoint GET list menggunakan pagination by default**.

## Features

- ✅ **DTO Pattern**: Request dan Response menggunakan DTO
- ✅ **Affiliate Links**: Setiap product include affiliate links saat GET detail
- ✅ **UUIDv7**: Menggunakan UUIDv7 untuk semua ID
- ✅ **Pagination**: Semua GET list menggunakan pagination by default
- ✅ **Validation**: Comprehensive validation untuk semua input
- ✅ **Error Handling**: Proper error handling dengan context timeout

## Pagination Parameters

Semua endpoint GET list mendukung query parameters berikut:

- `page` (optional): Page number, default = 1
- `limit` (optional): Items per page, default = 10, max = 100

## Pagination Response Format

```json
{
  "success": true,
  "data": {
    "items": [...], // Array of items
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_items": 50,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

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

#### 2. Get All Product Categories (With Pagination)

```http
GET /api/v1/product-categories?page=1&limit=10
```

**Response:**

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "01234567-89ab-cdef-0123-456789abcdef",
        "name": "Laptop Gaming",
        "slug": "laptop-gaming",
        "description": "Kategori untuk laptop gaming terbaik",
        "created_at": "2025-01-20T10:00:00Z",
        "updated_at": "2025-01-20T10:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_items": 25,
      "total_pages": 3,
      "has_next": true,
      "has_prev": false
    }
  }
}
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

#### 2. Get All Products (List View - No Affiliate Links, With Pagination)

```http
GET /api/v1/products?page=1&limit=10
```

**Response:**

```json
{
  "success": true,
  "data": {
    "items": [
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
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_items": 50,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
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

#### 4. Get Products by Category (With Affiliate Links & Pagination)

```http
GET /api/v1/products/category/{category_id}?page=1&limit=10
```

**Response:**

```json
{
  "success": true,
  "data": {
    "items": [
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
        "affiliate_links": [
          {
            "id": "01234567-89ab-cdef-0123-456789abcdef",
            "platform_name": "Tokopedia",
            "url": "https://tokopedia.com/product/asus-rog-strix-g15",
            "created_at": "2025-01-20T10:00:00Z",
            "updated_at": "2025-01-20T10:00:00Z"
          }
        ],
        "created_at": "2025-01-20T10:00:00Z",
        "updated_at": "2025-01-20T10:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_items": 15,
      "total_pages": 2,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

#### 5. Update Product

```http
PUT /api/v1/products/{id}
Authorization: Bearer <admin_token>
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

#### 2. Get Affiliate Links by Product ID

```http
GET /api/v1/products/{product_id}/affiliate
Authorization: Bearer <admin_token>
```

#### 3. Update Affiliate Link

```http
PUT /api/v1/products/{product_id}/affiliate/{affiliate_id}
Authorization: Bearer <admin_token>
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

#### 3. Get Products by Article ID (With Affiliate Links & Pagination)

```http
GET /api/v1/products/article/{article_id}?page=1&limit=10
```

**Response:**

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "01234567-89ab-cdef-0123-456789abcdef",
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
          }
        ],
        "created_at": "2025-01-20T10:00:00Z",
        "updated_at": "2025-01-20T10:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_items": 5,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

## Key Changes from Previous Version

### 🔄 **Pagination by Default**

- **All GET list endpoints** sekarang menggunakan pagination
- **Default values**: page=1, limit=10
- **Maximum limit**: 100 items per page
- **Query parameters**: `?page=1&limit=10`

### 📊 **Response Format Changes**

- **List endpoints**: Response wrapped dalam pagination object
- **Detail endpoints**: Tetap menggunakan format lama
- **Consistent structure**: Semua pagination menggunakan format yang sama

### 🚀 **Performance Benefits**

- **Reduced memory usage**: Tidak load semua data sekaligus
- **Faster response time**: Pagination mengurangi data transfer
- **Better scalability**: Dapat handle dataset yang besar

### 💡 **Usage Examples**

#### Frontend Integration

```javascript
// Get products with pagination
const getProducts = async (page = 1, limit = 10) => {
  const response = await fetch(`/api/v1/products?page=${page}&limit=${limit}`);
  const data = await response.json();

  return {
    products: data.data.items,
    pagination: data.data.pagination,
  };
};

// Get product detail with affiliate links
const getProductDetail = async (productId) => {
  const response = await fetch(`/api/v1/products/${productId}`);
  const data = await response.json();

  return data.data.product; // Includes affiliate_links
};

// Get products by category with pagination
const getProductsByCategory = async (categoryId, page = 1, limit = 10) => {
  const response = await fetch(
    `/api/v1/products/category/${categoryId}?page=${page}&limit=${limit}`
  );
  const data = await response.json();

  return {
    products: data.data.items, // Each product includes affiliate_links
    pagination: data.data.pagination,
  };
};
```

#### Pagination Navigation

```javascript
const PaginationComponent = ({ pagination, onPageChange }) => {
  return (
    <div className="pagination">
      {pagination.has_prev && (
        <button onClick={() => onPageChange(pagination.current_page - 1)}>
          Previous
        </button>
      )}

      <span>
        Page {pagination.current_page} of {pagination.total_pages}
      </span>

      {pagination.has_next && (
        <button onClick={() => onPageChange(pagination.current_page + 1)}>
          Next
        </button>
      )}
    </div>
  );
};
```

## Migration Guide

### From Non-Pagination to Pagination

#### Before:

```javascript
// Old way - no pagination
const products = await fetch("/api/v1/products").then((res) => res.json());
console.log(products.data.products); // Array of products
```

#### After:

```javascript
// New way - with pagination
const response = await fetch("/api/v1/products?page=1&limit=10").then((res) =>
  res.json()
);
console.log(response.data.items); // Array of products
console.log(response.data.pagination); // Pagination info
```

### Response Structure Changes

#### Before:

```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": {
    "products": [...]
  }
}
```

#### After:

```json
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_items": 50,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

Implementasi pagination sudah **production-ready** dan siap digunakan! 🎉
