// Bookmark API Service
// This file provides methods for bookmark management operations with status checking

import type { AxiosResponse } from 'axios';
import { apiClient } from './apiClient';
import type {
    ApiResponse,
    Bookmark,
    CreateBookmarkRequest,
    BookmarkStatus,
    PaginatedResponse,
    PaginationParams
} from './types';

// ============================================================================
// Response Handler Utility
// ============================================================================

class BookmarkApiError extends Error {
    constructor(
        public code: string,
        message: string,
        public statusCode?: number,
        public details?: Record<string, any>
    ) {
        super(message);
        this.name = 'BookmarkApiError';
    }
}

function handleResponse<T>(response: AxiosResponse<ApiResponse<T>>): T {
    if (response.data.success && response.data.data !== undefined) {
        return response.data.data;
    }

    const error = response.data.error;
    throw new BookmarkApiError(
        error?.code || 'UNKNOWN',
        error?.message || 'Unknown error occurred',
        response.status,
        error?.details
    );
}

function handlePaginatedResponse<T>(response: AxiosResponse<ApiResponse<T>>): PaginatedResponse<T> {
    if (response.data.success && response.data.data !== undefined) {
        return {
            data: response.data.data,
            pagination: response.data.pagination
        };
    }

    const error = response.data.error;
    throw new BookmarkApiError(
        error?.code || 'UNKNOWN',
        error?.message || 'Unknown error occurred',
        response.status,
        error?.details
    );
}

// ============================================================================
// Bookmark API Service Class
// ============================================================================

export class BookmarkApi {

    // ========================================================================
    // Add Bookmark
    // ========================================================================

    /**
     * Add a bookmark to an article
     * @param bookmarkData - Bookmark creation data with article ID
     * @returns Promise<Bookmark> - Created bookmark
     * @throws BookmarkApiError - If bookmark creation fails or already exists
     */
    async addBookmark(bookmarkData: CreateBookmarkRequest): Promise<Bookmark> {
        try {
            const response = await apiClient.post<ApiResponse<Bookmark>>(
                '/bookmarks',
                bookmarkData
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'CREATE_FAILED',
                    error.response.data?.error?.message || 'Failed to add bookmark',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                'Network error occurred while adding bookmark'
            );
        }
    }

    // ========================================================================
    // Remove Bookmark
    // ========================================================================

    /**
     * Remove a bookmark from an article
     * @param articleId - Article ID to remove bookmark from
     * @returns Promise<void> - Resolves when bookmark is removed
     * @throws BookmarkApiError - If bookmark removal fails or bookmark doesn't exist
     */
    async removeBookmark(articleId: number): Promise<void> {
        try {
            await apiClient.delete(`/articles/${articleId}/bookmarks`);
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'DELETE_FAILED',
                    error.response.data?.error?.message || `Failed to remove bookmark from article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                `Network error occurred while removing bookmark from article ${articleId}`
            );
        }
    }

    /**
     * Remove a specific bookmark by ID (admin use)
     * @param bookmarkId - Bookmark ID to remove
     * @returns Promise<void> - Resolves when bookmark is removed
     * @throws BookmarkApiError - If bookmark removal fails or bookmark doesn't exist
     */
    async removeBookmarkById(bookmarkId: number): Promise<void> {
        try {
            await apiClient.delete(`/bookmarks/${bookmarkId}`);
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'DELETE_FAILED',
                    error.response.data?.error?.message || `Failed to remove bookmark with ID ${bookmarkId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                `Network error occurred while removing bookmark with ID ${bookmarkId}`
            );
        }
    }

    // ========================================================================
    // Check Bookmark Status
    // ========================================================================

    /**
     * Check if current user has bookmarked an article
     * @param articleId - Article ID to check
     * @returns Promise<BookmarkStatus> - Bookmark status information
     * @throws BookmarkApiError - If status check fails
     */
    async checkBookmarkStatus(articleId: number): Promise<BookmarkStatus> {
        try {
            const response = await apiClient.get<ApiResponse<BookmarkStatus>>(
                `/articles/${articleId}/bookmarks/status`
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to check bookmark status for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                `Network error occurred while checking bookmark status for article ${articleId}`
            );
        }
    }

    // ========================================================================
    // Get Bookmarks by User
    // ========================================================================

