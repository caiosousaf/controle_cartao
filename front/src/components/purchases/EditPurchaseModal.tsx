import React, { useState, useEffect } from 'react';
import { X, ShoppingBag, AlertTriangle } from 'lucide-react';
import { categoryService } from '../../services/categoryService';
import { purchaseService } from '../../services/purchaseService';
import { format } from 'date-fns';
import { Category } from '../../types/category';
import { Purchase } from '../../types/purchase';

interface EditPurchaseModalProps {
  isOpen: boolean;
  onClose: () => void;
  purchase: Purchase;
  onPurchaseUpdated: () => void;
}

const EditPurchaseModal: React.FC<EditPurchaseModalProps> = ({
  isOpen,
  onClose,
  purchase,
  onPurchaseUpdated
}) => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [formData, setFormData] = useState({
    nome: '',
    descricao: '',
    local_compra: '',
    categoria_id: '',
    valor_parcela: '',
    data_compra: ''
  });
  const [updateAllInstallments, setUpdateAllInstallments] = useState(false);
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (isOpen) {
      fetchCategories();
      // Populate form with current purchase data
      setFormData({
        nome: purchase.nome || '',
        descricao: purchase.descricao || '',
        local_compra: purchase.local_compra || '',
        categoria_id: purchase.categoria_id || '',
        valor_parcela: purchase.valor_parcela?.toString() || '',
        data_compra: purchase.data_compra ? format(new Date(purchase.data_compra), 'yyyy-MM-dd') : ''
      });
    }
  }, [isOpen, purchase]);

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
    if (!formData.data_compra) return 'Data da compra é obrigatória';
    return null;
  };

  const formatDateForAPI = (dateStr: string) => {
    return dateStr;
  };

  const formatDateForDisplay = (dateStr: string) => {
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
    if (!purchase.id) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const updateData = {
        nome: formData.nome,
        descricao: formData.descricao || undefined,
        local_compra: formData.local_compra || undefined,
        categoria_id: formData.categoria_id,
        valor_parcela: parseFloat(formData.valor_parcela),
        data_compra: formatDateForAPI(formData.data_compra)
      };

      await purchaseService.updatePurchase(purchase.id, updateAllInstallments, updateData);
      
      onPurchaseUpdated();
      onClose();
    } catch (err) {
      setError('Erro ao atualizar compra');
    } finally {
      setIsLoading(false);
      setShowConfirmation(false);
    }
  };

  const handleClose = () => {
    onClose();
    setShowConfirmation(false);
    setUpdateAllInstallments(false);
    setError(null);
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
              <h2 className="text-2xl font-bold text-gray-900">Editar Compra</h2>
            </div>
            <button
              onClick={handleClose}
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
              <h3 className="text-lg font-medium text-gray-900">Confirmar Atualização</h3>
              
              {/* Update scope selection */}
              <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                <div className="flex items-start space-x-3">
                  <AlertTriangle className="h-5 w-5 text-yellow-600 mt-0.5" />
                  <div className="flex-1">
                    <h4 className="text-sm font-medium text-yellow-800 mb-2">
                      Escolha o escopo da atualização
                    </h4>
                    <div className="space-y-2">
                      <label className="flex items-center">
                        <input
                          type="radio"
                          name="updateScope"
                          checked={!updateAllInstallments}
                          onChange={() => setUpdateAllInstallments(false)}
                          className="mr-2"
                        />
                        <span className="text-sm text-yellow-700">
                          Atualizar apenas esta parcela ({purchase.parcela_atual}/{purchase.quantidade_parcelas})
                        </span>
                      </label>
                      <label className="flex items-center">
                        <input
                          type="radio"
                          name="updateScope"
                          checked={updateAllInstallments}
                          onChange={() => setUpdateAllInstallments(true)}
                          className="mr-2"
                        />
                        <span className="text-sm text-yellow-700">
                          Atualizar todas as parcelas deste agrupamento
                        </span>
                      </label>
                    </div>
                  </div>
                </div>
              </div>

              <div className="bg-gray-50 p-4 rounded-lg space-y-2">
                <p><strong>Nome da Compra:</strong> {formData.nome}</p>
                <p><strong>Local:</strong> {formData.local_compra || 'Não informado'}</p>
                <p><strong>Categoria:</strong> {selectedCategory?.nome || 'Não encontrada'}</p>
                <p><strong>Valor da Parcela:</strong> R$ {formData.valor_parcela}</p>
                <p><strong>Data da Compra:</strong> {formatDateForDisplay(formData.data_compra)}</p>
                <p><strong>Escopo:</strong> {updateAllInstallments ? 'Todas as parcelas' : 'Apenas esta parcela'}</p>
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
                  {isLoading ? 'Atualizando...' : 'Confirmar Atualização'}
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
                  Descrição
                </label>
                <textarea
                  id="descricao"
                  name="descricao"
                  value={formData.descricao}
                  onChange={handleChange}
                  rows={3}
                  className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                />
              </div>

              <div>
                <label htmlFor="local_compra" className="block text-sm font-medium text-gray-700">
                  Local da Compra
                </label>
                <input
                  type="text"
                  id="local_compra"
                  name="local_compra"
                  value={formData.local_compra}
                  onChange={handleChange}
                  className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
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

              {/* Read-only fields */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <h4 className="text-sm font-medium text-gray-700 mb-3">Informações não editáveis</h4>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <label className="block text-xs font-medium text-gray-500">Parcela Atual</label>
                    <p className="text-sm text-gray-900">{purchase.parcela_atual}</p>
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-gray-500">Total de Parcelas</label>
                    <p className="text-sm text-gray-900">{purchase.quantidade_parcelas}</p>
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-gray-500">Fatura</label>
                    <p className="text-sm text-gray-900">{purchase.nome_fatura || 'N/A'}</p>
                  </div>
                </div>
              </div>

              <div className="flex justify-end space-x-3">
                <button
                  type="button"
                  onClick={handleClose}
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

export default EditPurchaseModal;