// Simple validation test for BookmarkApi implementation
// This file validates that the BookmarkApi class has all required methods and proper types

import { bookmarkApi, BookmarkApi } from './bookmarkApi';
import type {
    Bookmark,
    CreateBookmarkRequest,
    BookmarkStatus,
    PaginationParams,
    PaginatedResponse
} from './types';

// ============================================================================
// Type Validation Tests
// ============================================================================

// Test that bookmarkApi is an instance of BookmarkApi
const isBookmarkApiInstance: boolean = bookmarkApi instanceof BookmarkApi;
console.log('bookmarkApi is instance of BookmarkApi:', isBookmarkApiInstance);

// ============================================================================
// Method Signature Validation
// ============================================================================

// Validate that all required methods exist with correct signatures
const validateBookmarkApiMethods = () => {
    // Basic bookmark operations
    const addBookmarkMethod: (bookmarkData: CreateBookmarkRequest) => Promise<Bookmark> = bookmarkApi.addBookmark.bind(bookmarkApi);
    const removeBookmarkMethod: (articleId: number) => Promise<void> = bookmarkApi.removeBookmark.bind(bookmarkApi);
    const removeBookmarkByIdMethod: (bookmarkId: number) => Promise<void> = bookmarkApi.removeBookmarkById.bind(bookmarkApi);

    // Status checking
    const checkBookmarkStatusMethod: (articleId: number) => Promise<BookmarkStatus> = bookmarkApi.checkBookmarkStatus.bind(bookmarkApi);

    // Get user bookmarks
    const getUserBookmarksMethod: () => Promise<Bookmark[]> = bookmarkApi.getUserBookmarks.bind(bookmarkApi);
    const getUserBookmarksPaginatedMethod: (params?: PaginationParams) => Promise<PaginatedResponse<Bookmark[]>> = bookmarkApi.getUserBookmarksPaginated.bind(bookmarkApi);

    // Get bookmarks by user (admin)
    const getBookmarksByUserMethod: (userId: number) => Promise<Bookmark[]> = bookmarkApi.getBookmarksByUser.bind(bookmarkApi);
    const getBookmarksByUserPaginatedMethod: (userId: number, params?: PaginationParams) => Promise<PaginatedResponse<Bookmark[]>> = bookmarkApi.getBookmarksByUserPaginated.bind(bookmarkApi);

    // Get bookmarks by article (admin)
    const getBookmarksByArticleMethod: (articleId: number) => Promise<Bookmark[]> = bookmarkApi.getBookmarksByArticle.bind(bookmarkApi);
    const getBookmarksByArticlePaginatedMethod: (articleId: number, params?: PaginationParams) => Promise<PaginatedResponse<Bookmark[]>> = bookmarkApi.getBookmarksByArticlePaginated.bind(bookmarkApi);

    // Statistics methods
    const getBookmarkCountByArticleMethod: (articleId: number) => Promise<number> = bookmarkApi.getBookmarkCountByArticle.bind(bookmarkApi);
    const getBookmarkCountByUserMethod: (userId: number) => Promise<number> = bookmarkApi.getBookmarkCountByUser.bind(bookmarkApi);

    // Toggle method
    const toggleBookmarkMethod: (articleId: number) => Promise<BookmarkStatus> = bookmarkApi.toggleBookmark.bind(bookmarkApi);

    // Get specific bookmark
    const getBookmarkByIdMethod: (bookmarkId: number) => Promise<Bookmark> = bookmarkApi.getBookmarkById.bind(bookmarkApi);

    console.log('All BookmarkApi methods have correct signatures');

    return {
        addBookmarkMethod,
        removeBookmarkMethod,
        removeBookmarkByIdMethod,
        checkBookmarkStatusMethod,
        getUserBookmarksMethod,
        getUserBookmarksPaginatedMethod,
        getBookmarksByUserMethod,
        getBookmarksByUserPaginatedMethod,
        getBookmarksByArticleMethod,
        getBookmarksByArticlePaginatedMethod,
        getBookmarkCountByArticleMethod,
        getBookmarkCountByUserMethod,
        toggleBookmarkMethod,
        getBookmarkByIdMethod
    };
};

// ============================================================================
// Type Interface Validation
// ============================================================================

// Validate request/response types
const validateBookmarkTypes = () => {
    // CreateBookmarkRequest validation
    const createBookmarkRequest: CreateBookmarkRequest = {
        article_id: 1
    };

    // Bookmark validation
    const bookmark: Bookmark = {
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
            name: 'Bookmarker',
            email: 'bookmarker@example.com',
            role: 'user',
            created_at: '2023-01-01T00:00:00Z',
            updated_at: '2023-01-01T00:00:00Z'
        },
        created_at: '2023-01-01T00:00:00Z'
    };

    // BookmarkStatus validation
    const bookmarkStatus: BookmarkStatus = {
        is_bookmarked: true
    };

    // PaginationParams validation
    const paginationParams: PaginationParams = {
        page: 1,
        limit: 10
    };

    console.log('All BookmarkApi types are properly defined');

    return {
        createBookmarkRequest,
        bookmark,
        bookmarkStatus,
        paginationParams
    };
};

// ============================================================================
// Run Validations
// ============================================================================

try {
    validateBookmarkApiMethods();
    validateBookmarkTypes();
    console.log('✅ BookmarkApi validation completed successfully');
} catch (error) {
    console.error('❌ BookmarkApi validation failed:', error);
}

// ============================================================================
// Export for potential use in other tests
// ============================================================================

export {
    isBookmarkApiInstance,
    validateBookmarkApiMethods,
    validateBookmarkTypes
};