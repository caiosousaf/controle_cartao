import React, { createContext, useState, useContext, ReactNode } from 'react';
import { Category } from '../types/category';
import { categoryService } from '../services/categoryService';
import { useAuth } from './AuthContext';

interface CategoryContextType {
  categories: Category[];
  isLoading: boolean;
  error: string | null;
  hasMore: boolean;
  totalCategories: number;
  fetchCategories: (limit?: number, offset?: number, ativo?: boolean) => Promise<void>;
  createCategory: (name: string) => Promise<void>;
  updateCategory: (categoryId: string, name: string) => Promise<void>;
  deactivateCategory: (categoryId: string) => Promise<void>;
  reactivateCategory: (categoryId: string) => Promise<void>;
  clearError: () => void;
}

const CategoryContext = createContext<CategoryContextType | undefined>(undefined);

export const CategoryProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [totalCategories, setTotalCategories] = useState(0);
  
  const { isAuthenticated } = useAuth();

  const fetchCategories = async (limit = 10, offset = 0, ativo?: boolean) => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await categoryService.getCategories(limit, offset, false, ativo);
      
      if (offset === 0) {
        setCategories(response.dados);
      } else {
        setCategories(prev => [...prev, ...response.dados]);
      }
      
      setHasMore(response.prox || false);
      
      // Fetch total count if it's the first page
      if (offset === 0) {
        fetchTotalCategories(ativo);
      }
    } catch (err) {
      setError('Erro ao carregar categorias');
      console.error('Error fetching categories:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchTotalCategories = async (ativo?: boolean) => {
    if (!isAuthenticated) return;
    
    try {
      const response = await categoryService.getCategories(1, 0, true, ativo);
      setTotalCategories(response.total || 0);
    } catch (err) {
      console.error('Error fetching total categories count:', err);
    }
  };

  const createCategory = async (name: string) => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      await categoryService.createCategory({ nome: name });
      // Refresh the list
      await fetchCategories();
    } catch (err) {
      setError('Erro ao criar categoria');
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  const updateCategory = async (categoryId: string, name: string) => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      await categoryService.updateCategory(categoryId, { nome: name });
      // Refresh the list
      await fetchCategories();
    } catch (err) {
      setError('Erro ao atualizar categoria');
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  const deactivateCategory = async (categoryId: string) => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      await categoryService.deactivateCategory(categoryId);
      // Refresh the list
      await fetchCategories();
    } catch (err) {
      setError('Erro ao desativar categoria');
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  const reactivateCategory = async (categoryId: string) => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      await categoryService.reactivateCategory(categoryId);
      // Refresh the list
      await fetchCategories();
    } catch (err) {
      setError('Erro ao reativar categoria');
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  const clearError = () => setError(null);

  const value = {
    categories,
    isLoading,
    error,
    hasMore,
    totalCategories,
    fetchCategories,
    createCategory,
    updateCategory,
    deactivateCategory,
    reactivateCategory,
    clearError
  };

  return <CategoryContext.Provider value={value}>{children}</CategoryContext.Provider>;
};

export const useCategories = (): CategoryContextType => {
  const context = useContext(CategoryContext);
  if (context === undefined) {
    throw new Error('useCategories must be used within a CategoryProvider');
  }
  return context;
};