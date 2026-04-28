package services

import (
	"testing"

	"merchant-service/internal/models"
)

func TestValidateWebhookRequestRejectsInvalidURL(t *testing.T) {
	err := validateWebhookRequest(models.CreateWebhookRequest{
		URL:        "ftp://example.test/hook",
		EventTypes: []string{"transaction.confirmed"},
		Secret:     "secret",
	})

	if err == nil {
		t.Fatal("expected invalid URL error")
	}
}

func TestValidateWebhookRequestRejectsUnsupportedEvent(t *testing.T) {
	err := validateWebhookRequest(models.CreateWebhookRequest{
		URL:        "https://example.test/hook",
		EventTypes: []string{"transaction.unknown"},
		Secret:     "secret",
	})

	if err == nil {
		t.Fatal("expected unsupported event error")
	}
}

func TestValidateWebhookRequestAcceptsSupportedEvents(t *testing.T) {
	err := validateWebhookRequest(models.CreateWebhookRequest{
		URL: "https://example.test/hook",
		EventTypes: []string{
			"invoice.created",
			"invoice.paid",
			"invoice.expired",
			"transaction.cancelled",
		},
		Secret: "secret",
	})

	if err != nil {
		t.Fatalf("expected supported events, got %v", err)
	}
}

func TestCreateAPIKeyRejectsUnsupportedScope(t *testing.T) {
	service := NewMerchantService(nil)

	_, err := service.CreateAPIKey("1", models.CreateAPIKeyRequest{
		Name:   "server",
		Scopes: []string{"wallet:write"},
	})

	if err == nil {
		t.Fatal("expected unsupported scope error")
	}
}
