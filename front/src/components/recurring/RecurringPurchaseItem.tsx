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
      <div className="bg-white rounded-lg shadow p-4 hover:shadow-md transition-shadow">
        <div className="flex justify-between items-start">
          <div>
            <h3 className="text-lg font-medium text-gray-900 flex items-center">
              <RepeatIcon className="h-5 w-5 mr-2 text-blue-500" />
              {purchase.nome || "Compra sem nome"}
            </h3>

            {purchase.descricao && (
              <p className="text-sm text-gray-600 mt-1">{purchase.descricao}</p>
            )}

            <div className="mt-3 flex flex-wrap gap-2">
              {purchase.local_compra && (
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                  <MapPin className="h-3 w-3 mr-1" />
                  {purchase.local_compra}
                </span>
              )}

              {purchase.data_criacao && (
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  <Calendar className="h-3 w-3 mr-1" />
                  {new Date(purchase.data_criacao).toLocaleDateString("pt-BR", {
                    day: "2-digit",
                    month: "2-digit",
                    year: "numeric",
                    timeZone: "UTC",
                  })}
                </span>
              )}

              <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                purchase.ativo 
                  ? 'bg-green-100 text-green-800' 
                  : 'bg-red-100 text-red-800'
              }`}>
                {purchase.ativo ? (
                  <CheckCircle className="h-3 w-3 mr-1" />
                ) : (
                  <XCircle className="h-3 w-3 mr-1" />
                )}
                {purchase.ativo ? 'Ativo' : 'Inativo'}
              </span>
            </div>
          </div>

          <div className="flex flex-col items-end">
            <span className="text-lg font-bold text-gray-900">
              R$ {purchase.valor_parcela?.toFixed(2) || "0.00"}
            </span>
            <span className="text-sm text-gray-600">por mês</span>
            
            <div className="mt-2 flex space-x-2">
              <button
                onClick={() => setIsEditing(true)}
                className="p-1 text-gray-500 hover:text-blue-600 transition-colors"
                title="Editar"
              >
                <Edit className="h-4 w-4" />
              </button>
              
              {isConfirmingAction === 'toggle' ? (
                <div className="flex items-center space-x-2">
                  <button
                    onClick={handleToggleActive}
                    disabled={isLoading}
                    className="text-xs bg-blue-600 text-white px-2 py-1 rounded hover:bg-blue-700"
                  >
                    Confirmar
                  </button>
                  <button
                    onClick={() => setIsConfirmingAction(null)}
                    disabled={isLoading}
                    className="text-xs bg-gray-300 text-gray-700 px-2 py-1 rounded hover:bg-gray-400"
                  >
                    Cancelar
                  </button>
                </div>
              ) : (
                <button
                  onClick={() => setIsConfirmingAction('toggle')}
                  className={`p-1 ${purchase.ativo ? 'text-red-500 hover:text-red-600' : 'text-green-500 hover:text-green-600'} transition-colors`}
                  title={purchase.ativo ? 'Desativar' : 'Reativar'}
                >
                  <Power className="h-4 w-4" />
                </button>
              )}

              {isConfirmingAction === 'remove' ? (
                <div className="flex items-center space-x-2">
                  <button
                    onClick={handleRemove}
                    disabled={isLoading}
                    className="text-xs bg-red-600 text-white px-2 py-1 rounded hover:bg-red-700"
                  >
                    Confirmar
                  </button>
                  <button
                    onClick={() => setIsConfirmingAction(null)}
                    disabled={isLoading}
                    className="text-xs bg-gray-300 text-gray-700 px-2 py-1 rounded hover:bg-gray-400"
                  >
                    Cancelar
                  </button>
                </div>
              ) : (
                <button
                  onClick={() => setIsConfirmingAction('remove')}
                  className="p-1 text-gray-500 hover:text-red-600 transition-colors"
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