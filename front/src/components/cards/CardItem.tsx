import React from 'react';
import { Card } from '../../types/card';
import { CreditCard } from 'lucide-react';
import { format } from 'date-fns';
import { ptBR } from 'date-fns/locale';

interface CardItemProps {
  card: Card;
  onClick: () => void;
}

const CardItem: React.FC<CardItemProps> = ({ card, onClick }) => {
  // Generate a random gradient for the card background
  const colors = [
    'from-blue-500 to-indigo-600',
    'from-purple-500 to-pink-600',
    'from-green-500 to-teal-600',
    'from-red-500 to-orange-600',
    'from-yellow-500 to-amber-600',
  ];
  const randomColor = colors[Math.floor(Math.random() * colors.length)];

  // Format creation date
  const formattedDate = card.data_criacao 
    ? format(new Date(card.data_criacao), "dd 'de' MMMM 'de' yyyy", { locale: ptBR })
    : 'Data desconhecida';

  return (
    <div 
      onClick={onClick}
      className={`bg-gradient-to-r ${randomColor} rounded-xl shadow-lg p-6 text-white hover:shadow-xl transform hover:-translate-y-1 transition-all cursor-pointer h-48 flex flex-col justify-between`}
    >
      <div className="flex justify-between items-start">
        <h3 className="text-xl font-bold truncate">{card.nome || 'Cartão sem nome'}</h3>
        <CreditCard className="h-8 w-8" />
      </div>
      
      <div className="mt-auto">
        <p className="text-xs opacity-75">Desde</p>
        <p className="text-sm">{formattedDate}</p>
      </div>
    </div>
  );
};

export default CardItem;