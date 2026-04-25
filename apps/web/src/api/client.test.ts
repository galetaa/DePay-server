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
});
