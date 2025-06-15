import React from "react";
import { Invoice } from "../../types/invoice";
import { Check, Clock, Calendar } from "lucide-react";

interface InvoiceItemProps {
  invoice: Invoice;
  onClick: () => void;
}

const InvoiceItem: React.FC<InvoiceItemProps> = ({ invoice, onClick }) => {
  // Determine status appearance
  const getStatusDetails = () => {
    if (!invoice.status)
      return { color: "gray", icon: Clock, text: "Status desconhecido" };

    const status = invoice.status.toLowerCase();
    if (status === "pago" || status === "paid") {
      return {
        color: "green",
        icon: Check,
        text: "Pago",
      };
    } else {
      return {
        color: "amber",
        icon: Clock,
        text: "Pendente",
      };
    }
  };

  const { color, icon: StatusIcon, text } = getStatusDetails();
  const colorClasses = {
    green: "bg-green-100 text-green-800",
    amber: "bg-amber-100 text-amber-800",
    gray: "bg-gray-100 text-gray-800",
  };

  return (
    <div
      onClick={onClick}
      className="bg-white rounded-lg shadow hover:shadow-md transition-shadow p-4 cursor-pointer border-l-4 border-blue-500"
    >
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-medium text-gray-900">
            {invoice.nome || "Fatura sem nome"}
          </h3>
          <p className="text-sm text-gray-600">
            {invoice.nome_cartao && `Cartão: ${invoice.nome_cartao}`}
          </p>
        </div>

        <div className="flex items-center space-x-4">
          {invoice.data_vencimento && (
            <div className="flex items-center text-sm text-gray-600">
              <Calendar className="h-4 w-4 mr-1" />
              <span>
                Vencimento:{" "}
                {new Date(invoice.data_vencimento).toLocaleDateString(
                  "pt-BR",
                  {
                    day: "2-digit",
                    month: "2-digit",
                    year: "numeric",
                    timeZone: "UTC", // Importante para evitar conversões de fuso indesejadas
                  }
                )}
              </span>
            </div>
          )}

          <span
            className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
              colorClasses[color as keyof typeof colorClasses]
            }`}
          >
            <StatusIcon className="h-3 w-3 mr-1" />
            {text}
          </span>
        </div>
      </div>
    </div>
  );
};

export default InvoiceItem;