export interface RecurringPurchase {
  id?: string;
  nome?: string;
  descricao?: string;
  compra_categoria_id?: string;
  local_compra?: string;
  valor_parcela?: number;
  ativo?: boolean;
  data_criacao?: string;
}

export interface RecurringPurchaseResponse {
  dados: RecurringPurchase[];
  prox?: boolean;
  total?: number;
}

export interface ExpenseEstimate {
  mes_ano: string;
  valor: number;
}

export interface ExpenseEstimateResponse {
  dados: ExpenseEstimate[];
}