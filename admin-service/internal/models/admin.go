package models

type TableListResponse struct {
	Tables []string `json:"tables"`
}

type TableRowsResponse struct {
	Columns []string         `json:"columns"`
	Rows    []map[string]any `json:"rows"`
	Limit   int              `json:"limit"`
}

type ExecuteFunctionRequest struct {
	Params []string `json:"params"`
}

type DemoPaymentRequest struct {
	InvoiceID string `json:"invoice_id"`
	UserID    string `json:"user_id"`
	WalletID  string `json:"wallet_id"`
}
