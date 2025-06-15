import React from "react";
import { Purchase } from "../../types/purchase";
import { ShoppingCart, MapPin, Tag, Calendar } from "lucide-react";

interface PurchaseItemProps {
  purchase: Purchase;
}

const PurchaseItem: React.FC<PurchaseItemProps> = ({ purchase }) => {
  return (
    <div className="bg-white rounded-lg shadow p-4 hover:shadow-md transition-shadow">
      <div className="flex justify-between items-start">
        <div>
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

        <div className="flex flex-col items-end">
          <span className="text-lg font-bold text-gray-900">
            R$ {purchase.valor_parcela?.toFixed(2) || "0.00"}
          </span>

          {purchase.parcela_atual !== undefined &&
            purchase.quantidade_parcelas !== undefined && (
              <span className="text-sm text-gray-600">
                Parcela {purchase.parcela_atual}/{purchase.quantidade_parcelas}
              </span>
            )}
        </div>
      </div>
    </div>
  );
};

export default PurchaseItem;
