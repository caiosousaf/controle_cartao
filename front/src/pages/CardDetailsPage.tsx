import React, { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useCards } from '../context/CardContext';
import { useInvoices } from '../context/InvoiceContext';
import InvoiceList from '../components/invoices/InvoiceList';
import { ArrowLeft, CreditCard } from 'lucide-react';

const CardDetailsPage: React.FC = () => {
  const { cardId } = useParams<{ cardId: string }>();
  const { cards, selectedCard, selectCard } = useCards();
  const { clearInvoices } = useInvoices();
  const navigate = useNavigate();
  
  useEffect(() => {
    if (!cardId) {
      navigate('/');
      return;
    }
    
    // If we already have cards loaded, find the matching card
    const card = cards.find(c => c.id === cardId);
    if (card) {
      selectCard(card);
    }
    
    // Cleanup when unmounting
    return () => {
      clearInvoices();
    };
  }, [cardId, cards]);

  const handleBack = () => {
    navigate('/');
  };

  return (
    <div className="space-y-8">
      <div className="relative">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-50 via-indigo-50 to-purple-50 dark:from-gray-800 dark:via-gray-800 dark:to-gray-900 rounded-3xl opacity-50"></div>
        <div className="relative px-8 py-12">
          <div className="flex items-center gap-3 mb-4">
            <button
              onClick={handleBack}
              className="p-2 rounded-xl text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-900 dark:hover:text-white transition-colors"
            >
              <ArrowLeft className="h-6 w-6" />
            </button>
            <div className="p-2 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl shadow-lg">
              <CreditCard className="h-6 w-6 text-white" />
            </div>
            <h1 className="text-4xl md:text-5xl font-bold bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900 dark:from-white dark:via-gray-100 dark:to-white bg-clip-text text-transparent">
              Detalhes do Cartão
            </h1>
          </div>
        </div>
      </div>
      
      {selectedCard && (
        <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-6">
          <div className="flex items-start justify-between">
            <div>
              <div className="flex items-center space-x-3">
                <div className="p-2 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
                  <CreditCard className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                </div>
                <h2 className="text-2xl font-bold text-gray-900 dark:text-white">{selectedCard.nome}</h2>
              </div>
              {selectedCard.data_criacao && (
                <p className="mt-2 text-gray-600 dark:text-gray-400 ml-11">
                  Criado em: {new Date(selectedCard.data_criacao).toLocaleDateString("pt-BR", {
                    day: "2-digit",
                    month: "2-digit",
                    year: "numeric",
                    timeZone: "UTC",
                  })}
                </p>
              )}
            </div>
          </div>
        </div>
      )}
      
      {cardId && (
        <InvoiceList cardId={cardId} />
      )}
    </div>
  );
};

export default CardDetailsPage;