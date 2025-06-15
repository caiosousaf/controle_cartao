import React, { useEffect, useState } from 'react';
import { RepeatIcon, Plus, FileDown } from 'lucide-react';
import { recurringService } from '../services/recurringService';
import { RecurringPurchase, ExpenseEstimate } from '../types/recurring';
import RecurringPurchaseItem from '../components/recurring/RecurringPurchaseItem';
import CreateRecurringModal from '../components/recurring/CreateRecurringModal';
import RegisterRecurringModal from '../components/recurring/RegisterRecurringModal';
import EstimateCard from '../components/recurring/EstimateCard';

const RecurringPurchasesPage: React.FC = () => {
  const [purchases, setPurchases] = useState<RecurringPurchase[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [totalPurchases, setTotalPurchases] = useState(0);
  const [page, setPage] = useState(1);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isRegisterModalOpen, setIsRegisterModalOpen] = useState(false);
  const [estimates, setEstimates] = useState<ExpenseEstimate[]>([]);
  const pageSize = 10;

  useEffect(() => {
    fetchPurchases();
    fetchEstimates();
  }, []);

  const fetchPurchases = async (offset = 0) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await recurringService.getRecurringPurchases(pageSize, offset);
      
      if (offset === 0) {
        setPurchases(response.dados);
      } else {
        setPurchases(prev => [...prev, ...response.dados]);
      }
      
      setHasMore(response.prox || false);
      setTotalPurchases(response.total || 0);
    } catch (err) {
      setError('Erro ao carregar compras recorrentes');
    } finally {
      setIsLoading(false);
    }
  };

  const fetchEstimates = async () => {
    try {
      const response = await recurringService.getEstimate();
      setEstimates(response.dados);
    } catch (err) {
      console.error('Error fetching estimates:', err);
    }
  };

  const loadMorePurchases = () => {
    if (hasMore && !isLoading) {
      const nextPage = page + 1;
      fetchPurchases((nextPage - 1) * pageSize);
      setPage(nextPage);
    }
  };

  const handlePurchaseCreated = () => {
    fetchPurchases(0);
    setPage(1);
    fetchEstimates();
  };

  return (
    <div className="space-y-8">
      <div className="flex items-center space-x-4">
        <RepeatIcon className="h-8 w-8 text-blue-600" />
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Compras Recorrentes</h1>
          <p className="mt-1 text-gray-600">Gerencie suas despesas mensais fixas</p>
        </div>
      </div>

      {estimates.length > 0 && <EstimateCard estimates={estimates} />}
      
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold text-gray-800">Lista de Compras Recorrentes</h2>
          <div className="flex items-center space-x-4">
            <p className="text-sm text-gray-500">
              {totalPurchases} {totalPurchases === 1 ? 'compra' : 'compras'} encontradas
            </p>
            <button
              onClick={() => setIsRegisterModalOpen(true)}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
            >
              <FileDown className="h-4 w-4 mr-2" />
              Registrar na Fatura
            </button>
            <button
              onClick={() => setIsCreateModalOpen(true)}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <Plus className="h-4 w-4 mr-2" />
              Nova Compra Recorrente
            </button>
          </div>
        </div>

        {error && (
          <div className="text-center p-6 bg-red-50 rounded-lg">
            <p className="text-red-600">{error}</p>
            <button
              onClick={() => fetchPurchases()}
              className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
            >
              Tentar novamente
            </button>
          </div>
        )}

        {!isLoading && purchases.length === 0 ? (
          <div className="flex flex-col items-center justify-center p-8 bg-gray-50 rounded-lg border border-gray-200">
            <RepeatIcon className="h-16 w-16 text-gray-400 mb-4" />
            <h3 className="text-lg font-medium text-gray-900">Nenhuma compra recorrente encontrada</h3>
            <p className="mt-1 text-gray-500">Comece cadastrando sua primeira compra recorrente.</p>
            <button
              onClick={() => setIsCreateModalOpen(true)}
              className="mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <Plus className="h-4 w-4 mr-2" />
              Criar Primeira Compra Recorrente
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            {purchases.map((purchase) => (
              <RecurringPurchaseItem 
                key={purchase.id} 
                purchase={purchase}
                onUpdate={() => {
                  fetchPurchases(0);
                  fetchEstimates();
                }}
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
              {isLoading ? 'Carregando...' : 'Carregar mais'}
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

      <CreateRecurringModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onPurchaseCreated={handlePurchaseCreated}
      />

      <RegisterRecurringModal
        isOpen={isRegisterModalOpen}
        onClose={() => setIsRegisterModalOpen(false)}
        onRegistered={handlePurchaseCreated}
      />
    </div>
  );
};

export default RecurringPurchasesPage;