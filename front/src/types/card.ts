export interface Card {
  id?: string;
  nome?: string;
  usuario_id?: string;
  data_criacao?: string;
  data_desativacao?: string;
}

export interface CardResponse {
  dados: Card[];
  prox?: boolean;
  total?: number;
}