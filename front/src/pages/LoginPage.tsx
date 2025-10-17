import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import LoginForm from '../components/auth/LoginForm';
import { useAuth } from '../context/AuthContext';
import { CreditCard } from 'lucide-react';

const LoginPage: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  
  useEffect(() => {
    if (isAuthenticated) {
      navigate('/');
    }
  }, [isAuthenticated, navigate]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-500 to-indigo-600 dark:from-blue-900 dark:to-indigo-900 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <div className="flex justify-center">
          <CreditCard className="mx-auto h-16 w-16 text-white" />
        </div>
        <h2 className="mt-6 text-center text-3xl font-extrabold text-white">
          Sistema de Gestão de Cartões
        </h2>
        <p className="mt-2 text-center text-white text-opacity-80">
          Faça login para gerenciar seus cartões, faturas e compras
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <LoginForm />
      </div>
    </div>
  );
};

export default LoginPage;