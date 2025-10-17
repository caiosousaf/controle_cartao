export interface Category {
  id?: string;
  nome?: string;
  usuario_id?: string;
  data_criacao?: string;
  data_desativacao?: string;
}

export interface CategoryResponse {
  dados: Category[];
  prox?: boolean;
  total?: number;
}

export interface CreateCategoryData {
  nome: string;
}

export interface UpdateCategoryData {
  nome: string;
}