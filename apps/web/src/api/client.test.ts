import { api } from "./client";

function jsonResponse(body: unknown, init?: ResponseInit) {
  return new Response(JSON.stringify(body), {
    status: init?.status ?? 200,
    headers: { "Content-Type": "application/json" },
    ...init
  });
}

describe("api client", () => {
  it("creates a demo invoice through the admin endpoint", async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ data: { invoice_id: "42", status: "issued" } }));
    vi.stubGlobal("fetch", fetchMock);

    await expect(api.createDemoInvoice()).resolves.toEqual({ invoice_id: "42", status: "issued" });
    expect(fetchMock).toHaveBeenCalledWith("/api/admin/demo/invoices", expect.objectContaining({ method: "POST" }));
  });

  it("submits a demo payment with invoice and wallet context", async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ data: { transaction_id: "84", status: "confirmed" } }));
    vi.stubGlobal("fetch", fetchMock);

    await expect(api.submitDemoPayment("42")).resolves.toEqual({ transaction_id: "84", status: "confirmed" });

    expect(fetchMock).toHaveBeenCalledWith(
      "/api/admin/demo/payments",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ invoice_id: "42", user_id: "1", wallet_id: "1" })
      })
    );
  });

  it("surfaces backend error messages", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse({ error: "boom" }, { status: 500 })));

    await expect(api.tables()).rejects.toThrow("boom");
  });

  it("loads merchant dashboard tables through typed table endpoints", async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ data: { columns: ["store_id"], rows: [{ store_id: 1 }] } }));
    vi.stubGlobal("fetch", fetchMock);

    await expect(api.merchantStores()).resolves.toEqual({ columns: ["store_id"], rows: [{ store_id: 1 }] });
    await api.merchantInvoices();
    await api.merchantTerminals();
    await api.merchantWebhooks();
    await api.merchantWebhookDeliveries();

    expect(fetchMock).toHaveBeenCalledWith("/api/admin/tables/stores?limit=100", expect.any(Object));
    expect(fetchMock).toHaveBeenCalledWith("/api/admin/tables/payment_invoices?limit=100", expect.any(Object));
    expect(fetchMock).toHaveBeenCalledWith("/api/admin/tables/terminals?limit=100", expect.any(Object));
    expect(fetchMock).toHaveBeenCalledWith("/api/admin/tables/merchant_webhooks?limit=100", expect.any(Object));
    expect(fetchMock).toHaveBeenCalledWith("/api/admin/tables/vw_webhook_delivery_status?limit=100", expect.any(Object));
  });

  it("loads product dashboard and system health endpoints", async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ data: { columns: ["id"], rows: [] } }));
    vi.stubGlobal("fetch", fetchMock);

    await api.users();
    await api.wallets();
    await api.userWalletBalances();
    await api.kycApplications();
    await api.complianceKycQueue();
    await api.merchantVerifications();
    await api.riskAlerts();
    await api.blacklistedWallets();

    fetchMock.mockResolvedValueOnce(jsonResponse({ data: { database: "ok", services: [] } }));
    await expect(api.systemHealth()).resolves.toEqual({ database: "ok", services: [] });

    expect(fetchMock).toHaveBeenCalledWith("/api/admin/tables/users?limit=100", expect.any(Object));
    expect(fetchMock).toHaveBeenCalledWith("/api/admin/tables/vw_user_wallet_balances?limit=100", expect.any(Object));
    expect(fetchMock).toHaveBeenCalledWith("/api/admin/tables/vw_compliance_kyc_queue?limit=100", expect.any(Object));
    expect(fetchMock).toHaveBeenCalledWith("/api/admin/tables/blacklisted_wallets?limit=100", expect.any(Object));
    expect(fetchMock).toHaveBeenCalledWith("/api/admin/system-health", expect.any(Object));
  });
});
