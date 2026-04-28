import { useEffect, useState } from "react";
import { personaById, type Persona } from "./auth/personas";
import { AdminSystemHealthPage } from "./pages/AdminSystemHealthPage";
import { CompliancePage } from "./pages/CompliancePage";
import { Layout } from "./components/Layout";
import { AnalyticsPage } from "./pages/AnalyticsPage";
import { DemoPage } from "./pages/DemoPage";
import { FunctionsPage } from "./pages/FunctionsPage";
import { LoginPage } from "./pages/LoginPage";
import { MerchantPage } from "./pages/MerchantPage";
import { TablesPage } from "./pages/TablesPage";
import { UserPage } from "./pages/UserPage";

type View =
  | "login"
  | "user-dashboard"
  | "user-profile"
  | "user-kyc"
  | "user-wallets"
  | "user-transactions"
  | "merchant-dashboard"
  | "merchant-invoices"
  | "merchant-webhooks"
  | "merchant-terminals"
  | "merchant-analytics"
  | "compliance-kyc"
  | "compliance-merchant-verifications"
  | "compliance-risk-alerts"
  | "compliance-blacklist"
  | "admin-system-health"
  | "tables"
  | "functions"
  | "analytics"
  | "demo";

const viewByPath: Record<string, View> = {
  "/": "login",
  "/login": "login",
  "/user/dashboard": "user-dashboard",
  "/user/profile": "user-profile",
  "/user/kyc": "user-kyc",
  "/user/wallets": "user-wallets",
  "/user/transactions": "user-transactions",
  "/merchant/dashboard": "merchant-dashboard",
  "/merchant/invoices": "merchant-invoices",
  "/merchant/webhooks": "merchant-webhooks",
  "/merchant/terminals": "merchant-terminals",
  "/merchant/analytics": "merchant-analytics",
  "/compliance/kyc": "compliance-kyc",
  "/compliance/merchant-verifications": "compliance-merchant-verifications",
  "/compliance/risk-alerts": "compliance-risk-alerts",
  "/compliance/blacklist": "compliance-blacklist",
  "/admin/system-health": "admin-system-health",
  "/admin/tables": "tables",
  "/admin/functions": "functions",
  "/admin/analytics": "analytics",
  "/admin/demo": "demo"
};

const pathByView: Record<View, string> = {
  login: "/login",
  "user-dashboard": "/user/dashboard",
  "user-profile": "/user/profile",
  "user-kyc": "/user/kyc",
  "user-wallets": "/user/wallets",
  "user-transactions": "/user/transactions",
  "merchant-dashboard": "/merchant/dashboard",
  "merchant-invoices": "/merchant/invoices",
  "merchant-webhooks": "/merchant/webhooks",
  "merchant-terminals": "/merchant/terminals",
  "merchant-analytics": "/merchant/analytics",
  "compliance-kyc": "/compliance/kyc",
  "compliance-merchant-verifications": "/compliance/merchant-verifications",
  "compliance-risk-alerts": "/compliance/risk-alerts",
  "compliance-blacklist": "/compliance/blacklist",
  "admin-system-health": "/admin/system-health",
  tables: "/admin/tables",
  functions: "/admin/functions",
  analytics: "/admin/analytics",
  demo: "/admin/demo"
};

function viewFromLocation() {
  return viewByPath[window.location.pathname] ?? "tables";
}

function personaFromStorage() {
  if (typeof window.localStorage?.getItem !== "function") {
    return personaById(undefined);
  }
  return personaById(window.localStorage.getItem("depay.persona"));
}

export function App() {
  const [view, setView] = useState<View>(viewFromLocation);
  const [persona, setPersona] = useState<Persona>(personaFromStorage);

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

  function navigatePath(path: string) {
    const nextView = viewByPath[path] ?? "login";
    navigate(nextView);
  }

function changePersona(nextPersona: Persona) {
    if (typeof window.localStorage?.setItem === "function") {
      window.localStorage.setItem("depay.persona", nextPersona.id);
    }
    setPersona(nextPersona);
  }

  if (view === "login") {
    return <LoginPage persona={persona} onPersonaChange={changePersona} onComplete={navigatePath} />;
  }

  return (
    <Layout view={view} persona={persona} onViewChange={navigate}>
      {view === "user-dashboard" && <UserPage section="dashboard" />}
      {view === "user-profile" && <UserPage section="profile" />}
      {view === "user-kyc" && <UserPage section="kyc" />}
      {view === "user-wallets" && <UserPage section="wallets" />}
      {view === "user-transactions" && <UserPage section="transactions" />}
      {view === "merchant-dashboard" && <MerchantPage section="dashboard" />}
      {view === "merchant-invoices" && <MerchantPage section="invoices" />}
      {view === "merchant-webhooks" && <MerchantPage section="webhooks" />}
      {view === "merchant-terminals" && <MerchantPage section="terminals" />}
      {view === "merchant-analytics" && <MerchantPage section="analytics" />}
      {view === "compliance-kyc" && <CompliancePage section="kyc" />}
      {view === "compliance-merchant-verifications" && <CompliancePage section="merchant-verifications" />}
      {view === "compliance-risk-alerts" && <CompliancePage section="risk-alerts" />}
      {view === "compliance-blacklist" && <CompliancePage section="blacklist" />}
      {view === "admin-system-health" && <AdminSystemHealthPage />}
      {view === "tables" && <TablesPage />}
      {view === "functions" && <FunctionsPage />}
      {view === "analytics" && <AnalyticsPage />}
      {view === "demo" && <DemoPage />}
    </Layout>
  );
}
