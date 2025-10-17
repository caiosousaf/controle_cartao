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
      <div className="bg-white rounded-lg shadow p-4 hover:shadow-md transition-shadow">
        <div className="flex justify-between items-start">
          <div>
            <h3 className="text-lg font-medium text-gray-900 flex items-center">
              <Tag className="h-5 w-5 mr-2 text-blue-500" />
              {category.nome || "Categoria sem nome"}
            </h3>

            <div className="mt-3 flex flex-wrap gap-2">
              {category.data_criacao && (
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  <Calendar className="h-3 w-3 mr-1" />
                  {new Date(category.data_criacao).toLocaleDateString("pt-BR", {
                    day: "2-digit",
                    month: "2-digit",
                    year: "numeric",
                    timeZone: "UTC",
                  })}
                </span>
              )}

              <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                isActive 
                  ? 'bg-green-100 text-green-800' 
                  : 'bg-red-100 text-red-800'
              }`}>
                {isActive ? (
                  <CheckCircle className="h-3 w-3 mr-1" />
                ) : (
                  <XCircle className="h-3 w-3 mr-1" />
                )}
                {isActive ? 'Ativa' : 'Inativa'}
              </span>
            </div>
          </div>

          <div className="flex items-center space-x-2">
            {isActive && (
              <button
                onClick={() => setIsEditing(true)}
                className="p-1 text-gray-500 hover:text-blue-600 transition-colors"
                title="Editar"
              >
                <Edit className="h-4 w-4" />
              </button>
            )}
            
            {isConfirmingAction ? (
              <div className="flex items-center space-x-2">
                <button
                  onClick={handleToggleActive}
                  disabled={isLoading}
                  className="text-xs bg-blue-600 text-white px-2 py-1 rounded hover:bg-blue-700"
                >
                  Confirmar
                </button>
                <button
                  onClick={() => setIsConfirmingAction(null)}
                  disabled={isLoading}
                  className="text-xs bg-gray-300 text-gray-700 px-2 py-1 rounded hover:bg-gray-400"
                >
                  Cancelar
                </button>
              </div>
            ) : (
              <button
                onClick={() => setIsConfirmingAction(isActive ? 'deactivate' : 'reactivate')}
                className={`p-1 ${isActive ? 'text-red-500 hover:text-red-600' : 'text-green-500 hover:text-green-600'} transition-colors`}
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