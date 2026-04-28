import { useQuery } from "@tanstack/react-query";
import { api } from "../api/client";
import { EmptyState } from "../components/EmptyState";
import { ErrorAlert } from "../components/ErrorAlert";
import { StatusBadge } from "../components/StatusBadge";

export function AdminSystemHealthPage() {
  const health = useQuery({ queryKey: ["admin", "system-health"], queryFn: api.systemHealth, refetchInterval: 30000 });
  const services = health.data?.services ?? [];
  const healthyCount = services.filter((service) => service.status === "ok").length;

  return (
    <section className="page">
      <div className="page-header">
        <div>
          <h1>System Health</h1>
          <p>Backend readiness</p>
        </div>
      </div>
      <ErrorAlert error={health.error} />
      {health.isLoading && <EmptyState title="Loading system health" />}
      {health.data && (
        <>
          <div className="metric-grid">
            <section className="metric-card">
              <span>Database</span>
              <strong>{health.data.database}</strong>
            </section>
            <section className="metric-card">
              <span>Services</span>
              <strong>
                {healthyCount}/{services.length}
              </strong>
            </section>
          </div>
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Service</th>
                  <th>Status</th>
                  <th>Endpoint</th>
                  <th>Error</th>
                </tr>
              </thead>
              <tbody>
                {services.map((service) => (
                  <tr key={service.name}>
                    <td>{service.name}</td>
                    <td>
                      <StatusBadge value={service.status} />
                    </td>
                    <td>{service.url}</td>
                    <td>{service.error ?? "NULL"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </>
      )}
    </section>
  );
}
