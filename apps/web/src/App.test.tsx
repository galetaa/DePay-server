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

function renderApp(path: string) {
  window.history.pushState({}, "", path);
  vi.stubGlobal("fetch", vi.fn().mockResolvedValue(jsonResponse({ data: { tables: ["users"], columns: [], rows: [] } })));

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
});
