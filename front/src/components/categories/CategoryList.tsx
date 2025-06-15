import React, { useEffect, useState } from 'react';
import { useCategories } from '../../context/CategoryContext';
import CategoryItem from './CategoryItem';
import CreateCategoryModal from './CreateCategoryModal';
import { Tag, Plus, Filter } from 'lucide-react';

const CategoryList: React.FC = () => {
  const { 
    categories, 
    isLoading, 
    error, 
    hasMore, 
    totalCategories, 
    fetchCategories 
  } = useCategories();
  const [page, setPage] = useState(1);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [statusFilter, setStatusFilter] = useState<'all' | 'active' | 'inactive'>('all');
  const pageSize = 10;
  
  useEffect(() => {
    // Initial load with filter
    const ativo = statusFilter === 'all' ? undefined : statusFilter === 'active';
    fetchCategories(pageSize, 0, ativo);
    setPage(1);
  }, [statusFilter]);

  const loadMoreCategories = () => {
    if (hasMore && !isLoading) {
      const nextPage = page + 1;
      const ativo = statusFilter === 'all' ? undefined : statusFilter === 'active';
      fetchCategories(pageSize, (nextPage - 1) * pageSize, ativo);
      setPage(nextPage);
    }
  };

  const handleCategoryCreated = () => {
    const ativo = statusFilter === 'all' ? undefined : statusFilter === 'active';
    fetchCategories(pageSize, 0, ativo);
    setPage(1);
  };

  const handleFilterChange = (filter: 'all' | 'active' | 'inactive') => {
    setStatusFilter(filter);
  };

  const getFilterLabel = () => {
    switch (statusFilter) {
      case 'active':
        return 'Ativas';
      case 'inactive':
        return 'Inativas';
      default:
        return 'Todas';
    }
  };

  // Show error state
  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold text-gray-800">Categorias</h2>
          <button
            onClick={() => setIsCreateModalOpen(true)}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            <Plus className="h-4 w-4 mr-2" />
            Nova Categoria
          </button>
        </div>
        
        <div className="text-center p-6 bg-red-50 rounded-lg">
          <p className="text-red-600">{error}</p>
          <button
            onClick={() => {
              const ativo = statusFilter === 'all' ? undefined : statusFilter === 'active';
              fetchCategories(pageSize, 0, ativo);
            }}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
          >
            Tentar novamente
          </button>
        </div>

        <CreateCategoryModal
          isOpen={isCreateModalOpen}
          onClose={() => setIsCreateModalOpen(false)}
        />
      </div>
    );
  }

  // Show loading state for initial load
  if (isLoading && (!categories || categories.length === 0)) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold text-gray-800">Categorias</h2>
          <button
            disabled
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-gray-400 cursor-not-allowed"
          >
            <Plus className="h-4 w-4 mr-2" />
            Nova Categoria
          </button>
        </div>

        <div className="space-y-4">
          {[...Array(3)].map((_, index) => (
            <div key={index} className="animate-pulse bg-white rounded-lg shadow p-4">
              <div className="flex justify-between items-start">
                <div>
                  <div className="h-6 bg-gray-200 rounded w-3/4 mb-3"></div>
                  <div className="flex space-x-2">
                    <div className="h-4 bg-gray-200 rounded w-20"></div>
                    <div className="h-4 bg-gray-200 rounded w-16"></div>
                  </div>
                </div>
                <div className="flex space-x-2">
                  <div className="h-8 w-8 bg-gray-200 rounded"></div>
                  <div className="h-8 w-8 bg-gray-200 rounded"></div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <h2 className="text-2xl font-bold text-gray-800">Categorias</h2>
        
        <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
          {/* Status Filter */}
          <div className="flex items-center space-x-2">
            <Filter className="h-4 w-4 text-gray-500" />
            <span className="text-sm text-gray-600">Filtrar:</span>
            <div className="flex rounded-md shadow-sm">
              <button
                onClick={() => handleFilterChange('all')}
                className={`px-3 py-1 text-xs font-medium rounded-l-md border ${
                  statusFilter === 'all'
                    ? 'bg-blue-600 text-white border-blue-600'
                    : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                }`}
              >
                Todas
              </button>
              <button
                onClick={() => handleFilterChange('active')}
                className={`px-3 py-1 text-xs font-medium border-t border-b ${
                  statusFilter === 'active'
                    ? 'bg-green-600 text-white border-green-600'
                    : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                }`}
              >
                Ativas
              </button>
              <button
                onClick={() => handleFilterChange('inactive')}
                className={`px-3 py-1 text-xs font-medium rounded-r-md border ${
                  statusFilter === 'inactive'
                    ? 'bg-red-600 text-white border-red-600'
                    : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                }`}
              >
                Inativas
              </button>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <p className="text-sm text-gray-500">
              {totalCategories} {getFilterLabel().toLowerCase()}
            </p>
            <button
              onClick={() => setIsCreateModalOpen(true)}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <Plus className="h-4 w-4 mr-2" />
              Nova Categoria
            </button>
          </div>
        </div>
      </div>

      {(!categories || categories.length === 0) && !isLoading ? (
        <div className="flex flex-col items-center justify-center p-8 bg-gray-50 rounded-lg border border-gray-200">
          <Tag className="h-16 w-16 text-gray-400 mb-4" />
          <h3 className="text-lg font-medium text-gray-900">
            {statusFilter === 'all' 
              ? 'Nenhuma categoria encontrada'
              : `Nenhuma categoria ${statusFilter === 'active' ? 'ativa' : 'inativa'} encontrada`
            }
          </h3>
          <p className="mt-1 text-gray-500">
            {statusFilter === 'all' 
              ? 'Você ainda não possui categorias cadastradas.'
              : `Não há categorias ${statusFilter === 'active' ? 'ativas' : 'inativas'} no momento.`
            }
          </p>
          {statusFilter === 'all' && (
            <button
              onClick={() => setIsCreateModalOpen(true)}
              className="mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <Plus className="h-4 w-4 mr-2" />
              Criar Primeira Categoria
            </button>
          )}
        </div>
      ) : (
        <div className="space-y-4">
          {categories.map((category) => (
            <CategoryItem 
              key={category.id} 
              category={category}
            />
          ))}
        </div>
      )}

      {hasMore && (
        <div className="mt-6 text-center">
          <button
            onClick={loadMoreCategories}
            disabled={isLoading}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50"
          >
            {isLoading ? 'Carregando...' : 'Carregar mais categorias'}
          </button>
        </div>
      )}

      <CreateCategoryModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
      />
    </div>
  );
};

export default CategoryList;