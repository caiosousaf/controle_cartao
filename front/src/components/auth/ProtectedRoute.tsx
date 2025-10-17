import React, { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import MainLayout from '../layout/MainLayout';

interface ProtectedRouteProps {
  children: ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const { isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    // User is not authenticated, redirect to login
    return <Navigate to="/login\" replace />;
  }

  // User is authenticated, render the children wrapped in the main layout
  return <MainLayout>{children}</MainLayout>;
};

export default ProtectedRoute;