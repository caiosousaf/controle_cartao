import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCards } from '../../context/CardContext';
import CardItem from './CardItem';
import CreateCardModal from './CreateCardModal';
import { CreditCard, Plus } from 'lucide-react';

const CardList: React.FC = () => {
  const { cards, isLoading, error, hasMore, totalCards, fetchCards } = useCards();
  const [page, setPage] = useState(1);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const pageSize = 10;
  const navigate = useNavigate();
  
  useEffect(() => {
    // Initial load
    fetchCards(pageSize, 0);
  }, []);

  const loadMoreCards = () => {
    if (hasMore && !isLoading) {
      const nextPage = page + 1;
      fetchCards(pageSize, (nextPage - 1) * pageSize);
      setPage(nextPage);
    }
  };

  const handleCardClick = (cardId: string) => {
    navigate(`/cards/${cardId}`);
  };

  const handleCardCreated = () => {
    fetchCards(pageSize, 0);
    setPage(1);
  };

  // Show error state
  if (error) {
    return (
      <div className="space-y-8">
        <div className="bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 rounded-2xl shadow-xl p-8 text-white relative overflow-hidden">
          <div className="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full -mr-32 -mt-32"></div>
          <div className="relative z-10 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div>
              <h2 className="text-3xl font-bold mb-2">Meus Cartões</h2>
              <p className="text-blue-100 text-sm">Gerencie seus cartões de crédito</p>
            </div>
            <button
              onClick={() => setIsCreateModalOpen(true)}
              className="inline-flex items-center px-6 py-3 bg-white text-indigo-600 font-semibold rounded-xl shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-200"
            >
              <Plus className="h-5 w-5 mr-2" />
              Novo Cartão
            </button>
          </div>
        </div>
        
        <div className="text-center p-8 bg-gradient-to-br from-red-50 to-rose-50 dark:from-red-900/20 dark:to-rose-900/20 rounded-2xl border-2 border-red-200 dark:border-red-800">
          <div className="max-w-md mx-auto">
            <div className="p-3 bg-red-100 dark:bg-red-900/30 rounded-full w-16 h-16 mx-auto mb-4 flex items-center justify-center">
              <CreditCard className="h-8 w-8 text-red-600 dark:text-red-400" />
            </div>
            <p className="text-red-800 dark:text-red-300 font-medium mb-4">{error}</p>
            <button
              onClick={() => fetchCards()}
              className="inline-flex items-center px-6 py-3 bg-red-600 text-white font-semibold rounded-xl hover:bg-red-700 transition-all duration-200 shadow-md hover:shadow-lg transform hover:scale-105"
            >
              Tentar novamente
            </button>
          </div>
        </div>

        <CreateCardModal
          isOpen={isCreateModalOpen}
          onClose={() => setIsCreateModalOpen(false)}
          onCardCreated={handleCardCreated}
        />
      </div>
    );
  }

  // Show loading state for initial load
  if (isLoading && (!cards || cards.length === 0)) {
    return (
      <div className="space-y-8">
        <div className="bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 rounded-2xl shadow-xl p-8 text-white relative overflow-hidden">
          <div className="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full -mr-32 -mt-32"></div>
          <div className="relative z-10 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div>
              <div className="h-8 bg-white/20 rounded-lg w-48 mb-3 animate-pulse"></div>
              <div className="h-4 bg-white/10 rounded w-32 animate-pulse"></div>
            </div>
            <div className="h-12 bg-white/20 rounded-xl w-36 animate-pulse"></div>
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {[...Array(4)].map((_, index) => (
            <div key={index} className="animate-pulse bg-gradient-to-br from-gray-200 to-gray-300 dark:from-gray-700 dark:to-gray-800 rounded-2xl shadow-xl p-6 min-h-[200px] flex flex-col justify-between">
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <div className="h-6 bg-gray-300 dark:bg-gray-600 rounded-lg w-3/4 mb-3"></div>
                  <div className="h-8 bg-gray-300 dark:bg-gray-600 rounded-lg w-1/2"></div>
                </div>
                <div className="h-12 w-12 bg-gray-300 dark:bg-gray-600 rounded-xl"></div>
              </div>
              <div className="mt-auto pt-4 border-t border-gray-300 dark:border-gray-600">
                <div className="h-3 bg-gray-300 dark:bg-gray-600 rounded w-1/4 mb-2"></div>
                <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-1/2"></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Header Section */}
      <div className="bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 rounded-2xl shadow-xl p-8 text-white relative overflow-hidden">
        {/* Decorative background elements */}
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full -mr-32 -mt-32"></div>
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full -ml-24 -mb-24"></div>
        
        <div className="relative z-10">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div>
              <h2 className="text-3xl font-bold mb-2">Meus Cartões</h2>
              <p className="text-blue-100 text-sm">
                {totalCards} {totalCards === 1 ? 'cartão cadastrado' : 'cartões cadastrados'}
              </p>
            </div>
            <button
              onClick={() => setIsCreateModalOpen(true)}
              className="inline-flex items-center px-6 py-3 bg-white text-indigo-600 font-semibold rounded-xl shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-200 hover:bg-blue-50 focus:outline-none focus:ring-4 focus:ring-white/50"
            >
              <Plus className="h-5 w-5 mr-2" />
              Novo Cartão
            </button>
          </div>
        </div>
      </div>

      {(!cards || cards.length === 0) && !isLoading ? (
        <div className="flex flex-col items-center justify-center p-12 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-800 dark:to-gray-900 rounded-2xl border-2 border-dashed border-gray-300 dark:border-gray-700">
          <div className="p-4 bg-blue-100 dark:bg-blue-900/30 rounded-full mb-4">
            <CreditCard className="h-12 w-12 text-blue-600 dark:text-blue-400" />
          </div>
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">Nenhum cartão encontrado</h3>
          <p className="text-gray-600 dark:text-gray-400 mb-6 text-center max-w-sm">
            Comece gerenciando suas finanças criando seu primeiro cartão de crédito.
          </p>
          <button
            onClick={() => setIsCreateModalOpen(true)}
            className="inline-flex items-center px-6 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-200 focus:outline-none focus:ring-4 focus:ring-blue-500/50"
          >
            <Plus className="h-5 w-5 mr-2" />
            Criar Primeiro Cartão
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {cards?.map((card) => (
            <CardItem 
              key={card.id} 
              card={card} 
              onClick={() => card.id && handleCardClick(card.id)} 
            />
          ))}
        </div>
      )}

      {hasMore && (
        <div className="mt-8 text-center">
          <button
            onClick={loadMoreCards}
            disabled={isLoading}
            className="inline-flex items-center px-6 py-3 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 font-medium rounded-xl shadow-md hover:shadow-lg border border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600 transform hover:scale-105 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
          >
            {isLoading ? (
              <>
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Carregando...
              </>
            ) : (
              'Carregar mais cartões'
            )}
          </button>
        </div>
      )}

      <CreateCardModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onCardCreated={handleCardCreated}
      />
    </div>
  );
};

export default CardList;