import { useState, useEffect } from 'react';
import { apiService } from '@/services/api';
import type { Article, Category, Tag } from '@/services/api';

// Generic hook for API calls
export function useApiCall<T>(
    apiCall: () => Promise<T>,
    dependencies: any[] = []
) {
    const [data, setData] = useState<T | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        let isMounted = true;

        const fetchData = async () => {
            try {
                setLoading(true);
                setError(null);
                const result = await apiCall();

                if (isMounted) {
                    setData(result);
                }
            } catch (err) {
                if (isMounted) {
                    setError(err instanceof Error ? err.message : 'An error occurred');
                    console.error('API call error:', err);
                }
            } finally {
                if (isMounted) {
                    setLoading(false);
                }
            }
        };

        fetchData();

        return () => {
            isMounted = false;
        };
    }, dependencies);

    const refetch = async () => {
        try {
            setLoading(true);
            setError(null);
            const result = await apiCall();
            setData(result);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            console.error('API refetch error:', err);
        } finally {
            setLoading(false);
        }
    };

    return { data, loading, error, refetch };
}

// Specific hooks for different API calls
export function useArticles() {
    return useApiCall(() => apiService.getAllArticles());
}

export function useCategories() {
    return useApiCall(() => apiService.getAllCategories());
}

export function useTags() {
    return useApiCall(() => apiService.getAllTags());
}

export function useArticlesByCategory(categoryName: string | undefined) {
    return useApiCall(
        () => categoryName ? apiService.getArticlesByCategory(categoryName) : Promise.resolve([]),
        [categoryName]
    );
}

export function useArticlesByTag(tagId: number | undefined) {
    return useApiCall(
        () => tagId ? apiService.getArticlesByTag(tagId) : Promise.resolve([]),
        [tagId]
    );
}

export function useArticleBySlug(slug: string | undefined) {
    return useApiCall(
        () => slug ? apiService.getArticleBySlug(slug) : Promise.reject(new Error('No slug provided')),
        [slug]
    );
}

export function useSearchArticles(query: string) {
    return useApiCall(
        () => query.trim() ? apiService.searchArticles(query) : Promise.resolve([]),
        [query]
    );
}

// Hook for manual search (doesn't auto-execute)
export function useManualSearch() {
    const [data, setData] = useState<Article[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const search = async (query: string) => {
        if (!query.trim()) {
            setData([]);
            return;
        }

        try {
            setLoading(true);
            setError(null);
            const results = await apiService.searchArticles(query);
            setData(results);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Search failed');
            console.error('Search error:', err);
        } finally {
            setLoading(false);
        }
    };

    const clear = () => {
        setData([]);
        setError(null);
    };

    return { data, loading, error, search, clear };
}