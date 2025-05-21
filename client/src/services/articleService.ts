
import axios from "axios";

const API_URL = "/api/v1"


export const getAllArticles = async () => {
    try {
        const res = await axios.get(`${API_URL}/article`);
        return {
            articles: res.data.data,
            message: res.data.message,
        };
    } catch (error) {
        console.error('Error di getAllArticles:', error);
        throw error;
    }
};


export const getArticleBySlug = async (slug: string) => {
    const res = await axios.get(`${API_URL}/articles/${slug}`)
    return res.data
}