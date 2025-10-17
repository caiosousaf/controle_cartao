import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { CardProvider } from './context/CardContext';
import { InvoiceProvider } from './context/InvoiceContext';
import { PurchaseProvider } from './context/PurchaseContext';
import { CategoryProvider } from './context/CategoryContext';
import { UserProvider } from './context/UserContext';
import ProtectedRoute from './components/auth/ProtectedRoute';
import LoginPage from './pages/LoginPage';
import DashboardPage from './pages/DashboardPage';
import CardDetailsPage from './pages/CardDetailsPage';
import InvoiceDetailsPage from './pages/InvoiceDetailsPage';
import AnalyticsPage from './pages/AnalyticsPage';
import RecurringPurchasesPage from './pages/RecurringPurchasesPage';
import CategoriesPage from './pages/CategoriesPage';
import UserProfilePage from './pages/UserProfilePage';

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <UserProvider>
          <CardProvider>
            <InvoiceProvider>
              <PurchaseProvider>
                <CategoryProvider>
                  <Routes>
                    <Route path="/login" element={<LoginPage />} />
                    <Route path="/" element={<ProtectedRoute><DashboardPage /></ProtectedRoute>} />
                    <Route path="/cards/:cardId" element={<ProtectedRoute><CardDetailsPage /></ProtectedRoute>} />
                    <Route path="/invoices/:invoiceId" element={<ProtectedRoute><InvoiceDetailsPage /></ProtectedRoute>} />
                    <Route path="/analytics" element={<ProtectedRoute><AnalyticsPage /></ProtectedRoute>} />
                    <Route path="/recurring" element={<ProtectedRoute><RecurringPurchasesPage /></ProtectedRoute>} />
                    <Route path="/categories" element={<ProtectedRoute><CategoriesPage /></ProtectedRoute>} />
                    <Route path="/profile" element={<ProtectedRoute><UserProfilePage /></ProtectedRoute>} />
                    <Route path="*" element={<Navigate to="/\" replace />} />
                  </Routes>
                </CategoryProvider>
              </PurchaseProvider>
            </InvoiceProvider>
          </CardProvider>
        </UserProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;