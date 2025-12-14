import React, { createContext, useState, useContext, ReactNode } from 'react';
import { Purchase, PurchaseFilters } from '../types/purchase';
import { purchaseService } from '../services/purchaseService';
import { useAuth } from './AuthContext';

interface PurchaseContextType {
  purchases: Purchase[];
  isLoading: boolean;
  error: string | null;
  hasMore: boolean;
  totalPurchases: number;
  totalAmount: number | null;
  fetchPurchasesByInvoiceId: (invoiceId: string, limit?: number, offset?: number) => Promise<void>;
  fetchTotalPurchases: (filters?: PurchaseFilters) => Promise<void>;
  clearPurchases: () => void;
  clearError: () => void;
}

const PurchaseContext = createContext<PurchaseContextType | undefined>(undefined);

export const PurchaseProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [purchases, setPurchases] = useState<Purchase[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [totalPurchases, setTotalPurchases] = useState(0);
  const [totalAmount, setTotalAmount] = useState<number | null>(null);
  
  const { isAuthenticated } = useAuth();

  const fetchPurchasesByInvoiceId = async (invoiceId: string, limit = 10, offset = 0) => {
    if (!isAuthenticated || !invoiceId) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await purchaseService.getPurchasesByInvoiceId(invoiceId, limit, offset);
      
      if (offset === 0) {
        setPurchases(response.dados !== undefined ? response.dados : []);
      } else {
        setPurchases(prev => [...prev, ...response.dados]);
      }
      
      setHasMore(response.prox || false);
      
      // Fetch total count if it's the first page
      if (offset === 0) {
        fetchTotalPurchasesCount(invoiceId);
      }
    } catch (err) {
      setError('Erro ao carregar compras');
      console.error(`Error fetching purchases for invoice ${invoiceId}:`, err);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchTotalPurchasesCount = async (invoiceId: string) => {
    if (!isAuthenticated || !invoiceId) return;
    
    try {
      const response = await purchaseService.getPurchasesByInvoiceId(invoiceId, 1, 0, true);
      setTotalPurchases(response.total || 0);
    } catch (err) {
      console.error(`Error fetching total purchases count for invoice ${invoiceId}:`, err);
    }
  };

  const fetchTotalPurchases = async (filters?: PurchaseFilters) => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await purchaseService.getTotalPurchases(
        filters?.dataEspecifica,
        filters?.ultimaParcela,
        filters?.pago,
        filters?.categoria_id
      );
      
      setTotalAmount(response.total);
    } catch (err) {
      setError('Erro ao calcular total de compras');
      console.error('Error fetching total purchases amount:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const clearPurchases = () => {
    setPurchases([]);
    setTotalPurchases(0);
  };

  const clearError = () => setError(null);

  const value = {
    purchases,
    isLoading,
    error,
    hasMore,
    totalPurchases,
    totalAmount,
    fetchPurchasesByInvoiceId,
    fetchTotalPurchases,
    clearPurchases,
    clearError
  };

  return <PurchaseContext.Provider value={value}>{children}</PurchaseContext.Provider>;
};

export const usePurchases = (): PurchaseContextType => {
  const context = useContext(PurchaseContext);
  if (context === undefined) {
    throw new Error('usePurchases must be used within a PurchaseProvider');
  }
  return context;
};