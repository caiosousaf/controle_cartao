import React, { useState, useEffect } from 'react';
import { X, RepeatIcon } from 'lucide-react';
import { recurringService } from '../../services/recurringService';
import { categoryService } from '../../services/categoryService';
import { RecurringPurchase } from '../../types/recurring';
import { Category } from '../../types/category';

interface CreateRecurringModalProps {
  isOpen: boolean;
  onClose: () => void;
  onPurchaseCreated: () => void;
  editPurchase?: RecurringPurchase;
}

const CreateRecurringModal: React.FC<CreateRecurringModalProps> = ({
  isOpen,
  onClose,
  onPurchaseCreated,
  editPurchase
}) => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [formData, setFormData] = useState({
    nome: '',
    descricao: '',
    compra_categoria_id: '',
    local_compra: '',
    valor_parcela: ''
  });
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (isOpen) {
      fetchCategories();
      if (editPurchase) {
        setFormData({
          nome: editPurchase.nome || '',
          descricao: editPurchase.descricao || '',
          compra_categoria_id: editPurchase.compra_categoria_id || '',
          local_compra: editPurchase.local_compra || '',
          valor_parcela: editPurchase.valor_parcela?.toString() || ''
        });
      }
    }
  }, [isOpen, editPurchase]);

  const fetchCategories = async () => {
    try {
      const response = await categoryService.getCategories(100000, 0);
      // Filter only active categories
      const activeCategories = response.dados.filter(category => !category.data_desativacao);
      setCategories(activeCategories);
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
    if (!formData.compra_categoria_id) return 'Categoria é obrigatória';
    if (!formData.valor_parcela || parseFloat(formData.valor_parcela) <= 0) return 'Valor da parcela inválido';
    return null;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const validationError = validateForm();
    if (validationError) {
      setError(validationError);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const data = {
        nome: formData.nome,
        descricao: formData.descricao || undefined,
        compra_categoria_id: formData.compra_categoria_id,
        local_compra: formData.local_compra || undefined,
        valor_parcela: parseFloat(formData.valor_parcela)
      };

      if (editPurchase?.id) {
        await recurringService.updateRecurringPurchase(editPurchase.id, data);
      } else {
        await recurringService.createRecurringPurchase(data);
      }
      
      onPurchaseCreated();
      onClose();
      setFormData({
        nome: '',
        descricao: '',
        compra_categoria_id: '',
        local_compra: '',
        valor_parcela: ''
      });
    } catch (err) {
      setError(editPurchase ? 'Erro ao atualizar compra recorrente' : 'Erro ao cadastrar compra recorrente');
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center space-x-3">
              <RepeatIcon className="h-6 w-6 text-blue-600" />
              <h2 className="text-2xl font-bold text-gray-900">
                {editPurchase ? 'Editar Compra Recorrente' : 'Nova Compra Recorrente'}
              </h2>
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

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="nome" className="block text-sm font-medium text-gray-700">
                Nome *
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
              <label htmlFor="compra_categoria_id" className="block text-sm font-medium text-gray-700">
                Categoria *
              </label>
              <select
                id="compra_categoria_id"
                name="compra_categoria_id"
                value={formData.compra_categoria_id}
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

            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={onClose}
                className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                disabled={isLoading}
              >
                Cancelar
              </button>
              <button
                type="submit"
                className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50"
                disabled={isLoading}
              >
                {isLoading ? (editPurchase ? 'Atualizando...' : 'Cadastrando...') : (editPurchase ? 'Atualizar' : 'Cadastrar')}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default CreateRecurringModal