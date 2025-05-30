// src/types/index.ts

export interface User {
    id: number;
    name: string;
    email: string;
    role: string; // Misalnya 'user', 'admin'
    // Hindari memasukkan password di tipe data frontend
    created_at?: string; // Opsional tergantung API
    updated_at?: string; // Opsional tergantung API
}

export interface Category {
    id: number;
    name: string;
    slug?: string; // Berguna untuk URL kategori
}

export interface Tag {
    id: number;
    name: string;
    slug?: string; // Berguna untuk URL tag
}

export interface Article {
    id: number;
    title: string;
    slug: string;
    content: string;
    user: User; // Penulis artikel
    category: Category | null; // Kategori bisa null
    tags?: Tag[]; // Artikel bisa punya banyak tag, opsional
    views: number;
    created_at: string;
    updated_at: string;
}

// Untuk respons API yang terstruktur (jika semua endpoint sama)
export interface ApiResponse<T> {
    data: T;
    message?: string;
    // Bisa ditambahkan metadata lain jika API mengembalikannya, e.g., pagination
}