    /**
     * Get all bookmarks by the current user (non-paginated)
     * @returns Promise<Bookmark[]> - Array of user's bookmarks
     * @throws BookmarkApiError - If retrieval fails
     */
    async getUserBookmarks(): Promise<Bookmark[]> {
        try {
            const response = await apiClient.get<ApiResponse<Bookmark[]>>('/bookmarks');
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || 'Failed to fetch user bookmarks',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                'Network error occurred while fetching user bookmarks'
            );
        }
    }

    /**
     * Get bookmarks by the current user with pagination
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<Bookmark[]>> - Paginated user bookmarks
     * @throws BookmarkApiError - If retrieval fails
     */
    async getUserBookmarksPaginated(params: PaginationParams = {}): Promise<PaginatedResponse<Bookmark[]>> {
        try {
            const { page = 1, limit = 10 } = params;
            const response = await apiClient.get<ApiResponse<Bookmark[]>>(
                '/bookmarks/paginated',
                {
                    params: { page, limit }
                }
            );
            return handlePaginatedResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || 'Failed to fetch paginated user bookmarks',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                'Network error occurred while fetching paginated user bookmarks'
            );
        }
    }

    /**
     * Get all bookmarks by a specific user (admin use)
     * @param userId - User ID
     * @returns Promise<Bookmark[]> - Array of bookmarks by the user
     * @throws BookmarkApiError - If retrieval fails
     */
    async getBookmarksByUser(userId: number): Promise<Bookmark[]> {
        try {
            const response = await apiClient.get<ApiResponse<Bookmark[]>>(
                `/users/${userId}/bookmarks`
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch bookmarks by user ${userId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching bookmarks by user ${userId}`
            );
        }
    }

    /**
     * Get bookmarks by a specific user with pagination (admin use)
     * @param userId - User ID
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<Bookmark[]>> - Paginated bookmarks by the user
     * @throws BookmarkApiError - If retrieval fails
     */
    async getBookmarksByUserPaginated(
        userId: number,
        params: PaginationParams = {}
    ): Promise<PaginatedResponse<Bookmark[]>> {
        try {
            const { page = 1, limit = 10 } = params;
            const response = await apiClient.get<ApiResponse<Bookmark[]>>(
                `/users/${userId}/bookmarks/paginated`,
                {
                    params: { page, limit }
                }
            );
            return handlePaginatedResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch paginated bookmarks by user ${userId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching paginated bookmarks by user ${userId}`
            );
        }
    }

    // ========================================================================
    // Get Bookmarks by Article
    // ========================================================================

    /**
     * Get all bookmarks for a specific article (admin use)
     * @param articleId - Article ID
     * @returns Promise<Bookmark[]> - Array of bookmarks for the article
     * @throws BookmarkApiError - If retrieval fails
     */
    async getBookmarksByArticle(articleId: number): Promise<Bookmark[]> {
        try {
            const response = await apiClient.get<ApiResponse<Bookmark[]>>(
                `/articles/${articleId}/bookmarks`
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch bookmarks for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching bookmarks for article ${articleId}`
            );
        }
    }

    /**
     * Get bookmarks for a specific article with pagination (admin use)
     * @param articleId - Article ID
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<Bookmark[]>> - Paginated bookmarks for the article
     * @throws BookmarkApiError - If retrieval fails
     */
    async getBookmarksByArticlePaginated(
        articleId: number,
        params: PaginationParams = {}
    ): Promise<PaginatedResponse<Bookmark[]>> {
        try {
            const { page = 1, limit = 10 } = params;
            const response = await apiClient.get<ApiResponse<Bookmark[]>>(
                `/articles/${articleId}/bookmarks/paginated`,
                {
                    params: { page, limit }
                }
            );
            return handlePaginatedResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch paginated bookmarks for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching paginated bookmarks for article ${articleId}`
            );
        }
    }

    // ========================================================================
    // Bookmark Statistics
    // ========================================================================

    /**
     * Get bookmark count for a specific article
     * @param articleId - Article ID
     * @returns Promise<number> - Number of bookmarks on the article
     * @throws BookmarkApiError - If retrieval fails
     */
    async getBookmarkCountByArticle(articleId: number): Promise<number> {
        try {
            const response = await apiClient.get<ApiResponse<{ count: number }>>(
                `/articles/${articleId}/bookmarks/count`
            );
            const result = handleResponse(response);
            return result.count;
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch bookmark count for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching bookmark count for article ${articleId}`
            );
        }
    }

    /**
     * Get bookmark count by a specific user
     * @param userId - User ID
     * @returns Promise<number> - Number of bookmarks by the user
     * @throws BookmarkApiError - If retrieval fails
     */
    async getBookmarkCountByUser(userId: number): Promise<number> {
        try {
            const response = await apiClient.get<ApiResponse<{ count: number }>>(
                `/users/${userId}/bookmarks/count`
            );
            const result = handleResponse(response);
            return result.count;
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch bookmark count for user ${userId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching bookmark count for user ${userId}`
            );
        }
    }

    // ========================================================================
    // Toggle Bookmark (Convenience Method)
    // ========================================================================

    /**
     * Toggle bookmark status for an article (add if not bookmarked, remove if bookmarked)
     * @param articleId - Article ID
     * @returns Promise<BookmarkStatus> - Updated bookmark status
     * @throws BookmarkApiError - If toggle operation fails
     */
    async toggleBookmark(articleId: number): Promise<BookmarkStatus> {
        try {
            // First check current status
            const currentStatus = await this.checkBookmarkStatus(articleId);

            if (currentStatus.is_bookmarked) {
                // Remove bookmark if currently bookmarked
                await this.removeBookmark(articleId);
            } else {
                // Add bookmark if not currently bookmarked
                await this.addBookmark({ article_id: articleId });
            }

            // Return updated status
            return await this.checkBookmarkStatus(articleId);
        } catch (error: any) {
            if (error instanceof BookmarkApiError) {
                throw error;
            }
            throw new BookmarkApiError(
                'TOGGLE_FAILED',
                `Failed to toggle bookmark for article ${articleId}`
            );
        }
    }

    // ========================================================================
    // Get Specific Bookmark
    // ========================================================================

    /**
     * Get a specific bookmark by ID
     * @param bookmarkId - Bookmark ID
     * @returns Promise<Bookmark> - Bookmark data
     * @throws BookmarkApiError - If retrieval fails or bookmark not found
     */
    async getBookmarkById(bookmarkId: number): Promise<Bookmark> {
        try {
            const response = await apiClient.get<ApiResponse<Bookmark>>(`/bookmarks/${bookmarkId}`);
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new BookmarkApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch bookmark with ID ${bookmarkId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new BookmarkApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching bookmark with ID ${bookmarkId}`
            );
        }
    }
}

// ============================================================================
// Export Singleton Instance
// ============================================================================

export const bookmarkApi = new BookmarkApi();
export default bookmarkApi;