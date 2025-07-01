
import axios from 'axios';

const API_URL = '/api/v1/tags';

export const getTags = async () => {
  try {
    const response = await axios.get(API_URL);
    return response.data;
  } catch (error) {
    console.error('Error fetching tags:', error);
    throw error;
  }
};
