package repositories

import (
	"strings"
	"testing"
	"time"

	"merchant-service/internal/models"
)

func approvedMemoryMerchantRepo() *memoryMerchantRepo {
	repo := NewMerchantRepository().(*memoryMerchantRepo)
	repo.merchants["1"] = models.Merchant{
		ID:                 "1",
		OwnerEmail:         "merchant@example.test",
		StoreName:          "Test Store",
		VerificationStatus: "approved",
		CreatedAt:          time.Now().UTC(),
	}
	return repo
}

func TestMemoryWebhookRepositoryRejectsDuplicateURL(t *testing.T) {
	repo := approvedMemoryMerchantRepo()
	req := models.CreateWebhookRequest{
		URL:    "https://example.test/hook",
		Secret: "secret",
	}

	_, err := repo.CreateWebhook("1", req)
	if err != nil {
		t.Fatalf("expected first webhook to be created, got %v", err)
	}
	_, err = repo.CreateWebhook("1", req)

	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected duplicate URL error, got %v", err)
	}
}

func TestMemoryWebhookRepositoryTestDelivery(t *testing.T) {
	repo := approvedMemoryMerchantRepo()
	webhook, err := repo.CreateWebhook("1", models.CreateWebhookRequest{
		URL:    "https://example.test/hook",
		Secret: "secret",
	})
	if err != nil {
		t.Fatalf("expected webhook to be created, got %v", err)
	}

	delivery, err := repo.TestWebhook("1", webhook.ID)

	if err != nil {
		t.Fatalf("expected test delivery, got %v", err)
	}
	if delivery.WebhookID != webhook.ID {
		t.Fatalf("expected webhook id %s, got %s", webhook.ID, delivery.WebhookID)
	}
	if delivery.EventType != "webhook.test" || delivery.Status != "delivered" {
		t.Fatalf("unexpected delivery state: %#v", delivery)
	}
	if delivery.Attempts != 1 || delivery.ResponseStatus != 204 {
		t.Fatalf("unexpected delivery attempt metadata: %#v", delivery)
	}
	if delivery.DeliveredAt == nil {
		t.Fatal("expected delivered_at")
	}
}

func TestMemoryAPIKeyRepositoryShowsSecretOnlyOnCreate(t *testing.T) {
	repo := approvedMemoryMerchantRepo()
	key, err := repo.CreateAPIKey("1", models.CreateAPIKeyRequest{
		Name:   "server",
		Scopes: []string{"invoice:read", "webhook:write"},
	})
	if err != nil {
		t.Fatalf("expected api key, got %v", err)
	}
	if key.Secret == "" || !strings.HasPrefix(key.Secret, "depay_") {
		t.Fatalf("expected one-time raw secret, got %#v", key)
	}

	keys, err := repo.ListAPIKeys("1")
	if err != nil {
		t.Fatalf("expected api key list, got %v", err)
	}
	if len(keys) != 1 || keys[0].Secret != "" {
		t.Fatalf("expected listed key without secret, got %#v", keys)
	}
}

func TestMemoryAPIKeyRepositoryRevokesKey(t *testing.T) {
	repo := approvedMemoryMerchantRepo()
	key, err := repo.CreateAPIKey("1", models.CreateAPIKeyRequest{
		Name:   "server",
		Scopes: []string{"invoice:read"},
	})
	if err != nil {
		t.Fatalf("expected api key, got %v", err)
	}
	if err := repo.RevokeAPIKey("1", key.ID); err != nil {
		t.Fatalf("expected revoke, got %v", err)
	}

	keys, err := repo.ListAPIKeys("1")
	if err != nil {
		t.Fatalf("expected api key list, got %v", err)
	}
	if keys[0].RevokedAt == nil {
		t.Fatal("expected revoked_at")
	}
}
