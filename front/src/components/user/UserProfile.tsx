import React, { useEffect, useState } from 'react';
import { useUser } from '../../context/UserContext';
import { User, Calendar, Mail, Edit, Save, X } from 'lucide-react';

const UserProfile: React.FC = () => {
  const { user, isLoading, error, fetchUser, updateUser, clearError } = useUser();
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({
    emailNovo: '',
    senhaAtual: '',
    senhaNova: '',
    confirmarSenha: ''
  });
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    fetchUser();
  }, []);

  useEffect(() => {
    if (user?.email) {
      setFormData(prev => ({
        ...prev,
        emailNovo: user.email || ''
      }));
    }
  }, [user]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    
    // Clear specific field error when user starts typing
    if (formErrors[name]) {
      setFormErrors(prev => ({ ...prev, [name]: '' }));
    }
  };

  const validateForm = () => {
    const errors: Record<string, string> = {};
    
    if (!formData.emailNovo.trim()) {
      errors.emailNovo = 'Email é obrigatório';
    } else if (!/\S+@\S+\.\S+/.test(formData.emailNovo)) {
      errors.emailNovo = 'Email inválido';
    }
    
    if (!formData.senhaAtual) {
      errors.senhaAtual = 'Senha atual é obrigatória';
    }
    
    if (!formData.senhaNova) {
      errors.senhaNova = 'Nova senha é obrigatória';
    } else if (formData.senhaNova.length < 6) {
      errors.senhaNova = 'Nova senha deve ter pelo menos 6 caracteres';
    }
    
    if (formData.senhaNova !== formData.confirmarSenha) {
      errors.confirmarSenha = 'Senhas não coincidem';
    }
    
    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm() || !user?.email) return;
    
    setIsSubmitting(true);
    clearError();
    
    try {
      await updateUser(
        user.email,
        formData.emailNovo,
        formData.senhaAtual,
        formData.senhaNova
      );
      
      setIsEditing(false);
      setFormData({
        emailNovo: formData.emailNovo,
        senhaAtual: '',
        senhaNova: '',
        confirmarSenha: ''
      });
    } catch (err) {
      console.error('Error updating user:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCancel = () => {
    setIsEditing(false);
    setFormData({
      emailNovo: user?.email || '',
      senhaAtual: '',
      senhaNova: '',
      confirmarSenha: ''
    });
    setFormErrors({});
    clearError();
  };

  if (isLoading && !user) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/3 mb-6"></div>
          <div className="space-y-4">
            <div className="h-4 bg-gray-200 rounded w-1/4"></div>
            <div className="h-10 bg-gray-200 rounded"></div>
            <div className="h-4 bg-gray-200 rounded w-1/4"></div>
            <div className="h-10 bg-gray-200 rounded"></div>
          </div>
        </div>
      </div>
    );
  }

  if (error && !user) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <div className="text-center p-6 bg-red-50 rounded-lg">
          <p className="text-red-600">{error}</p>
          <button
            onClick={fetchUser}
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
          >
            Tentar novamente
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-3">
          <User className="h-8 w-8 text-blue-600" />
          <h2 className="text-2xl font-bold text-gray-900">Perfil do Usuário</h2>
        </div>
        {!isEditing && (
          <button
            onClick={() => setIsEditing(true)}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            <Edit className="h-4 w-4 mr-2" />
            Editar Dados
          </button>
        )}
      </div>

      {error && (
        <div className="mb-4 p-4 text-sm text-red-700 bg-red-100 rounded-lg">
          {error}
        </div>
      )}

      {!isEditing ? (
        <div className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Nome
            </label>
            <p className="text-lg text-gray-900">{user?.nome || 'Não informado'}</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <div className="flex items-center space-x-2">
              <Mail className="h-5 w-5 text-gray-400" />
              <p className="text-lg text-gray-900">{user?.email || 'Não informado'}</p>
            </div>
          </div>

          {user?.data_criacao && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Membro desde
              </label>
              <div className="flex items-center space-x-2">
                <Calendar className="h-5 w-5 text-gray-400" />
                <p className="text-lg text-gray-900">
                  {new Date(user.data_criacao).toLocaleDateString('pt-BR', {
                    day: '2-digit',
                    month: 'long',
                    year: 'numeric',
                    timeZone: 'UTC'
                  })}
                </p>
              </div>
            </div>
          )}
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="emailAtual" className="block text-sm font-medium text-gray-700">
              Email Atual
            </label>
            <input
              type="email"
              id="emailAtual"
              value={user?.email || ''}
              disabled
              className="mt-1 block w-full rounded-md border-gray-300 bg-gray-100 shadow-sm cursor-not-allowed"
            />
          </div>

          <div>
            <label htmlFor="emailNovo" className="block text-sm font-medium text-gray-700">
              Novo Email *
            </label>
            <input
              type="email"
              id="emailNovo"
              name="emailNovo"
              value={formData.emailNovo}
              onChange={handleInputChange}
              className={`mt-1 block w-full rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 ${
                formErrors.emailNovo ? 'border-red-500' : 'border-gray-300'
              }`}
              required
            />
            {formErrors.emailNovo && (
              <p className="mt-1 text-sm text-red-600">{formErrors.emailNovo}</p>
            )}
          </div>

          <div>
            <label htmlFor="senhaAtual" className="block text-sm font-medium text-gray-700">
              Senha Atual *
            </label>
            <input
              type="password"
              id="senhaAtual"
              name="senhaAtual"
              value={formData.senhaAtual}
              onChange={handleInputChange}
              className={`mt-1 block w-full rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 ${
                formErrors.senhaAtual ? 'border-red-500' : 'border-gray-300'
              }`}
              required
            />
            {formErrors.senhaAtual && (
              <p className="mt-1 text-sm text-red-600">{formErrors.senhaAtual}</p>
            )}
          </div>

          <div>
            <label htmlFor="senhaNova" className="block text-sm font-medium text-gray-700">
              Nova Senha *
            </label>
            <input
              type="password"
              id="senhaNova"
              name="senhaNova"
              value={formData.senhaNova}
              onChange={handleInputChange}
              className={`mt-1 block w-full rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 ${
                formErrors.senhaNova ? 'border-red-500' : 'border-gray-300'
              }`}
              required
            />
            {formErrors.senhaNova && (
              <p className="mt-1 text-sm text-red-600">{formErrors.senhaNova}</p>
            )}
          </div>

          <div>
            <label htmlFor="confirmarSenha" className="block text-sm font-medium text-gray-700">
              Confirmar Nova Senha *
            </label>
            <input
              type="password"
              id="confirmarSenha"
              name="confirmarSenha"
              value={formData.confirmarSenha}
              onChange={handleInputChange}
              className={`mt-1 block w-full rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 ${
                formErrors.confirmarSenha ? 'border-red-500' : 'border-gray-300'
              }`}
              required
            />
            {formErrors.confirmarSenha && (
              <p className="mt-1 text-sm text-red-600">{formErrors.confirmarSenha}</p>
            )}
          </div>

          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={handleCancel}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
              disabled={isSubmitting}
            >
              <X className="h-4 w-4 mr-2 inline" />
              Cancelar
            </button>
            <button
              type="submit"
              className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50"
              disabled={isSubmitting}
            >
              <Save className="h-4 w-4 mr-2 inline" />
              {isSubmitting ? 'Salvando...' : 'Salvar Alterações'}
            </button>
          </div>
        </form>
      )}
    </div>
  );
};

export default UserProfile;