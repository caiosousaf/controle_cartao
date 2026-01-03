import React, { useState, useEffect } from 'react';
import { X, Clock, AlertTriangle, CheckCircle2 } from 'lucide-react';
import { purchaseService } from '../../services/purchaseService';
import { Purchase } from '../../types/purchase';

interface AnticipatePurchaseModalProps {
  isOpen: boolean;
  onClose: () => void;
  purchase: Purchase;
  invoiceId: string;
  onPurchaseUpdated: () => void;
}

const AnticipatePurchaseModal: React.FC<AnticipatePurchaseModalProps> = ({
  isOpen,
  onClose,
  purchase,
  invoiceId,
  onPurchaseUpdated
}) => {
  const [availableInstallments, setAvailableInstallments] = useState<number[]>([]);
  const [selectedInstallments, setSelectedInstallments] = useState<number[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingInstallments, setIsLoadingInstallments] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen) {
      if (purchase.agrupamento_id) {
        fetchAvailableInstallments();
        setSelectedInstallments([]);
        setError(null);
      } else {
        setError('Identificador da compra não disponível');
      }
    }
  }, [isOpen, purchase.agrupamento_id]);

  const fetchAvailableInstallments = async () => {
    if (!purchase.agrupamento_id || !invoiceId) return;
    
    setIsLoadingInstallments(true);
    setError(null);
    
    try {
      const installments = await purchaseService.getAvailableInstallments(invoiceId, purchase.agrupamento_id);
      setAvailableInstallments(installments);
    } catch (err) {
      setError('Erro ao carregar parcelas disponíveis');
      console.error('Error fetching available installments:', err);
    } finally {
      setIsLoadingInstallments(false);
    }
  };

  const handleInstallmentToggle = (installment: number) => {
    setSelectedInstallments(prev => {
      if (prev.includes(installment)) {
        return prev.filter(i => i !== installment);
      } else {
        return [...prev, installment].sort((a, b) => a - b);
      }
    });
  };

  const handleConfirm = async () => {
    if (!purchase.agrupamento_id || selectedInstallments.length === 0) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      await purchaseService.anticipateInstallments(
        invoiceId,
        selectedInstallments,
        purchase.agrupamento_id
      );
      
      onPurchaseUpdated();
      onClose();
    } catch (err) {
      setError('Erro ao antecipar parcelas');
      console.error('Error anticipating installments:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    onClose();
    setSelectedInstallments([]);
    setError(null);
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4 overflow-y-auto">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md my-8">
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center space-x-3">
              <Clock className="h-6 w-6 text-blue-600" />
              <h2 className="text-2xl font-bold text-gray-900">Antecipar Parcelas</h2>
            </div>
            <button
              onClick={handleClose}
              className="text-gray-400 hover:text-gray-500"
              disabled={isLoading}
            >
              <X className="h-6 w-6" />
            </button>
          </div>

          <div className="mb-4">
            <p className="text-sm text-gray-600 mb-2">
              <strong>Compra:</strong> {purchase.nome || 'Compra sem nome'}
            </p>
            <p className="text-sm text-gray-600">
              <strong>Parcela atual:</strong> {purchase.parcela_atual}/{purchase.quantidade_parcelas}
            </p>
          </div>

          {error && (
            <div className="mb-4 p-4 text-sm text-red-700 bg-red-100 rounded-lg flex items-start">
              <AlertTriangle className="h-5 w-5 mr-2 mt-0.5 flex-shrink-0" />
              <span>{error}</span>
            </div>
          )}

          {isLoadingInstallments ? (
            <div className="text-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
              <p className="mt-4 text-sm text-gray-600">Carregando parcelas disponíveis...</p>
            </div>
          ) : availableInstallments.length === 0 ? (
            <div className="text-center py-8">
              <AlertTriangle className="h-12 w-12 text-yellow-500 mx-auto mb-4" />
              <p className="text-gray-600">Nenhuma parcela disponível para antecipação</p>
            </div>
          ) : (
            <>
              <div className="mb-4">
                <h3 className="text-sm font-medium text-gray-700 mb-3">
                  Selecione as parcelas que deseja antecipar:
                </h3>
                <div className="max-h-64 overflow-y-auto border border-gray-200 rounded-lg p-3 space-y-2">
                  {availableInstallments.map((installment) => (
                    <label
                      key={installment}
                      className="flex items-center p-2 hover:bg-gray-50 rounded cursor-pointer"
                    >
                      <input
                        type="checkbox"
                        checked={selectedInstallments.includes(installment)}
                        onChange={() => handleInstallmentToggle(installment)}
                        className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                      />
                      <span className="ml-3 text-sm text-gray-700">
                        Parcela {installment}
                      </span>
                    </label>
                  ))}
                </div>
              </div>

              {selectedInstallments.length > 0 && (
                <div className="mb-4 p-3 bg-blue-50 rounded-lg">
                  <div className="flex items-center mb-2">
                    <CheckCircle2 className="h-5 w-5 text-blue-600 mr-2" />
                    <span className="text-sm font-medium text-blue-900">
                      {selectedInstallments.length} parcela(s) selecionada(s):
                    </span>
                  </div>
                  <div className="flex flex-wrap gap-2 mt-2">
                    {selectedInstallments.map((installment) => (
                      <span
                        key={installment}
                        className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                      >
                        Parcela {installment}
                      </span>
                    ))}
                  </div>
                </div>
              )}

              <div className="flex justify-end space-x-3 mt-6">
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
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                  disabled={isLoading || selectedInstallments.length === 0}
                >
                  {isLoading ? 'Antecipando...' : 'Confirmar Antecipação'}
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default AnticipatePurchaseModal;
