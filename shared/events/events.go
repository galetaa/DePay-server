package events

import "time"

const (
	TransactionCreated   = "transaction.created"
	TransactionSubmitted = "transaction.submitted"
	TransactionValidated = "transaction.validated"
	TransactionBroadcast = "transaction.broadcasted"
	TransactionConfirmed = "transaction.confirmed"
	TransactionFailed    = "transaction.failed"
	BalanceUpdated       = "balance.updated"
	RiskAlertCreated     = "risk_alert.created"
)

type Event struct {
	Type       string         `json:"type"`
	OccurredAt time.Time      `json:"occurred_at"`
	Payload    map[string]any `json:"payload"`
}

func New(eventType string, payload map[string]any) Event {
	return Event{
		Type:       eventType,
		OccurredAt: time.Now().UTC(),
		Payload:    payload,
	}
}
