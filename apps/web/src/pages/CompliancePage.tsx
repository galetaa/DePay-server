import type { ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import { api, BlacklistedWalletRow, KycQueueRow, MerchantVerificationRow, RiskAlertRow } from "../api/client";
import { EmptyState } from "../components/EmptyState";
import { ErrorAlert } from "../components/ErrorAlert";
import { StatusBadge } from "../components/StatusBadge";

type ComplianceSection = "kyc" | "merchant-verifications" | "risk-alerts" | "blacklist";

export function CompliancePage({ section }: { section: ComplianceSection }) {
  const kyc = useQuery({ queryKey: ["compliance", "kyc"], queryFn: api.complianceKycQueue });
  const merchants = useQuery({ queryKey: ["compliance", "merchant-verifications"], queryFn: api.merchantVerifications });
  const risks = useQuery({ queryKey: ["compliance", "risk-alerts"], queryFn: api.riskAlerts });
  const blacklist = useQuery({ queryKey: ["compliance", "blacklist"], queryFn: api.blacklistedWallets });
  const error = kyc.error ?? merchants.error ?? risks.error ?? blacklist.error;
  const isLoading = [kyc, merchants, risks, blacklist].some((query) => query.isLoading);

  return (
    <section className="page">
      <div className="page-header">
        <div>
          <h1>{title(section)}</h1>
          <p>Compliance workspace</p>
        </div>
      </div>
      <ErrorAlert error={error} />
      {isLoading && <EmptyState title="Loading compliance data" />}
      {!isLoading && (
        <>
          {section === "kyc" && <KycQueue rows={kyc.data?.rows ?? []} />}
          {section === "merchant-verifications" && <MerchantVerifications rows={merchants.data?.rows ?? []} />}
          {section === "risk-alerts" && <RiskAlerts rows={risks.data?.rows ?? []} />}
          {section === "blacklist" && <Blacklist rows={blacklist.data?.rows ?? []} />}
        </>
      )}
    </section>
  );
}

function KycQueue({ rows }: { rows: KycQueueRow[] }) {
  if (rows.length === 0) {
    return <EmptyState title="No KYC queue items" />;
  }
  return <SimpleTable columns={["User", "Email", "Status", "Submitted", "Documents"]} rows={rows.map((row) => [row.username ?? stringValue(row.user_id), row.email ?? "NULL", <StatusBadge value={row.kyc_status ?? "pending"} />, formatDate(row.submitted_at), stringValue(row.documents_count ?? 0)])} />;
}

function MerchantVerifications({ rows }: { rows: MerchantVerificationRow[] }) {
  if (rows.length === 0) {
    return <EmptyState title="No merchant verification applications" />;
  }
  return <SimpleTable columns={["Application", "Store", "Status", "Submitted", "Reviewed", "Reason"]} rows={rows.map((row) => [`#${row.merchant_verification_application_id}`, `#${row.store_id}`, <StatusBadge value={row.status} />, formatDate(row.submitted_at), formatDate(row.reviewed_at), row.rejection_reason ?? "NULL"])} />;
}

function RiskAlerts({ rows }: { rows: RiskAlertRow[] }) {
  if (rows.length === 0) {
    return <EmptyState title="No risk alerts" />;
  }
  return <SimpleTable columns={["Alert", "Type", "Level", "Transaction", "Created"]} rows={rows.map((row) => [`#${row.risk_alert_id}`, row.alert_type, <StatusBadge value={row.risk_level} />, row.transaction_id ? `#${row.transaction_id}` : "NULL", formatDate(row.created_at)])} />;
}

function Blacklist({ rows }: { rows: BlacklistedWalletRow[] }) {
  if (rows.length === 0) {
    return <EmptyState title="No blacklisted wallets" />;
  }
  return <SimpleTable columns={["Wallet", "Chain", "Reason", "Created"]} rows={rows.map((row) => [row.wallet_address, stringValue(row.chain_id), row.reason ?? "NULL", formatDate(row.created_at)])} />;
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

function title(section: ComplianceSection) {
  switch (section) {
    case "kyc":
      return "KYC Queue";
    case "merchant-verifications":
      return "Merchant Verifications";
    case "risk-alerts":
      return "Risk Alerts";
    case "blacklist":
      return "Blacklist";
  }
}

function stringValue(value: unknown) {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

function formatDate(value: unknown) {
  if (!value) {
    return "NULL";
  }
  const date = new Date(String(value));
  return Number.isNaN(date.getTime()) ? String(value) : date.toLocaleDateString("en-US", { month: "short", day: "2-digit", year: "numeric" });
}
