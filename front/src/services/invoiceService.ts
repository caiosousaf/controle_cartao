import api from './api';
import { InvoiceResponse } from '../types/invoice';

export const invoiceService = {
  getInvoicesByCardId: async (
    cardId: string, 
    limit = 10, 
    offset = 0, 
    getTotal = false,
    pago?: boolean
  ): Promise<InvoiceResponse> => {
    try {
      const params = new URLSearchParams();
      
      if (limit) params.append('limit', limit.toString());
      if (offset) params.append('offset', offset.toString());
      if (getTotal) params.append('total', 'true');
      if (pago !== undefined) params.append('pago', pago.toString());
      
      const response = await api.get<InvoiceResponse>(
        `/cadastros/cartao/${cardId}/faturas`, 
        { params }
      );
      return response.data;
    } catch (error) {
      console.error(`Error fetching invoices for card ${cardId}:`, error);
      throw error;
    }
  },

  payInvoice: async (invoiceId: string): Promise<void> => {
    try {
      await api.put(`/cadastros/fatura/${invoiceId}/status`, {
        status: "Pago"
      });
    } catch (error) {
      console.error('Error paying invoice:', error);
      throw error;
    }
  }
};