type ErrorAlertProps = {
  error?: unknown;
};

export function ErrorAlert({ error }: ErrorAlertProps) {
  if (!error) {
    return null;
  }
  const message = error instanceof Error ? error.message : "Request failed";
  return <div className="error-alert">{message}</div>;
}
