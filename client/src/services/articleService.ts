
import axios from "axios";

const API_URL = "/api/v1"

export const getAllArticles = async () => {
    try {
        const res = await axios.get(`${API_URL}/article`);
        console.log("Response API:", res.data);
        return {
            articles: res.data.data,
            message: res.data.message,
        };
    } catch (error) {
        console.error("Error di getAllArticles:", error);
        throw error; // agar bisa di-catch di tempat lain
    }
};


export const getArticleBySlug = async (slug: string) => {
    const res = await axios.get(`${API_URL}/articles/${slug}`)
    return res.data
}