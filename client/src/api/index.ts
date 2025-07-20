// Main API Export File
// This file exports all API services and utilities for easy importing

// Export types
export type {
    // Core types
    User,
    Article,
    Category,
    Tag,
    Comment,
    Like,
    Bookmark,
    ArticleTag,

    // Response DTOs
    UserResponse,
    ArticleResponse,
    CommentResponse,

    // Authentication types
    LoginRequest,
    LoginResponse,
    RegisterRequest,

    // Request types
    UpdateUserRequest,
    CreateArticleRequest,
    UpdateArticleRequest,
    CreateCategoryRequest,
    UpdateCategoryRequest,
    CreateTagRequest,
    UpdateTagRequest,
    AssignTagsByNameRequest,
    CreateCommentRequest,
    UpdateCommentRequest,
    CreateLikeRequest,
    CreateBookmarkRequest,

    // API response types
    ApiResponse,
    ErrorResponse,
    ResponseMetadata,
    PaginationMetadata,
    PaginatedResponse,

    // Utility types
    PaginationParams,
    ArticleFilters,
    ArticleSearchParams,
    LikeStatus,
    BookmarkStatus,
    ArticleWithTags,
    ApiErrorDetails,
    RequestConfig,
    ApiMethod
} from './types';

// Export constants
export { ApiErrorCode } from './types';

// Export API client
export { apiClient, TokenManager } from './apiClient';

// Export error handling utilities
export {
    ApiError,
    ErrorProcessor,
    ResponseHandler,
    NetworkDetector
} from './errorHandler';

// Export default API client instance
export { default as api } from './apiClient';