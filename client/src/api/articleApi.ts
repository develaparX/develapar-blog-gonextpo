
import axios from 'axios';

const API_URL = '/api/v1/articles';

export const getArticles = async () => {
  try {
    const response = await axios.get(API_URL);
    return response.data;
  } catch (error) {
    console.error('Error fetching articles:', error);
    throw error;
  }
};

// Add other article-related API calls here (e.g., getArticleById, createArticle, updateArticle, deleteArticle)
