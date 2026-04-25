type StatusBadgeProps = {
  value: string;
};

export function StatusBadge({ value }: StatusBadgeProps) {
  return <span className={`status status-${value.replace(/_/g, "-")}`}>{value}</span>;
}
