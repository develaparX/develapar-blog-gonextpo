// API Base Configuration
const API_BASE_URL = 'http://localhost:4300/api/v1';

// Types for API responses
export interface APIResponse<T> {
    success: boolean;
    data: T;
    error?: {
        code: string;
        message: string;
        details?: Record<string, string>;
    };
    meta?: {
        request_id: string;
        processing_time_ms: number;
        version: string;
        timestamp: string;
    };
    pagination?: {
        page: number;
        limit: number;
        total: number;
        total_pages: number;
        has_next: boolean;
        has_prev: boolean;
    };
}

// Entity Types
export interface Article {
    id: number;
    title: string;
    slug: string;
    content: string;
    user: {
        id: number;
        name: string;
        email: string;
    };
    category: {
        id: number;
        name: string;
    };
    views: number;
    created_at: string;
    updated_at: string;
}

export interface Category {
    id: number;
    name: string;
    description?: string;
}

export interface Tag {
    id: number;
    name: string;
    description?: string;
}

// Generic API fetch function
async function apiFetch<T>(endpoint: string): Promise<APIResponse<T>> {
    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data: APIResponse<T> = await response.json();
        return data;
    } catch (error) {
        console.error(`API Error for ${endpoint}:`, error);
        throw error;
    }
}

// API Service Class
class ApiService {
    // Articles API
    async getAllArticles(): Promise<Article[]> {
        const response = await apiFetch<{ articles: Article[]; message: string }>('/articles');
        if (response.success) {
            return response.data.articles || [];
        }
        throw new Error('Failed to fetch articles');
    }

    async getArticlesByCategory(categoryName: string): Promise<Article[]> {
        const response = await apiFetch<{ articles: Article[]; message: string }>(
            `/articles/category/${encodeURIComponent(categoryName)}`
        );
        if (response.success) {
            return response.data.articles || [];
        }
        throw new Error(`Failed to fetch articles for category: ${categoryName}`);
    }

    async getArticlesByTag(tagId: number): Promise<Article[]> {
        const response = await apiFetch<{ articles: Article[]; message: string }>(
            `/tags/${tagId}/articles`
        );
        if (response.success) {
            return response.data.articles || [];
        }
        throw new Error(`Failed to fetch articles for tag ID: ${tagId}`);
    }

    async getArticleBySlug(slug: string): Promise<Article> {
        const response = await apiFetch<{ article: Article; message: string }>(
            `/articles/${slug}`
        );
        if (response.success) {
            return response.data.article;
        }
        throw new Error(`Failed to fetch article: ${slug}`);
    }

    // Categories API
    async getAllCategories(): Promise<Category[]> {
        const response = await apiFetch<{ categories: Category[]; message: string }>('/categories/');
        if (response.success) {
            return response.data.categories || [];
        }
        throw new Error('Failed to fetch categories');
    }

    async getCategoryById(categoryId: number): Promise<Category> {
        const response = await apiFetch<{ category: Category; message: string }>(
            `/categories/${categoryId}`
        );
        if (response.success) {
            return response.data.category;
        }
        throw new Error(`Failed to fetch category ID: ${categoryId}`);
    }

    // Tags API
    async getAllTags(): Promise<Tag[]> {
        const response = await apiFetch<{ tags: Tag[]; message: string }>('/tags');
        if (response.success) {
            return response.data.tags || [];
        }
        throw new Error('Failed to fetch tags');
    }

    async getTagById(tagId: number): Promise<Tag> {
        const response = await apiFetch<{ tag: Tag; message: string }>(`/tags/${tagId}`);
        if (response.success) {
            return response.data.tag;
        }
        throw new Error(`Failed to fetch tag ID: ${tagId}`);
    }

    // Search functionality (client-side filtering)
    async searchArticles(query: string): Promise<Article[]> {
        const articles = await this.getAllArticles();
        const searchQuery = query.toLowerCase().trim();

        return articles.filter(article =>
            article.title.toLowerCase().includes(searchQuery) ||
            article.content.toLowerCase().includes(searchQuery)
        );
    }
}

// Export singleton instance
export const apiService = new ApiService();

// Export individual functions for convenience
export const {
    getAllArticles,
    getArticlesByCategory,
    getArticlesByTag,
    getArticleBySlug,
    getAllCategories,
    getCategoryById,
    getAllTags,
    getTagById,
    searchArticles
} = apiService;