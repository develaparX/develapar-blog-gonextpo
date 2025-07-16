# Improved Article Creation Flow

## Overview

Artikel sekarang dapat dibuat dengan category dan tags dalam satu request yang atomic. Slug juga dibuat otomatis dari judul artikel.

## Major Changes Made

### 1. New DTO for Article Creation

- Added `CreateArticleRequest` dengan support untuk tags
- Updated `UpdateArticleRequest` dengan support untuk tags
- Slug tidak lagi diterima sebagai input dari user

### 2. Enhanced Service Layer

- `CreateArticleWithTags`: Method baru untuk create artikel dengan tags sekaligus
- `UpdateArticle`: Sekarang support update tags juga
- Automatic slug generation dari title
- Automatic tag creation jika tag belum ada

### 3. Improved Controller

- `CreateArticleHandler`: Menggunakan `CreateArticleWithTags` method
- Single request untuk create artikel dengan category dan tags
- Better error handling dan response

### 4. Slug Generator Improvements

- Menangani multiple spaces dengan benar
- Menghilangkan karakter khusus
- Menangani consecutive hyphens
- Menangani leading/trailing hyphens
- Fallback ke "untitled" jika slug kosong

## API Flow Comparison

### Before (Multiple Requests Required)

```bash
# Step 1: Create Article
POST /api/v1/article/
{
  "title": "How to Build REST API",
  "slug": "how-to-build-rest-api",  // Manual input
  "content": "Article content...",
  "category": {"id": 1}
}

# Step 2: Assign Tags (Separate Request)
POST /api/v1/article-to-tag/
{
  "article_id": 1,
  "tags": ["golang", "api", "tutorial"]
}
```

### After (Single Atomic Request)

```bash
# Single Request - Everything Together
POST /api/v1/article/
{
  "title": "How to Build REST API",
  "content": "Article content...",
  "category_id": 1,
  "tags": ["golang", "api", "tutorial"]  // Optional
}
```

## New API Examples

### Create Article with Tags

```json
{
  "title": "How to Build REST API with Go",
  "content": "Complete tutorial on building REST API...",
  "category_id": 1,
  "tags": ["golang", "api", "tutorial", "backend"]
}
```

### Create Article without Tags

```json
{
  "title": "Simple Article",
  "content": "Article content...",
  "category_id": 2
}
```

### Update Article with Tags

```json
{
  "title": "Updated Title",
  "content": "Updated content...",
  "category_id": 3,
  "tags": ["updated", "new-tags"]
}
```

## Slug Generation Rules

1. **Lowercase conversion**: "Hello World" → "hello world"
2. **Special character replacement**: "Hello World!" → "hello world"
3. **Space to hyphen**: "hello world" → "hello-world"
4. **Multiple spaces normalization**: "hello world" → "hello-world"
5. **Consecutive hyphen removal**: "hello--world" → "hello-world"
6. **Trim hyphens**: "-hello-world-" → "hello-world"
7. **Empty fallback**: "" → "untitled"

## Examples

| Title                                             | Generated Slug                                  |
| ------------------------------------------------- | ----------------------------------------------- |
| "Hello World"                                     | "hello-world"                                   |
| "How to Build a REST API with Go & Gin Framework" | "how-to-build-a-rest-api-with-go-gin-framework" |
| "Top 10 Programming Languages in 2024"            | "top-10-programming-languages-in-2024"          |
| "Hello World!!!"                                  | "hello-world"                                   |
| "!@#$%"                                           | "untitled"                                      |

## Benefits

### 1. Improved User Experience

✅ **Single Request**: Create artikel dengan category dan tags sekaligus  
✅ **Automatic Slug**: User tidak perlu memikirkan slug  
✅ **Atomic Operation**: Semua berhasil atau semua gagal

### 2. Better Developer Experience

✅ **Simplified Frontend**: Hanya perlu 1 API call  
✅ **Better Error Handling**: Tidak ada partial state  
✅ **Consistent Response**: Format response yang konsisten

### 3. Technical Benefits

✅ **Atomic Transactions**: Create artikel dan assign tags dalam satu operasi  
✅ **Auto Tag Creation**: Tag baru otomatis dibuat jika belum ada  
✅ **SEO Friendly**: Slug otomatis SEO-friendly  
✅ **Consistency**: Semua slug mengikuti format yang sama

### 4. Backward Compatibility

✅ **Old Endpoints**: Endpoint lama masih berfungsi  
✅ **Gradual Migration**: Frontend bisa migrate secara bertahap

## Migration Guide

### For Frontend Developers

**Old Way:**

```javascript
// Step 1: Create article
const article = await createArticle({
  title: "My Article",
  slug: "my-article",
  content: "Content...",
  category: { id: 1 },
});

// Step 2: Assign tags
await assignTags({
  article_id: article.id,
  tags: ["tag1", "tag2"],
});
```

**New Way:**

```javascript
// Single request
const article = await createArticle({
  title: "My Article",
  content: "Content...",
  category_id: 1,
  tags: ["tag1", "tag2"],
});
```

### Response Format

```json
{
  "message": "Success create article with tags",
  "data": {
    "id": 1,
    "title": "How to Build REST API with Go",
    "slug": "how-to-build-rest-api-with-go",
    "content": "Complete tutorial...",
    "user": { "id": 1, "name": "John Doe" },
    "category": { "id": 1, "name": "Technology" },
    "views": 0,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

## Technical Implementation Details

### Service Layer Changes

```go
// New method for creating article with tags
func (a *articleService) CreateArticleWithTags(req dto.CreateArticleRequest, userID int) (model.Article, error) {
    // 1. Generate slug from title
    // 2. Create article
    // 3. Assign tags (create new tags if needed)
    // 4. Return complete article
}

// Enhanced update method
func (a *articleService) UpdateArticle(id int, req dto.UpdateArticleRequest) (model.Article, error) {
    // 1. Update article fields
    // 2. Update slug if title changed
    // 3. Update tags if provided
    // 4. Return updated article
}
```

### Controller Changes

```go
// Updated to use new DTO and service method
func (c *ArticleController) CreateArticleHandler(ctx *gin.Context) {
    var req dto.CreateArticleRequest
    // ... validation ...
    data, err := c.service.CreateArticleWithTags(req, userId)
    // ... response ...
}
```

### DTO Structure

```go
type CreateArticleRequest struct {
    Title      string   `json:"title" binding:"required"`
    Content    string   `json:"content" binding:"required"`
    CategoryID int      `json:"category_id" binding:"required"`
    Tags       []string `json:"tags,omitempty"`
}

type UpdateArticleRequest struct {
    Title      *string  `json:"title"`
    Content    *string  `json:"content"`
    CategoryID *int     `json:"category_id"`
    Tags       []string `json:"tags,omitempty"`
}
```

## Testing

All existing tests have been updated to work with the new flow:

- ✅ Unit tests for service layer
- ✅ Integration tests for controller layer
- ✅ Mock services updated
- ✅ All tests passing

## Conclusion

Flow baru ini memberikan pengalaman yang jauh lebih baik untuk frontend developers dan end users. Dengan satu request, artikel dapat dibuat lengkap dengan category dan tags, sementara slug dibuat otomatis untuk memastikan konsistensi dan SEO-friendly URLs.
