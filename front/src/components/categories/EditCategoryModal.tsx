import React, { useState, useEffect } from 'react';
import { X, Tag } from 'lucide-react';
import { useCategories } from '../../context/CategoryContext';
import { Category } from '../../types/category';

interface EditCategoryModalProps {
  isOpen: boolean;
  onClose: () => void;
  category: Category;
}

const EditCategoryModal: React.FC<EditCategoryModalProps> = ({
  isOpen,
  onClose,
  category
}) => {
  const { updateCategory } = useCategories();
  const [categoryName, setCategoryName] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (isOpen && category.nome) {
      setCategoryName(category.nome);
    }
  }, [isOpen, category]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!categoryName.trim()) {
      setError('O nome da categoria é obrigatório');
      return;
    }

    if (!category.id) {
      setError('ID da categoria não encontrado');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      await updateCategory(category.id, categoryName.trim());
      onClose();
    } catch (err) {
      setError('Erro ao atualizar categoria');
      console.error('Error updating category:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    onClose();
    setCategoryName(category.nome || '');
    setError(null);
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center space-x-3">
              <Tag className="h-6 w-6 text-blue-600" />
              <h2 className="text-2xl font-bold text-gray-900">Editar Categoria</h2>
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

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="categoryName" className="block text-sm font-medium text-gray-700">
                Nome da Categoria
              </label>
              <input
                type="text"
                id="categoryName"
                value={categoryName}
                onChange={(e) => setCategoryName(e.target.value)}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                placeholder="Digite o nome da categoria"
                required
              />
            </div>

            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={handleClose}
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
                {isLoading ? 'Atualizando...' : 'Atualizar Categoria'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default EditCategoryModal;