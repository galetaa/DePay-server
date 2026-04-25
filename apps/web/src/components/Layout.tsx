import { BarChart3, Database, FlaskConical, PlayCircle } from "lucide-react";
import type { ReactNode } from "react";

type View = "tables" | "functions" | "analytics" | "demo";

type LayoutProps = {
  view: View;
  onViewChange: (view: View) => void;
  children: ReactNode;
};

const navItems = [
  { id: "tables" as const, label: "Tables", icon: Database },
  { id: "functions" as const, label: "Functions", icon: FlaskConical },
  { id: "analytics" as const, label: "Analytics", icon: BarChart3 },
  { id: "demo" as const, label: "Demo", icon: PlayCircle }
];

export function Layout({ view, onViewChange, children }: LayoutProps) {
  return (
    <div className="shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-mark">D</span>
          <div>
            <strong>DePay</strong>
            <span>Admin</span>
          </div>
        </div>
        <nav className="nav">
          {navItems.map((item) => {
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
        </nav>
      </aside>
      <main className="main">{children}</main>
    </div>
  );
}
