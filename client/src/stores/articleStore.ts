import { create } from 'zustand';
import { getAllArticles } from '../services/articleService';

export interface Article {
    id: number;
    title: string;
    slug: string;
    content: string;
    user: {
        id: number;
        name: string;
        // tambahkan field lain jika perlu
    };
    category: {
        id: number;
        name: string;
        // tambahkan field lain jika perlu
    };
    views: number;
    created_at: string;  // karena dari API biasanya date jadi string
    updated_at: string;
}


interface ArticleState {
    articles: Article[];
    fetchArticles: () => Promise<void>;
    notification?: string
}

export const useArticleStore = create<ArticleState>((set) => ({
    articles: [],
    notification: undefined,
    fetchArticles: async () => {
        try {
            const { articles, message } = await getAllArticles();
            set({ articles, notification: message });
        } catch (error) {
            console.error("Error fetchArticles:", error);
        }
    },
}));
