import React from 'react';
import { TrendingUp } from 'lucide-react';
import { ExpenseEstimate } from '../../types/recurring';

interface EstimateCardProps {
  estimates: ExpenseEstimate[];
}

const EstimateCard: React.FC<EstimateCardProps> = ({ estimates }) => {
  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <div className="flex items-center space-x-3 mb-6">
        <TrendingUp className="h-6 w-6 text-blue-600" />
        <h2 className="text-xl font-bold text-gray-900">Previsão de Gastos</h2>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {estimates.map((estimate, index) => (
          <div
            key={index}
            className="bg-gray-50 rounded-lg p-4 border border-gray-200"
          >
            <p className="text-sm text-gray-600">{estimate.mes_ano}</p>
            <p className="text-2xl font-bold text-gray-900 mt-1">
              R$ {estimate.valor.toFixed(2)}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default EstimateCard