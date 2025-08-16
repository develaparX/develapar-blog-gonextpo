import { useMutation, useQueryClient } from '@tanstack/react-query';
import { apiService } from '@/services/api';
import { queryKeys } from '@/lib/queryKeys';
import type { Article, Category, Tag } from '@/services/api';

// Optimistic updates for better UX
export function useOptimisticArticleUpdate() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ articleId, updates }: { articleId: number; updates: Partial<Article> }) => {
            // This would be your actual API call
            return Promise.resolve({ ...updates, id: articleId });
        },
        onMutate: async ({ articleId, updates }) => {
            // Cancel any outgoing refetches
            await queryClient.cancelQueries({ queryKey: queryKeys.articles.all });

            // Snapshot the previous value
            const previousArticles = queryClient.getQueryData(queryKeys.articles.lists());

            // Optimistically update to the new value
            queryClient.setQueryData(queryKeys.articles.lists(), (old: Article[] | undefined) => {
                if (!old) return [];
                return old.map(article =>
                    article.id === articleId ? { ...article, ...updates } : article
                );
            });

            // Return a context object with the snapshotted value
            return { previousArticles };
        },
        onError: (err, variables, context) => {
            // If the mutation fails, use the context returned from onMutate to roll back
            if (context?.previousArticles) {
                queryClient.setQueryData(queryKeys.articles.lists(), context.previousArticles);
            }
        },
        onSettled: () => {
            // Always refetch after error or success
            queryClient.invalidateQueries({ queryKey: queryKeys.articles.all });
        },
    });
}

// Optimistic like/unlike functionality
export function useOptimisticLike() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ articleId, isLiked }: { articleId: number; isLiked: boolean }) => {
            // Your actual like/unlike API call would go here
            return Promise.resolve({ articleId, isLiked });
        },
        onMutate: async ({ articleId, isLiked }) => {
            // Optimistically update the UI immediately
            const updateArticleLikes = (articles: Article[] | undefined) => {
                if (!articles) return [];
                return articles.map(article => {
                    if (article.id === articleId) {
                        return {
                            ...article,
                            // Assuming you have a likes count field
                            likes: isLiked ? (article.likes || 0) + 1 : Math.max(0, (article.likes || 0) - 1),
                            isLiked
                        };
                    }
                    return article;
                });
            };

            // Update all relevant queries
            queryClient.setQueryData(queryKeys.articles.lists(), updateArticleLikes);

            // Update category-specific queries if they exist
            const queryCache = queryClient.getQueryCache();
            queryCache.getAll().forEach(query => {
                if (query.queryKey[0] === 'articles' && query.queryKey[1] === 'category') {
                    queryClient.setQueryData(query.queryKey, updateArticleLikes);
                }
            });
        },
        onError: () => {
            // Revert optimistic updates on error
            queryClient.invalidateQueries({ queryKey: queryKeys.articles.all });
        },
    });
}