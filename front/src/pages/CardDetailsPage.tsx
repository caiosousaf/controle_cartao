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
      <div className="flex items-center space-x-4">
        <button
          onClick={handleBack}
          className="p-2 rounded-full text-gray-600 hover:bg-gray-100 hover:text-gray-900 transition-colors"
        >
          <ArrowLeft className="h-6 w-6" />
        </button>
        <h1 className="text-3xl font-bold text-gray-900">Detalhes do Cartão</h1>
      </div>
      
      {selectedCard && (
        <div className="bg-white rounded-lg shadow-lg p-6">
          <div className="flex items-start justify-between">
            <div>
              <div className="flex items-center space-x-3">
                <CreditCard className="h-8 w-8 text-blue-600" />
                <h2 className="text-2xl font-bold text-gray-900">{selectedCard.nome}</h2>
              </div>
              {selectedCard.data_criacao && (
                <p className="mt-2 text-gray-600">
                  Criado em: {new Date(selectedCard.data_criacao).toLocaleDateString()}
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