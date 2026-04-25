import { Play } from "lucide-react";
import { FormEvent, useMemo, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { api } from "../api/client";
import { DataTable } from "../components/DataTable";
import { ErrorAlert } from "../components/ErrorAlert";

const functions = [
  { name: "get_user_kyc_wallet_summary", params: [] },
  { name: "get_user_wallet_balances", params: ["1"] },
  { name: "get_wallet_asset_distribution", params: ["1"] },
  { name: "get_transaction_card", params: ["1"] },
  { name: "get_user_transaction_history", params: ["1", "2020-01-01", "2100-01-01"] },
  { name: "get_store_transaction_history", params: ["1", "2020-01-01", "2100-01-01"] },
  { name: "get_blockchain_asset_activity", params: ["1", "2020-01-01", "2100-01-01"] },
  { name: "get_rpc_nodes_activity", params: ["1", "2020-01-01", "2100-01-01"] },
  { name: "get_store_turnover", params: ["1", "2020-01-01", "2100-01-01"] },
  { name: "get_store_success_rate", params: ["1", "2020-01-01", "2100-01-01"] },
  { name: "get_unverified_active_users", params: ["1", "1", "2020-01-01", "2100-01-01"] },
  { name: "get_failed_transactions_analytics", params: ["2020-01-01", "2100-01-01"] }
];

export function FunctionsPage() {
  const [selectedName, setSelectedName] = useState(functions[0].name);
  const selected = useMemo(() => functions.find((item) => item.name === selectedName) ?? functions[0], [selectedName]);
  const [params, setParams] = useState<string[]>(selected.params);
  const mutation = useMutation({
    mutationFn: () => api.executeFunction(selectedName, params)
  });

  function onSelect(name: string) {
    const next = functions.find((item) => item.name === name) ?? functions[0];
    setSelectedName(next.name);
    setParams(next.params);
  }

  function onSubmit(event: FormEvent) {
    event.preventDefault();
    mutation.mutate();
  }

  return (
    <section className="page">
      <div className="page-header">
        <div>
          <h1>Functions</h1>
          <p>SQL runner</p>
        </div>
      </div>
      <form className="toolbar" onSubmit={onSubmit}>
        <label>
          Function
          <select value={selectedName} onChange={(event) => onSelect(event.target.value)}>
            {functions.map((item) => (
              <option key={item.name} value={item.name}>
                {item.name}
              </option>
            ))}
          </select>
        </label>
        {params.map((param, index) => (
          <label key={index}>
            Param {index + 1}
            <input
              value={param}
              onChange={(event) => {
                const next = [...params];
                next[index] = event.target.value;
                setParams(next);
              }}
            />
          </label>
        ))}
        <button type="submit" className="primary-button" disabled={mutation.isPending}>
          <Play size={16} />
          Run
        </button>
      </form>
      <ErrorAlert error={mutation.error} />
      <DataTable data={mutation.data} emptyLabel="Run a function" />
    </section>
  );
}
