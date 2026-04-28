type EmptyStateProps = {
  title: string;
};

export function EmptyState({ title }: EmptyStateProps) {
  return <div className="empty">{title}</div>;
}
