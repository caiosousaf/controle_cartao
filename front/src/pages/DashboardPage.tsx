import React from 'react';
import CardList from '../components/cards/CardList';
import { Sparkles } from 'lucide-react';

const DashboardPage: React.FC = () => {
  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <div className="relative">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-50 via-indigo-50 to-purple-50 dark:from-gray-800 dark:via-gray-800 dark:to-gray-900 rounded-3xl opacity-50"></div>
        <div className="relative px-8 py-12">
          <div className="flex items-center gap-3 mb-4">
            <div className="p-2 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl shadow-lg">
              <Sparkles className="h-6 w-6 text-white" />
            </div>
            <h1 className="text-4xl md:text-5xl font-bold bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900 dark:from-white dark:via-gray-100 dark:to-white bg-clip-text text-transparent">
              Dashboard
            </h1>
          </div>
          <p className="text-lg text-gray-600 dark:text-gray-400 ml-14">
            Gerencie seus cartões, faturas e compras de forma inteligente
          </p>
        </div>
      </div>
      
      <CardList />
    </div>
  );
};

export default DashboardPage;