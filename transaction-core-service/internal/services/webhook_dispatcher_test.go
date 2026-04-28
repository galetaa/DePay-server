package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSignPayloadUsesTimestampAndBody(t *testing.T) {
	body := []byte(`{"type":"transaction.confirmed"}`)
	timestamp := "1714242424"
	secret := "webhook-secret"

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	require.Equal(t, expected, signPayload(timestamp, body, secret))
}

func TestDeliveryFailureStateSchedulesRetry(t *testing.T) {
	status, nextAttemptAt := deliveryFailureState(1, http.StatusInternalServerError)

	require.Equal(t, "retry_scheduled", status)
	require.True(t, nextAttemptAt.Valid)
	require.WithinDuration(t, time.Now().UTC().Add(30*time.Second), nextAttemptAt.Time, 2*time.Second)
}

func TestDeliveryFailureStateDeadLettersAfterMaxAttempts(t *testing.T) {
	status, nextAttemptAt := deliveryFailureState(5, http.StatusTooManyRequests)

	require.Equal(t, "dead_letter", status)
	require.False(t, nextAttemptAt.Valid)
}

func TestDeliveryFailureStateDoesNotRetryNonRetryableStatus(t *testing.T) {
	status, nextAttemptAt := deliveryFailureState(1, http.StatusBadRequest)

	require.Equal(t, "failed", status)
	require.False(t, nextAttemptAt.Valid)
}
