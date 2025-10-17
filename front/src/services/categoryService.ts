import api from './api';
import { CategoryResponse, CreateCategoryData, UpdateCategoryData } from '../types/category';

export const categoryService = {
  getCategories: async (limit = 10, offset = 0, getTotal = false, ativo?: boolean): Promise<CategoryResponse> => {
    try {
      const params = new URLSearchParams();
      
      if (limit) params.append('limit', limit.toString());
      if (offset) params.append('offset', offset.toString());
      if (getTotal) params.append('total', 'true');
      if (ativo !== undefined) params.append('ativo', ativo.toString());
      
      const response = await api.get<CategoryResponse>('/cadastros/categorias', { params });
      return response.data;
    } catch (error) {
      console.error('Error fetching categories:', error);
      throw error;
    }
  },

  createCategory: async (data: CreateCategoryData): Promise<void> => {
    try {
      await api.post('/cadastros/categorias', data);
    } catch (error) {
      console.error('Error creating category:', error);
      throw error;
    }
  },

  updateCategory: async (categoryId: string, data: UpdateCategoryData): Promise<void> => {
    try {
      await api.put(`/cadastros/categoria/${categoryId}`, data);
    } catch (error) {
      console.error('Error updating category:', error);
      throw error;
    }
  },

  deactivateCategory: async (categoryId: string): Promise<void> => {
    try {
      await api.put(`/cadastros/categoria/${categoryId}/remover`);
    } catch (error) {
      console.error('Error deactivating category:', error);
      throw error;
    }
  },

  reactivateCategory: async (categoryId: string): Promise<void> => {
    try {
      await api.put(`/cadastros/categoria/${categoryId}/reativar`);
    } catch (error) {
      console.error('Error reactivating category:', error);
      throw error;
    }
  }
};