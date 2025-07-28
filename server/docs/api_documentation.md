# Develapar API Documentation

## Overview

Develapar API adalah REST API untuk aplikasi blog yang menyediakan fitur lengkap untuk manajemen artikel, komentar, kategori, tag, bookmark, dan like.

## Base URL

```
http://localhost:4300/api/v1
```

## Authentication

API menggunakan Bearer Token authentication untuk endpoint yang memerlukan autentikasi.

```
Authorization: Bearer <your-jwt-token>
```

## Response Format

### Standard Response Structure

Semua response API menggunakan format standar berikut:

```json
{
  "success": true,
  "data": {
    // Response data varies by endpoint
  },
  "error": null,
  "pagination": {
    // Only present for paginated responses
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10,
    "has_next": true,
    "has_prev": false,
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
  },
  "meta": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "processing_time_ms": 15000000,
    "version": "1.0.0",
    "timestamp": "2025-07-24T20:43:16.123456789+07:00"
  }
}
```

### Error Response Structure

Untuk error responses:

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "field_name": "error description"
    },
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2025-07-24T20:43:16.123456789+07:00"
  },
  "meta": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "processing_time_ms": 15000000,
    "version": "1.0.0",
    "timestamp": "2025-07-24T20:43:16.123456789+07:00"
  }
}
```

## Response Metadata

### Meta Object

- `request_id`: Unique identifier untuk tracking request
- `processing_time_ms`: Waktu pemrosesan request dalam nanoseconds
- `version`: Versi API
- `timestamp`: Waktu response dibuat

### Pagination Object

- `page`: Halaman saat ini (1-based)
- `limit`: Jumlah item per halaman
- `total`: Total jumlah item
- `total_pages`: Total jumlah halaman
- `has_next`: Apakah ada halaman selanjutnya
- `has_prev`: Apakah ada halaman sebelumnya
- `request_id`: Request ID untuk tracking

## HTTP Status Codes

- `200 OK`: Request berhasil
- `201 Created`: Resource berhasil dibuat
- `400 Bad Request`: Request tidak valid
- `401 Unauthorized`: Authentication diperlukan
- `403 Forbidden`: Access denied
- `404 Not Found`: Resource tidak ditemukan
- `408 Request Timeout`: Request timeout
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

## Rate Limiting

API menggunakan rate limiting:

- Anonymous users: 50 requests per minute
- Authenticated users: 200 requests per minute

Rate limit headers:

- `X-RateLimit-Limit`: Limit maksimum
- `X-RateLimit-Remaining`: Sisa request
- `Retry-After`: Waktu tunggu jika limit terlampaui

## Request Tracking

Setiap request memiliki unique `request_id` yang dapat digunakan untuk tracking dan debugging. Request ID akan muncul di:

- Response header `X-Request-ID`
- Response body di `meta.request_id`
- Error response di `error.request_id`

## Endpoints

### Authentication

#### POST /auth/register

Register user baru

**Request Body:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "message": "User registered successfully",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2025-07-24T20:43:16+07:00",
      "updated_at": "2025-07-24T20:43:16+07:00"
    }
  },
  "meta": {...}
}
```

#### POST /auth/login

Login user

**Request Body:**

```json
{
  "identifier": "john@example.com",
  "password": "password123"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "message": "Login successful",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "meta": {...}
}
```

#### POST /auth/refresh

Refresh access token

**Request Body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Articles

#### GET /articles

Get all articles

**Response:**

```json
{
  "success": true,
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "Sample Article",
        "slug": "sample-article",
        "content": "Article content...",
        "user": {
          "id": 1,
          "name": "John Doe",
          "email": "john@example.com"
        },
        "category": {
          "id": 1,
          "name": "Technology"
        },
        "views": 100,
        "created_at": "2025-07-24T20:43:16+07:00",
        "updated_at": "2025-07-24T20:43:16+07:00"
      }
    ],
    "message": "Articles retrieved successfully"
  },
  "meta": {...}
}
```

#### GET /articles/paginated

Get articles with pagination

**Query Parameters:**

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10)

**Response:**

```json
{
  "success": true,
  "data": {
    "articles": [...],
    "message": "Articles retrieved successfully"
  },
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10,
    "has_next": true,
    "has_prev": false,
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
  },
  "meta": {...}
}
```

#### GET /articles/{slug}

Get article by slug

**Response:**

```json
{
  "success": true,
  "data": {
    "article": {
      "id": 1,
      "title": "Sample Article",
      "slug": "sample-article",
      "content": "Article content...",
      "user": {...},
      "category": {...},
      "views": 100,
      "created_at": "2025-07-24T20:43:16+07:00",
      "updated_at": "2025-07-24T20:43:16+07:00"
    },
    "message": "Article retrieved successfully"
  },
  "meta": {...}
}
```

#### POST /articles

Create new article (Authentication required)

**Request Body:**

```json
{
  "title": "New Article",
  "content": "Article content...",
  "category_id": 1
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "message": "Article created successfully",
    "article": {
      "id": 2,
      "title": "New Article",
      "slug": "new-article",
      "content": "Article content...",
      "user": {...},
      "category": {...},
      "views": 0,
      "created_at": "2025-07-24T20:43:16+07:00",
      "updated_at": "2025-07-24T20:43:16+07:00"
    }
  },
  "meta": {...}
}
```

#### PUT /articles/{article_id}

