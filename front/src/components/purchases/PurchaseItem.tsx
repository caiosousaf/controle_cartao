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
      <div className="bg-white rounded-lg shadow p-4 hover:shadow-md transition-shadow">
        <div className="flex justify-between items-start">
          <div className="flex-1">
            <h3 className="text-lg font-medium text-gray-900 flex items-center">
              <ShoppingCart className="h-5 w-5 mr-2 text-blue-500" />
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

              {purchase.categoria_nome && (
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                  <Tag className="h-3 w-3 mr-1" />
                  {purchase.categoria_nome}
                </span>
              )}

              {purchase.data_compra && (
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  <Calendar className="h-3 w-3 mr-1" />
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

          <div className="flex flex-col items-end ml-4">
            <span className="text-lg font-bold text-gray-900">
              R$ {purchase.valor_parcela?.toFixed(2) || "0.00"}
            </span>

            {purchase.parcela_atual !== undefined &&
              purchase.quantidade_parcelas !== undefined && (
                <span className="text-sm text-gray-600">
                  Parcela {purchase.parcela_atual}/{purchase.quantidade_parcelas}
                </span>
              )}

            <div className="flex items-center space-x-2 mt-2">
              <button
                onClick={() => setIsEditModalOpen(true)}
                className="p-1 text-gray-500 hover:text-blue-600 transition-colors"
                title="Editar compra"
              >
                <Edit className="h-4 w-4" />
              </button>
              <button
                onClick={() => canAnticipate && setIsAnticipateModalOpen(true)}
                disabled={!canAnticipate}
                className={`p-1 transition-colors ${
                  canAnticipate
                    ? 'text-gray-500 hover:text-purple-600 cursor-pointer'
                    : 'text-gray-300 cursor-not-allowed opacity-50'
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
                onClick={() => setIsDeleteModalOpen(true)}
                className="p-1 text-gray-500 hover:text-red-600 transition-colors"
                title="Remover compra"
              >
                <Trash2 className="h-4 w-4" />
              </button>
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