// Simple validation test for CommentApi implementation
// This file validates that the CommentApi class has all required methods and proper types

import { commentApi, CommentApi } from './commentApi';
import type {
    Comment,
    CommentResponse,
    CreateCommentRequest,
    UpdateCommentRequest,
    PaginationParams,
    PaginatedResponse
} from './types';

// ============================================================================
// Type Validation Tests
// ============================================================================

// Test that commentApi is an instance of CommentApi
const isCommentApiInstance: boolean = commentApi instanceof CommentApi;
console.log('commentApi is instance of CommentApi:', isCommentApiInstance);

// ============================================================================
// Method Signature Validation
// ============================================================================

// Validate that all required methods exist with correct signatures
const validateCommentApiMethods = () => {
    // Basic CRUD operations
    const createCommentMethod: (commentData: CreateCommentRequest) => Promise<Comment> = commentApi.createComment.bind(commentApi);
    const getCommentByIdMethod: (id: number) => Promise<Comment> = commentApi.getCommentById.bind(commentApi);
    const updateCommentMethod: (id: number, commentData: UpdateCommentRequest) => Promise<Comment> = commentApi.updateComment.bind(commentApi);
    const deleteCommentMethod: (id: number) => Promise<void> = commentApi.deleteComment.bind(commentApi);

    // Comments by article methods
    const getCommentsByArticleMethod: (articleId: number) => Promise<CommentResponse[]> = commentApi.getCommentsByArticle.bind(commentApi);
    const getCommentsByArticlePaginatedMethod: (articleId: number, params?: PaginationParams) => Promise<PaginatedResponse<CommentResponse[]>> = commentApi.getCommentsByArticlePaginated.bind(commentApi);

    // Comments by user methods
    const getCommentsByUserMethod: (userId: number) => Promise<CommentResponse[]> = commentApi.getCommentsByUser.bind(commentApi);
    const getCommentsByUserPaginatedMethod: (userId: number, params?: PaginationParams) => Promise<PaginatedResponse<CommentResponse[]>> = commentApi.getCommentsByUserPaginated.bind(commentApi);

    // All comments method
    const getAllCommentsMethod: (params?: PaginationParams) => Promise<PaginatedResponse<CommentResponse[]>> = commentApi.getAllComments.bind(commentApi);

    // Statistics methods
    const getCommentCountByArticleMethod: (articleId: number) => Promise<number> = commentApi.getCommentCountByArticle.bind(commentApi);
    const getCommentCountByUserMethod: (userId: number) => Promise<number> = commentApi.getCommentCountByUser.bind(commentApi);

    // Recent comments method
    const getRecentCommentsMethod: (limit?: number) => Promise<CommentResponse[]> = commentApi.getRecentComments.bind(commentApi);

    console.log('All CommentApi methods have correct signatures');

    return {
        createCommentMethod,
        getCommentByIdMethod,
        updateCommentMethod,
        deleteCommentMethod,
        getCommentsByArticleMethod,
        getCommentsByArticlePaginatedMethod,
        getCommentsByUserMethod,
        getCommentsByUserPaginatedMethod,
        getAllCommentsMethod,
        getCommentCountByArticleMethod,
        getCommentCountByUserMethod,
        getRecentCommentsMethod
    };
};

// ============================================================================
// Type Interface Validation
// ============================================================================

// Validate request/response types
const validateCommentTypes = () => {
    // CreateCommentRequest validation
    const createCommentRequest: CreateCommentRequest = {
        article_id: 1,
        content: 'This is a test comment'
    };

    // UpdateCommentRequest validation
    const updateCommentRequest: UpdateCommentRequest = {
        content: 'This is an updated comment'
    };

    // Comment validation
    const comment: Comment = {
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
            name: 'Commenter',
            email: 'commenter@example.com',
            role: 'user',
            created_at: '2023-01-01T00:00:00Z',
            updated_at: '2023-01-01T00:00:00Z'
        },
        content: 'This is a comment',
        created_at: '2023-01-01T00:00:00Z'
    };

    // CommentResponse validation
    const commentResponse: CommentResponse = {
        id: 1,
        content: 'This is a comment',
        created_at: '2023-01-01T00:00:00Z',
        user: {
            id: 2,
            name: 'Commenter',
            email: 'commenter@example.com'
        },
        article: {
            id: 1,
            title: 'Test Article',
            slug: 'test-article'
        }
    };

    // PaginationParams validation
    const paginationParams: PaginationParams = {
        page: 1,
        limit: 10
    };

    console.log('All CommentApi types are properly defined');

    return {
        createCommentRequest,
        updateCommentRequest,
        comment,
        commentResponse,
        paginationParams
    };
};

// ============================================================================
// Run Validations
// ============================================================================

try {
    validateCommentApiMethods();
    validateCommentTypes();
    console.log('✅ CommentApi validation completed successfully');
} catch (error) {
    console.error('❌ CommentApi validation failed:', error);
}

// ============================================================================
// Export for potential use in other tests
// ============================================================================

export {
    isCommentApiInstance,
    validateCommentApiMethods,
    validateCommentTypes
};