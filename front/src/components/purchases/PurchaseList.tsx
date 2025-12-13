import React, { useEffect, useState } from 'react';
import { usePurchases } from '../../context/PurchaseContext';
import PurchaseItem from './PurchaseItem';
import { ShoppingBag } from 'lucide-react';

interface PurchaseListProps {
  invoiceId: string;
}

const PurchaseList: React.FC<PurchaseListProps> = ({ invoiceId }) => {
  const { 
    purchases, 
    isLoading, 
    error, 
    hasMore, 
    totalPurchases,
    fetchPurchasesByInvoiceId 
  } = usePurchases();
  const [page, setPage] = useState(1);
  const pageSize = 10;
  
  useEffect(() => {
    if (invoiceId) {
      // Initial load
      fetchPurchasesByInvoiceId(invoiceId, pageSize, 0);
      setPage(1);
    }
  }, [invoiceId]);

  const loadMorePurchases = () => {
    if (hasMore && !isLoading && invoiceId) {
      const nextPage = page + 1;
      fetchPurchasesByInvoiceId(invoiceId, pageSize, (nextPage - 1) * pageSize);
      setPage(nextPage);
    }
  };

  const handlePurchaseUpdated = () => {
    if (invoiceId) {
      // Refresh the list
      fetchPurchasesByInvoiceId(invoiceId, pageSize, 0);
      setPage(1);
    }
  };
  
  // Calculate total amount
  const totalAmount = purchases.reduce((total, purchase) => {
    return total + (purchase.valor_parcela || 0);
  }, 0);

  if (error) {
    return (
      <div className="text-center p-6 bg-red-50 rounded-lg">
        <p className="text-red-600">{error}</p>
        <button
          onClick={() => invoiceId && fetchPurchasesByInvoiceId(invoiceId)}
          className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
        >
          Tentar novamente
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-800">Compras</h2>
        <div className="flex flex-col items-end">
          <p className="text-sm text-gray-500">
            {totalPurchases} {totalPurchases === 1 ? 'compra' : 'compras'} encontradas
          </p>
          {purchases.length > 0 && (
            <p className="text-lg font-semibold text-blue-600">
              Total: R$ {totalAmount.toFixed(2)}
            </p>
          )}
        </div>
      </div>

      {purchases.length === 0 && !isLoading ? (
        <div className="flex flex-col items-center justify-center p-8 bg-gray-50 rounded-lg border border-gray-200">
          <ShoppingBag className="h-16 w-16 text-gray-400 mb-4" />
          <h3 className="text-lg font-medium text-gray-900">Nenhuma compra encontrada</h3>
          <p className="mt-1 text-gray-500">Esta fatura ainda não possui compras registradas.</p>
        </div>
      ) : (
        <div className="space-y-4">
          {purchases.map((purchase) => (
            <PurchaseItem 
              key={purchase.id} 
              purchase={purchase}
              onPurchaseUpdated={handlePurchaseUpdated}
            />
          ))}
        </div>
      )}

      {hasMore && (
        <div className="mt-6 text-center">
          <button
            onClick={loadMorePurchases}
            disabled={isLoading}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50"
          >
            {isLoading ? 'Carregando...' : 'Carregar mais compras'}
          </button>
        </div>
      )}

      {isLoading && purchases.length === 0 && (
        <div className="space-y-4">
          {[...Array(3)].map((_, index) => (
            <div key={index} className="animate-pulse bg-white rounded-lg shadow p-4">
              <div className="h-5 bg-gray-200 rounded w-1/4 mb-3"></div>
              <div className="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
              <div className="flex justify-between">
                <div className="h-4 bg-gray-200 rounded w-1/4"></div>
                <div className="h-4 bg-gray-200 rounded w-1/6"></div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default PurchaseList;