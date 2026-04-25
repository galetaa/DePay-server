import { useQuery } from "@tanstack/react-query";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Line,
  LineChart,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis
} from "recharts";
import { api } from "../api/client";
import { ChartCard } from "../components/ChartCard";
import { ErrorAlert } from "../components/ErrorAlert";

const palette = ["#2563eb", "#f97316", "#16a34a", "#dc2626", "#7c3aed", "#0891b2"];

export function AnalyticsPage() {
  const turnover = useQuery({ queryKey: ["analytics", "turnover"], queryFn: api.storeTurnover });
  const statuses = useQuery({ queryKey: ["analytics", "statuses"], queryFn: api.statuses });
  const failed = useQuery({ queryKey: ["analytics", "failed"], queryFn: api.failedTransactions });
  const rpc = useQuery({ queryKey: ["analytics", "rpc"], queryFn: api.rpcHealth });

  const turnoverRows = (turnover.data?.rows ?? []).map((row) => ({
    name: String(row.asset_symbol ?? row.status ?? "asset"),
    value: Number(row.total_amount_usdt ?? 0)
  }));

  const statusRows = Object.entries(
    (statuses.data?.rows ?? []).reduce<Record<string, number>>((acc, row) => {
      const status = String(row.status ?? "unknown");
      acc[status] = (acc[status] ?? 0) + 1;
      return acc;
    }, {})
  ).map(([name, value]) => ({ name, value }));

  const failedRows = (failed.data?.rows ?? []).map((row) => ({
    name: String(row.chain_name ?? "chain"),
    value: Number(row.failed_count ?? 0)
  }));

  const rpcRows = (rpc.data?.rows ?? []).map((row) => ({
    name: String(row.node_name ?? "rpc"),
    latency: Number(row.avg_latency_ms ?? 0)
  }));

  return (
    <section className="page">
      <div className="page-header">
        <div>
          <h1>Analytics</h1>
          <p>Charts</p>
        </div>
      </div>
      <ErrorAlert error={turnover.error ?? statuses.error ?? failed.error ?? rpc.error} />
      <div className="charts-grid">
        <ChartCard title="Store Turnover">
          <ResponsiveContainer width="100%" height={260}>
            <BarChart data={turnoverRows}>
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
            <PieChart>
              <Pie data={statusRows} dataKey="value" nameKey="name" outerRadius={92} label>
                {statusRows.map((_, index) => (
                  <Cell key={index} fill={palette[index % palette.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </ChartCard>
        <ChartCard title="Failed Transactions">
          <ResponsiveContainer width="100%" height={260}>
            <BarChart data={failedRows}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="value" fill="#dc2626" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </ChartCard>
        <ChartCard title="RPC Latency">
          <ResponsiveContainer width="100%" height={260}>
            <LineChart data={rpcRows}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Line type="monotone" dataKey="latency" stroke="#16a34a" strokeWidth={2} />
            </LineChart>
          </ResponsiveContainer>
        </ChartCard>
      </div>
    </section>
  );
}
