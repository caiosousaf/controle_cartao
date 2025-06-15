import React from 'react';
import CategoryList from '../components/categories/CategoryList';
import { Tag } from 'lucide-react';

const CategoriesPage: React.FC = () => {
  return (
    <div className="space-y-8">
      <div className="flex items-center space-x-4">
        <Tag className="h-8 w-8 text-blue-600" />
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Categorias</h1>
          <p className="mt-1 text-gray-600">Gerencie as categorias das suas compras</p>
        </div>
      </div>
      
      <CategoryList />
    </div>
  );
};

export default CategoriesPage;