import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { App } from "./App";

function jsonResponse(body: unknown) {
  return new Response(JSON.stringify(body), {
    status: 200,
    headers: { "Content-Type": "application/json" }
  });
}

function merchantTableResponse(url: string) {
  if (url.includes("system-health")) {
    return { data: { database: "ok", services: [{ name: "admin-service", url: "local", status: "ok" }] } };
  }
  if (url.includes("vw_compliance_kyc_queue")) {
    return { data: { columns: ["user_id", "username", "kyc_status"], rows: [{ user_id: 2, username: "boris", kyc_status: "pending" }] } };
  }
  if (url.includes("merchant_verification_applications")) {
    return { data: { columns: ["merchant_verification_application_id", "store_id", "status"], rows: [{ merchant_verification_application_id: 3, store_id: 1, status: "approved" }] } };
  }
  if (url.includes("risk_alerts")) {
    return { data: { columns: ["risk_alert_id", "alert_type", "risk_level"], rows: [{ risk_alert_id: 8, alert_type: "large_payment", risk_level: "high" }] } };
  }
  if (url.includes("blacklisted_wallets")) {
    return { data: { columns: ["blacklisted_wallet_id", "chain_id", "wallet_address"], rows: [{ blacklisted_wallet_id: 9, chain_id: 1, wallet_address: "0xblocked" }] } };
  }
  if (url.includes("users")) {
    return {
      data: {
        columns: ["user_id", "username", "first_name", "last_name", "kyc_status"],
        rows: [{ user_id: 1, username: "alice", first_name: "Alice", last_name: "Stone", kyc_status: "approved" }]
      }
    };
  }
  if (url.includes("vw_user_wallet_balances")) {
    return { data: { columns: ["user_id", "wallet_id", "asset_symbol", "balance"], rows: [{ user_id: 1, wallet_id: 1, asset_symbol: "ETH", balance: "10.5" }] } };
  }
  if (url.includes("kyc_applications")) {
    return { data: { columns: ["kyc_application_id", "user_id", "status"], rows: [{ kyc_application_id: 2, user_id: 1, status: "approved" }] } };
  }
  if (url.includes("wallets")) {
    return { data: { columns: ["wallet_id", "user_id", "wallet_address"], rows: [{ wallet_id: 1, user_id: 1, wallet_address: "0xabc", chain_id: 1 }] } };
  }
  if (url.includes("stores")) {
    return { data: { columns: ["store_id", "store_name"], rows: [{ store_id: 1, store_name: "North Coffee", verification_status: "approved" }] } };
  }
  if (url.includes("payment_invoices")) {
    return {
      data: {
        columns: ["invoice_id", "store_id", "external_order_id", "amount_usdt", "status"],
        rows: [{ invoice_id: 11, store_id: 1, external_order_id: "order-11", amount_usdt: "120.50", status: "paid" }]
      }
    };
  }
  if (url.includes("terminals")) {
    return { data: { columns: ["terminal_id", "store_id", "serial_number", "status"], rows: [{ terminal_id: 3, store_id: 1, serial_number: "TERM-0001", status: "active" }] } };
  }
  if (url.includes("merchant_webhooks")) {
    return { data: { columns: ["webhook_id", "store_id", "url", "is_active"], rows: [{ webhook_id: 4, store_id: 1, url: "https://example.test/hook", is_active: true }] } };
  }
  if (url.includes("vw_webhook_delivery_status")) {
    return { data: { columns: ["webhook_delivery_id", "store_id", "event_type", "status"], rows: [{ webhook_delivery_id: 5, store_id: 1, event_type: "transaction.confirmed", status: "delivered" }] } };
  }
  if (url.includes("payment_transactions")) {
    return { data: { columns: ["transaction_id", "store_id", "amount_in_usdt", "status"], rows: [{ transaction_id: 6, store_id: 1, amount_in_usdt: "120.50", status: "confirmed" }] } };
  }
  return { data: { tables: ["users"], columns: [], rows: [] } };
}

function renderApp(path: string, responseForUrl?: (url: string) => unknown) {
  window.history.pushState({}, "", path);
  vi.stubGlobal(
    "fetch",
    vi.fn((input: RequestInfo | URL) => Promise.resolve(jsonResponse(responseForUrl ? responseForUrl(String(input)) : { data: { tables: ["users"], columns: [], rows: [] } })))
  );

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false }
    }
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  );
}

describe("App routing", () => {
  it("renders the page matching the current URL", () => {
    renderApp("/admin/functions");

    expect(screen.getByRole("heading", { name: "Functions" })).toBeInTheDocument();
  });

  it("updates the URL when navigating from the sidebar", async () => {
    renderApp("/admin/functions");

    await userEvent.click(screen.getByRole("button", { name: "Demo" }));

    await waitFor(() => {
      expect(window.location.pathname).toBe("/admin/demo");
      expect(screen.getByRole("heading", { name: "Demo" })).toBeInTheDocument();
    });
  });

  it("renders merchant routes from the URL", async () => {
    renderApp("/merchant/invoices", merchantTableResponse);

    expect(screen.getByRole("heading", { name: "Invoices" })).toBeInTheDocument();
    expect(await screen.findByText("order-11")).toBeInTheDocument();
  });

  it("lets a demo persona enter the user dashboard", async () => {
    renderApp("/login", merchantTableResponse);

    await userEvent.click(screen.getByRole("radio", { name: /Alice Stone/i }));
    await userEvent.click(screen.getByRole("button", { name: "Continue" }));

    await waitFor(() => {
      expect(window.location.pathname).toBe("/user/dashboard");
      expect(screen.getByRole("heading", { name: "User Dashboard" })).toBeInTheDocument();
    });
    expect(await screen.findByText("10.5000")).toBeInTheDocument();
  });

  it("renders compliance product routes", async () => {
    renderApp("/compliance/kyc", merchantTableResponse);

    expect(screen.getByRole("heading", { name: "KYC Queue" })).toBeInTheDocument();
    expect(await screen.findByText("boris")).toBeInTheDocument();
  });

  it("renders admin system health", async () => {
    renderApp("/admin/system-health", merchantTableResponse);

    expect(screen.getByRole("heading", { name: "System Health" })).toBeInTheDocument();
    expect(await screen.findByText("admin-service")).toBeInTheDocument();
  });
});
