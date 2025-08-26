// Query keys for TanStack Query
// Using hierarchical keys for better cache invalidation

export const queryKeys = {
    // Articles
    articles: {
        all: ['articles'] as const,
        lists: () => [...queryKeys.articles.all, 'list'] as const,
        list: (filters: Record<string, any>) => [...queryKeys.articles.lists(), { filters }] as const,
        details: () => [...queryKeys.articles.all, 'detail'] as const,
        detail: (slug: string) => [...queryKeys.articles.details(), slug] as const,
        byCategory: (categoryName: string) => [...queryKeys.articles.all, 'category', categoryName] as const,
        byTag: (tagId: number) => [...queryKeys.articles.all, 'tag', tagId] as const,
        search: (query: string) => [...queryKeys.articles.all, 'search', query] as const,
    },

    // Categories
    categories: {
        all: ['categories'] as const,
        lists: () => [...queryKeys.categories.all, 'list'] as const,
        list: (filters: Record<string, any>) => [...queryKeys.categories.lists(), { filters }] as const,
        details: () => [...queryKeys.categories.all, 'detail'] as const,
        detail: (id: number) => [...queryKeys.categories.details(), id] as const,
    },

    // Tags
    tags: {
        all: ['tags'] as const,
        lists: () => [...queryKeys.tags.all, 'list'] as const,
        list: (filters: Record<string, any>) => [...queryKeys.tags.lists(), { filters }] as const,
        details: () => [...queryKeys.tags.all, 'detail'] as const,
        detail: (id: number) => [...queryKeys.tags.details(), id] as const,
    },

    // Products
    products: {
        all: ['products'] as const,
        lists: () => [...queryKeys.products.all, 'list'] as const,
        list: (filters: Record<string, any>) => [...queryKeys.products.lists(), { filters }] as const,
        details: () => [...queryKeys.products.all, 'detail'] as const,
        detail: (id: string) => [...queryKeys.products.details(), id] as const,
        byArticle: (articleId: string) => [...queryKeys.products.all, 'article', articleId] as const,
    },
} as const;