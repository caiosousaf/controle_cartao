export interface User {
  nome?: string;
  email?: string;
  data_criacao?: string;
}

export interface UserResponse {
  nome?: string;
  email?: string;
  data_criacao?: string;
}

export interface UpdateUserData {
  email: string;
  email_novo: string;
  senha_atual: string;
  senha_nova: string;
}