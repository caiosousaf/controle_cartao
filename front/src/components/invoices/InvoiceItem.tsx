import React from "react";
import { Invoice } from "../../types/invoice";
import { Check, Clock, Calendar, CreditCard } from "lucide-react";

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
      className="group bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-xl border border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-600 transition-all duration-300 cursor-pointer overflow-hidden relative"
    >
      {/* Gradient accent bar */}
      <div className={`absolute left-0 top-0 bottom-0 w-1 bg-gradient-to-b ${
        color === 'green' ? 'from-green-400 to-emerald-600' :
        color === 'amber' ? 'from-amber-400 to-orange-600' :
        'from-gray-400 to-gray-600'
      }`}></div>
      
      <div className="pl-5 pr-6 py-5">
        <div className="flex items-center justify-between">
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-3 mb-2">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white truncate">
                {invoice.nome || "Fatura sem nome"}
              </h3>
              <span
                className={`inline-flex items-center px-3 py-1 rounded-lg text-xs font-semibold shadow-sm ${
                  color === 'green' 
                    ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300' 
                    : color === 'amber'
                    ? 'bg-amber-100 dark:bg-amber-900/30 text-amber-800 dark:text-amber-300'
                    : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300'
                }`}
              >
                <StatusIcon className="h-3 w-3 mr-1.5" />
                {text}
              </span>
            </div>
            
            <div className="flex flex-wrap items-center gap-4 text-sm text-gray-600 dark:text-gray-400">
              {invoice.nome_cartao && (
                <div className="flex items-center">
                  <CreditCard className="h-4 w-4 mr-1.5 text-gray-400" />
                  <span className="font-medium">{invoice.nome_cartao}</span>
                </div>
              )}
              
              {invoice.data_vencimento && (
                <div className="flex items-center">
                  <Calendar className="h-4 w-4 mr-1.5 text-gray-400" />
                  <span>
                    Vence em {new Date(invoice.data_vencimento).toLocaleDateString(
                      "pt-BR",
                      {
                        day: "2-digit",
                        month: "2-digit",
                        year: "numeric",
                        timeZone: "UTC",
                      }
                    )}
                  </span>
                </div>
              )}
            </div>
          </div>

          {/* Arrow indicator */}
          <div className="ml-4 flex-shrink-0">
            <div className="p-2 rounded-lg bg-gray-100 dark:bg-gray-700 group-hover:bg-blue-100 dark:group-hover:bg-blue-900/30 transition-colors">
              <svg className="w-5 h-5 text-gray-400 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default InvoiceItem;