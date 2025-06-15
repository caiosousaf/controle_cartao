import api from './api';
import { UserResponse, UpdateUserData } from '../types/user';

export const userService = {
  getUser: async (): Promise<UserResponse> => {
    try {
      const response = await api.get<UserResponse>('/cadastros/usuario');
      return response.data;
    } catch (error) {
      console.error('Error fetching user:', error);
      throw error;
    }
  },

  updateUser: async (data: UpdateUserData): Promise<void> => {
    try {
      await api.put('/cadastros/usuarios/alterar/senha', data);
    } catch (error) {
      console.error('Error updating user:', error);
      throw error;
    }
  }
};