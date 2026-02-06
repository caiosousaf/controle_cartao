import React, { useState } from 'react';
import { RecurringPurchase } from '../../types/recurring';
import { RepeatIcon, MapPin, Tag, Calendar, CheckCircle, XCircle, Edit, Power, Trash2 } from 'lucide-react';
import { recurringService } from '../../services/recurringService';
import CreateRecurringModal from './CreateRecurringModal';

interface RecurringPurchaseItemProps {
  purchase: RecurringPurchase;
  onUpdate: () => void;
}

const RecurringPurchaseItem: React.FC<RecurringPurchaseItemProps> = ({ purchase, onUpdate }) => {
  const [isConfirmingAction, setIsConfirmingAction] = useState<'toggle' | 'remove' | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const handleToggleActive = async () => {
    if (!purchase.id) return;
    
    setIsLoading(true);
    try {
      if (purchase.ativo) {
        await recurringService.deactivateRecurringPurchase(purchase.id);
      } else {
        await recurringService.reactivateRecurringPurchase(purchase.id);
      }
      onUpdate();
    } catch (error) {
      console.error('Error toggling recurring purchase status:', error);
    } finally {
      setIsLoading(false);
      setIsConfirmingAction(null);
    }
  };

  const handleRemove = async () => {
    if (!purchase.id) return;
    
    setIsLoading(true);
    try {
      await recurringService.removeRecurringPurchase(purchase.id);
      onUpdate();
    } catch (error) {
      console.error('Error removing recurring purchase:', error);
    } finally {
      setIsLoading(false);
      setIsConfirmingAction(null);
    }
  };

  return (
    <>
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-lg border border-gray-200 dark:border-gray-700 transition-all duration-300 p-5">
        <div className="flex justify-between items-start">
          <div className="flex-1">
            <div className="flex items-center gap-3 mb-2">
              <div className="p-2 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
                <RepeatIcon className="h-5 w-5 text-blue-600 dark:text-blue-400" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                {purchase.nome || "Compra sem nome"}
              </h3>
            </div>

            {purchase.descricao && (
              <p className="text-sm text-gray-600 dark:text-gray-400 ml-11 mb-3">{purchase.descricao}</p>
            )}

            <div className="ml-11 flex flex-wrap gap-2">
              {purchase.local_compra && (
                <span className="inline-flex items-center px-3 py-1 rounded-lg text-xs font-medium bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300 shadow-sm">
                  <MapPin className="h-3 w-3 mr-1.5" />
                  {purchase.local_compra}
                </span>
              )}

              {purchase.data_criacao && (
                <span className="inline-flex items-center px-3 py-1 rounded-lg text-xs font-medium bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 shadow-sm">
                  <Calendar className="h-3 w-3 mr-1.5" />
                  {new Date(purchase.data_criacao).toLocaleDateString("pt-BR", {
                    day: "2-digit",
                    month: "2-digit",
                    year: "numeric",
                    timeZone: "UTC",
                  })}
                </span>
              )}

              <span className={`inline-flex items-center px-3 py-1 rounded-lg text-xs font-medium shadow-sm ${
                purchase.ativo 
                  ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300' 
                  : 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300'
              }`}>
                {purchase.ativo ? (
                  <CheckCircle className="h-3 w-3 mr-1.5" />
                ) : (
                  <XCircle className="h-3 w-3 mr-1.5" />
                )}
                {purchase.ativo ? 'Ativo' : 'Inativo'}
              </span>
            </div>
          </div>

          <div className="flex flex-col items-end ml-4">
            <span className="text-2xl font-bold text-gray-900 dark:text-white">
              R$ {purchase.valor_parcela?.toFixed(2) || "0.00"}
            </span>
            <span className="text-sm text-gray-600 dark:text-gray-400">por mês</span>
            
            <div className="mt-3 flex items-center space-x-1 p-1 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
              <button
                onClick={() => setIsEditing(true)}
                className="p-2 text-gray-500 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-all duration-200"
                title="Editar"
              >
                <Edit className="h-4 w-4" />
              </button>
              
              {isConfirmingAction === 'toggle' ? (
                <div className="flex items-center space-x-1">
                  <button
                    onClick={handleToggleActive}
                    disabled={isLoading}
                    className="text-xs bg-blue-600 dark:bg-blue-500 text-white px-2 py-1 rounded-md hover:bg-blue-700 dark:hover:bg-blue-600 transition-colors"
                  >
                    Confirmar
                  </button>
                  <button
                    onClick={() => setIsConfirmingAction(null)}
                    disabled={isLoading}
                    className="text-xs bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-300 px-2 py-1 rounded-md hover:bg-gray-400 dark:hover:bg-gray-500 transition-colors"
                  >
                    Cancelar
                  </button>
                </div>
              ) : (
                <button
                  onClick={() => setIsConfirmingAction('toggle')}
                  className={`p-2 rounded-lg transition-all duration-200 ${
                    purchase.ativo 
                      ? 'text-red-500 dark:text-red-400 hover:text-red-600 dark:hover:text-red-300 hover:bg-red-50 dark:hover:bg-red-900/20' 
                      : 'text-green-500 dark:text-green-400 hover:text-green-600 dark:hover:text-green-300 hover:bg-green-50 dark:hover:bg-green-900/20'
                  }`}
                  title={purchase.ativo ? 'Desativar' : 'Reativar'}
                >
                  <Power className="h-4 w-4" />
                </button>
              )}

              {isConfirmingAction === 'remove' ? (
                <div className="flex items-center space-x-1">
                  <button
                    onClick={handleRemove}
                    disabled={isLoading}
                    className="text-xs bg-red-600 dark:bg-red-500 text-white px-2 py-1 rounded-md hover:bg-red-700 dark:hover:bg-red-600 transition-colors"
                  >
                    Confirmar
                  </button>
                  <button
                    onClick={() => setIsConfirmingAction(null)}
                    disabled={isLoading}
                    className="text-xs bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-300 px-2 py-1 rounded-md hover:bg-gray-400 dark:hover:bg-gray-500 transition-colors"
                  >
                    Cancelar
                  </button>
                </div>
              ) : (
                <button
                  onClick={() => setIsConfirmingAction('remove')}
                  className="p-2 text-gray-500 dark:text-gray-400 hover:text-red-600 dark:hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-all duration-200"
                  title="Remover"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              )}
            </div>
          </div>
        </div>
      </div>

      <CreateRecurringModal
        isOpen={isEditing}
        onClose={() => setIsEditing(false)}
        onPurchaseCreated={onUpdate}
        editPurchase={purchase}
      />
    </>
  );
};

export default RecurringPurchaseItem;