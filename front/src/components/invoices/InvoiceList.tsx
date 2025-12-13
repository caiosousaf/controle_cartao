import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useInvoices } from '../../context/InvoiceContext';
import InvoiceItem from './InvoiceItem';
import CreateInvoiceModal from './CreateInvoiceModal';
import { Receipt, Plus, Filter } from 'lucide-react';

interface InvoiceListProps {
  cardId: string;
}

const InvoiceList: React.FC<InvoiceListProps> = ({ cardId }) => {
  const { 
    invoices, 
    isLoading, 
    error, 
    hasMore, 
    totalInvoices, 
    fetchInvoicesByCardId 
  } = useInvoices();
  const [page, setPage] = useState(1);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [paymentFilter, setPaymentFilter] = useState<'all' | 'paid' | 'pending'>('all');
  const pageSize = 30;
  const navigate = useNavigate();
  
  useEffect(() => {
    if (cardId) {
      // Initial load
      const pagoParam = paymentFilter === 'all' ? undefined : paymentFilter === 'pending';
      fetchInvoicesByCardId(cardId, pageSize, 0, pagoParam);
      setPage(1);
    }
  }, [cardId, paymentFilter]);

  const loadMoreInvoices = () => {
    if (hasMore && !isLoading && cardId) {
      const nextPage = page + 1;
      const pagoParam = paymentFilter === 'all' ? undefined : paymentFilter === 'pending';
      fetchInvoicesByCardId(cardId, pageSize, (nextPage - 1) * pageSize, pagoParam);
      setPage(nextPage);
    }
  };

  const handleInvoiceClick = (invoiceId: string) => {
    navigate(`/invoices/${invoiceId}`);
  };

  const handleInvoiceCreated = () => {
    if (cardId) {
      const pagoParam = paymentFilter === 'all' ? undefined : paymentFilter === 'pending';
      fetchInvoicesByCardId(cardId, pageSize, 0, pagoParam);
      setPage(1);
    }
  };

  const handleFilterChange = (filter: 'all' | 'paid' | 'pending') => {
    setPaymentFilter(filter);
  };

  const getFilterLabel = () => {
    switch (paymentFilter) {
      case 'paid':
        return 'Pagas';
      case 'pending':
        return 'Pendentes';
      default:
        return 'Todas';
    }
  };

  if (error) {
    return (
      <div className="text-center p-6 bg-red-50 rounded-lg">
        <p className="text-red-600">{error}</p>
        <button
          onClick={() => {
            if (cardId) {
              const pagoParam = paymentFilter === 'all' ? undefined : paymentFilter === 'pending';
              fetchInvoicesByCardId(cardId, pageSize, 0, pagoParam);
            }
          }}
          className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
        >
          Tentar novamente
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <h2 className="text-2xl font-bold text-gray-800">Faturas</h2>
        
        <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
          {/* Payment Status Filter */}
          <div className="flex items-center space-x-2">
            <Filter className="h-4 w-4 text-gray-500" />
            <span className="text-sm text-gray-600">Filtrar:</span>
            <div className="flex rounded-md shadow-sm">
              <button
                onClick={() => handleFilterChange('all')}
                className={`px-3 py-1 text-xs font-medium rounded-l-md border ${
                  paymentFilter === 'all'
                    ? 'bg-blue-600 text-white border-blue-600'
                    : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                }`}
              >
                Todas
              </button>
              <button
                onClick={() => handleFilterChange('paid')}
                className={`px-3 py-1 text-xs font-medium border-t border-b ${
                  paymentFilter === 'paid'
                    ? 'bg-green-600 text-white border-green-600'
                    : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                }`}
              >
                Pagas
              </button>
              <button
                onClick={() => handleFilterChange('pending')}
                className={`px-3 py-1 text-xs font-medium rounded-r-md border ${
                  paymentFilter === 'pending'
                    ? 'bg-amber-600 text-white border-amber-600'
                    : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                }`}
              >
                Pendentes
              </button>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <p className="text-sm text-gray-500">
              {totalInvoices || 0} {(totalInvoices || 0) === 1 ? 'fatura' : 'faturas'} {getFilterLabel().toLowerCase()}
            </p>
            <button
              onClick={() => setIsCreateModalOpen(true)}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <Plus className="h-4 w-4 mr-2" />
              Nova Fatura
            </button>
          </div>
        </div>
      </div>

      {(!invoices || invoices.length === 0) && !isLoading ? (
        <div className="flex flex-col items-center justify-center p-8 bg-gray-50 rounded-lg border border-gray-200">
          <Receipt className="h-16 w-16 text-gray-400 mb-4" />
          <h3 className="text-lg font-medium text-gray-900">
            {paymentFilter === 'all' 
              ? 'Nenhuma fatura encontrada'
              : `Nenhuma fatura ${paymentFilter === 'paid' ? 'paga' : 'pendente'} encontrada`
            }
          </h3>
          <p className="mt-1 text-gray-500">
            {paymentFilter === 'all' 
              ? 'Este cartão ainda não possui faturas.'
              : `Este cartão não possui faturas ${paymentFilter === 'paid' ? 'pagas' : 'pendentes'}.`
            }
          </p>
          {paymentFilter === 'all' && (
            <button
              onClick={() => setIsCreateModalOpen(true)}
              className="mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <Plus className="h-4 w-4 mr-2" />
              Criar Primeira Fatura
            </button>
          )}
        </div>
      ) : (
        <div className="space-y-4">
          {invoices?.map((invoice) => (
            <InvoiceItem 
              key={invoice.id} 
              invoice={invoice}
              onClick={() => invoice.id && handleInvoiceClick(invoice.id)}
            />
          ))}
        </div>
      )}

      {hasMore && (
        <div className="mt-6 text-center">
          <button
            onClick={loadMoreInvoices}
            disabled={isLoading}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50"
          >
            {isLoading ? 'Carregando...' : 'Carregar mais faturas'}
          </button>
        </div>
      )}

      {isLoading && (!invoices || invoices.length === 0) && (
        <div className="space-y-4">
          {[...Array(3)].map((_, index) => (
            <div key={index} className="animate-pulse bg-white rounded-lg shadow p-4">
              <div className="h-5 bg-gray-200 rounded w-1/4 mb-3"></div>
              <div className="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
              <div className="h-4 bg-gray-200 rounded w-3/4"></div>
            </div>
          ))}
        </div>
      )}

      <CreateInvoiceModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onInvoiceCreated={handleInvoiceCreated}
        cardId={cardId}
      />
    </div>
  );
};

export default InvoiceList;