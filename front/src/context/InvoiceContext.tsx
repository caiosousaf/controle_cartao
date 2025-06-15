import React, { createContext, useState, useContext, ReactNode } from 'react';
import { Invoice } from '../types/invoice';
import { invoiceService } from '../services/invoiceService';
import { useAuth } from './AuthContext';

interface InvoiceContextType {
  invoices: Invoice[];
  selectedInvoice: Invoice | null;
  isLoading: boolean;
  error: string | null;
  hasMore: boolean;
  totalInvoices: number;
  fetchInvoicesByCardId: (cardId: string, limit?: number, offset?: number, pago?: boolean) => Promise<void>;
  selectInvoice: (invoice: Invoice) => void;
  clearSelectedInvoice: () => void;
  clearInvoices: () => void;
  clearError: () => void;
}

const InvoiceContext = createContext<InvoiceContextType | undefined>(undefined);

export const InvoiceProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [totalInvoices, setTotalInvoices] = useState(0);
  
  const { isAuthenticated } = useAuth();

  const fetchInvoicesByCardId = async (cardId: string, limit = 10, offset = 0, pago?: boolean) => {
    if (!isAuthenticated || !cardId) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await invoiceService.getInvoicesByCardId(cardId, limit, offset, false, pago);
      
      if (offset === 0) {
        setInvoices(response.dados);
      } else {
        setInvoices(prev => [...prev, ...response.dados]);
      }
      
      setHasMore(response.prox || false);
      
      // Fetch total count if it's the first page
      if (offset === 0) {
        fetchTotalInvoices(cardId, pago);
      }
    } catch (err) {
      setError('Erro ao carregar faturas');
      console.error(`Error fetching invoices for card ${cardId}:`, err);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchTotalInvoices = async (cardId: string, pago?: boolean) => {
    if (!isAuthenticated || !cardId) return;
    
    try {
      const response = await invoiceService.getInvoicesByCardId(cardId, 1, 0, true, pago);
      setTotalInvoices(response.total || 0);
    } catch (err) {
      console.error(`Error fetching total invoices count for card ${cardId}:`, err);
    }
  };

  const selectInvoice = (invoice: Invoice) => {
    setSelectedInvoice(invoice);
  };

  const clearSelectedInvoice = () => {
    setSelectedInvoice(null);
  };

  const clearInvoices = () => {
    setInvoices([]);
    setSelectedInvoice(null);
    setTotalInvoices(0);
  };

  const clearError = () => setError(null);

  const value = {
    invoices,
    selectedInvoice,
    isLoading,
    error,
    hasMore,
    totalInvoices,
    fetchInvoicesByCardId,
    selectInvoice,
    clearSelectedInvoice,
    clearInvoices,
    clearError
  };

  return <InvoiceContext.Provider value={value}>{children}</InvoiceContext.Provider>;
};

export const useInvoices = (): InvoiceContextType => {
  const context = useContext(InvoiceContext);
  if (context === undefined) {
    throw new Error('useInvoices must be used within an InvoiceProvider');
  }
  return context;
};