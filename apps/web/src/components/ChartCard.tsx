import type { ReactNode } from "react";

type ChartCardProps = {
  title: string;
  children: ReactNode;
};

export function ChartCard({ title, children }: ChartCardProps) {
  return (
    <section className="chart-card">
      <header>{title}</header>
      <div className="chart-body">{children}</div>
    </section>
  );
}
