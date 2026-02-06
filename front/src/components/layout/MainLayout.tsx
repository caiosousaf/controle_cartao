import React, { ReactNode, useState, useEffect } from 'react';
import { useNavigate, Link, useLocation } from 'react-router-dom';
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
  Tag,
  ChevronLeft,
  ChevronRight
} from 'lucide-react';

interface MainLayoutProps {
  children: ReactNode;
}

const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const { logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(() => {
    const saved = localStorage.getItem('sidebarCollapsed');
    return saved === 'true';
  });
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

  const toggleSidebarCollapse = () => {
    const newState = !sidebarCollapsed;
    setSidebarCollapsed(newState);
    localStorage.setItem('sidebarCollapsed', newState.toString());
  };

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/';
    }
    return location.pathname.startsWith(path);
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex">
      {/* Sidebar - Fixed on desktop */}
      <div className={`hidden lg:flex lg:flex-shrink-0 transition-all duration-300 ${sidebarCollapsed ? 'w-20' : 'w-72'}`}>
        <div className="flex flex-col w-full relative">
          {/* Sidebar Header */}
          <div className="flex items-center h-16 flex-shrink-0 px-4 border-b border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 relative">
            {!sidebarCollapsed && (
              <div className="p-2 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl shadow-lg">
                <CreditCard className="h-6 w-6 text-white" />
              </div>
            )}
            {!sidebarCollapsed ? (
              <span className="ml-3 text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 dark:from-white dark:to-gray-300 bg-clip-text text-transparent">
                CardManager
              </span>
            ) : (
              <div className="p-2 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl shadow-lg mx-auto">
                <CreditCard className="h-6 w-6 text-white" />
              </div>
            )}
            
            {/* Collapse Toggle Button - Now in Header */}
            <button
              onClick={toggleSidebarCollapse}
              className={`absolute ${sidebarCollapsed ? 'right-2' : 'right-3'} p-2 rounded-lg text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-700 dark:hover:text-gray-300 transition-all duration-200 z-10`}
              title={sidebarCollapsed ? 'Expandir menu' : 'Minimizar menu'}
            >
              {sidebarCollapsed ? (
                <ChevronRight className="h-5 w-5" />
              ) : (
                <ChevronLeft className="h-5 w-5" />
              )}
            </button>
          </div>
          
          {/* Sidebar Navigation */}
          <div className="flex-1 flex flex-col overflow-y-auto bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800">
            <nav className="flex-1 px-3 py-6 space-y-2">
              <NavLink
                to="/"
                icon={Home}
                label="Início"
                isActive={isActive('/')}
                onClick={() => setSidebarOpen(false)}
                gradient="from-blue-500 to-indigo-600"
                collapsed={sidebarCollapsed}
              />
              <NavLink
                to="/categories"
                icon={Tag}
                label="Categorias"
                isActive={isActive('/categories')}
                onClick={() => setSidebarOpen(false)}
                gradient="from-purple-500 to-pink-600"
                collapsed={sidebarCollapsed}
              />
              <NavLink
                to="/analytics"
                icon={BarChart3}
                label="Análise de Compras"
                isActive={isActive('/analytics')}
                onClick={() => setSidebarOpen(false)}
                gradient="from-emerald-500 to-teal-600"
                collapsed={sidebarCollapsed}
              />
              <NavLink
                to="/recurring"
                icon={RepeatIcon}
                label="Compras Recorrentes"
                isActive={isActive('/recurring')}
                onClick={() => setSidebarOpen(false)}
                gradient="from-amber-500 to-orange-600"
                collapsed={sidebarCollapsed}
              />
            </nav>
            
            {/* User Section */}
            {!sidebarCollapsed && (
              <div className="flex-shrink-0 p-4 border-t border-gray-200 dark:border-gray-800">
                <div className="flex items-center px-4 py-3 rounded-xl bg-gray-50 dark:bg-gray-800/50">
                  <div className="flex-shrink-0 p-2 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-lg">
                    <User className="h-5 w-5 text-white" />
                  </div>
                  <div className="ml-3 flex-1">
                    <div className="text-sm font-semibold text-gray-800 dark:text-gray-200">Usuário</div>
                    <Link
                      to="/profile"
                      className="text-xs text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400"
                    >
                      Ver perfil
                    </Link>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Mobile sidebar backdrop */}
      <div 
        className={`fixed inset-0 bg-gray-900/50 backdrop-blur-sm z-40 transition-opacity duration-300 ease-in-out lg:hidden ${
          sidebarOpen ? 'opacity-100' : 'opacity-0 pointer-events-none'
        }`}
        onClick={toggleSidebar}
      />

      {/* Mobile Sidebar */}
      <div 
        className={`fixed inset-y-0 left-0 w-72 bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800 shadow-2xl z-50 transform transition duration-300 ease-in-out lg:hidden ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        <div className="h-16 flex items-center justify-between px-6 border-b border-gray-200 dark:border-gray-800">
          <div className="flex items-center">
            <div className="p-2 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl">
              <CreditCard className="h-5 w-5 text-white" />
            </div>
            <span className="ml-3 text-lg font-bold bg-gradient-to-r from-gray-900 to-gray-700 dark:from-white dark:to-gray-300 bg-clip-text text-transparent">
              CardManager
            </span>
          </div>
          <button
            type="button"
            onClick={toggleSidebar}
            className="p-2 rounded-xl text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500 transition-all"
          >
            <span className="sr-only">Close sidebar</span>
            <X className="h-6 w-6" />
          </button>
        </div>
        
        <nav className="flex-1 px-4 py-6 space-y-2">
          <NavLink
            to="/"
            icon={Home}
            label="Início"
            isActive={isActive('/')}
            onClick={() => setSidebarOpen(false)}
            gradient="from-blue-500 to-indigo-600"
          />
          <NavLink
            to="/categories"
            icon={Tag}
            label="Categorias"
            isActive={isActive('/categories')}
            onClick={() => setSidebarOpen(false)}
            gradient="from-purple-500 to-pink-600"
          />
          <NavLink
            to="/analytics"
            icon={BarChart3}
            label="Análise de Compras"
            isActive={isActive('/analytics')}
            onClick={() => setSidebarOpen(false)}
            gradient="from-emerald-500 to-teal-600"
          />
          <NavLink
            to="/recurring"
            icon={RepeatIcon}
            label="Compras Recorrentes"
            isActive={isActive('/recurring')}
            onClick={() => setSidebarOpen(false)}
            gradient="from-amber-500 to-orange-600"
          />
        </nav>
        
        <div className="p-4 border-t border-gray-200 dark:border-gray-800">
          <div className="flex items-center px-4 py-3 rounded-xl bg-gray-50 dark:bg-gray-800/50">
            <div className="flex-shrink-0 p-2 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-lg">
              <User className="h-5 w-5 text-white" />
            </div>
            <div className="ml-3 flex-1">
              <div className="text-sm font-semibold text-gray-800 dark:text-gray-200">Usuário</div>
              <Link
                to="/profile"
                className="text-xs text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400"
                onClick={() => setSidebarOpen(false)}
              >
                Ver perfil
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Top Header */}
        <header className="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 shadow-sm sticky top-0 z-10">
          <div className="px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between h-16 items-center">
              <div className="flex items-center">
                <button
                  type="button"
                  className="p-2 rounded-xl text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500 transition-all lg:hidden"
                  onClick={toggleSidebar}
                >
                  <span className="sr-only">Open sidebar</span>
                  <Menu className="h-6 w-6" />
                </button>
              </div>
              
              <div className="flex items-center space-x-3">
                <button
                  onClick={toggleTheme}
                  className="p-2.5 text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-800 rounded-xl transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  {isDark ? (
                    <Sun className="h-5 w-5" />
                  ) : (
                    <Moon className="h-5 w-5" />
                  )}
                </button>
                
                <div className="hidden md:flex md:items-center md:space-x-3">
                  <Link
                    to="/profile"
                    className="inline-flex items-center px-4 py-2 text-sm font-medium rounded-xl text-gray-700 dark:text-gray-300 bg-gray-50 dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none transition-all duration-200 hover:shadow-md"
                  >
                    <User className="h-4 w-4 mr-2" />
                    Perfil
                  </Link>
                  <button
                    type="button"
                    onClick={handleLogout}
                    className="inline-flex items-center px-4 py-2 text-sm font-medium rounded-xl text-white bg-gradient-to-r from-red-500 to-rose-600 hover:from-red-600 hover:to-rose-700 focus:outline-none transition-all duration-200 shadow-md hover:shadow-lg transform hover:scale-105"
                  >
                    <LogOut className="h-4 w-4 mr-2" />
                    Sair
                  </button>
                </div>
              </div>
            </div>
          </div>
        </header>

        {/* Main content */}
        <main className="flex-1 overflow-y-auto bg-gray-50 dark:bg-gray-900">
          <div className="py-8 px-4 sm:px-6 lg:px-8">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
};

// NavLink Component
interface NavLinkProps {
  to: string;
  icon: React.ElementType;
  label: string;
  isActive: boolean;
  onClick: () => void;
  gradient: string;
  collapsed?: boolean;
}

const NavLink: React.FC<NavLinkProps> = ({ to, icon: Icon, label, isActive, onClick, gradient, collapsed = false }) => {
  return (
    <Link
      to={to}
      onClick={onClick}
      className={`group flex items-center ${collapsed ? 'justify-center px-3' : 'px-4'} py-3 text-sm font-medium rounded-xl transition-all duration-200 relative ${
        isActive
          ? `bg-gradient-to-r ${gradient} text-white shadow-lg shadow-blue-500/25`
          : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800'
      }`}
      title={collapsed ? label : undefined}
    >
      <Icon className={`${collapsed ? '' : 'mr-3'} h-5 w-5 transition-colors flex-shrink-0 ${
        isActive ? 'text-white' : 'text-gray-500 dark:text-gray-400 group-hover:text-gray-700 dark:group-hover:text-gray-200'
      }`} />
      {!collapsed && (
        <span className="truncate">{label}</span>
      )}
      {collapsed && !isActive && (
        <div className="absolute left-full ml-2 px-3 py-2 bg-gray-900 dark:bg-gray-800 text-white text-sm rounded-lg opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity whitespace-nowrap z-50 shadow-lg">
          {label}
          <div className="absolute right-full top-1/2 -translate-y-1/2 border-4 border-transparent border-r-gray-900 dark:border-r-gray-800"></div>
        </div>
      )}
    </Link>
  );
};

export default MainLayout;