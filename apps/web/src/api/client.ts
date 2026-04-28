export type TableRows<Row = Record<string, unknown>> = {
  columns: string[];
  rows: Row[];
  limit?: number;
};

export type StoreRow = {
  store_id: string | number;
  store_name: string;
  legal_name?: string;
  contact_email?: string;
  verification_status: string;
  created_at?: string;
};

export type UserRow = {
  user_id: string | number;
  username: string;
  email?: string | null;
  first_name?: string | null;
  last_name?: string | null;
  kyc_status: string;
  created_at?: string;
};

export type WalletRow = {
  wallet_id: string | number;
  user_id?: string | number | null;
  store_id?: string | number | null;
  chain_id: string | number;
  wallet_address: string;
  is_store_wallet: boolean | string;
  wallet_label?: string | null;
  created_at?: string;
};

export type WalletBalanceRow = {
  user_id?: string | number | null;
  wallet_id?: string | number;
  wallet_label?: string | null;
  wallet_address?: string;
  chain_name?: string;
  asset_symbol?: string;
  balance?: string | number;
};

export type KycApplicationRow = {
  kyc_application_id: string | number;
  user_id: string | number;
  status: string;
  reviewer_id?: string | number | null;
  rejection_reason?: string | null;
  submitted_at?: string;
  reviewed_at?: string | null;
};

export type KycQueueRow = {
  user_id?: string | number;
  username?: string;
  email?: string | null;
  kyc_status?: string;
  submitted_at?: string;
  documents_count?: string | number;
};

export type MerchantVerificationRow = {
  merchant_verification_application_id: string | number;
  store_id: string | number;
  status: string;
  reviewer_id?: string | number | null;
  rejection_reason?: string | null;
  submitted_at?: string;
  reviewed_at?: string | null;
};

export type RiskAlertRow = {
  risk_alert_id: string | number;
  user_id?: string | number | null;
  store_id?: string | number | null;
  transaction_id?: string | number | null;
  alert_type: string;
  risk_level: string;
  details?: unknown;
  created_at?: string;
};

export type BlacklistedWalletRow = {
  blacklisted_wallet_id: string | number;
  chain_id: string | number;
  wallet_address: string;
  reason?: string | null;
  created_by?: string | number | null;
  created_at?: string;
};

export type MerchantInvoiceRow = {
  invoice_id: string | number;
  store_id: string | number;
  user_id?: string | number | null;
  external_order_id?: string | null;
  amount_usdt: string | number;
  status: string;
  expires_at?: string;
  paid_at?: string | null;
  created_at?: string;
};

export type MerchantTerminalRow = {
  terminal_id: string | number;
  store_id: string | number;
  serial_number: string;
  status: string;
  last_seen_at?: string | null;
  created_at?: string;
};

export type MerchantWebhookRow = {
  webhook_id: string | number;
  store_id: string | number;
  url: string;
  event_types?: string[] | string;
  is_active: boolean | string;
  failure_count: string | number;
  last_success_at?: string | null;
  last_failure_at?: string | null;
  created_at?: string;
};

export type MerchantWebhookDeliveryRow = {
  webhook_delivery_id: string | number;
  webhook_id?: string | number | null;
  store_id: string | number;
  store_name?: string;
  transaction_id?: string | number | null;
  event_type: string;
  payload?: unknown;
  status: string;
  attempts: string | number;
  response_status?: string | number | null;
  error_message?: string | null;
  created_at?: string;
  last_attempt_at?: string | null;
  next_attempt_at?: string | null;
  delivered_at?: string | null;
};

export type MerchantTransactionRow = {
  transaction_id: string | number;
  store_id: string | number;
  invoice_id?: string | number | null;
  amount_in_usdt: string | number;
  status: string;
  created_at?: string;
  completed_at?: string | null;
};

export type ServiceHealth = {
  name: string;
  url: string;
  status: string;
  error?: string;
};

export type SystemHealth = {
  database: string;
  services: ServiceHealth[];
};

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {})
    },
    ...init
  });

  const payload = await response.json().catch(() => ({}));
  if (!response.ok) {
    const message = payload?.error?.message ?? payload?.error ?? "Request failed";
    throw new Error(message);
  }
  return payload.data as T;
}

export const api = {
  tables: () => request<{ tables: string[] }>("/api/admin/tables"),
  tableRows: (table: string, limit = 50) => request<TableRows>(`/api/admin/tables/${table}?limit=${limit}`),
  executeFunction: (name: string, params: string[]) =>
    request<TableRows>(`/api/admin/functions/${name}/execute`, {
      method: "POST",
      body: JSON.stringify({ params })
    }),
  storeTurnover: () => request<TableRows>("/api/analytics/store-turnover?store_id=1&date_from=2020-01-01&date_to=2100-01-01"),
  statuses: () => request<TableRows>("/api/analytics/transaction-statuses"),
  failedTransactions: () => request<TableRows>("/api/analytics/failed-transactions?date_from=2020-01-01&date_to=2100-01-01"),
  rpcHealth: () => request<TableRows>("/api/analytics/rpc-health"),
  users: () => request<TableRows<UserRow>>("/api/admin/tables/users?limit=100"),
  wallets: () => request<TableRows<WalletRow>>("/api/admin/tables/wallets?limit=100"),
  userWalletBalances: () => request<TableRows<WalletBalanceRow>>("/api/admin/tables/vw_user_wallet_balances?limit=100"),
  kycApplications: () => request<TableRows<KycApplicationRow>>("/api/admin/tables/kyc_applications?limit=100"),
  userTransactions: () => request<TableRows<MerchantTransactionRow>>("/api/admin/tables/payment_transactions?limit=100"),
  merchantStores: () => request<TableRows<StoreRow>>("/api/admin/tables/stores?limit=100"),
  merchantInvoices: () => request<TableRows<MerchantInvoiceRow>>("/api/admin/tables/payment_invoices?limit=100"),
  merchantTerminals: () => request<TableRows<MerchantTerminalRow>>("/api/admin/tables/terminals?limit=100"),
  merchantWebhooks: () => request<TableRows<MerchantWebhookRow>>("/api/admin/tables/merchant_webhooks?limit=100"),
  merchantWebhookDeliveries: () => request<TableRows<MerchantWebhookDeliveryRow>>("/api/admin/tables/vw_webhook_delivery_status?limit=100"),
  merchantTransactions: () => request<TableRows<MerchantTransactionRow>>("/api/admin/tables/payment_transactions?limit=100"),
  complianceKycQueue: () => request<TableRows<KycQueueRow>>("/api/admin/tables/vw_compliance_kyc_queue?limit=100"),
  merchantVerifications: () => request<TableRows<MerchantVerificationRow>>("/api/admin/tables/merchant_verification_applications?limit=100"),
  riskAlerts: () => request<TableRows<RiskAlertRow>>("/api/admin/tables/risk_alerts?limit=100"),
  blacklistedWallets: () => request<TableRows<BlacklistedWalletRow>>("/api/admin/tables/blacklisted_wallets?limit=100"),
  systemHealth: () => request<SystemHealth>("/api/admin/system-health"),
  createDemoInvoice: () => request<Record<string, unknown>>("/api/admin/demo/invoices", { method: "POST" }),
  submitDemoPayment: (invoiceId: string) =>
    request<Record<string, unknown>>("/api/admin/demo/payments", {
      method: "POST",
      body: JSON.stringify({ invoice_id: invoiceId, user_id: "1", wallet_id: "1" })
    })
};
