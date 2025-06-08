import React, { createContext, useState, useContext, ReactNode, useEffect } from 'react';
import { Card, CardResponse } from '../types/card';
import { cardService } from '../services/cardService';
import { useAuth } from './AuthContext';

interface CardContextType {
  cards: Card[];
  selectedCard: Card | null;
  isLoading: boolean;
  error: string | null;
  hasMore: boolean;
  totalCards: number;
  fetchCards: (limit?: number, offset?: number) => Promise<void>;
  selectCard: (card: Card) => void;
  clearSelectedCard: () => void;
  clearError: () => void;
}

const CardContext = createContext<CardContextType | undefined>(undefined);

export const CardProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [cards, setCards] = useState<Card[]>([]);
  const [selectedCard, setSelectedCard] = useState<Card | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [totalCards, setTotalCards] = useState(0);
  
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    if (isAuthenticated) {
      fetchCards();
      fetchTotalCards();
    }
  }, [isAuthenticated]);

  const fetchCards = async (limit = 10, offset = 0) => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await cardService.getCards(limit, offset);
      
      if (offset === 0) {
        setCards(response.dados);
      } else {
        setCards(prev => [...prev, ...response.dados]);
      }
      
      setHasMore(response.prox || false);
    } catch (err) {
      setError('Erro ao carregar cartões');
      console.error('Error fetching cards:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchTotalCards = async () => {
    if (!isAuthenticated) return;
    
    try {
      const response = await cardService.getCards(1, 0, true);
      setTotalCards(response.total || 0);
    } catch (err) {
      console.error('Error fetching total cards count:', err);
    }
  };

  const selectCard = (card: Card) => {
    setSelectedCard(card);
  };

  const clearSelectedCard = () => {
    setSelectedCard(null);
  };

  const clearError = () => setError(null);

  const value = {
    cards,
    selectedCard,
    isLoading,
    error,
    hasMore,
    totalCards,
    fetchCards,
    selectCard,
    clearSelectedCard,
    clearError
  };

  return <CardContext.Provider value={value}>{children}</CardContext.Provider>;
};

export const useCards = (): CardContextType => {
  const context = useContext(CardContext);
  if (context === undefined) {
    throw new Error('useCards must be used within a CardProvider');
  }
  return context;
};