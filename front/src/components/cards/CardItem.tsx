import React from 'react';
import { Card } from '../../types/card';
import { CreditCard, Calendar } from 'lucide-react';
import { format } from 'date-fns';
import { ptBR } from 'date-fns/locale';

interface CardItemProps {
  card: Card;
  onClick: () => void;
}

const CardItem: React.FC<CardItemProps> = ({ card, onClick }) => {
  // Generate consistent gradient based on card name for better UX
  const colors = [
    { gradient: 'from-blue-500 via-blue-600 to-indigo-700', accent: 'bg-blue-400/20' },
    { gradient: 'from-purple-500 via-purple-600 to-pink-700', accent: 'bg-purple-400/20' },
    { gradient: 'from-emerald-500 via-teal-600 to-cyan-700', accent: 'bg-emerald-400/20' },
    { gradient: 'from-rose-500 via-red-600 to-orange-700', accent: 'bg-rose-400/20' },
    { gradient: 'from-amber-500 via-yellow-600 to-orange-600', accent: 'bg-amber-400/20' },
    { gradient: 'from-violet-500 via-purple-600 to-fuchsia-700', accent: 'bg-violet-400/20' },
  ];
  
  // Use card name hash for consistent color per card
  const colorIndex = card.id 
    ? parseInt(card.id.slice(-1), 16) % colors.length 
    : Math.floor(Math.random() * colors.length);
  const selectedColor = colors[colorIndex];

  // Format creation date
  const formattedDate = card.data_criacao 
    ? format(new Date(card.data_criacao), "MMM yyyy", { locale: ptBR })
    : 'N/A';

  return (
    <div 
      onClick={onClick}
      className={`relative bg-gradient-to-br ${selectedColor.gradient} rounded-2xl shadow-xl hover:shadow-2xl transform hover:-translate-y-2 hover:scale-[1.02] transition-all duration-300 cursor-pointer overflow-hidden group`}
    >
      {/* Decorative elements */}
      <div className="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full -mr-16 -mt-16 group-hover:scale-150 transition-transform duration-500"></div>
      <div className="absolute bottom-0 left-0 w-24 h-24 bg-white/5 rounded-full -ml-12 -mb-12 group-hover:scale-125 transition-transform duration-500"></div>
      
      {/* Content */}
      <div className="relative p-6 h-full min-h-[200px] flex flex-col justify-between">
        <div className="flex justify-between items-start">
          <div className="flex-1">
            <div className={`inline-block px-3 py-1 rounded-full ${selectedColor.accent} backdrop-blur-sm mb-3`}>
              <p className="text-xs font-semibold text-white/90 uppercase tracking-wide">Cartão</p>
            </div>
            <h3 className="text-2xl font-bold text-white mb-1 leading-tight">
              {card.nome || 'Cartão sem nome'}
            </h3>
          </div>
          <div className="ml-4 p-3 bg-white/10 backdrop-blur-sm rounded-xl group-hover:bg-white/20 transition-colors">
            <CreditCard className="h-6 w-6 text-white" />
          </div>
        </div>
        
        <div className="mt-auto pt-4 border-t border-white/20">
          <div className="flex items-center text-white/80">
            <Calendar className="h-4 w-4 mr-2" />
            <div>
              <p className="text-xs opacity-75">Criado em</p>
              <p className="text-sm font-medium">{formattedDate}</p>
            </div>
          </div>
        </div>
      </div>
      
      {/* Hover effect overlay */}
      <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
    </div>
  );
};

export default CardItem;