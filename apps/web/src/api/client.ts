export type TableRows = {
  columns: string[];
  rows: Record<string, unknown>[];
  limit?: number;
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
  createDemoInvoice: () => request<Record<string, unknown>>("/api/admin/demo/invoices", { method: "POST" }),
  submitDemoPayment: (invoiceId: string) =>
    request<Record<string, unknown>>("/api/admin/demo/payments", {
      method: "POST",
      body: JSON.stringify({ invoice_id: invoiceId, user_id: "1", wallet_id: "1" })
    })
};
