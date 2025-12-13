import React, { useState, useEffect } from 'react';
import { DollarSign, Calculator, TrendingUp, TrendingDown, Edit, Save, X } from 'lucide-react';
import { ExpenseEstimate } from '../../types/recurring';

interface SalaryBalanceCardProps {
  estimates: ExpenseEstimate[];
}

const SalaryBalanceCard: React.FC<SalaryBalanceCardProps> = ({ estimates }) => {
  const [salary, setSalary] = useState<number | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [tempSalary, setTempSalary] = useState('');

  useEffect(() => {
    // Load salary from localStorage on component mount
    const savedSalary = localStorage.getItem('userSalary');
    if (savedSalary) {
      setSalary(parseFloat(savedSalary));
    }
  }, []);

  const handleSaveSalary = () => {
    const salaryValue = parseFloat(tempSalary);
    if (!isNaN(salaryValue) && salaryValue > 0) {
      setSalary(salaryValue);
      localStorage.setItem('userSalary', salaryValue.toString());
      setIsEditing(false);
      setTempSalary('');
    }
  };

  const handleCancelEdit = () => {
    setIsEditing(false);
    setTempSalary('');
  };

  const handleStartEdit = () => {
    setIsEditing(true);
    setTempSalary(salary?.toString() || '');
  };

  const calculateBalance = () => {
    if (!salary || estimates.length === 0) return null;

    // Take only the first 3 months of estimates
    const next3Months = estimates.slice(0, 3);
    const totalExpenses = next3Months.reduce((sum, estimate) => sum + estimate.valor, 0);
    const totalIncome = salary * next3Months.length;
    const balance = totalIncome - totalExpenses;

    return {
      totalIncome,
      totalExpenses,
      balance,
      months: next3Months
    };
  };

  const balanceData = calculateBalance();

  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-3">
          <Calculator className="h-6 w-6 text-blue-600" />
          <h2 className="text-xl font-bold text-gray-900">Análise Financeira</h2>
        </div>
        
        {!isEditing && (
          <button
            onClick={handleStartEdit}
            className="inline-flex items-center px-3 py-1 text-sm text-blue-600 hover:text-blue-700 transition-colors"
          >
            <Edit className="h-4 w-4 mr-1" />
            {salary ? 'Editar Salário' : 'Definir Salário'}
          </button>
        )}
      </div>

      {/* Salary Input Section */}
      <div className="mb-6 p-4 bg-gray-50 rounded-lg">
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Salário Mensal
        </label>
        
        {isEditing ? (
          <div className="flex items-center space-x-3">
            <div className="flex-1 relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <span className="text-gray-500 sm:text-sm">R$</span>
              </div>
              <input
                type="number"
                value={tempSalary}
                onChange={(e) => setTempSalary(e.target.value)}
                step="0.01"
                min="0"
                className="pl-8 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                placeholder="0,00"
                autoFocus
              />
            </div>
            <button
              onClick={handleSaveSalary}
              className="inline-flex items-center px-3 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
            >
              <Save className="h-4 w-4 mr-1" />
              Salvar
            </button>
            <button
              onClick={handleCancelEdit}
              className="inline-flex items-center px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
            >
              <X className="h-4 w-4 mr-1" />
              Cancelar
            </button>
          </div>
        ) : (
          <div className="flex items-center space-x-2">
            <DollarSign className="h-5 w-5 text-gray-400" />
            <span className="text-lg font-semibold text-gray-900">
              {salary ? `R$ ${salary.toFixed(2)}` : 'Não definido'}
            </span>
          </div>
        )}
      </div>

      {/* Balance Calculation */}
      {balanceData && (
        <div className="space-y-4">
          <h3 className="text-lg font-medium text-gray-900">
            Projeção dos Próximos {balanceData.months.length} Meses
          </h3>
          
          {/* Monthly breakdown */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
            {balanceData.months.map((month, index) => {
              const monthBalance = salary! - month.valor;
              return (
                <div key={index} className="bg-gray-50 rounded-lg p-3 border">
                  <p className="text-sm font-medium text-gray-600">{month.mes_ano}</p>
                  <div className="mt-1 space-y-1">
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">Receita:</span>
                      <span className="text-green-600">R$ {salary!.toFixed(2)}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">Gastos:</span>
                      <span className="text-red-600">R$ {month.valor.toFixed(2)}</span>
                    </div>
                    <div className="flex justify-between text-sm font-medium border-t pt-1">
                      <span>Saldo:</span>
                      <span className={monthBalance >= 0 ? 'text-green-600' : 'text-red-600'}>
                        R$ {monthBalance.toFixed(2)}
                      </span>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>

          {/* Total Summary */}
          <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg p-4 border border-blue-200">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="text-center">
                <div className="flex items-center justify-center space-x-2 mb-1">
                  <TrendingUp className="h-5 w-5 text-green-500" />
                  <span className="text-sm font-medium text-gray-600">Receita Total</span>
                </div>
                <p className="text-xl font-bold text-green-600">
                  R$ {balanceData.totalIncome.toFixed(2)}
                </p>
              </div>
              
              <div className="text-center">
                <div className="flex items-center justify-center space-x-2 mb-1">
                  <TrendingDown className="h-5 w-5 text-red-500" />
                  <span className="text-sm font-medium text-gray-600">Gastos Totais</span>
                </div>
                <p className="text-xl font-bold text-red-600">
                  R$ {balanceData.totalExpenses.toFixed(2)}
                </p>
              </div>
              
              <div className="text-center">
                <div className="flex items-center justify-center space-x-2 mb-1">
                  {balanceData.balance >= 0 ? (
                    <TrendingUp className="h-5 w-5 text-green-500" />
                  ) : (
                    <TrendingDown className="h-5 w-5 text-red-500" />
                  )}
                  <span className="text-sm font-medium text-gray-600">Saldo Final</span>
                </div>
                <p className={`text-2xl font-bold ${
                  balanceData.balance >= 0 ? 'text-green-600' : 'text-red-600'
                }`}>
                  R$ {balanceData.balance.toFixed(2)}
                </p>
                <p className="text-xs text-gray-500 mt-1">
                  {balanceData.balance >= 0 ? 'Saldo positivo' : 'Saldo negativo'}
                </p>
              </div>
            </div>
          </div>

          {/* Alert for negative balance */}
          {balanceData.balance < 0 && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <div className="flex items-center space-x-2">
                <TrendingDown className="h-5 w-5 text-red-500" />
                <span className="text-sm font-medium text-red-800">
                  Atenção: Seus gastos recorrentes excedem sua receita em R$ {Math.abs(balanceData.balance).toFixed(2)}
                </span>
              </div>
              <p className="text-xs text-red-600 mt-1">
                Considere revisar suas despesas ou aumentar sua receita.
              </p>
            </div>
          )}
        </div>
      )}

      {/* No salary message */}
      {!salary && (
        <div className="text-center py-6">
          <Calculator className="h-12 w-12 text-gray-400 mx-auto mb-3" />
          <p className="text-gray-600 mb-2">
            Defina seu salário mensal para ver a análise financeira
          </p>
          <p className="text-sm text-gray-500">
            Calcularemos automaticamente seu saldo após deduzir os gastos recorrentes
          </p>
        </div>
      )}
    </div>
  );
};

export default SalaryBalanceCard;