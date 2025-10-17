import React, { useState } from 'react';
import { X, Receipt } from 'lucide-react';
import api from '../../services/api';
import { format } from 'date-fns';

interface CreateInvoiceModalProps {
  isOpen: boolean;
  onClose: () => void;
  onInvoiceCreated: () => void;
  cardId: string;
}

const CreateInvoiceModal: React.FC<CreateInvoiceModalProps> = ({
  isOpen,
  onClose,
  onInvoiceCreated,
  cardId
}) => {
  const [dueDate, setDueDate] = useState(format(new Date(), 'yyyy-MM-dd'));
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!dueDate) {
      setError('A data de vencimento é obrigatória');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      await api.post(`/cadastros/cartao/${cardId}/faturas`, {
        data_vencimento: dueDate
      });
      
      onInvoiceCreated();
      onClose();
      setDueDate(format(new Date(), 'yyyy-MM-dd'));
    } catch (err) {
      setError('Erro ao criar fatura');
      console.error('Error creating invoice:', err);
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
              <Receipt className="h-6 w-6 text-blue-600" />
              <h2 className="text-2xl font-bold text-gray-900">Nova Fatura</h2>
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
              <label htmlFor="dueDate" className="block text-sm font-medium text-gray-700">
                Data de Vencimento
              </label>
              <input
                type="date"
                id="dueDate"
                value={dueDate}
                onChange={(e) => setDueDate(e.target.value)}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
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
                {isLoading ? 'Criando...' : 'Criar Fatura'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default CreateInvoiceModal;