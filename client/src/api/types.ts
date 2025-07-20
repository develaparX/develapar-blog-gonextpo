// API Type Definitions
// This file contains all TypeScript interfaces for API requests and responses

// ============================================================================
// Base API Response Structure
// ============================================================================

export interface ApiResponse<T = any> {
    success: boolean;
    data?: T;
    error?: ErrorResponse;
    pagination?: PaginationMetadata;
    meta?: ResponseMetadata;
}

export interface ErrorResponse {
    code: string;
    message: string;
    details?: Record<string, any>;
    request_id?: string;
    timestamp: string;
}

export interface ResponseMetadata {
    request_id: string;
    processing_time_ms: number;
    version: string;
    timestamp: string;
}

export interface PaginationMetadata {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
    request_id?: string;
}

// ============================================================================
// Core Entity Types
// ============================================================================

export interface User {
    id: number;
    name: string;
    email: string;
    role: string;
    created_at: string;
    updated_at: string;
}

export interface Category {
    id: number;
    name: string;
}

export interface Tag {
    id: number;
    name: string;
}

export interface Article {
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

export interface Comment {
    id: number;
    article: Article;
    user: User;
    content: string;
    created_at: string;
}

export interface Like {
    id: number;
    article: Article;
    user: User;
    created_at: string;
}

export interface Bookmark {
    id: number;
    article: Article;
    user: User;
    created_at: string;
}

export interface ArticleTag {
    article: Article;
    tag: Tag;
}

// ============================================================================
// Response DTOs (for API responses)
// ============================================================================

export interface UserResponse {
    id: number;
    name: string;
    email: string;
}

export interface ArticleResponse {
    id: number;
    title: string;
    slug: string;
}

export interface CommentResponse {
    id: number;
    content: string;
    created_at: string;
    user: UserResponse;
    article: ArticleResponse;
}

// ============================================================================
// Authentication Types
// ============================================================================

export interface LoginRequest {
    identifier: string;
    password: string;
}

export interface LoginResponse {
    access_token: string;
    refresh_token: string;
}

export interface RegisterRequest {
    name: string;
    email: string;
    password: string;
}

// ============================================================================
// User Management Request Types
// ============================================================================

export interface UpdateUserRequest {
    name?: string;
    email?: string;
    password?: string;
}

// ============================================================================
// Article Management Request Types
// ============================================================================

export interface CreateArticleRequest {
    title: string;
    content: string;
    category_id: number;
    tags?: string[];
}

export interface UpdateArticleRequest {
    title?: string;
    content?: string;
    category_id?: number;
    tags?: string[];
}

// ============================================================================
// Category Management Request Types
// ============================================================================

export interface CreateCategoryRequest {
    name: string;
}

export interface UpdateCategoryRequest {
    name?: string;
}

// ============================================================================
// Tag Management Request Types
// ============================================================================

export interface CreateTagRequest {
    name: string;
}

export interface UpdateTagRequest {
    name?: string;
}

export interface AssignTagsByNameRequest {
    article_id: number;
    tags: string[];
}

// ============================================================================
// Comment Management Request Types
// ============================================================================

export interface CreateCommentRequest {
    article_id: number;
    content: string;
}

export interface UpdateCommentRequest {
    content?: string;
}

// ============================================================================
// Like Management Request Types
// ============================================================================

export interface CreateLikeRequest {
    article_id: number;
}

// ============================================================================
// Bookmark Management Request Types
// ============================================================================

export interface CreateBookmarkRequest {
    article_id: number;
}

// ============================================================================
// Pagination and Query Types
// ============================================================================

export interface PaginationParams {
    page?: number;
    limit?: number;
}

export interface PaginatedResponse<T> {
    data: T;
    pagination: PaginationMetadata;
}

// ============================================================================
// API Error Types
// ============================================================================

export const ApiErrorCode = {
    // Network errors
    NETWORK_ERROR: 'NETWORK_ERROR',
    TIMEOUT: 'TIMEOUT',

    // Authentication errors
    UNAUTHORIZED: 'UNAUTHORIZED',
    FORBIDDEN: 'FORBIDDEN',
    TOKEN_EXPIRED: 'TOKEN_EXPIRED',

    // Validation errors
    VALIDATION_ERROR: 'VALIDATION_ERROR',
    INVALID_INPUT: 'INVALID_INPUT',

    // Resource errors
    NOT_FOUND: 'NOT_FOUND',
    CONFLICT: 'CONFLICT',

    // Server errors
    INTERNAL_ERROR: 'INTERNAL_ERROR',
    SERVICE_UNAVAILABLE: 'SERVICE_UNAVAILABLE',

    // Rate limiting
    RATE_LIMIT_EXCEEDED: 'RATE_LIMIT_EXCEEDED',

    // Request errors
    REQUEST_CANCELLED: 'REQUEST_CANCELLED',

    // Unknown
    UNKNOWN: 'UNKNOWN'
} as const;

export type ApiErrorCode = typeof ApiErrorCode[keyof typeof ApiErrorCode];

export interface ApiErrorDetails {
    field?: string;
    code?: string;
    message?: string;
    [key: string]: any;
}

// ============================================================================
// Utility Types
// ============================================================================

export type ApiMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

export interface RequestConfig {
    method: ApiMethod;
    url: string;
    data?: any;
    params?: Record<string, any>;
    headers?: Record<string, string>;
}

// ============================================================================
// Extended Article Types with Tags
// ============================================================================

export interface ArticleWithTags extends Omit<Article, 'category'> {
    category: Category | null;
    tags?: Tag[];
}

// ============================================================================
// Search and Filter Types
// ============================================================================

export interface ArticleFilters {
    category?: string;
    author?: number;
    tags?: string[];
    search?: string;
}

export interface ArticleSearchParams extends PaginationParams, ArticleFilters { }

// ============================================================================
// Status Check Types
// ============================================================================

export interface LikeStatus {
    is_liked: boolean;
    like_count: number;
}

export interface BookmarkStatus {
    is_bookmarked: boolean;
}