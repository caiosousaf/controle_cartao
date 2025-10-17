import api from './api';
import { Card, CardResponse } from '../types/card';

export const cardService = {
  getCards: async (limit = 10, offset = 0, getTotal = false): Promise<CardResponse> => {
    try {
      const params = new URLSearchParams();
      
      if (limit) params.append('limit', limit.toString());
      if (offset) params.append('offset', offset.toString());
      if (getTotal) params.append('total', 'true');
      
      const response = await api.get<CardResponse>('/cadastros/cartoes', { params });
      return response.data;
    } catch (error) {
      console.error('Error fetching cards:', error);
      throw error;
    }
  },
};