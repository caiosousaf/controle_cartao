import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useInvoices } from "../context/InvoiceContext";
import { usePurchases } from "../context/PurchaseContext";
import PurchaseList from "../components/purchases/PurchaseList";
import CreatePurchaseModal from "../components/purchases/CreatePurchaseModal";
import PayInvoiceModal from "../components/invoices/PayInvoiceModal";
import { ArrowLeft, Receipt, Calendar, Check, Clock, Plus } from "lucide-react";

const InvoiceDetailsPage: React.FC = () => {
  const { invoiceId } = useParams<{ invoiceId: string }>();
  const { invoices, selectedInvoice, selectInvoice } = useInvoices();
  const { clearPurchases, fetchPurchasesByInvoiceId } = usePurchases();
  const navigate = useNavigate();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isPayModalOpen, setIsPayModalOpen] = useState(false);

  useEffect(() => {
    if (!invoiceId) {
      navigate("/");
      return;
    }

    // If we already have invoices loaded, find the matching invoice
    const invoice = invoices.find((i) => i.id === invoiceId);
    if (invoice) {
      selectInvoice(invoice);
    }

    // Cleanup when unmounting
    return () => {
      clearPurchases();
    };
  }, [invoiceId, invoices]);

  const handleBack = () => {
    if (selectedInvoice?.fatura_cartao_id) {
      navigate(`/cards/${selectedInvoice.fatura_cartao_id}`);
    } else {
      navigate("/");
    }
  };

  // Determine status appearance
  const getStatusDetails = () => {
    if (!selectedInvoice?.status)
      return { color: "gray", icon: Clock, text: "Status desconhecido" };

    const status = selectedInvoice.status.toLowerCase();
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

  const handlePurchaseCreated = () => {
    if (invoiceId) {
      fetchPurchasesByInvoiceId(invoiceId);
    }
  };

  const handleInvoicePaid = () => {
    if (invoiceId) {
      // Refresh the invoice data
      const invoice = invoices.find((i) => i.id === invoiceId);
      if (invoice) {
        selectInvoice({ ...invoice, status: "Pago" });
      }
    }
  };

  const isPaid =
    selectedInvoice?.status?.toLowerCase() === "pago" ||
    selectedInvoice?.status?.toLowerCase() === "paid";

  return (
    <div className="space-y-8">
      <div className="flex items-center space-x-4">
        <button
          onClick={handleBack}
          className="p-2 rounded-full text-gray-600 hover:bg-gray-100 hover:text-gray-900 transition-colors"
        >
          <ArrowLeft className="h-6 w-6" />
        </button>
        <h1 className="text-3xl font-bold text-gray-900">Detalhes da Fatura</h1>
      </div>

      {selectedInvoice && (
        <div className="bg-white rounded-lg shadow-lg p-6">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between">
            <div>
              <div className="flex items-center space-x-3">
                <Receipt className="h-8 w-8 text-blue-600" />
                <h2 className="text-2xl font-bold text-gray-900">
                  {selectedInvoice.nome}
                </h2>
                <span
                  className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    colorClasses[color as keyof typeof colorClasses]
                  }`}
                >
                  <StatusIcon className="h-3 w-3 mr-1" />
                  {text}
                </span>
              </div>
              <p className="mt-2 text-gray-600">
                Cartão: {selectedInvoice.nome_cartao || "Não especificado"}
              </p>
            </div>

            <div className="mt-4 md:mt-0 flex items-center space-x-4">
              {selectedInvoice.data_vencimento && (
                <div className="flex items-center space-x-2 text-gray-700">
                  <Calendar className="h-5 w-5 text-gray-500" />
                  <span>
                    Vencimento:{" "}
                    {new Date(
                      selectedInvoice.data_vencimento
                    ).toLocaleDateString("pt-BR", {
                      day: "2-digit",
                      month: "2-digit",
                      year: "numeric",
                      timeZone: "UTC",
                    })}
                  </span>
                </div>
              )}

              <div className="flex space-x-2">
                {!isPaid && (
                  <>
                    <button
                      onClick={() => setIsCreateModalOpen(true)}
                      className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                    >
                      <Plus className="h-4 w-4 mr-2" />
                      Nova Compra
                    </button>
                    <button
                      onClick={() => setIsPayModalOpen(true)}
                      className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                    >
                      <Check className="h-4 w-4 mr-2" />
                      Pagar Fatura
                    </button>
                  </>
                )}
              </div>
            </div>
          </div>
        </div>
      )}

      {invoiceId && (
        <>
          <PurchaseList invoiceId={invoiceId} />
          {selectedInvoice && (
            <>
              <CreatePurchaseModal
                isOpen={isCreateModalOpen}
                onClose={() => setIsCreateModalOpen(false)}
                invoiceId={invoiceId}
                invoiceName={selectedInvoice.nome || ""}
                onPurchaseCreated={handlePurchaseCreated}
              />
              <PayInvoiceModal
                isOpen={isPayModalOpen}
                onClose={() => setIsPayModalOpen(false)}
                invoiceId={invoiceId}
                invoiceName={selectedInvoice.nome || ""}
                onInvoicePaid={handleInvoicePaid}
              />
            </>
          )}
        </>
      )}
    </div>
  );
};

export default InvoiceDetailsPage;