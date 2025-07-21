// Article API Module
// This file provides comprehensive article management operations

import type { AxiosResponse } from 'axios';
import { apiClient } from './apiClient';
import { ResponseHandler } from './errorHandler';
import type {
    ApiResponse,
    ArticleWithTags,
    CreateArticleRequest,
    UpdateArticleRequest,
    PaginationParams,
    PaginatedResponse,
    ArticleSearchParams,
    ArticleFilters
} from './types';

// ============================================================================
// Article API Class
// ============================================================================

export class ArticleApi {
    // ========================================================================
    // CRUD Operations (Task 4.1)
    // ========================================================================

    /**
     * Create a new article with tag association
     * Requirements: 4.1, 4.2, 4.3, 4.4, 4.5
     */
    async createArticle(articleData: CreateArticleRequest): Promise<ArticleWithTags> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags>> = await apiClient.post(
                '/articles',
                articleData
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get article by slug for public access
     * Requirements: 4.1, 4.2, 4.3, 4.4, 4.5
     */
    async getArticleBySlug(slug: string): Promise<ArticleWithTags> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags>> = await apiClient.get(
                `/articles/${slug}`
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Update article with ownership verification
     * Requirements: 4.1, 4.2, 4.3, 4.4, 4.5
     */
    async updateArticle(id: number, articleData: UpdateArticleRequest): Promise<ArticleWithTags> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags>> = await apiClient.put(
                `/articles/${id}`,
                articleData
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Delete article with proper authorization
     * Requirements: 4.1, 4.2, 4.3, 4.4, 4.5
     */
    async deleteArticle(id: number): Promise<void> {
        try {
            await apiClient.delete(`/articles/${id}`);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    // ========================================================================
    // Filtering and Pagination Operations (Task 4.2)
    // ========================================================================

    /**
     * Get all articles with pagination support
     * Requirements: 4.2, 4.3
     */
    async getAllArticles(params?: PaginationParams): Promise<ArticleWithTags[]> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                '/articles',
                { params }
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get all articles with pagination metadata
     * Requirements: 4.2, 4.3
     */
    async getAllArticlesPaginated(params?: PaginationParams): Promise<PaginatedResponse<ArticleWithTags[]>> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                '/articles/paginated',
                { params }
            );
            return ResponseHandler.handlePaginatedSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get articles by category with filtering
     * Requirements: 4.2, 4.3
     */
    async getArticlesByCategory(categoryName: string, params?: PaginationParams): Promise<ArticleWithTags[]> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                `/articles/category/${encodeURIComponent(categoryName)}`,
                { params }
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get articles by category with pagination metadata
     * Requirements: 4.2, 4.3
     */
    async getArticlesByCategoryPaginated(
        categoryName: string,
        params?: PaginationParams
    ): Promise<PaginatedResponse<ArticleWithTags[]>> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                `/articles/category/${encodeURIComponent(categoryName)}/paginated`,
                { params }
            );
            return ResponseHandler.handlePaginatedSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get articles by author
     * Requirements: 4.2, 4.3
     */
    async getArticlesByAuthor(userId: number, params?: PaginationParams): Promise<ArticleWithTags[]> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                `/articles/author/${userId}`,
                { params }
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get articles by author with pagination metadata
     * Requirements: 4.2, 4.3
     */
    async getArticlesByAuthorPaginated(
        userId: number,
        params?: PaginationParams
    ): Promise<PaginatedResponse<ArticleWithTags[]>> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                `/articles/author/${userId}/paginated`,
                { params }
            );
            return ResponseHandler.handlePaginatedSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Search articles with comprehensive filtering
     * Requirements: 4.2, 4.3
     */
    async searchArticles(searchParams: ArticleSearchParams): Promise<ArticleWithTags[]> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                '/articles/search',
                { params: searchParams }
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Search articles with pagination metadata
     * Requirements: 4.2, 4.3
     */
    async searchArticlesPaginated(searchParams: ArticleSearchParams): Promise<PaginatedResponse<ArticleWithTags[]>> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                '/articles/search/paginated',
                { params: searchParams }
            );
            return ResponseHandler.handlePaginatedSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    // ========================================================================
    // Additional Utility Methods
    // ========================================================================

    /**
     * Get article by ID (for authenticated access)
     */
    async getArticleById(id: number): Promise<ArticleWithTags> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags>> = await apiClient.get(
                `/articles/id/${id}`
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get articles by tag
     */
    async getArticlesByTag(tagName: string, params?: PaginationParams): Promise<ArticleWithTags[]> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                `/articles/tag/${encodeURIComponent(tagName)}`,
                { params }
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get articles by multiple tags
     */
    async getArticlesByTags(tags: string[], params?: PaginationParams): Promise<ArticleWithTags[]> {
        try {
            const queryParams = {
                ...params,
                tags: tags.join(',')
            };
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                '/articles/tags',
                { params: queryParams }
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get recent articles
     */
    async getRecentArticles(limit: number = 10): Promise<ArticleWithTags[]> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                '/articles/recent',
                { params: { limit } }
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get popular articles (by views)
     */
    async getPopularArticles(limit: number = 10): Promise<ArticleWithTags[]> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                '/articles/popular',
                { params: { limit } }
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Increment article view count
     */
    async incrementViews(slug: string): Promise<void> {
        try {
            await apiClient.post(`/articles/${slug}/view`);
        } catch (error) {
            // Don't throw error for view tracking failures
            console.warn('Failed to increment article views:', error);
        }
    }

    // ========================================================================
    // Advanced Filtering Methods
    // ========================================================================

    /**
     * Get articles with advanced filtering options
     */
    async getArticlesWithFilters(filters: ArticleFilters & PaginationParams): Promise<PaginatedResponse<ArticleWithTags[]>> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.get(
                '/articles/filter',
                { params: filters }
            );
            return ResponseHandler.handlePaginatedSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get article count by filters
     */
    async getArticleCount(filters?: ArticleFilters): Promise<number> {
        try {
            const response: AxiosResponse<ApiResponse<{ count: number }>> = await apiClient.get(
                '/articles/count',
                { params: filters }
            );
            const result = ResponseHandler.handleSuccess(response);
            return result.count;
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    // ========================================================================
    // Batch Operations
    // ========================================================================

    /**
     * Get multiple articles by IDs
     */
    async getArticlesByIds(ids: number[]): Promise<ArticleWithTags[]> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.post(
                '/articles/batch',
                { ids }
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Bulk update articles (admin only)
     */
    async bulkUpdateArticles(updates: Array<{ id: number; data: UpdateArticleRequest }>): Promise<ArticleWithTags[]> {
        try {
            const response: AxiosResponse<ApiResponse<ArticleWithTags[]>> = await apiClient.put(
                '/articles/bulk',
                { updates }
            );
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Bulk delete articles (admin only)
     */
    async bulkDeleteArticles(ids: number[]): Promise<void> {
        try {
            await apiClient.delete('/articles/bulk', {
                data: { ids }
            });
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }
}

// ============================================================================
// Export Singleton Instance
// ============================================================================

export const articleApi = new ArticleApi();
export default articleApi;