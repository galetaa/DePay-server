import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { MerchantPage } from "./MerchantPage";

function jsonResponse(body: unknown) {
  return new Response(JSON.stringify(body), {
    status: 200,
    headers: { "Content-Type": "application/json" }
  });
}

function responseFor(url: string) {
  if (url.includes("stores")) {
    return {
      data: {
        columns: ["store_id", "store_name", "verification_status"],
        rows: [{ store_id: 1, store_name: "North Coffee", legal_name: "North Coffee LLC", contact_email: "merchant@example.com", verification_status: "approved" }]
      }
    };
  }
  if (url.includes("payment_invoices")) {
    return {
      data: {
        columns: ["invoice_id", "store_id", "external_order_id", "amount_usdt", "status"],
        rows: [
          { invoice_id: 10, store_id: 1, external_order_id: "order-paid", amount_usdt: "120.50", status: "paid" },
          { invoice_id: 11, store_id: 1, external_order_id: "order-open", amount_usdt: "30.00", status: "issued" }
        ]
      }
    };
  }
  if (url.includes("terminals")) {
    return { data: { columns: ["terminal_id", "store_id", "serial_number", "status"], rows: [{ terminal_id: 1, store_id: 1, serial_number: "TERM-0001", status: "active" }] } };
  }
  if (url.includes("merchant_webhooks")) {
    return { data: { columns: ["webhook_id", "store_id", "url", "is_active"], rows: [{ webhook_id: 1, store_id: 1, url: "https://example.test/hook", is_active: true, failure_count: 0 }] } };
  }
  if (url.includes("vw_webhook_delivery_status")) {
    return {
      data: {
        columns: ["webhook_delivery_id", "store_id", "event_type", "status"],
        rows: [{ webhook_delivery_id: 1, store_id: 1, event_type: "transaction.confirmed", status: "delivered", attempts: 1 }]
      }
    };
  }
  return { data: { columns: ["transaction_id", "store_id", "amount_in_usdt", "status"], rows: [{ transaction_id: 1, store_id: 1, amount_in_usdt: "120.50", status: "confirmed" }] } };
}

function renderMerchant(section: "dashboard" | "invoices" | "webhooks" | "terminals" | "analytics") {
  vi.stubGlobal("fetch", vi.fn((input: RequestInfo | URL) => Promise.resolve(jsonResponse(responseFor(String(input))))));
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false }
    }
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <MerchantPage section={section} />
    </QueryClientProvider>
  );
}

describe("MerchantPage", () => {
  it("summarizes merchant dashboard data", async () => {
    renderMerchant("dashboard");

    expect(screen.getByRole("heading", { name: "Merchant Dashboard" })).toBeInTheDocument();
    expect(await screen.findByRole("heading", { name: "North Coffee" })).toBeInTheDocument();
    expect(screen.getAllByText("$120.50").length).toBeGreaterThan(0);
    expect(screen.getAllByText("1/1").length).toBeGreaterThan(0);
  });

  it("renders invoice rows for the selected store", async () => {
    renderMerchant("invoices");

    expect(screen.getByRole("heading", { name: "Invoices" })).toBeInTheDocument();
    expect(await screen.findByText("order-paid")).toBeInTheDocument();
    expect(screen.getByText("order-open")).toBeInTheDocument();
  });
});
