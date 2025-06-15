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
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold text-gray-800">Meus Cartões</h2>
          <button
            onClick={() => setIsCreateModalOpen(true)}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            <Plus className="h-4 w-4 mr-2" />
            Novo Cartão
          </button>
        </div>
        
        <div className="text-center p-6 bg-red-50 rounded-lg">
          <p className="text-red-600">{error}</p>
          <button
            onClick={() => fetchCards()}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
          >
            Tentar novamente
          </button>
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
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold text-gray-800">Meus Cartões</h2>
          <button
            disabled
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-gray-400 cursor-not-allowed"
          >
            <Plus className="h-4 w-4 mr-2" />
            Novo Cartão
          </button>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {[...Array(3)].map((_, index) => (
            <div key={index} className="animate-pulse bg-gradient-to-r from-gray-200 to-gray-300 rounded-xl shadow-lg p-6 h-48 flex flex-col justify-between">
              <div className="flex justify-between items-start">
                <div className="h-6 bg-gray-300 rounded w-3/4"></div>
                <div className="h-8 w-8 bg-gray-300 rounded"></div>
              </div>
              <div className="mt-auto">
                <div className="h-3 bg-gray-300 rounded w-1/4 mb-1"></div>
                <div className="h-4 bg-gray-300 rounded w-1/2"></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-800">Meus Cartões</h2>
        <div className="flex items-center space-x-4">
          <p className="text-sm text-gray-500">
            {totalCards} {totalCards === 1 ? 'cartão' : 'cartões'} encontrados
          </p>
          <button
            onClick={() => setIsCreateModalOpen(true)}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            <Plus className="h-4 w-4 mr-2" />
            Novo Cartão
          </button>
        </div>
      </div>

      {(!cards || cards.length === 0) && !isLoading ? (
        <div className="flex flex-col items-center justify-center p-8 bg-gray-50 rounded-lg border border-gray-200">
          <CreditCard className="h-16 w-16 text-gray-400 mb-4" />
          <h3 className="text-lg font-medium text-gray-900">Nenhum cartão encontrado</h3>
          <p className="mt-1 text-gray-500">Você ainda não possui cartões cadastrados.</p>
          <button
            onClick={() => setIsCreateModalOpen(true)}
            className="mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            <Plus className="h-4 w-4 mr-2" />
            Criar Primeiro Cartão
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
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
        <div className="mt-6 text-center">
          <button
            onClick={loadMoreCards}
            disabled={isLoading}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50"
          >
            {isLoading ? 'Carregando...' : 'Carregar mais cartões'}
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