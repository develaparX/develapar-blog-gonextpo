// Simple validation test for TagApi implementation
// This file validates that the TagApi class has all required methods and proper types

import { tagApi, TagApi } from './tagApi';
import type {
    Tag,
    CreateTagRequest,
    UpdateTagRequest,
    AssignTagsByNameRequest,
    PaginationParams,
    PaginatedResponse
} from './types';

// ============================================================================
// Type Validation Tests
// ============================================================================

// Test that tagApi is an instance of TagApi
const isTagApiInstance: boolean = tagApi instanceof TagApi;
console.log('tagApi is instance of TagApi:', isTagApiInstance);

// ============================================================================
// Method Signature Validation
// ============================================================================

// Validate that all required methods exist with correct signatures
const validateTagApiMethods = () => {
    // Basic CRUD operations
    const createTagMethod: (tagData: CreateTagRequest) => Promise<Tag> = tagApi.createTag.bind(tagApi);
    const getAllTagsMethod: () => Promise<Tag[]> = tagApi.getAllTags.bind(tagApi);
    const getAllTagsPaginatedMethod: (params?: PaginationParams) => Promise<PaginatedResponse<Tag[]>> = tagApi.getAllTagsPaginated.bind(tagApi);
    const getTagByIdMethod: (id: number) => Promise<Tag> = tagApi.getTagById.bind(tagApi);
    const updateTagMethod: (id: number, tagData: UpdateTagRequest) => Promise<Tag> = tagApi.updateTag.bind(tagApi);
    const deleteTagMethod: (id: number) => Promise<void> = tagApi.deleteTag.bind(tagApi);

    // Tag-Article association methods
    const assignTagsByNameMethod: (assignmentData: AssignTagsByNameRequest) => Promise<Tag[]> = tagApi.assignTagsByName.bind(tagApi);
    const getTagsByArticleMethod: (articleId: number) => Promise<Tag[]> = tagApi.getTagsByArticle.bind(tagApi);
    const removeAllTagsFromArticleMethod: (articleId: number) => Promise<void> = tagApi.removeAllTagsFromArticle.bind(tagApi);
    const removeTagFromArticleMethod: (articleId: number, tagId: number) => Promise<void> = tagApi.removeTagFromArticle.bind(tagApi);

    // Search and filter methods
    const searchTagsMethod: (query: string, limit?: number) => Promise<Tag[]> = tagApi.searchTags.bind(tagApi);
    const getPopularTagsMethod: (limit?: number) => Promise<Tag[]> = tagApi.getPopularTags.bind(tagApi);

    console.log('All TagApi methods have correct signatures');

    return {
        createTagMethod,
        getAllTagsMethod,
        getAllTagsPaginatedMethod,
        getTagByIdMethod,
        updateTagMethod,
        deleteTagMethod,
        assignTagsByNameMethod,
        getTagsByArticleMethod,
        removeAllTagsFromArticleMethod,
        removeTagFromArticleMethod,
        searchTagsMethod,
        getPopularTagsMethod
    };
};

// ============================================================================
// Type Interface Validation
// ============================================================================

// Validate request/response types
const validateTagTypes = () => {
    // CreateTagRequest validation
    const createTagRequest: CreateTagRequest = {
        name: 'Technology'
    };

    // UpdateTagRequest validation
    const updateTagRequest: UpdateTagRequest = {
        name: 'Updated Technology'
    };

    // AssignTagsByNameRequest validation
    const assignTagsByNameRequest: AssignTagsByNameRequest = {
        article_id: 1,
        tags: ['tech', 'programming']
    };

    // Tag validation
    const tag: Tag = {
        id: 1,
        name: 'Technology'
    };

    // PaginationParams validation
    const paginationParams: PaginationParams = {
        page: 1,
        limit: 10
    };

    console.log('All TagApi types are properly defined');

    return {
        createTagRequest,
        updateTagRequest,
        assignTagsByNameRequest,
        tag,
        paginationParams
    };
};

// ============================================================================
// Run Validations
// ============================================================================

try {
    validateTagApiMethods();
    validateTagTypes();
    console.log('✅ TagApi validation completed successfully');
} catch (error) {
    console.error('❌ TagApi validation failed:', error);
}

// ============================================================================
// Export for potential use in other tests
// ============================================================================

export {
    isTagApiInstance,
    validateTagApiMethods,
    validateTagTypes
};