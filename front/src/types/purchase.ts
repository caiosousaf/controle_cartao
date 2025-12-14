export interface Purchase {
  id?: string;
  nome?: string;
  descricao?: string;
  local_compra?: string;
  categoria_id?: string;
  categoria_nome?: string;
  valor_parcela?: number;
  parcela_atual?: number;
  quantidade_parcelas?: number;
  fatura_id?: string;
  nome_fatura?: string;
  data_compra?: string;
  data_criacao?: string;
}

export interface PurchaseResponse {
  dados: Purchase[];
  prox?: boolean;
  total?: number;
}

export interface TotalPurchaseResponse {
  total: number;
}

export interface PurchaseFilters {
  dataEspecifica?: string;
  ultimaParcela?: boolean;
  pago?: boolean;
  categoria_id?: string;
}