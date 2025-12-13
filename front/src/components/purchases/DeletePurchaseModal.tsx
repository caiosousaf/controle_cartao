import React, { useState } from 'react';
import { X, Trash2, AlertTriangle } from 'lucide-react';
import { purchaseService } from '../../services/purchaseService';
import { Purchase } from '../../types/purchase';

interface DeletePurchaseModalProps {
  isOpen: boolean;
  onClose: () => void;
  purchase: Purchase;
  onPurchaseDeleted: () => void;
}

const DeletePurchaseModal: React.FC<DeletePurchaseModalProps> = ({
  isOpen,
  onClose,
  purchase,
  onPurchaseDeleted
}) => {
  const [removeAllInstallments, setRemoveAllInstallments] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleConfirm = async () => {
    if (!purchase.id) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      await purchaseService.deletePurchase(purchase.id, removeAllInstallments);
      onPurchaseDeleted();
      onClose();
    } catch (err) {
      setError('Erro ao remover compra');
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    onClose();
    setRemoveAllInstallments(false);
    setError(null);
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center space-x-3">
              <Trash2 className="h-6 w-6 text-red-600" />
              <h2 className="text-2xl font-bold text-gray-900">Remover Compra</h2>
            </div>
            <button
              onClick={handleClose}
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
            <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-4">
              <div className="flex items-start space-x-3">
                <AlertTriangle className="h-5 w-5 text-red-600 mt-0.5" />
                <div>
                  <h4 className="text-sm font-medium text-red-800 mb-1">
                    Atenção: Esta ação não pode ser desfeita
                  </h4>
                  <p className="text-sm text-red-700">
                    Você está prestes a remover a compra <strong>{purchase.nome}</strong>.
                  </p>
                </div>
              </div>
            </div>

            <div className="space-y-3">
              <h4 className="text-sm font-medium text-gray-700">Escolha o escopo da remoção:</h4>
              
              <label className="flex items-start space-x-3 p-3 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer">
                <input
                  type="radio"
                  name="removeScope"
                  checked={!removeAllInstallments}
                  onChange={() => setRemoveAllInstallments(false)}
                  className="mt-1"
                />
                <div>
                  <span className="text-sm font-medium text-gray-900">
                    Remover apenas esta parcela
                  </span>
                  <p className="text-xs text-gray-600 mt-1">
                    Parcela {purchase.parcela_atual} de {purchase.quantidade_parcelas}
                  </p>
                </div>
              </label>

              <label className="flex items-start space-x-3 p-3 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer">
                <input
                  type="radio"
                  name="removeScope"
                  checked={removeAllInstallments}
                  onChange={() => setRemoveAllInstallments(true)}
                  className="mt-1"
                />
                <div>
                  <span className="text-sm font-medium text-gray-900">
                    Remover todas as parcelas deste agrupamento
                  </span>
                  <p className="text-xs text-gray-600 mt-1">
                    Todas as {purchase.quantidade_parcelas} parcelas serão removidas
                  </p>
                </div>
              </label>
            </div>
          </div>

          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={handleClose}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
              disabled={isLoading}
            >
              Cancelar
            </button>
            <button
              type="button"
              onClick={handleConfirm}
              className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 disabled:opacity-50"
              disabled={isLoading}
            >
              {isLoading ? 'Removendo...' : 'Confirmar Remoção'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DeletePurchaseModal;