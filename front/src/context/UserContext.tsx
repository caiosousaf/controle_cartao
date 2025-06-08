import React, { createContext, useState, useContext, ReactNode } from 'react';
import { User } from '../types/user';
import { userService } from '../services/userService';
import { useAuth } from './AuthContext';

interface UserContextType {
  user: User | null;
  isLoading: boolean;
  error: string | null;
  fetchUser: () => Promise<void>;
  updateUser: (email: string, emailNovo: string, senhaAtual: string, senhaNova: string) => Promise<void>;
  clearError: () => void;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export const UserProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const { isAuthenticated } = useAuth();

  const fetchUser = async () => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const userData = await userService.getUser();
      setUser(userData);
    } catch (err) {
      setError('Erro ao carregar dados do usuário');
      console.error('Error fetching user:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const updateUser = async (email: string, emailNovo: string, senhaAtual: string, senhaNova: string) => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      await userService.updateUser({
        email,
        email_novo: emailNovo,
        senha_atual: senhaAtual,
        senha_nova: senhaNova
      });
      
      // Update local user data with new email
      if (user) {
        setUser({ ...user, email: emailNovo });
      }
    } catch (err) {
      setError('Erro ao atualizar dados do usuário');
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  const clearError = () => setError(null);

  const value = {
    user,
    isLoading,
    error,
    fetchUser,
    updateUser,
    clearError
  };

  return <UserContext.Provider value={value}>{children}</UserContext.Provider>;
};

export const useUser = (): UserContextType => {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
};