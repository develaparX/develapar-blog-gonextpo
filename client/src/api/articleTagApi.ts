
import axios from 'axios';

const API_URL = '/api/v1/article-tags';

export const getArticleTags = async () => {
  try {
    const response = await axios.get(API_URL);
    return response.data;
  } catch (error) {
    console.error('Error fetching article tags:', error);
    throw error;
  }
};
