
import axios from 'axios';

const API_URL = '/api/v1/bookmarks';

export const getBookmarks = async () => {
  try {
    const response = await axios.get(API_URL);
    return response.data;
  } catch (error) {
    console.error('Error fetching bookmarks:', error);
    throw error;
  }
};
