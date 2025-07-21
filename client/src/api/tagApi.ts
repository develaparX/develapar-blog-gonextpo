// Tag API Service
// This file provides methods for tag management operations and article-tag associations

import type { AxiosResponse } from 'axios';
import { apiClient } from './apiClient';
import type {
    ApiResponse,
    Tag,
    CreateTagRequest,
    UpdateTagRequest,
    AssignTagsByNameRequest,
    PaginatedResponse,
    PaginationParams
} from './types';

// ============================================================================
// Response Handler Utility
// ============================================================================

class TagApiError extends Error {
    constructor(
        public code: string,
        message: string,
        public statusCode?: number,
        public details?: Record<string, any>
    ) {
        super(message);
        this.name = 'TagApiError';
    }
}

function handleResponse<T>(response: AxiosResponse<ApiResponse<T>>): T {
    if (response.data.success && response.data.data !== undefined) {
        return response.data.data;
    }

    const error = response.data.error;
    throw new TagApiError(
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
    throw new TagApiError(
        error?.code || 'UNKNOWN',
        error?.message || 'Unknown error occurred',
        response.status,
        error?.details
    );
}

// ============================================================================
// Tag API Service Class
// ============================================================================

export class TagApi {

    // ========================================================================
    // Create Tag
    // ========================================================================

    /**
     * Create a new tag
     * @param tagData - Tag creation data
     * @returns Promise<Tag> - Created tag
     * @throws TagApiError - If creation fails
     */
    async createTag(tagData: CreateTagRequest): Promise<Tag> {
        try {
            const response = await apiClient.post<ApiResponse<Tag>>(
                '/tags',
                tagData
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'CREATE_FAILED',
                    error.response.data?.error?.message || 'Failed to create tag',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                'Network error occurred while creating tag'
            );
        }
    }

    // ========================================================================
    // Get All Tags
    // ========================================================================

    /**
     * Get all tags (non-paginated) - useful for tag selection dropdowns
     * @returns Promise<Tag[]> - Array of all tags
     * @throws TagApiError - If retrieval fails
     */
    async getAllTags(): Promise<Tag[]> {
        try {
            const response = await apiClient.get<ApiResponse<Tag[]>>('/tags');
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || 'Failed to fetch tags',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                'Network error occurred while fetching tags'
            );
        }
    }

    /**
     * Get all tags with pagination
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<Tag[]>> - Paginated tags
     * @throws TagApiError - If retrieval fails
     */
    async getAllTagsPaginated(params: PaginationParams = {}): Promise<PaginatedResponse<Tag[]>> {
        try {
            const { page = 1, limit = 10 } = params;
            const response = await apiClient.get<ApiResponse<Tag[]>>(
                '/tags/paginated',
                {
                    params: { page, limit }
                }
            );
            return handlePaginatedResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || 'Failed to fetch paginated tags',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                'Network error occurred while fetching paginated tags'
            );
        }
    }

    // ========================================================================
    // Get Tag by ID
    // ========================================================================

    /**
     * Get a specific tag by ID
     * @param id - Tag ID
     * @returns Promise<Tag> - Tag data
     * @throws TagApiError - If retrieval fails or tag not found
     */
    async getTagById(id: number): Promise<Tag> {
        try {
            const response = await apiClient.get<ApiResponse<Tag>>(`/tags/${id}`);
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch tag with ID ${id}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching tag with ID ${id}`
            );
        }
    }

    // ========================================================================
    // Update Tag
    // ========================================================================

    /**
     * Update an existing tag
     * @param id - Tag ID to update
     * @param tagData - Updated tag data
     * @returns Promise<Tag> - Updated tag
     * @throws TagApiError - If update fails or tag not found
     */
    async updateTag(id: number, tagData: UpdateTagRequest): Promise<Tag> {
        try {
            const response = await apiClient.put<ApiResponse<Tag>>(
                `/tags/${id}`,
                tagData
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'UPDATE_FAILED',
                    error.response.data?.error?.message || `Failed to update tag with ID ${id}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                `Network error occurred while updating tag with ID ${id}`
            );
        }
    }

    // ========================================================================
    // Delete Tag
    // ========================================================================

    /**
     * Delete a tag
     * @param id - Tag ID to delete
     * @returns Promise<void> - Resolves when deletion is complete
     * @throws TagApiError - If deletion fails or tag not found
     */
    async deleteTag(id: number): Promise<void> {
        try {
            await apiClient.delete(`/tags/${id}`);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'DELETE_FAILED',
                    error.response.data?.error?.message || `Failed to delete tag with ID ${id}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                `Network error occurred while deleting tag with ID ${id}`
            );
        }
    }

    // ========================================================================
    // Tag-Article Association Methods
    // ========================================================================

    /**
     * Assign tags to an article by tag names
     * @param assignmentData - Article ID and tag names
     * @returns Promise<Tag[]> - Array of assigned tags
     * @throws TagApiError - If assignment fails
     */
    async assignTagsByName(assignmentData: AssignTagsByNameRequest): Promise<Tag[]> {
        try {
            const response = await apiClient.post<ApiResponse<Tag[]>>(
                '/tags/assign-by-name',
                assignmentData
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'ASSIGNMENT_FAILED',
                    error.response.data?.error?.message || 'Failed to assign tags to article',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                'Network error occurred while assigning tags to article'
            );
        }
    }

    /**
     * Get all tags associated with a specific article
     * @param articleId - Article ID
     * @returns Promise<Tag[]> - Array of tags associated with the article
     * @throws TagApiError - If retrieval fails
     */
    async getTagsByArticle(articleId: number): Promise<Tag[]> {
        try {
            const response = await apiClient.get<ApiResponse<Tag[]>>(`/articles/${articleId}/tags`);
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch tags for article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching tags for article ${articleId}`
            );
        }
    }

    /**
     * Remove all tags from an article
     * @param articleId - Article ID
     * @returns Promise<void> - Resolves when removal is complete
     * @throws TagApiError - If removal fails
     */
    async removeAllTagsFromArticle(articleId: number): Promise<void> {
        try {
            await apiClient.delete(`/articles/${articleId}/tags`);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'REMOVAL_FAILED',
                    error.response.data?.error?.message || `Failed to remove tags from article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                `Network error occurred while removing tags from article ${articleId}`
            );
        }
    }

    /**
     * Remove a specific tag from an article
     * @param articleId - Article ID
     * @param tagId - Tag ID to remove
     * @returns Promise<void> - Resolves when removal is complete
     * @throws TagApiError - If removal fails
     */
    async removeTagFromArticle(articleId: number, tagId: number): Promise<void> {
        try {
            await apiClient.delete(`/articles/${articleId}/tags/${tagId}`);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'REMOVAL_FAILED',
                    error.response.data?.error?.message || `Failed to remove tag ${tagId} from article ${articleId}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                `Network error occurred while removing tag ${tagId} from article ${articleId}`
            );
        }
    }

    // ========================================================================
    // Search and Filter Methods
    // ========================================================================

    /**
     * Search tags by name (useful for autocomplete)
     * @param query - Search query
     * @param limit - Maximum number of results (default: 10)
     * @returns Promise<Tag[]> - Array of matching tags
     * @throws TagApiError - If search fails
     */
    async searchTags(query: string, limit: number = 10): Promise<Tag[]> {
        try {
            const response = await apiClient.get<ApiResponse<Tag[]>>('/tags/search', {
                params: { q: query, limit }
            });
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'SEARCH_FAILED',
                    error.response.data?.error?.message || 'Failed to search tags',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                'Network error occurred while searching tags'
            );
        }
    }

    /**
     * Get popular tags (most used tags)
     * @param limit - Maximum number of results (default: 10)
     * @returns Promise<Tag[]> - Array of popular tags
     * @throws TagApiError - If retrieval fails
     */
    async getPopularTags(limit: number = 10): Promise<Tag[]> {
        try {
            const response = await apiClient.get<ApiResponse<Tag[]>>('/tags/popular', {
                params: { limit }
            });
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new TagApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || 'Failed to fetch popular tags',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new TagApiError(
                'NETWORK_ERROR',
                'Network error occurred while fetching popular tags'
            );
        }
    }
}

// ============================================================================
// Export Singleton Instance
// ============================================================================

export const tagApi = new TagApi();
export default tagApi;