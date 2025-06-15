export interface Invoice {
  id?: string;
  nome?: string;
  fatura_cartao_id?: string;
  nome_cartao?: string;
  status?: string;
  data_criacao?: string;
  data_vencimento?: string;
}

export interface InvoiceResponse {
  dados: Invoice[];
  prox?: boolean;
  total?: number;
}