import React from 'react';
import UserProfile from '../components/user/UserProfile';
import { User } from 'lucide-react';

const UserProfilePage: React.FC = () => {
  return (
    <div className="space-y-8">
      <div className="flex items-center space-x-4">
        <User className="h-8 w-8 text-blue-600" />
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Perfil do Usuário</h1>
          <p className="mt-1 text-gray-600">Visualize e edite suas informações pessoais</p>
        </div>
      </div>
      
      <UserProfile />
    </div>
  );
};

export default UserProfilePage;