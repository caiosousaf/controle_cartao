import api from './api';

interface LoginCredentials {
  email: string;
  senha: string;
}

interface LoginResponse {
  token: string;
}

export const authService = {
  login: async (credentials: LoginCredentials): Promise<string> => {
    try {
      const response = await api.post<LoginResponse>('/usuarios/login', credentials);
      
      // Store the token in local storage
      const token = response.data.token;
      localStorage.setItem('authToken', token);
      
      return token;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  },

  logout: (): void => {
    localStorage.removeItem('authToken');
  },

  isAuthenticated: (): boolean => {
    return !!localStorage.getItem('authToken');
  }
};