import React from 'react';
import CardList from '../components/cards/CardList';

const DashboardPage: React.FC = () => {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="mt-1 text-gray-600">Gerencie seus cartões, faturas e compras</p>
      </div>
      
      <CardList />
    </div>
  );
};

export default DashboardPage;