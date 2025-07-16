# Automatic Slug Generation

## Overview

Slug sekarang dibuat otomatis dari judul artikel. User tidak perlu lagi menginput slug secara manual.

## Changes Made

### 1. DTO Updates

- Removed `Slug` field from `UpdateArticleRequest`
- Slug tidak lagi diterima sebagai input dari user

### 2. Service Layer Updates

- `CreateArticle`: Slug dibuat otomatis dari `payload.Title`
- `UpdateArticle`: Ketika title diupdate, slug juga otomatis diupdate

### 3. Slug Generator Improvements

- Menangani multiple spaces dengan benar
- Menghilangkan karakter khusus
- Menangani consecutive hyphens
- Menangani leading/trailing hyphens
- Fallback ke "untitled" jika slug kosong

## API Examples

### Create Article (Before)

```json
{
  "title": "How to Build REST API",
  "slug": "how-to-build-rest-api", // Manual input
  "content": "Article content...",
  "category": { "id": 1 }
}
```

### Create Article (After)

```json
{
  "title": "How to Build REST API",
  // slug akan otomatis menjadi: "how-to-build-rest-api"
  "content": "Article content...",
  "category": { "id": 1 }
}
```

### Update Article (Before)

```json
{
  "title": "Updated Title",
  "slug": "updated-title", // Manual input
  "content": "Updated content..."
}
```

### Update Article (After)

```json
{
  "title": "Updated Title",
  // slug akan otomatis menjadi: "updated-title"
  "content": "Updated content..."
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

1. **Consistency**: Semua slug mengikuti format yang sama
2. **User Experience**: User tidak perlu memikirkan slug
3. **SEO Friendly**: Slug otomatis SEO-friendly
4. **Error Prevention**: Menghindari slug yang tidak valid
5. **Maintenance**: Lebih mudah maintain karena otomatis
