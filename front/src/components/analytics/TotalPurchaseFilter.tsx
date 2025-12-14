import React, { useState, useEffect } from 'react';
import { usePurchases } from '../../context/PurchaseContext';
import { PurchaseFilters } from '../../types/purchase';
import { Calendar, DollarSign, Check, Tag } from 'lucide-react';
import { categoryService } from '../../services/categoryService';
import { Category } from '../../types/category';

const TotalPurchaseFilter: React.FC = () => {
  const { isLoading, error, totalAmount, fetchTotalPurchases } = usePurchases();
  const [filters, setFilters] = useState<PurchaseFilters>({});
  const [categories, setCategories] = useState<Category[]>([]);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const response = await categoryService.getCategories(100000, 0, false, true);
        setCategories(response.dados);
      } catch (err) {
        console.error('Erro ao carregar categorias para o filtro', err);
      }
    };

    fetchCategories();
  }, []);
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await fetchTotalPurchases(filters);
  };
  
  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target as HTMLInputElement;
    
    if (type === 'checkbox') {
      const checked = (e.target as HTMLInputElement).checked;
      setFilters(prev => ({ ...prev, [name]: checked }));
    } else if (name === 'dataEspecifica') {
      setFilters(prev => ({ ...prev, [name]: value || undefined }));
    } else {
      const booleanValue = value === 'true' ? true : value === 'false' ? false : undefined;
      setFilters(prev => ({ ...prev, [name]: booleanValue !== undefined ? booleanValue : value || undefined }));
    }
  };

  return (
    <div className="bg-white rounded-lg shadow p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-800">Total de Compras</h2>
      </div>
      
      {error && (
        <div className="p-4 text-sm text-red-700 bg-red-100 rounded-lg" role="alert">
          {error}
        </div>
      )}
      
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div>
            <label htmlFor="dataEspecifica" className="block text-sm font-medium text-gray-700 mb-1">
              Mês/Ano
            </label>
            <div className="relative mt-1 rounded-md shadow-sm">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <Calendar className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type="text"
                name="dataEspecifica"
                id="dataEspecifica"
                placeholder="MM/AAAA"
                className="focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 pr-3 py-2 sm:text-sm border-gray-300 rounded-md"
                onChange={handleChange}
                value={filters.dataEspecifica || ''}
              />
            </div>
          </div>

          <div>
            <label htmlFor="categoria_id" className="block text-sm font-medium text-gray-700 mb-1">
              Categoria
            </label>
            <div className="relative mt-1 rounded-md shadow-sm">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Tag className="h-5 w-5 text-gray-400" />
                </div>
                <select
                id="categoria_id"
                name="categoria_id"
                className="focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 pr-3 py-2 sm:text-sm border-gray-300 rounded-md"
                onChange={handleChange}
                value={filters.categoria_id ? String(filters.categoria_id) : ''}
                >
                <option value="">Todas as Categorias</option>
                {categories.map((category) => (
                    <option key={category.id} value={category.id}>
                    {category.nome}
                    </option>
                ))}
                </select>
            </div>
          </div>
          
          <div>
            <label htmlFor="ultimaParcela" className="block text-sm font-medium text-gray-700 mb-1">
              Última Parcela
            </label>
            <select
              id="ultimaParcela"
              name="ultimaParcela"
              className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
              onChange={handleChange}
              value={filters.ultimaParcela !== undefined ? String(filters.ultimaParcela) : ''}
            >
              <option value="">Todos</option>
              <option value="true">Sim</option>
              <option value="false">Não</option>
            </select>
          </div>
          
          <div>
            <label htmlFor="pago" className="block text-sm font-medium text-gray-700 mb-1">
              Status de Pagamento
            </label>
            <select
              id="pago"
              name="pago"
              className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
              onChange={handleChange}
              value={filters.pago !== undefined ? String(filters.pago) : ''}
            >
              <option value="">Todos</option>
              <option value="false">Pago</option>
              <option value="true">Pendente</option>
            </select>
          </div>
        </div>
        
        <div className="flex justify-end">
          <button
            type="submit"
            disabled={isLoading}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
          >
            {isLoading ? 'Calculando...' : 'Calcular Total'}
          </button>
        </div>
      </form>
      
      {totalAmount && (
        <div className="mt-6 p-4 bg-green-50 rounded-lg border border-green-200">
          <div className="flex items-center">
            <DollarSign className="h-8 w-8 text-green-500 mr-3" />
            <div>
              <h3 className="text-lg font-medium text-green-800">Total de Compras</h3>
              <p className="text-2xl font-bold text-green-600">{totalAmount}</p>
              <p className="text-sm text-green-700 mt-1">
                <Check className="h-4 w-4 inline mr-1" />
                Cálculo concluído com base nos filtros selecionados
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default TotalPurchaseFilter;