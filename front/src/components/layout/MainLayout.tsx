import React, { ReactNode, useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { 
  CreditCard, 
  LogOut, 
  Menu, 
  X, 
  BarChart3, 
  User, 
  Home,
  RepeatIcon,
  Moon,
  Sun,
  Tag
} from 'lucide-react';

interface MainLayoutProps {
  children: ReactNode;
}

const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const { logout } = useAuth();
  const navigate = useNavigate();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [isDark, setIsDark] = useState(() => {
    if (typeof window !== 'undefined') {
      return document.documentElement.classList.contains('dark');
    }
    return false;
  });

  useEffect(() => {
    if (isDark) {
      document.documentElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
    } else {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
    }
  }, [isDark]);

  useEffect(() => {
    const savedTheme = localStorage.getItem('theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    
    if (savedTheme === 'dark' || (!savedTheme && prefersDark)) {
      setIsDark(true);
    }
  }, []);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  const toggleTheme = () => {
    setIsDark(!isDark);
  };

  return (
    <div className="min-h-screen bg-gray-100 dark:bg-gray-900 flex flex-col">
      {/* Header */}
      <header className="bg-white dark:bg-gray-800 shadow-sm sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <div className="flex">
              <button
                type="button"
                className="p-2 rounded-md text-gray-400 hover:text-gray-500 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500 md:hidden"
                onClick={toggleSidebar}
              >
                <span className="sr-only">Open sidebar</span>
                <Menu className="h-6 w-6" />
              </button>
              <div className="flex items-center">
                <CreditCard className="h-8 w-8 text-blue-600" />
                <span className="ml-2 text-xl font-bold text-gray-900 dark:text-white">CardManager</span>
              </div>
            </div>
            
            <div className="flex items-center space-x-4">
              <button
                onClick={toggleTheme}
                className="p-2 text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white rounded-lg"
              >
                {isDark ? (
                  <Sun className="h-5 w-5" />
                ) : (
                  <Moon className="h-5 w-5" />
                )}
              </button>
              
              {/* Profile link in header for desktop */}
              <div className="hidden md:flex md:items-center md:space-x-4">
                <Link
                  to="/profile"
                  className="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-gray-500 dark:text-gray-400 bg-white dark:bg-gray-800 hover:text-gray-700 dark:hover:text-gray-200 focus:outline-none transition"
                >
                  <User className="h-4 w-4 mr-1" />
                  Perfil
                </Link>
                <button
                  type="button"
                  onClick={handleLogout}
                  className="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-gray-500 dark:text-gray-400 bg-white dark:bg-gray-800 hover:text-gray-700 dark:hover:text-gray-200 focus:outline-none transition"
                >
                  <LogOut className="h-4 w-4 mr-1" />
                  Sair
                </button>
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Mobile sidebar backdrop */}
      <div 
        className={`fixed inset-0 bg-gray-600 bg-opacity-75 z-20 transition-opacity duration-300 ease-in-out ${
          sidebarOpen ? 'opacity-100' : 'opacity-0 pointer-events-none'
        }`}
        onClick={toggleSidebar}
      />

      {/* Sidebar */}
      <div 
        className={`fixed inset-y-0 left-0 w-64 bg-white dark:bg-gray-800 shadow-lg z-30 transform transition duration-300 ease-in-out ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        } md:translate-x-0 md:relative md:z-0`}
      >
        <div className="h-16 flex items-center justify-between px-4 md:hidden">
          <div className="flex items-center">
            <CreditCard className="h-8 w-8 text-blue-600" />
            <span className="ml-2 text-xl font-bold text-gray-900 dark:text-white">CardManager</span>
          </div>
          <button
            type="button"
            onClick={toggleSidebar}
            className="p-2 rounded-md text-gray-400 hover:text-gray-500 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
          >
            <span className="sr-only">Close sidebar</span>
            <X className="h-6 w-6" />
          </button>
        </div>
        
        <nav className="mt-4 px-2 space-y-1">
          {/* Profile link at the top for mobile */}
          <div className="md:hidden mb-4 pb-4 border-b border-gray-200 dark:border-gray-700">
            <Link
              to="/profile"
              className="group flex items-center px-2 py-2 text-base font-medium rounded-md text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-white"
              onClick={() => setSidebarOpen(false)}
            >
              <User className="mr-3 h-6 w-6 text-gray-500 dark:text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300" />
              Perfil
            </Link>
          </div>

          <Link
            to="/"
            className="group flex items-center px-2 py-2 text-base font-medium rounded-md text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-white"
            onClick={() => setSidebarOpen(false)}
          >
            <Home className="mr-3 h-6 w-6 text-gray-500 dark:text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300" />
            Início
          </Link>
          <Link
            to="/categories"
            className="group flex items-center px-2 py-2 text-base font-medium rounded-md text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-white"
            onClick={() => setSidebarOpen(false)}
          >
            <Tag className="mr-3 h-6 w-6 text-gray-500 dark:text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300" />
            Categorias
          </Link>
          <Link
            to="/analytics"
            className="group flex items-center px-2 py-2 text-base font-medium rounded-md text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-white"
            onClick={() => setSidebarOpen(false)}
          >
            <BarChart3 className="mr-3 h-6 w-6 text-gray-500 dark:text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300" />
            Análise de Compras
          </Link>
          <Link
            to="/recurring"
            className="group flex items-center px-2 py-2 text-base font-medium rounded-md text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-white"
            onClick={() => setSidebarOpen(false)}
          >
            <RepeatIcon className="mr-3 h-6 w-6 text-gray-500 dark:text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300" />
            Compras Recorrentes
          </Link>
          
          <div className="md:hidden pt-4 pb-3 border-t border-gray-200 dark:border-gray-700">
            <div className="flex items-center px-4">
              <div className="flex-shrink-0">
                <User className="h-10 w-10 rounded-full text-gray-400 dark:text-gray-500 bg-gray-100 dark:bg-gray-700 p-2" />
              </div>
              <div className="ml-3">
                <div className="text-base font-medium text-gray-800 dark:text-gray-200">Usuário</div>
                <button
                  type="button"
                  onClick={handleLogout}
                  className="mt-1 inline-flex items-center text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300"
                >
                  <LogOut className="h-4 w-4 mr-1" />
                  Sair
                </button>
              </div>
            </div>
          </div>
        </nav>
      </div>

      {/* Main content */}
      <main className="flex-1 py-6 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto w-full">
        {children}
      </main>
    </div>
  );
};

export default MainLayout;