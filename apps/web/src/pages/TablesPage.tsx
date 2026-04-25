import { RefreshCw } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { useMemo, useState } from "react";
import { api } from "../api/client";
import { DataTable } from "../components/DataTable";
import { ErrorAlert } from "../components/ErrorAlert";

export function TablesPage() {
  const [selectedTable, setSelectedTable] = useState("users");
  const tablesQuery = useQuery({ queryKey: ["tables"], queryFn: api.tables });
  const rowsQuery = useQuery({
    queryKey: ["table", selectedTable],
    queryFn: () => api.tableRows(selectedTable),
    enabled: selectedTable.length > 0
  });
  const tables = useMemo(() => tablesQuery.data?.tables ?? [], [tablesQuery.data]);

  return (
    <section className="page">
      <div className="page-header">
        <div>
          <h1>Tables</h1>
          <p>PostgreSQL browser</p>
        </div>
        <button type="button" className="icon-button" onClick={() => rowsQuery.refetch()} title="Refresh rows">
          <RefreshCw size={18} />
        </button>
      </div>
      <div className="toolbar">
        <label>
          Table
          <select value={selectedTable} onChange={(event) => setSelectedTable(event.target.value)}>
            {tables.map((table) => (
              <option key={table} value={table}>
                {table}
              </option>
            ))}
          </select>
        </label>
      </div>
      <ErrorAlert error={tablesQuery.error ?? rowsQuery.error} />
      <DataTable data={rowsQuery.data} />
    </section>
  );
}
