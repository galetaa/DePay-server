import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { DemoPage } from "./DemoPage";

function jsonResponse(body: unknown) {
  return new Response(JSON.stringify(body), {
    status: 200,
    headers: { "Content-Type": "application/json" }
  });
}

function renderDemoPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false }
    }
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <DemoPage />
    </QueryClientProvider>
  );
}

describe("DemoPage", () => {
  it("runs invoice and payment flow", async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(jsonResponse({ data: { invoice_id: "21", status: "issued" } }))
      .mockResolvedValueOnce(jsonResponse({ data: { transaction_id: "83", status: "confirmed" } }));
    vi.stubGlobal("fetch", fetchMock);

    renderDemoPage();

    await userEvent.click(screen.getByRole("button", { name: /create/i }));
    expect(await screen.findByText("#21")).toBeInTheDocument();

    await userEvent.click(screen.getByRole("button", { name: /submit/i }));

    await waitFor(() => {
      expect(screen.getByText("83")).toBeInTheDocument();
      expect(screen.getByText("confirmed")).toBeInTheDocument();
    });
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });
});
