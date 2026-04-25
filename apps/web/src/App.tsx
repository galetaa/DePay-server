import { useState } from "react";
import { Layout } from "./components/Layout";
import { AnalyticsPage } from "./pages/AnalyticsPage";
import { DemoPage } from "./pages/DemoPage";
import { FunctionsPage } from "./pages/FunctionsPage";
import { TablesPage } from "./pages/TablesPage";

type View = "tables" | "functions" | "analytics" | "demo";

export function App() {
  const [view, setView] = useState<View>("tables");

  return (
    <Layout view={view} onViewChange={setView}>
      {view === "tables" && <TablesPage />}
      {view === "functions" && <FunctionsPage />}
      {view === "analytics" && <AnalyticsPage />}
      {view === "demo" && <DemoPage />}
    </Layout>
  );
}
