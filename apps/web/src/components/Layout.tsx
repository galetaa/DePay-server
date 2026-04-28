import { BarChart3, CreditCard, Database, FlaskConical, Gauge, Link2, ListChecks, LogIn, MonitorSmartphone, PlayCircle, ReceiptText, ShieldAlert, UserCircle, WalletCards } from "lucide-react";
import type { ReactNode } from "react";
import type { Persona, Role } from "../auth/personas";

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

type LayoutProps = {
  view: View;
  persona: Persona;
  onViewChange: (view: View) => void;
  children: ReactNode;
};

const navSections = [
  {
    title: "Session",
    roles: ["user", "merchant", "compliance", "admin"] as Role[],
    items: [{ id: "login" as const, label: "Sign in", icon: LogIn }]
  },
  {
    title: "User",
    roles: ["user"] as Role[],
    items: [
      { id: "user-dashboard" as const, label: "Dashboard", icon: Gauge },
      { id: "user-profile" as const, label: "Profile", icon: UserCircle },
      { id: "user-kyc" as const, label: "KYC", icon: ListChecks },
      { id: "user-wallets" as const, label: "Wallets", icon: WalletCards },
      { id: "user-transactions" as const, label: "Transactions", icon: CreditCard }
    ]
  },
  {
    title: "Merchant",
    roles: ["merchant"] as Role[],
    items: [
      { id: "merchant-dashboard" as const, label: "Dashboard", icon: Gauge },
      { id: "merchant-invoices" as const, label: "Invoices", icon: ReceiptText },
      { id: "merchant-webhooks" as const, label: "Webhooks", icon: Link2 },
      { id: "merchant-terminals" as const, label: "Terminals", icon: MonitorSmartphone },
      { id: "merchant-analytics" as const, label: "Analytics", icon: BarChart3 }
    ]
  },
  {
    title: "Compliance",
    roles: ["compliance"] as Role[],
    items: [
      { id: "compliance-kyc" as const, label: "KYC", icon: ListChecks },
      { id: "compliance-merchant-verifications" as const, label: "Merchants", icon: UserCircle },
      { id: "compliance-risk-alerts" as const, label: "Risk Alerts", icon: ShieldAlert },
      { id: "compliance-blacklist" as const, label: "Blacklist", icon: WalletCards }
    ]
  },
  {
    title: "Admin",
    roles: ["admin"] as Role[],
    items: [
      { id: "admin-system-health" as const, label: "System Health", icon: Gauge },
      { id: "tables" as const, label: "Tables", icon: Database },
      { id: "functions" as const, label: "Functions", icon: FlaskConical },
      { id: "analytics" as const, label: "Analytics", icon: BarChart3 },
      { id: "demo" as const, label: "Demo", icon: PlayCircle }
    ]
  }
];

export function Layout({ view, persona, onViewChange, children }: LayoutProps) {
  const sections = navSections.filter((section) => section.roles.some((role) => persona.roles.includes(role)));

  return (
    <div className="shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-mark">D</span>
          <div>
            <strong>DePay</strong>
            <span>{persona.label}</span>
          </div>
        </div>
        <nav className="nav">
          {sections.map((section) => (
            <div className="nav-section" key={section.title}>
              <span className="nav-section-title">{section.title}</span>
              {section.items.map((item) => {
                const Icon = item.icon;
                return (
                  <button
                    key={item.id}
                    type="button"
                    className={view === item.id ? "nav-button active" : "nav-button"}
                    onClick={() => onViewChange(item.id)}
                    title={item.label}
                  >
                    <Icon size={18} />
                    <span>{item.label}</span>
                  </button>
                );
              })}
            </div>
          ))}
        </nav>
      </aside>
      <main className="main">{children}</main>
    </div>
  );
}
