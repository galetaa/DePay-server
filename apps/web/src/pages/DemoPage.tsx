import { CheckCircle2, FilePlus2, Send } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { api } from "../api/client";
import { ErrorAlert } from "../components/ErrorAlert";

export function DemoPage() {
  const [invoiceId, setInvoiceId] = useState("");
  const invoiceMutation = useMutation({
    mutationFn: api.createDemoInvoice,
    onSuccess: (data) => setInvoiceId(String(data.invoice_id ?? ""))
  });
  const paymentMutation = useMutation({
    mutationFn: () => api.submitDemoPayment(invoiceId)
  });

  return (
    <section className="page">
      <div className="page-header">
        <div>
          <h1>Demo</h1>
          <p>Payment flow</p>
        </div>
      </div>
      <ErrorAlert error={invoiceMutation.error ?? paymentMutation.error} />
      <div className="demo-flow">
        <div className="demo-step">
          <FilePlus2 size={22} />
          <div>
            <strong>Invoice</strong>
            <span>{invoiceId ? `#${invoiceId}` : "Ready"}</span>
          </div>
          <button type="button" className="primary-button" onClick={() => invoiceMutation.mutate()} disabled={invoiceMutation.isPending}>
            Create
          </button>
        </div>
        <div className="demo-step">
          <Send size={22} />
          <div>
            <strong>Transaction</strong>
            <span>{paymentMutation.data ? String(paymentMutation.data.transaction_id) : "Waiting"}</span>
          </div>
          <button type="button" className="primary-button" onClick={() => paymentMutation.mutate()} disabled={!invoiceId || paymentMutation.isPending}>
            Submit
          </button>
        </div>
        <div className="demo-step">
          <CheckCircle2 size={22} />
          <div>
            <strong>Status</strong>
            <span>{paymentMutation.data ? String(paymentMutation.data.status) : "Pending"}</span>
          </div>
        </div>
      </div>
    </section>
  );
}
