import React from 'react';
import TotalPurchaseFilter from '../components/analytics/TotalPurchaseFilter';
import { BarChart3 } from 'lucide-react';

const AnalyticsPage: React.FC = () => {
  return (
    <div className="space-y-8">
      <div className="flex items-center space-x-4">
        <BarChart3 className="h-8 w-8 text-blue-600" />
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Análise de Compras</h1>
          <p className="mt-1 text-gray-600">Visualize totais de compras com filtros personalizados</p>
        </div>
      </div>
      
      <TotalPurchaseFilter />
    </div>
  );
};

export default AnalyticsPage;