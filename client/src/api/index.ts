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

// Export API services
export { authApi, AuthApi } from './authApi';
export { userApi, UserApi } from './userApi';
export { articleApi, ArticleApi } from './articleApi';
export { categoryApi, CategoryApi } from './categoryApi';
export { tagApi, TagApi } from './tagApi';
export { commentApi, CommentApi } from './commentApi';
export { likeApi, LikeApi } from './likeApi';
export { bookmarkApi, BookmarkApi } from './bookmarkApi';

// Export token management utilities
export { tokenManager, EnhancedTokenManager, TokenRefreshManager } from './tokenManager';

// Export default API client instance
export { default as api } from './apiClient';