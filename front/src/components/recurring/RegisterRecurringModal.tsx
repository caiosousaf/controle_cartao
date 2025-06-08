import React, { useState } from 'react';
import { X, AlertCircle } from 'lucide-react';
import { recurringService } from '../../services/recurringService';

interface RegisterRecurringModalProps {
  isOpen: boolean;
  onClose: () => void;
  onRegistered: () => void;
}

const RegisterRecurringModal: React.FC<RegisterRecurringModalProps> = ({
  isOpen,
  onClose,
  onRegistered
}) => {
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleConfirm = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      await recurringService.registerRecurringPurchases();
      onRegistered();
      onClose();
    } catch (err) {
      setError('Erro ao registrar compras recorrentes');
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
              <AlertCircle className="h-6 w-6 text-blue-600" />
              <h2 className="text-2xl font-bold text-gray-900">Registrar Compras Recorrentes</h2>
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

          <div className="mb-6">
            <p className="text-gray-700">
              Você tem certeza que deseja registrar todas as compras recorrentes ativas na fatura?
            </p>
            <p className="mt-2 text-sm text-gray-500">
              Esta ação irá criar novas compras na fatura atual para todas as compras recorrentes ativas.
            </p>
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
              type="button"
              onClick={handleConfirm}
              className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50"
              disabled={isLoading}
            >
              {isLoading ? 'Processando...' : 'Confirmar'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default RegisterRecurringModal;