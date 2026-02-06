import React, { useState } from "react";
import { Purchase } from "../../types/purchase";
import { ShoppingCart, MapPin, Tag, Calendar, Edit, Trash2, Clock } from "lucide-react";
import EditPurchaseModal from "./EditPurchaseModal";
import DeletePurchaseModal from "./DeletePurchaseModal";
import AnticipatePurchaseModal from "./AnticipatePurchaseModal";

interface PurchaseItemProps {
  purchase: Purchase;
  invoiceId: string;
  onPurchaseUpdated: () => void;
}

const PurchaseItem: React.FC<PurchaseItemProps> = ({ purchase, invoiceId, onPurchaseUpdated }) => {
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isAnticipateModalOpen, setIsAnticipateModalOpen] = useState(false);

  // Verifica se é possível antecipar parcelas
  const canAnticipate = 
    purchase.quantidade_parcelas !== undefined &&
    purchase.quantidade_parcelas > 1 &&
    purchase.parcela_atual !== undefined &&
    purchase.parcela_atual < purchase.quantidade_parcelas;

  return (
    <>
      <div className="group bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-lg border border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-600 transition-all duration-300 overflow-hidden relative">
        {/* Gradient accent bar */}
        <div className="absolute left-0 top-0 bottom-0 w-1 bg-gradient-to-b from-blue-400 to-indigo-600"></div>
        
        <div className="pl-5 pr-6 py-5">
          <div className="flex justify-between items-start gap-4">
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-3 mb-2">
                <div className="p-2 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
                  <ShoppingCart className="h-5 w-5 text-blue-600 dark:text-blue-400" />
                </div>
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white truncate">
                  {purchase.nome || "Compra sem nome"}
                </h3>
              </div>

              {purchase.descricao && (
                <p className="text-sm text-gray-600 dark:text-gray-400 ml-11 mb-3 line-clamp-2">
                  {purchase.descricao}
                </p>
              )}

              <div className="ml-11 flex flex-wrap gap-2">
                {purchase.local_compra && (
                  <span className="inline-flex items-center px-3 py-1 rounded-lg text-xs font-medium bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 shadow-sm">
                    <MapPin className="h-3 w-3 mr-1.5" />
                    {purchase.local_compra}
                  </span>
                )}

                {purchase.categoria_nome && (
                  <span className="inline-flex items-center px-3 py-1 rounded-lg text-xs font-medium bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300 shadow-sm">
                    <Tag className="h-3 w-3 mr-1.5" />
                    {purchase.categoria_nome}
                  </span>
                )}

                {purchase.data_compra && (
                  <span className="inline-flex items-center px-3 py-1 rounded-lg text-xs font-medium bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 shadow-sm">
                    <Calendar className="h-3 w-3 mr-1.5" />
                    {new Date(purchase.data_compra).toLocaleDateString("pt-BR", {
                      day: "2-digit",
                      month: "2-digit",
                      year: "numeric",
                      timeZone: "UTC",
                    })}
                  </span>
                )}
              </div>
            </div>

            <div className="flex flex-col items-end ml-4 flex-shrink-0">
              <div className="text-right mb-2">
                <span className="text-2xl font-bold text-gray-900 dark:text-white">
                  R$ {purchase.valor_parcela?.toFixed(2) || "0.00"}
                </span>
                {purchase.parcela_atual !== undefined &&
                  purchase.quantidade_parcelas !== undefined && (
                    <div className="mt-1">
                      <span className="inline-flex items-center px-2.5 py-1 rounded-md text-xs font-semibold bg-indigo-100 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300">
                        Parcela {purchase.parcela_atual}/{purchase.quantidade_parcelas}
                      </span>
                    </div>
                  )}
              </div>

              <div className="flex items-center space-x-1 mt-3 p-1 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    setIsEditModalOpen(true);
                  }}
                  className="p-2 text-gray-500 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-all duration-200"
                  title="Editar compra"
                >
                  <Edit className="h-4 w-4" />
                </button>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    canAnticipate && setIsAnticipateModalOpen(true);
                  }}
                  disabled={!canAnticipate}
                  className={`p-2 rounded-lg transition-all duration-200 ${
                    canAnticipate
                      ? 'text-gray-500 dark:text-gray-400 hover:text-purple-600 dark:hover:text-purple-400 hover:bg-purple-50 dark:hover:bg-purple-900/20 cursor-pointer'
                      : 'text-gray-300 dark:text-gray-600 cursor-not-allowed opacity-50'
                  }`}
                  title={
                    canAnticipate
                      ? 'Antecipar parcelas'
                      : purchase.quantidade_parcelas === 1
                      ? 'Antecipação disponível apenas para compras parceladas'
                      : 'Não é possível antecipar a última parcela'
                  }
                >
                  <Clock className="h-4 w-4" />
                </button>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    setIsDeleteModalOpen(true);
                  }}
                  className="p-2 text-gray-500 dark:text-gray-400 hover:text-red-600 dark:hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-all duration-200"
                  title="Remover compra"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <EditPurchaseModal
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        purchase={purchase}
        onPurchaseUpdated={onPurchaseUpdated}
      />

      <DeletePurchaseModal
        isOpen={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        purchase={purchase}
        onPurchaseDeleted={onPurchaseUpdated}
      />

      <AnticipatePurchaseModal
        isOpen={isAnticipateModalOpen}
        onClose={() => setIsAnticipateModalOpen(false)}
        purchase={purchase}
        invoiceId={invoiceId}
        onPurchaseUpdated={onPurchaseUpdated}
      />
    </>
  );
};

export default PurchaseItem;