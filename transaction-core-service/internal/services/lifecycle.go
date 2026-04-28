package services

import "fmt"

var terminalTransactionStatuses = map[string]bool{
	"confirmed": true,
	"failed":    true,
	"cancelled": true,
}

var allowedTransactionTransitions = map[string]map[string]bool{
	"created": {
		"submitted": true,
		"cancelled": true,
	},
	"submitted": {
		"validated": true,
		"cancelled": true,
		"failed":    true,
	},
	"validated": {
		"broadcasted": true,
		"failed":      true,
	},
	"broadcasted": {
		"confirmed": true,
		"failed":    true,
	},
}

type transitionDecision struct {
	Idempotent bool
}

func validateTransactionTransition(from string, to string) (transitionDecision, error) {
	if from == "" {
		from = "created"
	}
	if to == "" {
		return transitionDecision{}, fmt.Errorf("transaction status is required")
	}
	if from == to {
		return transitionDecision{Idempotent: true}, nil
	}
	if terminalTransactionStatuses[from] {
		return transitionDecision{}, fmt.Errorf("transaction status %s is terminal", from)
	}
	if allowedTransactionTransitions[from][to] {
		return transitionDecision{}, nil
	}
	return transitionDecision{}, fmt.Errorf("invalid transaction status transition from %s to %s", from, to)
}
