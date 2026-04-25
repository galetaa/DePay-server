import { useEffect, useState } from "react";
import { Layout } from "./components/Layout";
import { AnalyticsPage } from "./pages/AnalyticsPage";
import { DemoPage } from "./pages/DemoPage";
import { FunctionsPage } from "./pages/FunctionsPage";
import { TablesPage } from "./pages/TablesPage";

type View = "tables" | "functions" | "analytics" | "demo";

const viewByPath: Record<string, View> = {
  "/": "tables",
  "/admin/tables": "tables",
  "/admin/functions": "functions",
  "/admin/analytics": "analytics",
  "/admin/demo": "demo"
};

const pathByView: Record<View, string> = {
  tables: "/admin/tables",
  functions: "/admin/functions",
  analytics: "/admin/analytics",
  demo: "/admin/demo"
};

function viewFromLocation() {
  return viewByPath[window.location.pathname] ?? "tables";
}

export function App() {
  const [view, setView] = useState<View>(viewFromLocation);

  useEffect(() => {
    function onPopState() {
      setView(viewFromLocation());
    }

    window.addEventListener("popstate", onPopState);
    return () => window.removeEventListener("popstate", onPopState);
  }, []);

  function navigate(nextView: View) {
    setView(nextView);
    const nextPath = pathByView[nextView];
    if (window.location.pathname !== nextPath) {
      window.history.pushState({}, "", nextPath);
    }
  }

  return (
    <Layout view={view} onViewChange={navigate}>
      {view === "tables" && <TablesPage />}
      {view === "functions" && <FunctionsPage />}
      {view === "analytics" && <AnalyticsPage />}
      {view === "demo" && <DemoPage />}
    </Layout>
  );
}
