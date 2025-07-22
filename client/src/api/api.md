# Client API Layer Documentation

## Overview

The Client API Layer provides a comprehensive, type-safe interface for all HTTP communication with the Go backend server. It includes authentication management, error handling, response processing, and full CRUD operations for all application entities.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Authentication](#authentication)
3. [Error Handling](#error-handling)
4. [API Services](#api-services)
5. [Type Definitions](#type-definitions)
6. [Best Practices](#best-practices)
7. [Examples](#examples)

## Getting Started

### Installation and Setup

The API layer is automatically configured when you import any API service. The base configuration uses `/api/v1` as the default endpoint.

```typescript
import { userApi, articleApi, authApi } from "@/api";
```

### Environment Configuration

Configure the API base URL using environment variables:

```bash
# .env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### Basic Usage

```typescript
import { userApi, ApiError } from "@/api";

try {
  const users = await userApi.getAllUsers();
  console.log("Users:", users);
} catch (error) {
  if (error instanceof ApiError) {
    console.error("API Error:", error.message);
  }
}
```

## Authentication

### Authentication Flow

The authentication system handles login, registration, token refresh, and logout operations automatically.

#### Login

```typescript
import { authApi } from "@/api";

try {
  const response = await authApi.login({
    identifier: "user@example.com", // Email or username
    password: "password123",
  });

  console.log("Login successful:", response);
  // Tokens are automatically stored
} catch (error) {
  console.error("Login failed:", error.message);
}
```

#### Registration

```typescript
import { authApi } from "@/api";

try {
  const user = await authApi.register({
    name: "John Doe",
    email: "john@example.com",
    password: "securePassword123",
  });

  console.log("Registration successful:", user);
} catch (error) {
  console.error("Registration failed:", error.message);
}
```

#### Token Management

Tokens are automatically managed by the system:

```typescript
import { authApi } from "@/api";

// Check authentication status
const isAuthenticated = authApi.isAuthenticated();

// Get current user from token
const currentUser = authApi.getCurrentUser();

// Check if token needs refresh
const needsRefresh = authApi.needsTokenRefresh();

// Manually ensure valid token
const isValid = await authApi.ensureValidToken();
```

#### Logout

```typescript
import { authApi } from "@/api";

try {
  await authApi.logout();
  console.log("Logout successful");
} catch (error) {
  console.error("Logout error:", error.message);
}
```

### Authentication Requirements

| Endpoint Category  | Authentication Required | Notes                                                 |
| ------------------ | ----------------------- | ----------------------------------------------------- |
| Public Articles    | No                      | Reading articles by slug                              |
| User Management    | Yes                     | All operations require authentication                 |
| Article Management | Yes                     | Create, update, delete operations                     |
| Comments           | Yes                     | All comment operations                                |
| Likes/Bookmarks    | Yes                     | All like and bookmark operations                      |
| Categories/Tags    | Partial                 | Read operations public, write operations require auth |

## Error Handling

### Error Types

The API layer categorizes errors into several types:

```typescript
import { ApiError, ApiErrorCode } from "@/api";

// Network errors
ApiErrorCode.NETWORK_ERROR;
ApiErrorCode.TIMEOUT;
ApiErrorCode.REQUEST_CANCELLED;

// Authentication errors
ApiErrorCode.UNAUTHORIZED;
ApiErrorCode.FORBIDDEN;
ApiErrorCode.TOKEN_EXPIRED;

// Validation errors
ApiErrorCode.VALIDATION_ERROR;
ApiErrorCode.INVALID_INPUT;

// Resource errors
ApiErrorCode.NOT_FOUND;
ApiErrorCode.CONFLICT;

// Server errors
ApiErrorCode.INTERNAL_ERROR;
ApiErrorCode.SERVICE_UNAVAILABLE;
```

### Error Handling Examples

```typescript
import { userApi, ApiError, ErrorProcessor } from "@/api";

try {
  const user = await userApi.getUserById(123);
} catch (error) {
  if (error instanceof ApiError) {
    // Check error type
    if (error.isAuthenticationError) {
      // Handle authentication errors
      console.log("Please log in to continue");
    } else if (error.isValidationError) {
      // Handle validation errors
      console.log("Validation errors:", error.details);
    } else if (error.isNetworkError) {
      // Handle network errors
      console.log("Network connection issue");
    }

    // Get user-friendly message
    const friendlyMessage = ErrorProcessor.getUserFriendlyMessage(error);
    console.log(friendlyMessage);
  }
}
```

### Automatic Error Recovery

The system automatically handles:

- Token refresh on 401 errors
- Request retry on network failures
- Exponential backoff for failed requests

## API Services

### Authentication API

#### Methods

| Method               | Description          | Authentication | Parameters        |
| -------------------- | -------------------- | -------------- | ----------------- |
| `login(credentials)` | User login           | No             | `LoginRequest`    |
| `register(userData)` | User registration    | No             | `RegisterRequest` |
| `refreshToken()`     | Refresh access token | No             | None              |
| `logout()`           | User logout          | Yes            | None              |
| `isAuthenticated()`  | Check auth status    | No             | None              |
| `getCurrentUser()`   | Get user from token  | No             | None              |

#### Example Usage

```typescript
import { authApi } from "@/api";

// Login example
const loginData = await authApi.login({
  identifier: "user@example.com",
  password: "password123",
});

// Registration example
const newUser = await authApi.register({
  name: "Jane Doe",
  email: "jane@example.com",
  password: "securePassword123",
});
```

### User Management API

#### Methods

| Method                         | Description              | Authentication | Parameters                  |
| ------------------------------ | ------------------------ | -------------- | --------------------------- |
| `getUserById(id)`              | Get user by ID           | Yes            | `number`                    |
| `getAllUsers()`                | Get all users            | Yes            | None                        |
| `getAllUsersPaginated(params)` | Get paginated users      | Yes            | `PaginationParams`          |
| `updateUser(id, data)`         | Update user              | Yes            | `number, UpdateUserRequest` |
| `deleteUser(id)`               | Delete user              | Yes            | `number`                    |
| `getCurrentUser()`             | Get current user profile | Yes            | None                        |
| `updateCurrentUser(data)`      | Update current user      | Yes            | `UpdateUserRequest`         |

#### Example Usage

```typescript
import { userApi } from "@/api";

// Get user by ID
const user = await userApi.getUserById(123);

// Get paginated users
const paginatedUsers = await userApi.getAllUsersPaginated({
  page: 1,
  limit: 10,
});

// Update user
const updatedUser = await userApi.updateUser(123, {
  name: "Updated Name",
  email: "updated@example.com",
});

// Get current user profile
const currentUser = await userApi.getCurrentUser();
```

### Article Management API

#### Methods

| Method                                    | Description              | Authentication | Parameters                     |
| ----------------------------------------- | ------------------------ | -------------- | ------------------------------ |
| `createArticle(data)`                     | Create new article       | Yes            | `CreateArticleRequest`         |
| `getArticleBySlug(slug)`                  | Get article by slug      | No             | `string`                       |
| `updateArticle(id, data)`                 | Update article           | Yes            | `number, UpdateArticleRequest` |
| `deleteArticle(id)`                       | Delete article           | Yes            | `number`                       |
| `getAllArticles(params)`                  | Get all articles         | No             | `PaginationParams`             |
| `getArticlesByCategory(category, params)` | Get articles by category | No             | `string, PaginationParams`     |
| `getArticlesByAuthor(userId, params)`     | Get articles by author   | No             | `number, PaginationParams`     |
| `searchArticles(params)`                  | Search articles          | No             | `ArticleSearchParams`          |

#### Example Usage

```typescript
import { articleApi } from "@/api";

// Create article
const newArticle = await articleApi.createArticle({
  title: "My New Article",
  content: "Article content here...",
  category_id: 1,
  tags: ["javascript", "react"],
});

// Get article by slug
const article = await articleApi.getArticleBySlug("my-article-slug");

// Search articles
const searchResults = await articleApi.searchArticles({
  search: "javascript",
  category: "technology",
  page: 1,
  limit: 10,
});

// Get articles by category
const categoryArticles = await articleApi.getArticlesByCategory("technology");
```

### Category Management API

#### Methods

| Method                     | Description        | Authentication | Parameters                      |
| -------------------------- | ------------------ | -------------- | ------------------------------- |
| `createCategory(data)`     | Create category    | Yes            | `CreateCategoryRequest`         |
| `getAllCategories()`       | Get all categories | No             | None                            |
| `getCategoryById(id)`      | Get category by ID | No             | `number`                        |
| `updateCategory(id, data)` | Update category    | Yes            | `number, UpdateCategoryRequest` |
| `deleteCategory(id)`       | Delete category    | Yes            | `number`                        |

#### Example Usage

```typescript
import { categoryApi } from "@/api";

// Create category
const newCategory = await categoryApi.createCategory({
  name: "Technology",
});

// Get all categories
const categories = await categoryApi.getAllCategories();

// Update category
const updatedCategory = await categoryApi.updateCategory(1, {
  name: "Updated Technology",
});
```

### Tag Management API

#### Methods

| Method                        | Description            | Authentication | Parameters                 |
| ----------------------------- | ---------------------- | -------------- | -------------------------- |
| `createTag(data)`             | Create tag             | Yes            | `CreateTagRequest`         |
| `getAllTags()`                | Get all tags           | No             | None                       |
| `getTagById(id)`              | Get tag by ID          | No             | `number`                   |
| `updateTag(id, data)`         | Update tag             | Yes            | `number, UpdateTagRequest` |
| `deleteTag(id)`               | Delete tag             | Yes            | `number`                   |
| `assignTagsByName(data)`      | Assign tags to article | Yes            | `AssignTagsByNameRequest`  |
| `getTagsByArticle(articleId)` | Get article tags       | No             | `number`                   |

#### Example Usage

```typescript
import { tagApi } from "@/api";

// Create tag
const newTag = await tagApi.createTag({
  name: "javascript",
});

// Get all tags for dropdown
const allTags = await tagApi.getAllTags();

// Assign tags to article
const assignedTags = await tagApi.assignTagsByName({
  article_id: 123,
  tags: ["javascript", "react", "typescript"],
});

// Search tags for autocomplete
const matchingTags = await tagApi.searchTags("java", 5);
```

### Comment Management API

#### Methods

| Method                                | Description          | Authentication | Parameters                     |
| ------------------------------------- | -------------------- | -------------- | ------------------------------ |
| `createComment(data)`                 | Create comment       | Yes            | `CreateCommentRequest`         |
| `getCommentsByArticle(articleId)`     | Get article comments | No             | `number`                       |
| `getCommentsByUser(userId)`           | Get user comments    | Yes            | `number`                       |
| `updateComment(id, data)`             | Update comment       | Yes            | `number, UpdateCommentRequest` |
| `deleteComment(id)`                   | Delete comment       | Yes            | `number`                       |
| `getCommentCountByArticle(articleId)` | Get comment count    | No             | `number`                       |

#### Example Usage

```typescript
import { commentApi } from "@/api";

// Create comment
const newComment = await commentApi.createComment({
  article_id: 123,
  content: "Great article! Thanks for sharing.",
});

// Get comments for article
const articleComments = await commentApi.getCommentsByArticle(123);

// Get paginated comments
const paginatedComments = await commentApi.getCommentsByArticlePaginated(123, {
  page: 1,
  limit: 10,
});

// Update comment
const updatedComment = await commentApi.updateComment(456, {
  content: "Updated comment content",
});
```

### Like Management API

#### Methods

| Method                         | Description         | Authentication | Parameters          |
| ------------------------------ | ------------------- | -------------- | ------------------- |
| `addLike(data)`                | Add like to article | Yes            | `CreateLikeRequest` |
| `removeLike(articleId)`        | Remove like         | Yes            | `number`            |
| `checkLikeStatus(articleId)`   | Check like status   | Yes            | `number`            |
| `getLikesByArticle(articleId)` | Get article likes   | No             | `number`            |
| `getLikesByUser(userId)`       | Get user likes      | Yes            | `number`            |
| `toggleLike(articleId)`        | Toggle like status  | Yes            | `number`            |

#### Example Usage

```typescript
import { likeApi } from "@/api";

// Add like
const newLike = await likeApi.addLike({
  article_id: 123,
});

// Check like status
const likeStatus = await likeApi.checkLikeStatus(123);
console.log("Is liked:", likeStatus.is_liked);
console.log("Like count:", likeStatus.like_count);

// Toggle like (convenience method)
const updatedStatus = await likeApi.toggleLike(123);

// Get user's likes
const userLikes = await likeApi.getLikesByUser(456);
```

### Bookmark Management API

#### Methods

| Method                           | Description            | Authentication | Parameters              |
| -------------------------------- | ---------------------- | -------------- | ----------------------- |
| `addBookmark(data)`              | Add bookmark           | Yes            | `CreateBookmarkRequest` |
| `removeBookmark(articleId)`      | Remove bookmark        | Yes            | `number`                |
| `checkBookmarkStatus(articleId)` | Check bookmark status  | Yes            | `number`                |
| `getUserBookmarks()`             | Get user bookmarks     | Yes            | None                    |
| `toggleBookmark(articleId)`      | Toggle bookmark status | Yes            | `number`                |

#### Example Usage

```typescript
import { bookmarkApi } from "@/api";

// Add bookmark
const newBookmark = await bookmarkApi.addBookmark({
  article_id: 123,
});

// Check bookmark status
const bookmarkStatus = await bookmarkApi.checkBookmarkStatus(123);
console.log("Is bookmarked:", bookmarkStatus.is_bookmarked);

// Get user's bookmarks
const userBookmarks = await bookmarkApi.getUserBookmarks();

// Toggle bookmark
const updatedStatus = await bookmarkApi.toggleBookmark(123);
```

## Type Definitions

### Core Entity Types

```typescript
interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  created_at: string;
  updated_at: string;
}

interface Article {
  id: number;
  title: string;
  slug: string;
  content: string;
  user: User;
  category: Category;
  views: number;
  created_at: string;
  updated_at: string;
}

interface Category {
  id: number;
  name: string;
}

interface Tag {
  id: number;
  name: string;
}

interface Comment {
  id: number;
  article: Article;
  user: User;
  content: string;
  created_at: string;
}
```

### Request Types

```typescript
interface LoginRequest {
  identifier: string; // Email or username
  password: string;
}

interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

interface CreateArticleRequest {
  title: string;
  content: string;
  category_id: number;
  tags?: string[];
}

interface UpdateArticleRequest {
  title?: string;
  content?: string;
  category_id?: number;
  tags?: string[];
}
```

### Response Types

```typescript
interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: ErrorResponse;
  pagination?: PaginationMetadata;
  meta?: ResponseMetadata;
}

interface PaginationMetadata {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

interface LikeStatus {
  is_liked: boolean;
  like_count: number;
}

interface BookmarkStatus {
  is_bookmarked: boolean;
}
```

## Best Practices

### 1. Error Handling

Always handle errors appropriately:

```typescript
import { userApi, ApiError, ErrorProcessor } from "@/api";

try {
  const user = await userApi.getUserById(123);
  // Handle success
} catch (error) {
  if (error instanceof ApiError) {
    // Handle specific error types
    if (error.isAuthenticationError) {
      // Redirect to login
    } else if (error.isValidationError) {
      // Show validation errors
      console.log(error.details);
    }

    // Show user-friendly message
    const message = ErrorProcessor.getUserFriendlyMessage(error);
    showNotification(message);
  }
}
```

### 2. Pagination

Use pagination for large datasets:

```typescript
import { articleApi } from "@/api";

const paginatedArticles = await articleApi.getAllArticlesPaginated({
  page: 1,
  limit: 20,
});

console.log("Articles:", paginatedArticles.data);
console.log("Pagination:", paginatedArticles.pagination);
```

### 3. Authentication Checks

Check authentication before making authenticated requests:

```typescript
import { authApi, userApi } from "@/api";

if (authApi.isAuthenticated()) {
  const currentUser = await userApi.getCurrentUser();
} else {
  // Redirect to login
}
```

### 4. Loading States

Handle loading states in your components:

```typescript
import { useState } from "react";
import { articleApi } from "@/api";

function ArticleList() {
  const [articles, setArticles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const loadArticles = async () => {
    setLoading(true);
    setError(null);

    try {
      const data = await articleApi.getAllArticles();
      setArticles(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Component JSX...
}
```

## Examples

### Complete Authentication Flow

```typescript
import { authApi, userApi } from "@/api";

class AuthService {
  async login(email: string, password: string) {
    try {
      // Attempt login
      const response = await authApi.login({
        identifier: email,
        password: password,
      });

      // Get user profile
      const user = await userApi.getCurrentUser();

      return { success: true, user };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }

  async register(name: string, email: string, password: string) {
    try {
      const user = await authApi.register({
        name,
        email,
        password,
      });

      return { success: true, user };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }

  async logout() {
    try {
      await authApi.logout();
      return { success: true };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }
}
```

### Article Management with Error Handling

```typescript
import { articleApi, ApiError } from "@/api";

class ArticleService {
  async createArticle(articleData) {
    try {
      const article = await articleApi.createArticle(articleData);
      return { success: true, article };
    } catch (error) {
      if (error instanceof ApiError) {
        if (error.isValidationError) {
          return {
            success: false,
            validationErrors: error.details,
          };
        }
      }
      return { success: false, error: error.message };
    }
  }

  async getArticleWithInteractions(slug) {
    try {
      // Get article
      const article = await articleApi.getArticleBySlug(slug);

      // Get like status if authenticated
      let likeStatus = null;
      let bookmarkStatus = null;

      if (authApi.isAuthenticated()) {
        [likeStatus, bookmarkStatus] = await Promise.all([
          likeApi.checkLikeStatus(article.id),
          bookmarkApi.checkBookmarkStatus(article.id),
        ]);
      }

      return {
        success: true,
        article,
        likeStatus,
        bookmarkStatus,
      };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }
}
```

### Search and Filtering

```typescript
import { articleApi, categoryApi, tagApi } from "@/api";

class SearchService {
  async searchArticles(query, filters = {}) {
    try {
      const searchParams = {
        search: query,
        category: filters.category,
        tags: filters.tags,
        author: filters.author,
        page: filters.page || 1,
        limit: filters.limit || 10,
      };

      const results = await articleApi.searchArticlesPaginated(searchParams);

      return {
        success: true,
        articles: results.data,
        pagination: results.pagination,
      };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }

  async getSearchFilters() {
    try {
      const [categories, popularTags] = await Promise.all([
        categoryApi.getAllCategories(),
        tagApi.getPopularTags(20),
      ]);

      return {
        success: true,
        categories,
        tags: popularTags,
      };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }
}
```

### Batch Operations

```typescript
import { articleApi, likeApi, bookmarkApi } from "@/api";

class BatchService {
  async getArticlesWithStatus(articleIds) {
    try {
      // Get articles
      const articles = await articleApi.getArticlesByIds(articleIds);

      // Get like and bookmark status for each article
      const statusPromises = articles.map(async (article) => {
        const [likeStatus, bookmarkStatus] = await Promise.all([
          likeApi.checkLikeStatus(article.id),
          bookmarkApi.checkBookmarkStatus(article.id),
        ]);

        return {
          ...article,
          likeStatus,
          bookmarkStatus,
        };
      });

      const articlesWithStatus = await Promise.all(statusPromises);

      return { success: true, articles: articlesWithStatus };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }
}
```

## Troubleshooting

### Common Issues

1. **401 Unauthorized Errors**

   - Check if user is logged in
   - Verify token hasn't expired
   - Ensure proper authentication headers

2. **Network Errors**

   - Check internet connection
   - Verify API server is running
   - Check CORS configuration

3. **Validation Errors**

   - Check request data format
   - Verify required fields are provided
   - Validate data types match API expectations

4. **Rate Limiting**
   - Implement request throttling
   - Add retry logic with exponential backoff
   - Cache frequently accessed data

### Debug Mode

Enable debug logging by setting environment variable:

```bash
VITE_API_DEBUG=true
```

This will log all API requests and responses to the console.

---

For more information or support, please refer to the backend API documentation or contact the development team.
