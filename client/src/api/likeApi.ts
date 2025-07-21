// Like API Service
// This file provides methods for like management operations with status checking

import type { AxiosResponse } from 'axios';
import { apiClient } from './apiClient';
import type {
    ApiResponse,
    Like,
    CreateLikeRequest,
    LikeStatus,
    PaginatedResponse,
    PaginationParams
} from './types';

// ============================================================================
// Response Handler Utility
// ============================================================================

class LikeApiError extends Error {
    constructor(
        public code: string,
        message: string,
        public statusCode?: number,
        public details?: Record<string, any>
    ) {
        super(message);
        this.name = 'LikeApiError';
    }
}

function handleResponse<T>(response: AxiosResponse<ApiResponse<T>>): T {
    if (response.data.success && response.data.data !== undefined) {
        return response.data.data;
    }

    const error = response.data.error;
    throw new LikeApiError(
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
    throw new LikeApiError(
        error?.code || 'UNKNOWN',
        error?.message || 'Unknown error occurred',
        response.status,
        error?.details
    );
}

// ============================================================================
// Like API Service Class
// ============================================================================

export class LikeApi {

    // ========================================================================
    // Add Like
    // ========================================================================

    /**
     * Add a like to an article
     * @param likeData - Like creation data with article ID
     * @returns Promise<Like> - Created like
     * @throws LikeApiError - If like creation fails or already exists
     */
    async addLike(likeData: CreateLikeRequest): Promise<Like> {
        try {
            const response = await apiClient.post<ApiResponse<Like>>(
                '/likes',
                likeData
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new LikeApiError(
                    error.response.data?.error?.code || 'CREATE_FAILED',
                    error.response.data?.error?.message || 'Failed to add like',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new LikeApiError(
                'NETWORK_ERROR',
                'Network error occurred while adding like'
            );
        }
    }

    // ========================================================================
    // Remove Like
    // ========================================================================

    /**
     * Remove a like from an article
     * @param articleId - Article ID to remove like from
     * @returns Promise<void> - Resolves when like is removed
     * @throws LikeApiError - If like removal fails or like doesn't exist
     */
    async removeLike(articleId: number): Promise<void> {
        try {
            await apiClient.delete(`/articles/${articleId}/likes`);
        } catch (error: any) {
            if (error.response) {
                throw new LikeApiError(
                    error.response.data?.error?.code || 'DELETE_FAILED',
                    error.response.data?.error?.message || `Failed to remove like from article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new LikeApiError(
                'NETWORK_ERROR',
                `Network error occurred while removing like from article ${articleId}`
            );
        }
    }

    /**
     * Remove a specific like by ID (admin use)
     * @param likeId - Like ID to remove
     * @returns Promise<void> - Resolves when like is removed
     * @throws LikeApiError - If like removal fails or like doesn't exist
     */
    async removeLikeById(likeId: number): Promise<void> {
        try {
            await apiClient.delete(`/likes/${likeId}`);
        } catch (error: any) {
            if (error.response) {
                throw new LikeApiError(
                    error.response.data?.error?.code || 'DELETE_FAILED',
                    error.response.data?.error?.message || `Failed to remove like with ID ${likeId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new LikeApiError(
                'NETWORK_ERROR',
                `Network error occurred while removing like with ID ${likeId}`
            );
        }
    }

    // ========================================================================
    // Check Like Status
    // ========================================================================

    /**
     * Check if current user has liked an article and get like count
     * @param articleId - Article ID to check
     * @returns Promise<LikeStatus> - Like status and count information
     * @throws LikeApiError - If status check fails
     */
    async checkLikeStatus(articleId: number): Promise<LikeStatus> {
        try {
            const response = await apiClient.get<ApiResponse<LikeStatus>>(
                `/articles/${articleId}/likes/status`
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new LikeApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to check like status for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new LikeApiError(
                'NETWORK_ERROR',
                `Network error occurred while checking like status for article ${articleId}`
            );
        }
    }

    // ========================================================================
    // Get Likes by Article
    // ========================================================================

    /**
     * Get all likes for a specific article (non-paginated)
     * @param articleId - Article ID
     * @returns Promise<Like[]> - Array of likes for the article
     * @throws LikeApiError - If retrieval fails
     */
    async getLikesByArticle(articleId: number): Promise<Like[]> {
        try {
            const response = await apiClient.get<ApiResponse<Like[]>>(
                `/articles/${articleId}/likes`
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new LikeApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch likes for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new LikeApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching likes for article ${articleId}`
            );
        }
    }

    /**
     * Get likes for a specific article with pagination
     * @param articleId - Article ID
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<Like[]>> - Paginated likes for the article
     * @throws LikeApiError - If retrieval fails
     */
    async getLikesByArticlePaginated(
        articleId: number,
        params: PaginationParams = {}
    ): Promise<PaginatedResponse<Like[]>> {
        try {
            const { page = 1, limit = 10 } = params;
            const response = await apiClient.get<ApiResponse<Like[]>>(
                `/articles/${articleId}/likes/paginated`,
                {
                    params: { page, limit }
                }
            );
            return handlePaginatedResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new LikeApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch paginated likes for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new LikeApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching paginated likes for article ${articleId}`
            );
        }
    }

    // ========================================================================
    // Get Likes by User
    // ========================================================================

    /**
     * Get all likes by a specific user (non-paginated)
     * @param userId - User ID
     * @returns Promise<Like[]> - Array of likes by the user
     * @throws LikeApiError - If retrieval fails
     */
    async getLikesByUser(userId: number): Promise<Like[]> {
        try {
            const response = await apiClient.get<ApiResponse<Like[]>>(
                `/users/${userId}/likes`
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new LikeApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch likes by user ${userId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new LikeApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching likes by user ${userId}`
            );
        }
    }

    /**
     * Get likes by a specific user with pagination
     * @param userId - User ID
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<Like[]>> - Paginated likes by the user
     * @throws LikeApiError - If retrieval fails
     */
    async getLikesByUserPaginated(
        userId: number,
        params: PaginationParams = {}
    ): Promise<PaginatedResponse<Like[]>> {
        try {
            const { page = 1, limit = 10 } = params;
            const response = await apiClient.get<ApiResponse<Like[]>>(
                `/users/${userId}/likes/paginated`,
                {
                    params: { page, limit }
                }
            );
            return handlePaginatedResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new LikeApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch paginated likes by user ${userId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new LikeApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching paginated likes by user ${userId}`
            );
        }
    }

    // ========================================================================
    // Like Statistics
    // ========================================================================

    /**
     * Get like count for a specific article
     * @param articleId - Article ID
     * @returns Promise<number> - Number of likes on the article
     * @throws LikeApiError - If retrieval fails
     */
    async getLikeCountByArticle(articleId: number): Promise<number> {
        try {
            const response = await apiClient.get<ApiResponse<{ count: number }>>(
                `/articles/${articleId}/likes/count`
            );
            const result = handleResponse(response);
            return result.count;
        } catch (error: any) {
            if (error.response) {
                throw new LikeApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch like count for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new LikeApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching like count for article ${articleId}`
            );
        }
    }

    /**
     * Get like count by a specific user
     * @param userId - User ID
     * @returns Promise<number> - Number of likes by the user
     * @throws LikeApiError - If retrieval fails
     */
    async getLikeCountByUser(userId: number): Promise<number> {
        try {
            const response = await apiClient.get<ApiResponse<{ count: number }>>(
                `/users/${userId}/likes/count`
            );
            const result = handleResponse(response);
            return result.count;
        } catch (error: any) {
            if (error.response) {
                throw new LikeApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch like count for user ${userId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new LikeApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching like count for user ${userId}`
            );
        }
    }

    // ========================================================================
    // Toggle Like (Convenience Method)
    // ========================================================================

    /**
     * Toggle like status for an article (add if not liked, remove if liked)
     * @param articleId - Article ID
     * @returns Promise<LikeStatus> - Updated like status
     * @throws LikeApiError - If toggle operation fails
     */
    async toggleLike(articleId: number): Promise<LikeStatus> {
        try {
            // First check current status
            const currentStatus = await this.checkLikeStatus(articleId);

            if (currentStatus.is_liked) {
                // Remove like if currently liked
                await this.removeLike(articleId);
            } else {
                // Add like if not currently liked
                await this.addLike({ article_id: articleId });
            }

            // Return updated status
            return await this.checkLikeStatus(articleId);
        } catch (error: any) {
            if (error instanceof LikeApiError) {
                throw error;
            }
            throw new LikeApiError(
                'TOGGLE_FAILED',
                `Failed to toggle like for article ${articleId}`
            );
        }
    }
}

// ============================================================================
// Export Singleton Instance
// ============================================================================

export const likeApi = new LikeApi();
export default likeApi;