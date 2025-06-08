import api from './api';
import { RecurringPurchase, RecurringPurchaseResponse, ExpenseEstimateResponse } from '../types/recurring';

interface CreateRecurringPurchaseData {
  nome: string;
  descricao?: string;
  compra_categoria_id: string;
  local_compra?: string;
  valor_parcela: number;
}

export const recurringService = {
  getRecurringPurchases: async (limit = 10, offset = 0, getTotal = false): Promise<RecurringPurchaseResponse> => {
    try {
      const params = new URLSearchParams();
      
      if (limit) params.append('limit', limit.toString());
      if (offset) params.append('offset', offset.toString());
      if (getTotal) params.append('total', 'true');
      
      const response = await api.get<RecurringPurchaseResponse>('/cadastros/compras/recorrente', { params });
      return response.data;
    } catch (error) {
      console.error('Error fetching recurring purchases:', error);
      throw error;
    }
  },

  createRecurringPurchase: async (data: CreateRecurringPurchaseData): Promise<void> => {
    try {
      await api.post('/cadastros/compras/recorrente/cadastro', data);
    } catch (error) {
      console.error('Error creating recurring purchase:', error);
      throw error;
    }
  },

  registerRecurringPurchases: async (): Promise<void> => {
    try {
      await api.post('/cadastros/compras/recorrente');
    } catch (error) {
      console.error('Error registering recurring purchases:', error);
      throw error;
    }
  },

  updateRecurringPurchase: async (id: string, data: CreateRecurringPurchaseData): Promise<void> => {
    try {
      await api.put(`/cadastros/compras/recorrente/${id}`, data);
    } catch (error) {
      console.error('Error updating recurring purchase:', error);
      throw error;
    }
  },

  deactivateRecurringPurchase: async (id: string): Promise<void> => {
    try {
      await api.put(`/cadastros/compras/recorrente/${id}/desativar`);
    } catch (error) {
      console.error('Error deactivating recurring purchase:', error);
      throw error;
    }
  },

  reactivateRecurringPurchase: async (id: string): Promise<void> => {
    try {
      await api.put(`/cadastros/compras/recorrente/${id}/reativar`);
    } catch (error) {
      console.error('Error reactivating recurring purchase:', error);
      throw error;
    }
  },

  removeRecurringPurchase: async (id: string): Promise<void> => {
    try {
      await api.delete(`/cadastros/compras/recorrente/${id}/remover`);
    } catch (error) {
      console.error('Error removing recurring purchase:', error);
      throw error;
    }
  },

  getEstimate: async (): Promise<ExpenseEstimateResponse> => {
    try {
      const response = await api.get<ExpenseEstimateResponse>('/cadastros/compras/recorrente/previsao');
      return response.data;
    } catch (error) {
      console.error('Error fetching estimate:', error);
      throw error;
    }
  }
};