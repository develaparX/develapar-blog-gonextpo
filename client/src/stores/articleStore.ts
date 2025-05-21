import { create } from 'zustand';

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
    notification?: string;
    setArticles: (articles: Article[]) => void;
    setNotification: (message: string) => void;
}

export const useArticleStore = create<ArticleState>((set) => ({
    articles: [],
    notification: undefined,
    setArticles: (articles) => set({ articles }),
    setNotification: (message) => set({ notification: message }),
}));
