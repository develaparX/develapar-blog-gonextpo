// Category API Tests
// This file contains unit tests for the category API service

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { categoryApi, CategoryApi } from './categoryApi';
import { apiClient } from './apiClient';
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from './types';

// Mock the apiClient
vi.mock('./apiClient', () => ({
    apiClient: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn(),
    }
}));

describe('CategoryApi', () => {
    let mockApiClient: any;

    beforeEach(() => {
        mockApiClient = apiClient as any;
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.resetAllMocks();
    });

    describe('createCategory', () => {
        it('should create a category successfully', async () => {
            const categoryData: CreateCategoryRequest = { name: 'Technology' };
            const expectedCategory: Category = { id: 1, name: 'Technology' };

            mockApiClient.post.mockResolvedValue({
                data: {
                    success: true,
                    data: expectedCategory
                }
            });

            const result = await categoryApi.createCategory(categoryData);

            expect(mockApiClient.post).toHaveBeenCalledWith('/categories', categoryData);
            expect(result).toEqual(expectedCategory);
        });

        it('should handle creation errors', async () => {
            const categoryData: CreateCategoryRequest = { name: 'Technology' };

            mockApiClient.post.mockRejectedValue({
                response: {
                    status: 400,
                    data: {
                        error: {
                            code: 'VALIDATION_ERROR',
                            message: 'Category name already exists'
                        }
                    }
                }
            });

            await expect(categoryApi.createCategory(categoryData)).rejects.toThrow('Category name already exists');
        });

        it('should handle network errors during creation', async () => {
            const categoryData: CreateCategoryRequest = { name: 'Technology' };

            mockApiClient.post.mockRejectedValue(new Error('Network error'));

            await expect(categoryApi.createCategory(categoryData)).rejects.toThrow('Network error occurred while creating category');
        });
    });

    describe('getAllCategories', () => {
        it('should fetch all categories successfully', async () => {
            const expectedCategories: Category[] = [
                { id: 1, name: 'Technology' },
                { id: 2, name: 'Science' }
            ];

            mockApiClient.get.mockResolvedValue({
                data: {
                    success: true,
                    data: expectedCategories
                }
            });

            const result = await categoryApi.getAllCategories();

            expect(mockApiClient.get).toHaveBeenCalledWith('/categories');
            expect(result).toEqual(expectedCategories);
        });

        it('should handle fetch errors', async () => {
            mockApiClient.get.mockRejectedValue({
                response: {
                    status: 500,
                    data: {
                        error: {
                            code: 'INTERNAL_ERROR',
                            message: 'Database connection failed'
                        }
                    }
                }
            });

            await expect(categoryApi.getAllCategories()).rejects.toThrow('Database connection failed');
        });
    });

    describe('getAllCategoriesPaginated', () => {
        it('should fetch paginated categories successfully', async () => {
            const expectedCategories: Category[] = [
                { id: 1, name: 'Technology' }
            ];
            const expectedPagination = {
                page: 1,
                limit: 10,
                total: 1,
                total_pages: 1,
                has_next: false,
                has_prev: false
            };

            mockApiClient.get.mockResolvedValue({
                data: {
                    success: true,
                    data: expectedCategories,
                    pagination: expectedPagination
                }
            });

            const result = await categoryApi.getAllCategoriesPaginated({ page: 1, limit: 10 });

            expect(mockApiClient.get).toHaveBeenCalledWith('/categories/paginated', {
                params: { page: 1, limit: 10 }
            });
            expect(result.data).toEqual(expectedCategories);
            expect(result.pagination).toEqual(expectedPagination);
        });

        it('should use default pagination parameters', async () => {
            const expectedCategories: Category[] = [];

            mockApiClient.get.mockResolvedValue({
                data: {
                    success: true,
                    data: expectedCategories
                }
            });

            await categoryApi.getAllCategoriesPaginated();

            expect(mockApiClient.get).toHaveBeenCalledWith('/categories/paginated', {
                params: { page: 1, limit: 10 }
            });
        });
    });

    describe('getCategoryById', () => {
        it('should fetch category by ID successfully', async () => {
            const expectedCategory: Category = { id: 1, name: 'Technology' };

            mockApiClient.get.mockResolvedValue({
                data: {
                    success: true,
                    data: expectedCategory
                }
            });

            const result = await categoryApi.getCategoryById(1);

            expect(mockApiClient.get).toHaveBeenCalledWith('/categories/1');
            expect(result).toEqual(expectedCategory);
        });

        it('should handle not found errors', async () => {
            mockApiClient.get.mockRejectedValue({
                response: {
                    status: 404,
                    data: {
                        error: {
                            code: 'NOT_FOUND',
                            message: 'Category not found'
                        }
                    }
                }
            });

            await expect(categoryApi.getCategoryById(999)).rejects.toThrow('Category not found');
        });
    });

    describe('updateCategory', () => {
        it('should update category successfully', async () => {
            const updateData: UpdateCategoryRequest = { name: 'Updated Technology' };
            const expectedCategory: Category = { id: 1, name: 'Updated Technology' };

            mockApiClient.put.mockResolvedValue({
                data: {
                    success: true,
                    data: expectedCategory
                }
            });

            const result = await categoryApi.updateCategory(1, updateData);

            expect(mockApiClient.put).toHaveBeenCalledWith('/categories/1', updateData);
            expect(result).toEqual(expectedCategory);
        });

        it('should handle update errors', async () => {
            const updateData: UpdateCategoryRequest = { name: 'Updated Technology' };

            mockApiClient.put.mockRejectedValue({
                response: {
                    status: 403,
                    data: {
                        error: {
                            code: 'FORBIDDEN',
                            message: 'Insufficient permissions'
                        }
                    }
                }
            });

            await expect(categoryApi.updateCategory(1, updateData)).rejects.toThrow('Insufficient permissions');
        });
    });

    describe('deleteCategory', () => {
        it('should delete category successfully', async () => {
            mockApiClient.delete.mockResolvedValue({});

            await categoryApi.deleteCategory(1);

            expect(mockApiClient.delete).toHaveBeenCalledWith('/categories/1');
        });

        it('should handle delete errors', async () => {
            mockApiClient.delete.mockRejectedValue({
                response: {
                    status: 409,
                    data: {
                        error: {
                            code: 'CONFLICT',
                            message: 'Cannot delete category with associated articles'
                        }
                    }
                }
            });

            await expect(categoryApi.deleteCategory(1)).rejects.toThrow('Cannot delete category with associated articles');
        });

        it('should handle network errors during deletion', async () => {
            mockApiClient.delete.mockRejectedValue(new Error('Network error'));

            await expect(categoryApi.deleteCategory(1)).rejects.toThrow('Network error occurred while deleting category with ID 1');
        });
    });
});