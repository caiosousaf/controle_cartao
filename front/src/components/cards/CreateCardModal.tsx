import React, { useState } from 'react';
import { X, CreditCard } from 'lucide-react';
import api from '../../services/api';

interface CreateCardModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCardCreated: () => void;
}

const CreateCardModal: React.FC<CreateCardModalProps> = ({
  isOpen,
  onClose,
  onCardCreated
}) => {
  const [cardName, setCardName] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!cardName.trim()) {
      setError('O nome do cartão é obrigatório');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      await api.post('/cadastros/cartoes', {
        nome: cardName
      });
      
      onCardCreated();
      onClose();
      setCardName('');
    } catch (err) {
      setError('Erro ao criar cartão');
      console.error('Error creating card:', err);
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center space-x-3">
              <CreditCard className="h-6 w-6 text-blue-600" />
              <h2 className="text-2xl font-bold text-gray-900">Novo Cartão</h2>
            </div>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-500"
            >
              <X className="h-6 w-6" />
            </button>
          </div>

          {error && (
            <div className="mb-4 p-4 text-sm text-red-700 bg-red-100 rounded-lg">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="cardName" className="block text-sm font-medium text-gray-700">
                Nome do Cartão
              </label>
              <input
                type="text"
                id="cardName"
                value={cardName}
                onChange={(e) => setCardName(e.target.value)}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                placeholder="Digite o nome do cartão"
                required
              />
            </div>

            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={onClose}
                className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                disabled={isLoading}
              >
                Cancelar
              </button>
              <button
                type="submit"
                className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50"
                disabled={isLoading}
              >
                {isLoading ? 'Criando...' : 'Criar Cartão'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default CreateCardModal;