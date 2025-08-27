import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiService } from '@/services/api';
import { queryKeys } from '@/lib/queryKeys';

// Articles Queries
export function useArticles() {
    return useQuery({
        queryKey: queryKeys.articles.lists(),
        queryFn: () => apiService.getAllArticles(),
        staleTime: 1000 * 60 * 5, // 5 minutes
    });
}

export function usePaginatedArticles(page: number = 1, limit: number = 10) {
    return useQuery({
        queryKey: [...queryKeys.articles.lists(), 'paginated', { page, limit }],
        queryFn: () => apiService.getPaginatedArticles({ page, limit }),
        staleTime: 1000 * 60 * 5, // 5 minutes
        placeholderData: (previousData) => previousData, // Keep previous data while loading new page
    });
}

export function useArticlesByCategory(categoryName: string | undefined) {
    return useQuery({
        queryKey: queryKeys.articles.byCategory(categoryName || ''),
        queryFn: () => categoryName ? apiService.getArticlesByCategory(categoryName) : Promise.resolve([]),
        enabled: !!categoryName, // Only run query if categoryName exists
        staleTime: 1000 * 60 * 3, // 3 minutes
    });
}

export function useArticlesByTag(tagId: number | undefined) {
    return useQuery({
        queryKey: queryKeys.articles.byTag(tagId || 0),
        queryFn: () => tagId ? apiService.getArticlesByTag(tagId) : Promise.resolve([]),
        enabled: !!tagId, // Only run query if tagId exists
        staleTime: 1000 * 60 * 3, // 3 minutes
    });
}

export function useArticleBySlug(slug: string | undefined) {
    return useQuery({
        queryKey: queryKeys.articles.detail(slug || ''),
        queryFn: () => {
            if (!slug) throw new Error('No slug provided');
            return apiService.getArticleBySlug(slug);
        },
        enabled: !!slug, // Only run query if slug exists
        staleTime: 1000 * 60 * 10, // 10 minutes - article details don't change often
    });
}

// Categories Queries
export function useCategories() {
    return useQuery({
        queryKey: queryKeys.categories.lists(),
        queryFn: () => apiService.getAllCategories(),
        staleTime: 1000 * 60 * 10, // 10 minutes - categories don't change often
    });
}

export function useCategoryById(categoryId: number | undefined) {
    return useQuery({
        queryKey: queryKeys.categories.detail(categoryId || 0),
        queryFn: () => {
            if (!categoryId) throw new Error('No category ID provided');
            return apiService.getCategoryById(categoryId);
        },
        enabled: !!categoryId,
        staleTime: 1000 * 60 * 10, // 10 minutes
    });
}

// Tags Queries
export function useTags() {
    return useQuery({
        queryKey: queryKeys.tags.lists(),
        queryFn: () => apiService.getAllTags(),
        staleTime: 1000 * 60 * 10, // 10 minutes - tags don't change often
    });
}

export function useTagById(tagId: number | undefined) {
    return useQuery({
        queryKey: queryKeys.tags.detail(tagId || 0),
        queryFn: () => {
            if (!tagId) throw new Error('No tag ID provided');
            return apiService.getTagById(tagId);
        },
        enabled: !!tagId,
        staleTime: 1000 * 60 * 10, // 10 minutes
    });
}

// Products Queries
export function useProductsByArticleId(articleId: string | undefined) {
    return useQuery({
        queryKey: queryKeys.products.byArticle(articleId || ''),
        queryFn: () => {
            if (!articleId) throw new Error('No article ID provided');
            return apiService.getProductsByArticleId(articleId);
        },
        enabled: !!articleId, // Only run query if articleId exists
        staleTime: 1000 * 60 * 5, // 5 minutes - products don't change often
    });
}

// Search functionality with manual trigger
export function useSearchArticles(query: string, enabled: boolean = false) {
    return useQuery({
        queryKey: queryKeys.articles.search(query),
        queryFn: () => apiService.searchArticles(query),
        enabled: enabled && !!query.trim(), // Only run when enabled and query exists
        staleTime: 1000 * 60 * 2, // 2 minutes - search results can be more dynamic
    });
}

// Manual search hook with imperative control
export function useManualSearch() {
    const queryClient = useQueryClient();

    const searchMutation = useMutation({
        mutationFn: (query: string) => apiService.searchArticles(query),
        onSuccess: (data, query) => {
            // Cache the search results
            queryClient.setQueryData(queryKeys.articles.search(query), data);
        },
    });

    const search = (query: string) => {
        if (!query.trim()) {
            return Promise.resolve([]);
        }
        return searchMutation.mutateAsync(query);
    };

    const clear = () => {
        searchMutation.reset();
    };

    return {
        data: searchMutation.data || [],
        isLoading: searchMutation.isPending,
        error: searchMutation.error?.message || null,
        search,
        clear,
        isSuccess: searchMutation.isSuccess,
        isError: searchMutation.isError,
    };
}

// Utility hooks for cache management
export function useInvalidateArticles() {
    const queryClient = useQueryClient();

    return {
        invalidateAll: () => queryClient.invalidateQueries({ queryKey: queryKeys.articles.all }),
        invalidateByCategory: (categoryName: string) =>
            queryClient.invalidateQueries({ queryKey: queryKeys.articles.byCategory(categoryName) }),
        invalidateByTag: (tagId: number) =>
            queryClient.invalidateQueries({ queryKey: queryKeys.articles.byTag(tagId) }),
        invalidateSearch: () =>
            queryClient.invalidateQueries({ queryKey: [...queryKeys.articles.all, 'search'] }),
    };
}

export function useInvalidateCategories() {
    const queryClient = useQueryClient();

    return {
        invalidateAll: () => queryClient.invalidateQueries({ queryKey: queryKeys.categories.all }),
        invalidateDetail: (categoryId: number) =>
            queryClient.invalidateQueries({ queryKey: queryKeys.categories.detail(categoryId) }),
    };
}

export function useInvalidateTags() {
    const queryClient = useQueryClient();

    return {
        invalidateAll: () => queryClient.invalidateQueries({ queryKey: queryKeys.tags.all }),
        invalidateDetail: (tagId: number) =>
            queryClient.invalidateQueries({ queryKey: queryKeys.tags.detail(tagId) }),
    };
}

export function useInvalidateProducts() {
    const queryClient = useQueryClient();

    return {
        invalidateAll: () => queryClient.invalidateQueries({ queryKey: queryKeys.products.all }),
        invalidateByArticle: (articleId: string) =>
            queryClient.invalidateQueries({ queryKey: queryKeys.products.byArticle(articleId) }),
        invalidateDetail: (productId: string) =>
            queryClient.invalidateQueries({ queryKey: queryKeys.products.detail(productId) }),
    };
}

// Prefetch hooks for performance optimization
export function usePrefetchArticlesByCategory() {
    const queryClient = useQueryClient();

    return (categoryName: string) => {
        queryClient.prefetchQuery({
            queryKey: queryKeys.articles.byCategory(categoryName),
            queryFn: () => apiService.getArticlesByCategory(categoryName),
            staleTime: 1000 * 60 * 3,
        });
    };
}

export function usePrefetchArticlesByTag() {
    const queryClient = useQueryClient();

    return (tagId: number) => {
        queryClient.prefetchQuery({
            queryKey: queryKeys.articles.byTag(tagId),
            queryFn: () => apiService.getArticlesByTag(tagId),
            staleTime: 1000 * 60 * 3,
        });
    };
}