Update article (Authentication required)

**Request Body:**

```json
{
  "title": "Updated Article",
  "content": "Updated content...",
  "category_id": 1
}
```

#### DELETE /articles/{article_id}

Delete article (Authentication required)

**Response:**

```json
{
  "success": true,
  "data": {
    "message": "Article deleted successfully"
  },
  "meta": {...}
}
```

### Categories

#### GET /categories

Get all categories

#### GET /categories/{category_id}

Get category by ID

#### POST /categories

Create new category (Authentication required)

#### PUT /categories/{category_id}

Update category (Authentication required)

#### DELETE /categories/{category_id}

Delete category (Authentication required)

### Tags

#### GET /tags

Get all tags

#### GET /tags/{tag_id}

Get tag by ID

#### POST /tags

Create new tag (Authentication required)

#### PUT /tags/{tag_id}

Update tag (Authentication required)

#### DELETE /tags/{tag_id}

Delete tag (Authentication required)

### Comments

#### GET /comments/article/{article_id}

Get comments for an article

#### GET /comments/user/{user_id}

Get comments by user

#### POST /comments

Create new comment (Authentication required)

#### PUT /comments/{comment_id}

Update comment (Authentication required)

#### DELETE /comments/{comment_id}

Delete comment (Authentication required)

### Bookmarks

#### GET /bookmarks/{user_id}

Get user's bookmarks

#### POST /bookmarks

Create bookmark (Authentication required)

#### DELETE /bookmarks

Remove bookmark (Authentication required)

#### GET /bookmarks/check

Check if article is bookmarked (Authentication required)

### Likes

#### GET /likes/article/{article_id}

Get likes for an article

#### GET /likes/user/{user_id}

Get likes by user

#### POST /likes

Add like (Authentication required)

#### DELETE /likes

Remove like (Authentication required)

#### GET /likes/check

Check if article is liked (Authentication required)

### Article Tags

#### GET /article-tags/{article_id}

Get tags for an article

#### POST /article-tags/{article_id}

Assign tags to article (Authentication required)

#### DELETE /article-tags/{article_id}/{tag_id}

Remove tag from article (Authentication required)

#### GET /tags/{tag_id}/articles

Get articles by tag

### Health & Monitoring

#### GET /health

Basic health check

#### GET /health/detailed

Detailed health check with database status

#### GET /health/database

Database connection statistics

#### GET /metrics

Get all metrics

#### GET /metrics/summary

Get metrics summary

#### GET /metrics/requests

Get request metrics

#### GET /metrics/database

Get database metrics

#### GET /metrics/application

Get application metrics

#### GET /metrics/errors

Get error metrics

#### POST /metrics/reset

Reset metrics

## Error Codes

### Common Error Codes

- `VALIDATION_ERROR`: Input validation failed
- `NOT_FOUND`: Resource not found
- `UNAUTHORIZED`: Authentication required
- `FORBIDDEN`: Access denied
- `INTERNAL_ERROR`: Internal server error
- `TIMEOUT_ERROR`: Request timeout
- `REQUEST_CANCELLED`: Request was cancelled
- `RATE_LIMIT_EXCEEDED`: Rate limit exceeded

### Validation Error Details

Validation errors include detailed field-level error information in the `details` object:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "title": "Title is required",
      "email": "Invalid email format",
      "password": "Password must be at least 8 characters"
    }
  }
}
```

## Examples

### Creating an Article with Tags

```bash
curl -X POST http://localhost:4300/api/v1/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "title": "My New Article",
    "content": "This is the article content...",
    "category_id": 1
  }'
```

### Getting Paginated Articles

```bash
curl "http://localhost:4300/api/v1/articles/paginated?page=2&limit=5"
```

### Checking Bookmark Status

```bash
curl -H "Authorization: Bearer your-jwt-token" \
  "http://localhost:4300/api/v1/bookmarks/check?article_id=1"
```

## SDK Examples

### JavaScript/Node.js

```javascript
const API_BASE = "http://localhost:4300/api/v1";

// Get articles with error handling
async function getArticles() {
  try {
    const response = await fetch(`${API_BASE}/articles`);
    const data = await response.json();

    if (data.success) {
      console.log("Articles:", data.data.articles);
      console.log("Request ID:", data.meta.request_id);
    } else {
      console.error("Error:", data.error.message);
    }
  } catch (error) {
    console.error("Network error:", error);
  }
}

// Create article with authentication
async function createArticle(token, articleData) {
  try {
    const response = await fetch(`${API_BASE}/articles`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(articleData),
    });

    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Error creating article:", error);
  }
}
```

### Python

```python
import requests

API_BASE = 'http://localhost:4300/api/v1'

def get_articles():
    response = requests.get(f'{API_BASE}/articles')
    data = response.json()

    if data['success']:
        print('Articles:', data['data']['articles'])
        print('Request ID:', data['meta']['request_id'])
    else:
        print('Error:', data['error']['message'])

def create_article(token, article_data):
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {token}'
    }

    response = requests.post(
        f'{API_BASE}/articles',
        json=article_data,
        headers=headers
    )

    return response.json()
```

## Changelog

### Version 1.0.0

- Initial API release
- Standard response format with metadata
- Request tracking with unique IDs
- Rate limiting implementation
- Comprehensive error handling
- Full CRUD operations for all resources
- Authentication and authorization
- Pagination support
- Health checks and monitoring
