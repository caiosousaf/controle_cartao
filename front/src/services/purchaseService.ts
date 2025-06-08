import api from './api';
import { PurchaseResponse, TotalPurchaseResponse } from '../types/purchase';

interface CreatePurchaseData {
  nome: string;
  descricao?: string;
  local_compra?: string;
  categoria_id: string;
  valor_parcela: number;
  parcela_atual: number;
  quantidade_parcelas: number;
  fatura_id: string;
  data_compra: string;
}

export const purchaseService = {
  getPurchasesByInvoiceId: async (
    invoiceId: string,
    limit = 10,
    offset = 0,
    getTotal = false
  ): Promise<PurchaseResponse> => {
    try {
      const params = new URLSearchParams();
      
      params.append('fatura_id', invoiceId);
      if (limit) params.append('limit', limit.toString());
      if (offset) params.append('offset', offset.toString());
      if (getTotal) params.append('total', 'true');
      
      const response = await api.get<PurchaseResponse>('/cadastros/compras', { params });
      return response.data;
    } catch (error) {
      console.error(`Error fetching purchases for invoice ${invoiceId}:`, error);
      throw error;
    }
  },

  getTotalPurchases: async (
    dataEspecifica?: string,
    ultimaParcela?: boolean,
    pago?: boolean
  ): Promise<TotalPurchaseResponse> => {
    try {
      const params = new URLSearchParams();
      
      if (dataEspecifica) params.append('data_especifica', dataEspecifica);
      if (ultimaParcela !== undefined) params.append('ultima_parcela', ultimaParcela.toString());
      if (pago !== undefined) params.append('pago', pago.toString());
      
      const response = await api.get<TotalPurchaseResponse>('/cadastros/compras/total', { params });
      return response.data;
    } catch (error) {
      console.error('Error fetching total purchases:', error);
      throw error;
    }
  },

  createPurchase: async (invoiceId: string, data: CreatePurchaseData): Promise<void> => {
    try {
      await api.post(`/cadastros/fatura/${invoiceId}/compras`, data);
    } catch (error) {
      console.error('Error creating purchase:', error);
      throw error;
    }
  },
};