import React, { useState, useEffect } from 'react';
import { X, ShoppingBag } from 'lucide-react';
import { categoryService } from '../../services/categoryService';
import { purchaseService } from '../../services/purchaseService';
import { format } from 'date-fns';
import { Category } from '../../types/category';

interface CreatePurchaseModalProps {
  isOpen: boolean;
  onClose: () => void;
  invoiceId: string;
  invoiceName: string;
  onPurchaseCreated: () => void;
}

const CreatePurchaseModal: React.FC<CreatePurchaseModalProps> = ({
  isOpen,
  onClose,
  invoiceId,
  invoiceName,
  onPurchaseCreated
}) => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [formData, setFormData] = useState({
    nome: '',
    descricao: '',
    local_compra: '',
    categoria_id: '',
    valor_parcela: '',
    parcela_atual: '1',
    quantidade_parcelas: '1',
    data_compra: format(new Date(), 'yyyy-MM-dd')
  });
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (isOpen) {
      fetchCategories();
    }
  }, [isOpen]);

  const fetchCategories = async () => {
    try {
      const response = await categoryService.getCategories(100000, 0, false, true);
      setCategories(response.dados);
    } catch (err) {
      setError('Erro ao carregar categorias');
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const validateForm = () => {
    if (!formData.nome.trim()) return 'Nome é obrigatório';
    if (!formData.categoria_id) return 'Categoria é obrigatória';
    if (!formData.valor_parcela || parseFloat(formData.valor_parcela) <= 0) return 'Valor da parcela inválido';
    if (parseInt(formData.parcela_atual) < 1) return 'Parcela atual inválida';
    if (parseInt(formData.quantidade_parcelas) < 1) return 'Quantidade de parcelas inválida';
    if (parseInt(formData.parcela_atual) > parseInt(formData.quantidade_parcelas)) {
      return 'Parcela atual não pode ser maior que a quantidade de parcelas';
    }
    if (!formData.data_compra) return 'Data da compra é obrigatória';
    return null;
  };

  const formatDateForAPI = (dateStr: string) => {
    // Simply return the date string as is, since it's already in YYYY-MM-DD format
    return dateStr;
  };

  const formatDateForDisplay = (dateStr: string) => {
    // Parse the date string and format for display
    const [year, month, day] = dateStr.split('-');
    return `${day}/${month}/${year}`;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const validationError = validateForm();
    if (validationError) {
      setError(validationError);
      return;
    }
    setShowConfirmation(true);
  };

  const handleConfirm = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      await purchaseService.createPurchase(invoiceId, {
        ...formData,
        valor_parcela: parseFloat(formData.valor_parcela),
        parcela_atual: parseInt(formData.parcela_atual),
        quantidade_parcelas: parseInt(formData.quantidade_parcelas),
        fatura_id: invoiceId,
        data_compra: formatDateForAPI(formData.data_compra)
      });
      
      onPurchaseCreated();
      onClose();
    } catch (err) {
      setError('Erro ao cadastrar compra');
    } finally {
      setIsLoading(false);
      setShowConfirmation(false);
    }
  };

  if (!isOpen) return null;

  const selectedCategory = categories.find(cat => cat.id === formData.categoria_id);

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4 overflow-y-auto">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-2xl my-8">
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center space-x-3">
              <ShoppingBag className="h-6 w-6 text-blue-600" />
              <h2 className="text-2xl font-bold text-gray-900">Nova Compra</h2>
            </div>
            <button
              onClick={onClose}
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

          {showConfirmation ? (
            <div className="space-y-6">
              <h3 className="text-lg font-medium text-gray-900">Confirmar Compra</h3>
              <div className="bg-gray-50 p-4 rounded-lg space-y-2">
                <p><strong>Fatura:</strong> {invoiceName}</p>
                <p><strong>Nome da Compra:</strong> {formData.nome}</p>
                <p><strong>Local:</strong> {formData.local_compra || 'Não informado'}</p>
                <p><strong>Categoria:</strong> {selectedCategory?.nome || 'Não encontrada'}</p>
                <p><strong>Valor da Parcela:</strong> R$ {formData.valor_parcela}</p>
                <p><strong>Parcelas:</strong> {formData.parcela_atual}/{formData.quantidade_parcelas}</p>
                <p><strong>Data da Compra:</strong> {formatDateForDisplay(formData.data_compra)}</p>
              </div>
              
              <div className="flex justify-end space-x-3">
                <button
                  type="button"
                  onClick={() => setShowConfirmation(false)}
                  className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                  disabled={isLoading}
                >
                  Voltar
                </button>
                <button
                  type="button"
                  onClick={handleConfirm}
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50"
                  disabled={isLoading}
                >
                  {isLoading ? 'Cadastrando...' : 'Confirmar'}
                </button>
              </div>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-6">
              <div>
                <label htmlFor="nome" className="block text-sm font-medium text-gray-700">
                  Nome da Compra *
                </label>
                <input
                  type="text"
                  id="nome"
                  name="nome"
                  value={formData.nome}
                  onChange={handleChange}
                  className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                  required
                />
              </div>

              <div>
                <label htmlFor="descricao" className="block text-sm font-medium text-gray-700">
                  Descrição *
                </label>
                <textarea
                  id="descricao"
                  name="descricao"
                  value={formData.descricao}
                  onChange={handleChange}
                  rows={3}
                  className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                  required
                />
              </div>

              <div>
                <label htmlFor="local_compra" className="block text-sm font-medium text-gray-700">
                  Local da Compra *
                </label>
                <input
                  type="text"
                  id="local_compra"
                  name="local_compra"
                  value={formData.local_compra}
                  onChange={handleChange}
                  className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                  required
                />
              </div>

              <div>
                <label htmlFor="categoria_id" className="block text-sm font-medium text-gray-700">
                  Categoria *
                </label>
                <select
                  id="categoria_id"
                  name="categoria_id"
                  value={formData.categoria_id}
                  onChange={handleChange}
                  className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                  required
                >
                  <option value="">Selecione uma categoria</option>
                  {categories.map(category => (
                    <option key={category.id} value={category.id}>
                      {category.nome}
                    </option>
                  ))}
                </select>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label htmlFor="valor_parcela" className="block text-sm font-medium text-gray-700">
                    Valor da Parcela *
                  </label>
                  <div className="mt-1 relative rounded-md shadow-sm">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <span className="text-gray-500 sm:text-sm">R$</span>
                    </div>
                    <input
                      type="number"
                      id="valor_parcela"
                      name="valor_parcela"
                      value={formData.valor_parcela}
                      onChange={handleChange}
                      step="0.01"
                      min="0.01"
                      className="pl-8 mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                      required
                    />
                  </div>
                </div>

                <div>
                  <label htmlFor="data_compra" className="block text-sm font-medium text-gray-700">
                    Data da Compra *
                  </label>
                  <input
                    type="date"
                    id="data_compra"
                    name="data_compra"
                    value={formData.data_compra}
                    onChange={handleChange}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                    required
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label htmlFor="parcela_atual" className="block text-sm font-medium text-gray-700">
                    Parcela Atual *
                  </label>
                  <input
                    type="number"
                    id="parcela_atual"
                    name="parcela_atual"
                    value={formData.parcela_atual}
                    onChange={handleChange}
                    min="1"
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                    required
                  />
                </div>

                <div>
                  <label htmlFor="quantidade_parcelas" className="block text-sm font-medium text-gray-700">
                    Quantidade de Parcelas *
                  </label>
                  <input
                    type="number"
                    id="quantidade_parcelas"
                    name="quantidade_parcelas"
                    value={formData.quantidade_parcelas}
                    onChange={handleChange}
                    min="1"
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                    required
                  />
                </div>
              </div>

              <div className="flex justify-end space-x-3">
                <button
                  type="button"
                  onClick={onClose}
                  className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                >
                  Prosseguir
                </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
};

export default CreatePurchaseModal;