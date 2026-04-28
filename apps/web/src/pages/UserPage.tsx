import { useMemo, useState } from "react";
import type { ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import { api, KycApplicationRow, MerchantTransactionRow, UserRow, WalletBalanceRow, WalletRow } from "../api/client";
import { EmptyState } from "../components/EmptyState";
import { ErrorAlert } from "../components/ErrorAlert";
import { StatusBadge } from "../components/StatusBadge";

type UserSection = "dashboard" | "profile" | "kyc" | "wallets" | "transactions";

export function UserPage({ section }: { section: UserSection }) {
  const [selectedUserId, setSelectedUserId] = useState("");
  const users = useQuery({ queryKey: ["user", "users"], queryFn: api.users });
  const wallets = useQuery({ queryKey: ["user", "wallets"], queryFn: api.wallets });
  const balances = useQuery({ queryKey: ["user", "balances"], queryFn: api.userWalletBalances });
  const kyc = useQuery({ queryKey: ["user", "kyc"], queryFn: api.kycApplications });
  const transactions = useQuery({ queryKey: ["user", "transactions"], queryFn: api.userTransactions });

  const userRows = users.data?.rows ?? [];
  const activeUserId = selectedUserId || stringValue(userRows[0]?.user_id);
  const activeUser = userRows.find((user) => sameId(user.user_id, activeUserId));
  const isLoading = [users, wallets, balances, kyc, transactions].some((query) => query.isLoading);
  const error = users.error ?? wallets.error ?? balances.error ?? kyc.error ?? transactions.error;

  const data = useMemo(
    () => ({
      wallets: byOptionalOwner(wallets.data?.rows ?? [], "user_id", activeUserId),
      balances: byOptionalOwner(balances.data?.rows ?? [], "user_id", activeUserId),
      kyc: byOptionalOwner(kyc.data?.rows ?? [], "user_id", activeUserId),
      transactions: byOptionalOwner(transactions.data?.rows ?? [], "user_id", activeUserId)
    }),
    [activeUserId, balances.data, kyc.data, transactions.data, wallets.data]
  );

  return (
    <section className="page">
      <div className="page-header">
        <div>
          <h1>{userTitle(section)}</h1>
          <p>{activeUser ? displayName(activeUser) : "User workspace"}</p>
        </div>
      </div>
      <div className="toolbar">
        <label>
          User
          <select value={activeUserId} onChange={(event) => setSelectedUserId(event.target.value)}>
            {userRows.map((user) => (
              <option key={stringValue(user.user_id)} value={stringValue(user.user_id)}>
                {displayName(user)}
              </option>
            ))}
          </select>
        </label>
        {activeUser && <StatusBadge value={activeUser.kyc_status} />}
      </div>
      <ErrorAlert error={error} />
      {isLoading && <EmptyState title="Loading user data" />}
      {!isLoading && !activeUser && <EmptyState title="No users" />}
      {!isLoading && activeUser && (
        <>
          {section === "dashboard" && <UserDashboard user={activeUser} data={data} />}
          {section === "profile" && <UserProfile user={activeUser} />}
          {section === "kyc" && <UserKyc rows={data.kyc} status={activeUser.kyc_status} />}
          {section === "wallets" && <UserWallets wallets={data.wallets} balances={data.balances} />}
          {section === "transactions" && <UserTransactions rows={data.transactions} />}
        </>
      )}
    </section>
  );
}

function UserDashboard({ user, data }: { user: UserRow; data: UserData }) {
  const totalBalance = data.balances.reduce((sum, row) => sum + numberValue(row.balance), 0);
  return (
    <>
      <div className="metric-grid">
        <Metric title="Wallets" value={String(data.wallets.length)} />
        <Metric title="Portfolio balance" value={totalBalance.toFixed(4)} />
        <Metric title="Transactions" value={String(data.transactions.length)} />
        <Metric title="KYC" value={user.kyc_status} />
      </div>
      <div className="merchant-overview">
        <UserProfile user={user} />
        <UserTransactions rows={data.transactions.slice(0, 5)} />
      </div>
    </>
  );
}

function UserProfile({ user }: { user: UserRow }) {
  return (
    <section className="panel">
      <header>
        <h2>Profile</h2>
        <StatusBadge value={user.kyc_status} />
      </header>
      <dl className="detail-list">
        <Detail label="Username" value={user.username} />
        <Detail label="Name" value={displayName(user)} />
        <Detail label="Email" value={user.email ?? "NULL"} />
        <Detail label="Created" value={formatDate(user.created_at)} />
      </dl>
    </section>
  );
}

function UserKyc({ rows, status }: { rows: KycApplicationRow[]; status: string }) {
  if (rows.length === 0) {
    return (
      <section className="panel">
        <header>
          <h2>KYC</h2>
          <StatusBadge value={status} />
        </header>
        <EmptyState title="No KYC applications" />
      </section>
    );
  }
  return <SimpleTable columns={["Application", "Status", "Submitted", "Reviewed", "Reason"]} rows={rows.map((row) => [`#${row.kyc_application_id}`, <StatusBadge value={row.status} />, formatDate(row.submitted_at), formatDate(row.reviewed_at), row.rejection_reason ?? "NULL"])} />;
}

function UserWallets({ wallets, balances }: { wallets: WalletRow[]; balances: WalletBalanceRow[] }) {
  if (wallets.length === 0) {
    return <EmptyState title="No wallets" />;
  }
  return (
    <div className="stack">
      <SimpleTable columns={["Wallet", "Label", "Address", "Chain"]} rows={wallets.map((row) => [`#${row.wallet_id}`, row.wallet_label ?? "NULL", row.wallet_address, String(row.chain_id)])} />
      <SimpleTable columns={["Wallet", "Chain", "Asset", "Balance"]} rows={balances.map((row) => [stringValue(row.wallet_id), row.chain_name ?? "chain", row.asset_symbol ?? "asset", stringValue(row.balance)])} />
    </div>
  );
}

function UserTransactions({ rows }: { rows: MerchantTransactionRow[] }) {
  if (rows.length === 0) {
    return <EmptyState title="No transactions" />;
  }
  return <SimpleTable columns={["Transaction", "Invoice", "Amount", "Status", "Created"]} rows={rows.map((row) => [`#${row.transaction_id}`, row.invoice_id ? `#${row.invoice_id}` : "NULL", formatCurrency(numberValue(row.amount_in_usdt)), <StatusBadge value={row.status} />, formatDate(row.created_at)])} />;
}

type UserData = {
  wallets: WalletRow[];
  balances: WalletBalanceRow[];
  kyc: KycApplicationRow[];
  transactions: MerchantTransactionRow[];
};

function userTitle(section: UserSection) {
  switch (section) {
    case "dashboard":
      return "User Dashboard";
    case "profile":
      return "Profile";
    case "kyc":
      return "KYC";
    case "wallets":
      return "Wallets";
    case "transactions":
      return "Transactions";
  }
}

function displayName(user: UserRow) {
  const name = [user.first_name, user.last_name].filter(Boolean).join(" ");
  return name || user.username;
}

function Detail({ label, value }: { label: string; value: ReactNode }) {
  return (
    <div>
      <dt>{label}</dt>
      <dd>{value}</dd>
    </div>
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

function byOptionalOwner<Row extends Record<string, unknown>>(rows: Row[], key: string, ownerId: string) {
  return rows.filter((row) => sameId(row[key], ownerId));
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
  return Number.isNaN(date.getTime()) ? String(value) : date.toLocaleDateString("en-US", { month: "short", day: "2-digit", year: "numeric" });
}
