import { useMemo, useState } from "react";
import type { ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import { Bar, BarChart, CartesianGrid, Cell, Pie, PieChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { api, MerchantInvoiceRow, MerchantTerminalRow, MerchantTransactionRow, MerchantWebhookDeliveryRow, MerchantWebhookRow, StoreRow } from "../api/client";
import { ChartCard } from "../components/ChartCard";
import { EmptyState } from "../components/EmptyState";
import { ErrorAlert } from "../components/ErrorAlert";
import { StatusBadge } from "../components/StatusBadge";

type MerchantSection = "dashboard" | "invoices" | "webhooks" | "terminals" | "analytics";

type MerchantPageProps = {
  section: MerchantSection;
};

const palette = ["#2563eb", "#f97316", "#16a34a", "#dc2626", "#0891b2", "#7c3aed"];

export function MerchantPage({ section }: MerchantPageProps) {
  const [selectedStoreId, setSelectedStoreId] = useState("");
  const stores = useQuery({ queryKey: ["merchant", "stores"], queryFn: api.merchantStores });
  const invoices = useQuery({ queryKey: ["merchant", "invoices"], queryFn: api.merchantInvoices });
  const terminals = useQuery({ queryKey: ["merchant", "terminals"], queryFn: api.merchantTerminals });
  const webhooks = useQuery({ queryKey: ["merchant", "webhooks"], queryFn: api.merchantWebhooks });
  const deliveries = useQuery({ queryKey: ["merchant", "webhook-deliveries"], queryFn: api.merchantWebhookDeliveries });
  const transactions = useQuery({ queryKey: ["merchant", "transactions"], queryFn: api.merchantTransactions });

  const storeRows = stores.data?.rows ?? [];
  const activeStoreId = selectedStoreId || stringValue(storeRows[0]?.store_id);
  const activeStore = storeRows.find((store) => sameId(store.store_id, activeStoreId));
  const isLoading = [stores, invoices, terminals, webhooks, deliveries, transactions].some((query) => query.isLoading);
  const error = stores.error ?? invoices.error ?? terminals.error ?? webhooks.error ?? deliveries.error ?? transactions.error;

  const data = useMemo(
    () => ({
      invoices: byStore(invoices.data?.rows ?? [], activeStoreId),
      terminals: byStore(terminals.data?.rows ?? [], activeStoreId),
      webhooks: byStore(webhooks.data?.rows ?? [], activeStoreId),
      deliveries: byStore(deliveries.data?.rows ?? [], activeStoreId),
      transactions: byStore(transactions.data?.rows ?? [], activeStoreId)
    }),
    [activeStoreId, deliveries.data, invoices.data, terminals.data, transactions.data, webhooks.data]
  );

  const summary = useMemo(() => summarizeMerchant(data), [data]);
  const heading = sectionTitle(section);

  return (
    <section className="page">
      <div className="page-header">
        <div>
          <h1>{heading}</h1>
          <p>{activeStore?.store_name ?? "Merchant workspace"}</p>
        </div>
      </div>

      <div className="toolbar">
        <label>
          Store
          <select value={activeStoreId} onChange={(event) => setSelectedStoreId(event.target.value)}>
            {storeRows.map((store) => (
              <option key={stringValue(store.store_id)} value={stringValue(store.store_id)}>
                {store.store_name}
              </option>
            ))}
          </select>
        </label>
        {activeStore && <StatusBadge value={String(activeStore.verification_status)} />}
      </div>

      <ErrorAlert error={error} />
      {isLoading && <EmptyState title="Loading merchant data" />}
      {!isLoading && !activeStore && <EmptyState title="No merchant stores" />}
      {!isLoading && activeStore && (
        <>
          {section === "dashboard" && <DashboardSection data={data} summary={summary} store={activeStore} />}
          {section === "invoices" && <InvoicesSection rows={data.invoices} />}
          {section === "webhooks" && <WebhooksSection rows={data.webhooks} deliveries={data.deliveries} />}
          {section === "terminals" && <TerminalsSection rows={data.terminals} />}
          {section === "analytics" && <MerchantAnalyticsSection data={data} summary={summary} />}
        </>
      )}
    </section>
  );
}

function DashboardSection({ data, summary, store }: { data: MerchantData; summary: MerchantSummary; store: StoreRow }) {
  return (
    <>
      <div className="metric-grid">
        <Metric title="Paid revenue" value={formatCurrency(summary.paidRevenue)} />
        <Metric title="Open invoices" value={String(summary.openInvoices)} />
        <Metric title="Active terminals" value={`${summary.activeTerminals}/${data.terminals.length}`} />
        <Metric title="Webhook success" value={`${summary.deliveredWebhooks}/${data.deliveries.length}`} />
      </div>
      <div className="merchant-overview">
        <section className="panel">
          <header>
            <h2>{store.store_name}</h2>
            <StatusBadge value={String(store.verification_status)} />
          </header>
          <dl className="detail-list">
            <div>
              <dt>Legal name</dt>
              <dd>{store.legal_name ?? "NULL"}</dd>
            </div>
            <div>
              <dt>Contact</dt>
              <dd>{store.contact_email ?? "NULL"}</dd>
            </div>
            <div>
              <dt>Invoices</dt>
              <dd>{data.invoices.length}</dd>
            </div>
          </dl>
        </section>
        <RecentInvoices rows={data.invoices.slice(0, 5)} />
      </div>
    </>
  );
}

function InvoicesSection({ rows }: { rows: MerchantInvoiceRow[] }) {
  if (rows.length === 0) {
    return <EmptyState title="No invoices for this store" />;
  }
  return (
    <SimpleTable
      columns={["Invoice", "Order", "Amount", "Status", "Created", "Paid"]}
      rows={rows.map((row) => [
        `#${row.invoice_id}`,
        row.external_order_id ?? "NULL",
        formatCurrency(numberValue(row.amount_usdt)),
        <StatusBadge value={row.status} />,
        formatDate(row.created_at),
        formatDate(row.paid_at)
      ])}
    />
  );
}

function WebhooksSection({ rows, deliveries }: { rows: MerchantWebhookRow[]; deliveries: MerchantWebhookDeliveryRow[] }) {
  if (rows.length === 0) {
    return <EmptyState title="No webhooks for this store" />;
  }
  return (
    <div className="stack">
      <SimpleTable
        columns={["Webhook", "URL", "Status", "Failures", "Last success", "Last failure"]}
        rows={rows.map((row) => [
          `#${row.webhook_id}`,
          row.url,
          truthy(row.is_active) ? "active" : "inactive",
          String(row.failure_count ?? 0),
          formatDate(row.last_success_at),
          formatDate(row.last_failure_at)
        ])}
      />
      <section className="panel">
        <header>
          <h2>Deliveries</h2>
        </header>
        {deliveries.length === 0 ? (
          <EmptyState title="No delivery attempts" />
        ) : (
          <SimpleTable
            columns={["Delivery", "Event", "Status", "Attempts", "Response", "Next attempt", "Payload"]}
            rows={deliveries.map((row) => [
              `#${row.webhook_delivery_id}`,
              row.event_type,
              <StatusBadge value={row.status} />,
              String(row.attempts ?? 0),
              row.response_status ? String(row.response_status) : "NULL",
              formatDate(row.next_attempt_at),
              payloadPreview(row.payload)
            ])}
          />
        )}
      </section>
    </div>
  );
}

function TerminalsSection({ rows }: { rows: MerchantTerminalRow[] }) {
  if (rows.length === 0) {
    return <EmptyState title="No terminals for this store" />;
  }
  return (
    <SimpleTable
      columns={["Terminal", "Serial", "Status", "Last seen", "Created"]}
      rows={rows.map((row) => [`#${row.terminal_id}`, row.serial_number, <StatusBadge value={row.status} />, formatDate(row.last_seen_at), formatDate(row.created_at)])}
    />
  );
}

function MerchantAnalyticsSection({ data, summary }: { data: MerchantData; summary: MerchantSummary }) {
  const invoiceStatusRows = groupCount(data.invoices, (row) => row.status);
  const webhookStatusRows = groupCount(data.deliveries, (row) => row.status);
  const transactionStatusRows = groupCount(data.transactions, (row) => row.status);
  const invoiceAmountRows = data.invoices.map((row) => ({
    name: `#${row.invoice_id}`,
    value: numberValue(row.amount_usdt)
  }));

  return (
    <>
      <div className="metric-grid">
        <Metric title="Total invoiced" value={formatCurrency(summary.totalInvoiced)} />
        <Metric title="Paid revenue" value={formatCurrency(summary.paidRevenue)} />
        <Metric title="Transactions" value={String(data.transactions.length)} />
        <Metric title="Webhook failures" value={String(summary.failedWebhooks)} />
      </div>
      <div className="charts-grid">
        <ChartCard title="Invoice Statuses">
          <ResponsiveContainer width="100%" height={260}>
            <PieChart>
              <Pie data={invoiceStatusRows} dataKey="value" nameKey="name" outerRadius={92} label>
                {invoiceStatusRows.map((_, index) => (
                  <Cell key={index} fill={palette[index % palette.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </ChartCard>
        <ChartCard title="Invoice Amounts">
          <ResponsiveContainer width="100%" height={260}>
            <BarChart data={invoiceAmountRows}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="value" fill="#2563eb" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </ChartCard>
        <ChartCard title="Transaction Statuses">
          <ResponsiveContainer width="100%" height={260}>
            <BarChart data={transactionStatusRows}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="value" fill="#16a34a" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </ChartCard>
        <ChartCard title="Webhook Deliveries">
          <ResponsiveContainer width="100%" height={260}>
            <BarChart data={webhookStatusRows}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="value" fill="#f97316" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </ChartCard>
      </div>
    </>
  );
}

function RecentInvoices({ rows }: { rows: MerchantInvoiceRow[] }) {
  return (
    <section className="panel">
      <header>
        <h2>Recent invoices</h2>
      </header>
      {rows.length === 0 ? (
        <EmptyState title="No recent invoices" />
      ) : (
        <SimpleTable
          columns={["Invoice", "Amount", "Status"]}
          rows={rows.map((row) => [`#${row.invoice_id}`, formatCurrency(numberValue(row.amount_usdt)), <StatusBadge value={row.status} />])}
        />
      )}
    </section>
  );
}

function Metric({ title, value }: { title: string; value: string }) {
  return (
    <section className="metric-card">
      <span>{title}</span>
      <strong>{value}</strong>
    </section>
  );
}

function SimpleTable({ columns, rows }: { columns: string[]; rows: ReactNode[][] }) {
  return (
    <div className="table-wrap">
      <table>
        <thead>
          <tr>
            {columns.map((column) => (
              <th key={column}>{column}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, rowIndex) => (
            <tr key={rowIndex}>
              {row.map((cell, cellIndex) => (
                <td key={cellIndex}>{cell}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

type MerchantData = {
  invoices: MerchantInvoiceRow[];
  terminals: MerchantTerminalRow[];
  webhooks: MerchantWebhookRow[];
  deliveries: MerchantWebhookDeliveryRow[];
  transactions: MerchantTransactionRow[];
};

type MerchantSummary = {
  totalInvoiced: number;
  paidRevenue: number;
  openInvoices: number;
  activeTerminals: number;
  deliveredWebhooks: number;
  failedWebhooks: number;
};

function summarizeMerchant(data: MerchantData): MerchantSummary {
  return {
    totalInvoiced: data.invoices.reduce((sum, invoice) => sum + numberValue(invoice.amount_usdt), 0),
    paidRevenue: data.invoices.filter((invoice) => invoice.status === "paid").reduce((sum, invoice) => sum + numberValue(invoice.amount_usdt), 0),
    openInvoices: data.invoices.filter((invoice) => invoice.status === "issued" || invoice.status === "draft").length,
    activeTerminals: data.terminals.filter((terminal) => terminal.status === "active").length,
    deliveredWebhooks: data.deliveries.filter((delivery) => delivery.status === "delivered").length,
    failedWebhooks: data.deliveries.filter((delivery) => ["failed", "retry_scheduled", "dead_letter"].includes(delivery.status)).length
  };
}

function sectionTitle(section: MerchantSection) {
  switch (section) {
    case "dashboard":
      return "Merchant Dashboard";
    case "invoices":
      return "Invoices";
    case "webhooks":
      return "Webhooks";
    case "terminals":
      return "Terminals";
    case "analytics":
      return "Merchant Analytics";
  }
}

function byStore<Row extends { store_id: string | number | null | undefined }>(rows: Row[], storeId: string) {
  return rows.filter((row) => sameId(row.store_id, storeId));
}

function sameId(left: unknown, right: unknown) {
  return stringValue(left) === stringValue(right);
}

function stringValue(value: unknown) {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

function numberValue(value: unknown) {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : 0;
}

function formatCurrency(value: number) {
  return `$${value.toFixed(2)}`;
}

function formatDate(value: unknown) {
  if (!value) {
    return "NULL";
  }
  const date = new Date(String(value));
  if (Number.isNaN(date.getTime())) {
    return String(value);
  }
  return date.toLocaleDateString("en-US", { month: "short", day: "2-digit", year: "numeric" });
}

function payloadPreview(value: unknown) {
  if (!value) {
    return "NULL";
  }
  const payload = typeof value === "string" ? value : JSON.stringify(value);
  return payload.length > 80 ? `${payload.slice(0, 77)}...` : payload;
}

function truthy(value: unknown) {
  return value === true || value === "true" || value === "t" || value === 1 || value === "1";
}

function groupCount<Row>(rows: Row[], getKey: (row: Row) => string) {
  return Object.entries(
    rows.reduce<Record<string, number>>((acc, row) => {
      const key = getKey(row) || "unknown";
      acc[key] = (acc[key] ?? 0) + 1;
      return acc;
    }, {})
  ).map(([name, value]) => ({ name, value }));
}
