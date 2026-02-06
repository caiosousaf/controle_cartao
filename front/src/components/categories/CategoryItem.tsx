import React, { useState } from 'react';
import { Category } from '../../types/category';
import { Tag, Calendar, CheckCircle, XCircle, Edit, Power, Trash2 } from 'lucide-react';
import { useCategories } from '../../context/CategoryContext';
import EditCategoryModal from './EditCategoryModal';

interface CategoryItemProps {
  category: Category;
}

const CategoryItem: React.FC<CategoryItemProps> = ({ category }) => {
  const { deactivateCategory, reactivateCategory } = useCategories();
  const [isConfirmingAction, setIsConfirmingAction] = useState<'deactivate' | 'reactivate' | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const isActive = !category.data_desativacao;

  const handleToggleActive = async () => {
    if (!category.id) return;
    
    setIsLoading(true);
    try {
      if (isActive) {
        await deactivateCategory(category.id);
      } else {
        await reactivateCategory(category.id);
      }
    } catch (error) {
      console.error('Error toggling category status:', error);
    } finally {
      setIsLoading(false);
      setIsConfirmingAction(null);
    }
  };

  return (
    <>
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-lg border border-gray-200 dark:border-gray-700 transition-all duration-300 p-5">
        <div className="flex justify-between items-start">
          <div className="flex-1">
            <div className="flex items-center gap-3 mb-3">
              <div className="p-2 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
                <Tag className="h-5 w-5 text-blue-600 dark:text-blue-400" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                {category.nome || "Categoria sem nome"}
              </h3>
            </div>

            <div className="ml-11 flex flex-wrap gap-2">
              {category.data_criacao && (
                <span className="inline-flex items-center px-3 py-1 rounded-lg text-xs font-medium bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 shadow-sm">
                  <Calendar className="h-3 w-3 mr-1.5" />
                  {new Date(category.data_criacao).toLocaleDateString("pt-BR", {
                    day: "2-digit",
                    month: "2-digit",
                    year: "numeric",
                    timeZone: "UTC",
                  })}
                </span>
              )}

              <span className={`inline-flex items-center px-3 py-1 rounded-lg text-xs font-medium shadow-sm ${
                isActive 
                  ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300' 
                  : 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300'
              }`}>
                {isActive ? (
                  <CheckCircle className="h-3 w-3 mr-1.5" />
                ) : (
                  <XCircle className="h-3 w-3 mr-1.5" />
                )}
                {isActive ? 'Ativa' : 'Inativa'}
              </span>
            </div>
          </div>

          <div className="flex items-center space-x-1 ml-4 p-1 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
            {isActive && (
              <button
                onClick={() => setIsEditing(true)}
                className="p-2 text-gray-500 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-all duration-200"
                title="Editar"
              >
                <Edit className="h-4 w-4" />
              </button>
            )}
            
            {isConfirmingAction ? (
              <div className="flex items-center space-x-1">
                <button
                  onClick={handleToggleActive}
                  disabled={isLoading}
                  className="text-xs bg-blue-600 dark:bg-blue-500 text-white px-2 py-1 rounded-md hover:bg-blue-700 dark:hover:bg-blue-600 transition-colors"
                >
                  Confirmar
                </button>
                <button
                  onClick={() => setIsConfirmingAction(null)}
                  disabled={isLoading}
                  className="text-xs bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-300 px-2 py-1 rounded-md hover:bg-gray-400 dark:hover:bg-gray-500 transition-colors"
                >
                  Cancelar
                </button>
              </div>
            ) : (
              <button
                onClick={() => setIsConfirmingAction(isActive ? 'deactivate' : 'reactivate')}
                className={`p-2 rounded-lg transition-all duration-200 ${
                  isActive 
                    ? 'text-red-500 dark:text-red-400 hover:text-red-600 dark:hover:text-red-300 hover:bg-red-50 dark:hover:bg-red-900/20' 
                    : 'text-green-500 dark:text-green-400 hover:text-green-600 dark:hover:text-green-300 hover:bg-green-50 dark:hover:bg-green-900/20'
                }`}
                title={isActive ? 'Desativar' : 'Reativar'}
              >
                {isActive ? <Trash2 className="h-4 w-4" /> : <Power className="h-4 w-4" />}
              </button>
            )}
          </div>
        </div>
      </div>

      <EditCategoryModal
        isOpen={isEditing}
        onClose={() => setIsEditing(false)}
        category={category}
      />
    </>
  );
};

export default CategoryItem;