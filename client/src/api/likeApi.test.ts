// Simple validation test for LikeApi implementation
// This file validates that the LikeApi class has all required methods and proper types

import { likeApi, LikeApi } from './likeApi';
import type {
    Like,
    CreateLikeRequest,
    LikeStatus,
    PaginationParams,
    PaginatedResponse
} from './types';

// ============================================================================
// Type Validation Tests
// ============================================================================

// Test that likeApi is an instance of LikeApi
const isLikeApiInstance: boolean = likeApi instanceof LikeApi;
console.log('likeApi is instance of LikeApi:', isLikeApiInstance);

// ============================================================================
// Method Signature Validation
// ============================================================================

// Validate that all required methods exist with correct signatures
const validateLikeApiMethods = () => {
    // Basic like operations
    const addLikeMethod: (likeData: CreateLikeRequest) => Promise<Like> = likeApi.addLike.bind(likeApi);
    const removeLikeMethod: (articleId: number) => Promise<void> = likeApi.removeLike.bind(likeApi);
    const removeLikeByIdMethod: (likeId: number) => Promise<void> = likeApi.removeLikeById.bind(likeApi);

    // Status checking
    const checkLikeStatusMethod: (articleId: number) => Promise<LikeStatus> = likeApi.checkLikeStatus.bind(likeApi);

    // Get likes by article
    const getLikesByArticleMethod: (articleId: number) => Promise<Like[]> = likeApi.getLikesByArticle.bind(likeApi);
    const getLikesByArticlePaginatedMethod: (articleId: number, params?: PaginationParams) => Promise<PaginatedResponse<Like[]>> = likeApi.getLikesByArticlePaginated.bind(likeApi);

    // Get likes by user
    const getLikesByUserMethod: (userId: number) => Promise<Like[]> = likeApi.getLikesByUser.bind(likeApi);
    const getLikesByUserPaginatedMethod: (userId: number, params?: PaginationParams) => Promise<PaginatedResponse<Like[]>> = likeApi.getLikesByUserPaginated.bind(likeApi);

    // Statistics methods
    const getLikeCountByArticleMethod: (articleId: number) => Promise<number> = likeApi.getLikeCountByArticle.bind(likeApi);
    const getLikeCountByUserMethod: (userId: number) => Promise<number> = likeApi.getLikeCountByUser.bind(likeApi);

    // Toggle method
    const toggleLikeMethod: (articleId: number) => Promise<LikeStatus> = likeApi.toggleLike.bind(likeApi);

    console.log('All LikeApi methods have correct signatures');

    return {
        addLikeMethod,
        removeLikeMethod,
        removeLikeByIdMethod,
        checkLikeStatusMethod,
        getLikesByArticleMethod,
        getLikesByArticlePaginatedMethod,
        getLikesByUserMethod,
        getLikesByUserPaginatedMethod,
        getLikeCountByArticleMethod,
        getLikeCountByUserMethod,
        toggleLikeMethod
    };
};

// ============================================================================
// Type Interface Validation
// ============================================================================

// Validate request/response types
const validateLikeTypes = () => {
    // CreateLikeRequest validation
    const createLikeRequest: CreateLikeRequest = {
        article_id: 1
    };

    // Like validation
    const like: Like = {
        id: 1,
        article: {
            id: 1,
            title: 'Test Article',
            slug: 'test-article',
            content: 'Test content',
            user: {
                id: 1,
                name: 'Test User',
                email: 'test@example.com',
                role: 'user',
                created_at: '2023-01-01T00:00:00Z',
                updated_at: '2023-01-01T00:00:00Z'
            },
            category: {
                id: 1,
                name: 'Technology'
            },
            views: 100,
            created_at: '2023-01-01T00:00:00Z',
            updated_at: '2023-01-01T00:00:00Z'
        },
        user: {
            id: 2,
            name: 'Liker',
            email: 'liker@example.com',
            role: 'user',
            created_at: '2023-01-01T00:00:00Z',
            updated_at: '2023-01-01T00:00:00Z'
        },
        created_at: '2023-01-01T00:00:00Z'
    };

    // LikeStatus validation
    const likeStatus: LikeStatus = {
        is_liked: true,
        like_count: 5
    };

    // PaginationParams validation
    const paginationParams: PaginationParams = {
        page: 1,
        limit: 10
    };

    console.log('All LikeApi types are properly defined');

    return {
        createLikeRequest,
        like,
        likeStatus,
        paginationParams
    };
};

// ============================================================================
// Run Validations
// ============================================================================

try {
    validateLikeApiMethods();
    validateLikeTypes();
    console.log('✅ LikeApi validation completed successfully');
} catch (error) {
    console.error('❌ LikeApi validation failed:', error);
}

// ============================================================================
// Export for potential use in other tests
// ============================================================================

export {
    isLikeApiInstance,
    validateLikeApiMethods,
    validateLikeTypes
};