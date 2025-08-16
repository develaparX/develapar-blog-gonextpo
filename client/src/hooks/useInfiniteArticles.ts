import { useInfiniteQuery } from '@tanstack/react-query';
import { apiService } from '@/services/api';
import { queryKeys } from '@/lib/queryKeys';
import type { Article } from '@/services/api';

interface PaginatedResponse<T> {
    data: T[];
    pagination: {
        page: number;
        limit: number;
        total: number;
        total_pages: number;
        has_next: boolean;
        has_prev: boolean;
    };
}

// Infinite scrolling for articles
export function useInfiniteArticles(limit: number = 10) {
    return useInfiniteQuery({
        queryKey: [...queryKeys.articles.lists(), { limit }],
        queryFn: async ({ pageParam = 1 }) => {
            // This assumes your API supports pagination
            // You'd need to update your API service to support this
            const response = await fetch(
                `http://localhost:4300/api/v1/articles/paginated?page=${pageParam}&limit=${limit}`
            );
            const data = await response.json();

            if (data.success) {
                return {
                    articles: data.data.articles || [],
                    pagination: data.pagination,
                };
            }
            throw new Error('Failed to fetch articles');
        },
        initialPageParam: 1,
        getNextPageParam: (lastPage) => {
            return lastPage.pagination.has_next
                ? lastPage.pagination.page + 1
                : undefined;
        },
        getPreviousPageParam: (firstPage) => {
            return firstPage.pagination.has_prev
                ? firstPage.pagination.page - 1
                : undefined;
        },
        staleTime: 1000 * 60 * 5, // 5 minutes
    });
}

// Infinite scrolling for articles by category
export function useInfiniteArticlesByCategory(categoryName: string, limit: number = 10) {
    return useInfiniteQuery({
        queryKey: [...queryKeys.articles.byCategory(categoryName), { limit }],
        queryFn: async ({ pageParam = 1 }) => {
            // This assumes your API supports pagination for category
            const response = await fetch(
                `http://localhost:4300/api/v1/articles/category/${encodeURIComponent(categoryName)}/paginated?page=${pageParam}&limit=${limit}`
            );
            const data = await response.json();

            if (data.success) {
                return {
                    articles: data.data.articles || [],
                    pagination: data.pagination,
                };
            }
            throw new Error(`Failed to fetch articles for category: ${categoryName}`);
        },
        initialPageParam: 1,
        getNextPageParam: (lastPage) => {
            return lastPage.pagination.has_next
                ? lastPage.pagination.page + 1
                : undefined;
        },
        enabled: !!categoryName,
        staleTime: 1000 * 60 * 3, // 3 minutes
    });
}

// Hook to get flattened articles from infinite query
export function useFlattenedArticles(infiniteQuery: ReturnType<typeof useInfiniteArticles>) {
    const articles = infiniteQuery.data?.pages.flatMap(page => page.articles) || [];

    return {
        articles,
        hasNextPage: infiniteQuery.hasNextPage,
        fetchNextPage: infiniteQuery.fetchNextPage,
        isFetchingNextPage: infiniteQuery.isFetchingNextPage,
        isLoading: infiniteQuery.isLoading,
        error: infiniteQuery.error,
    };
}