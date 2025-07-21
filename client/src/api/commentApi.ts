// Comment API Service
// This file provides methods for comment management operations with article and user filtering

import type { AxiosResponse } from 'axios';
import { apiClient } from './apiClient';
import type {
    ApiResponse,
    Comment,
    CommentResponse,
    CreateCommentRequest,
    UpdateCommentRequest,
    PaginatedResponse,
    PaginationParams
} from './types';

// ============================================================================
// Response Handler Utility
// ============================================================================

class CommentApiError extends Error {
    constructor(
        public code: string,
        message: string,
        public statusCode?: number,
        public details?: Record<string, any>
    ) {
        super(message);
        this.name = 'CommentApiError';
    }
}

function handleResponse<T>(response: AxiosResponse<ApiResponse<T>>): T {
    if (response.data.success && response.data.data !== undefined) {
        return response.data.data;
    }

    const error = response.data.error;
    throw new CommentApiError(
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
    throw new CommentApiError(
        error?.code || 'UNKNOWN',
        error?.message || 'Unknown error occurred',
        response.status,
        error?.details
    );
}

// ============================================================================
// Comment API Service Class
// ============================================================================

export class CommentApi {

    // ========================================================================
    // Create Comment
    // ========================================================================

    /**
     * Create a new comment on an article
     * @param commentData - Comment creation data with article association
     * @returns Promise<Comment> - Created comment
     * @throws CommentApiError - If creation fails
     */
    async createComment(commentData: CreateCommentRequest): Promise<Comment> {
        try {
            const response = await apiClient.post<ApiResponse<Comment>>(
                '/comments',
                commentData
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'CREATE_FAILED',
                    error.response.data?.error?.message || 'Failed to create comment',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                'Network error occurred while creating comment'
            );
        }
    }

    // ========================================================================
    // Get Comments by Article
    // ========================================================================

    /**
     * Get all comments for a specific article (non-paginated)
     * @param articleId - Article ID
     * @returns Promise<CommentResponse[]> - Array of comments for the article
     * @throws CommentApiError - If retrieval fails
     */
    async getCommentsByArticle(articleId: number): Promise<CommentResponse[]> {
        try {
            const response = await apiClient.get<ApiResponse<CommentResponse[]>>(
                `/articles/${articleId}/comments`
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch comments for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching comments for article ${articleId}`
            );
        }
    }

    /**
     * Get comments for a specific article with pagination
     * @param articleId - Article ID
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<CommentResponse[]>> - Paginated comments for the article
     * @throws CommentApiError - If retrieval fails
     */
    async getCommentsByArticlePaginated(
        articleId: number,
        params: PaginationParams = {}
    ): Promise<PaginatedResponse<CommentResponse[]>> {
        try {
            const { page = 1, limit = 10 } = params;
            const response = await apiClient.get<ApiResponse<CommentResponse[]>>(
                `/articles/${articleId}/comments/paginated`,
                {
                    params: { page, limit }
                }
            );
            return handlePaginatedResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch paginated comments for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching paginated comments for article ${articleId}`
            );
        }
    }

    // ========================================================================
    // Get Comments by User
    // ========================================================================

    /**
     * Get all comments by a specific user (non-paginated)
     * @param userId - User ID
     * @returns Promise<CommentResponse[]> - Array of comments by the user
     * @throws CommentApiError - If retrieval fails
     */
    async getCommentsByUser(userId: number): Promise<CommentResponse[]> {
        try {
            const response = await apiClient.get<ApiResponse<CommentResponse[]>>(
                `/users/${userId}/comments`
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch comments by user ${userId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching comments by user ${userId}`
            );
        }
    }

    /**
     * Get comments by a specific user with pagination
     * @param userId - User ID
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<CommentResponse[]>> - Paginated comments by the user
     * @throws CommentApiError - If retrieval fails
     */
    async getCommentsByUserPaginated(
        userId: number,
        params: PaginationParams = {}
    ): Promise<PaginatedResponse<CommentResponse[]>> {
        try {
            const { page = 1, limit = 10 } = params;
            const response = await apiClient.get<ApiResponse<CommentResponse[]>>(
                `/users/${userId}/comments/paginated`,
                {
                    params: { page, limit }
                }
            );
            return handlePaginatedResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch paginated comments by user ${userId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching paginated comments by user ${userId}`
            );
        }
    }

    // ========================================================================
    // Get All Comments
    // ========================================================================

    /**
     * Get all comments (admin/moderation use)
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<CommentResponse[]>> - Paginated comments
     * @throws CommentApiError - If retrieval fails
     */
    async getAllComments(params: PaginationParams = {}): Promise<PaginatedResponse<CommentResponse[]>> {
        try {
            const { page = 1, limit = 10 } = params;
            const response = await apiClient.get<ApiResponse<CommentResponse[]>>(
                '/comments',
                {
                    params: { page, limit }
                }
            );
            return handlePaginatedResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || 'Failed to fetch all comments',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                'Network error occurred while fetching all comments'
            );
        }
    }

    // ========================================================================
    // Get Comment by ID
    // ========================================================================

    /**
     * Get a specific comment by ID
     * @param id - Comment ID
     * @returns Promise<Comment> - Comment data
     * @throws CommentApiError - If retrieval fails or comment not found
     */
    async getCommentById(id: number): Promise<Comment> {
        try {
            const response = await apiClient.get<ApiResponse<Comment>>(`/comments/${id}`);
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch comment with ID ${id}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching comment with ID ${id}`
            );
        }
    }

    // ========================================================================
    // Update Comment (with ownership checks)
    // ========================================================================

    /**
     * Update an existing comment (requires ownership or admin privileges)
     * @param id - Comment ID to update
     * @param commentData - Updated comment data
     * @returns Promise<Comment> - Updated comment
     * @throws CommentApiError - If update fails, comment not found, or insufficient permissions
     */
    async updateComment(id: number, commentData: UpdateCommentRequest): Promise<Comment> {
        try {
            const response = await apiClient.put<ApiResponse<Comment>>(
                `/comments/${id}`,
                commentData
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'UPDATE_FAILED',
                    error.response.data?.error?.message || `Failed to update comment with ID ${id}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                `Network error occurred while updating comment with ID ${id}`
            );
        }
    }

    // ========================================================================
    // Delete Comment (with ownership checks)
    // ========================================================================

    /**
     * Delete a comment (requires ownership or admin privileges)
     * @param id - Comment ID to delete
     * @returns Promise<void> - Resolves when deletion is complete
     * @throws CommentApiError - If deletion fails, comment not found, or insufficient permissions
     */
    async deleteComment(id: number): Promise<void> {
        try {
            await apiClient.delete(`/comments/${id}`);
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'DELETE_FAILED',
                    error.response.data?.error?.message || `Failed to delete comment with ID ${id}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                `Network error occurred while deleting comment with ID ${id}`
            );
        }
    }

    // ========================================================================
    // Comment Statistics and Filtering
    // ========================================================================

    /**
     * Get comment count for a specific article
     * @param articleId - Article ID
     * @returns Promise<number> - Number of comments on the article
     * @throws CommentApiError - If retrieval fails
     */
    async getCommentCountByArticle(articleId: number): Promise<number> {
        try {
            const response = await apiClient.get<ApiResponse<{ count: number }>>(
                `/articles/${articleId}/comments/count`
            );
            const result = handleResponse(response);
            return result.count;
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch comment count for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching comment count for article ${articleId}`
            );
        }
    }

    /**
     * Get comment count by a specific user
     * @param userId - User ID
     * @returns Promise<number> - Number of comments by the user
     * @throws CommentApiError - If retrieval fails
     */
    async getCommentCountByUser(userId: number): Promise<number> {
        try {
            const response = await apiClient.get<ApiResponse<{ count: number }>>(
                `/users/${userId}/comments/count`
            );
            const result = handleResponse(response);
            return result.count;
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch comment count for user ${userId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching comment count for user ${userId}`
            );
        }
    }

    // ========================================================================
    // Recent Comments
    // ========================================================================

    /**
     * Get recent comments across all articles
     * @param limit - Maximum number of recent comments to fetch (default: 10)
     * @returns Promise<CommentResponse[]> - Array of recent comments
     * @throws CommentApiError - If retrieval fails
     */
    async getRecentComments(limit: number = 10): Promise<CommentResponse[]> {
        try {
            const response = await apiClient.get<ApiResponse<CommentResponse[]>>(
                '/comments/recent',
                {
                    params: { limit }
                }
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CommentApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || 'Failed to fetch recent comments',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CommentApiError(
                'NETWORK_ERROR',
                'Network error occurred while fetching recent comments'
            );
        }
    }
}

// ============================================================================
// Export Singleton Instance
// ============================================================================

export const commentApi = new CommentApi();
export default commentApi;