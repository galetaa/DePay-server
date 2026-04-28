package services

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"merchant-service/internal/models"
	"merchant-service/internal/repositories"
	"shared/auth"

	"golang.org/x/crypto/bcrypt"
)

type MerchantService interface {
	Register(req models.RegisterRequest) (models.Merchant, auth.TokenPair, error)
	Login(req models.LoginRequest) (models.Merchant, auth.TokenPair, error)
	GetMe(merchantID string) (models.Merchant, error)
	SubmitVerification(merchantID string, req models.VerificationRequest) (models.Merchant, error)
	CreateInvoice(merchantID string, req models.CreateInvoiceRequest) (models.Invoice, error)
	ListInvoices(merchantID string) ([]models.Invoice, error)
	CreateTerminal(merchantID string, req models.CreateTerminalRequest) (models.Terminal, error)
	ListTerminals(merchantID string) ([]models.Terminal, error)
	CreateWebhook(merchantID string, req models.CreateWebhookRequest) (models.Webhook, error)
	ListWebhooks(merchantID string) ([]models.Webhook, error)
	GetWebhook(merchantID string, webhookID string) (models.Webhook, error)
	DeleteWebhook(merchantID string, webhookID string) error
	TestWebhook(merchantID string, webhookID string) (models.WebhookDelivery, error)
	CreateAPIKey(merchantID string, req models.CreateAPIKeyRequest) (models.APIKey, error)
	ListAPIKeys(merchantID string) ([]models.APIKey, error)
	RevokeAPIKey(merchantID string, keyID string) error
}

type merchantService struct {
	repo repositories.MerchantRepository
}

func NewMerchantService(repo repositories.MerchantRepository) MerchantService {
	return &merchantService{repo: repo}
}

func (s *merchantService) Register(req models.RegisterRequest) (models.Merchant, auth.TokenPair, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.Merchant{}, auth.TokenPair{}, err
	}
	merchant, err := s.repo.Create(models.Merchant{
		OwnerEmail:   req.OwnerEmail,
		StoreName:    req.StoreName,
		LegalName:    req.LegalName,
		PasswordHash: string(hash),
	})
	if err != nil {
		return models.Merchant{}, auth.TokenPair{}, err
	}
	pair, err := auth.NewTokenPair(merchant.ID, []string{"merchant"}, 24*time.Hour, 30*24*time.Hour)
	return merchant, pair, err
}

func (s *merchantService) Login(req models.LoginRequest) (models.Merchant, auth.TokenPair, error) {
	merchant, err := s.repo.GetByEmail(req.OwnerEmail)
	if err != nil {
		return models.Merchant{}, auth.TokenPair{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(merchant.PasswordHash), []byte(req.Password)); err != nil {
		return models.Merchant{}, auth.TokenPair{}, errors.New("invalid credentials")
	}
	pair, err := auth.NewTokenPair(merchant.ID, []string{"merchant"}, 24*time.Hour, 30*24*time.Hour)
	return merchant, pair, err
}

func (s *merchantService) GetMe(merchantID string) (models.Merchant, error) {
	return s.repo.GetByID(merchantID)
}

func (s *merchantService) SubmitVerification(merchantID string, req models.VerificationRequest) (models.Merchant, error) {
	return s.repo.SubmitVerification(merchantID, req)
}

func (s *merchantService) CreateInvoice(merchantID string, req models.CreateInvoiceRequest) (models.Invoice, error) {
	return s.repo.CreateInvoice(merchantID, req)
}

func (s *merchantService) ListInvoices(merchantID string) ([]models.Invoice, error) {
	return s.repo.ListInvoices(merchantID)
}

func (s *merchantService) CreateTerminal(merchantID string, req models.CreateTerminalRequest) (models.Terminal, error) {
	return s.repo.CreateTerminal(merchantID, req)
}

func (s *merchantService) ListTerminals(merchantID string) ([]models.Terminal, error) {
	return s.repo.ListTerminals(merchantID)
}

func (s *merchantService) CreateWebhook(merchantID string, req models.CreateWebhookRequest) (models.Webhook, error) {
	if err := validateWebhookRequest(req); err != nil {
		return models.Webhook{}, err
	}
	return s.repo.CreateWebhook(merchantID, req)
}

func (s *merchantService) ListWebhooks(merchantID string) ([]models.Webhook, error) {
	return s.repo.ListWebhooks(merchantID)
}

func (s *merchantService) GetWebhook(merchantID string, webhookID string) (models.Webhook, error) {
	return s.repo.GetWebhook(merchantID, webhookID)
}

func (s *merchantService) DeleteWebhook(merchantID string, webhookID string) error {
	return s.repo.DeleteWebhook(merchantID, webhookID)
}

func (s *merchantService) TestWebhook(merchantID string, webhookID string) (models.WebhookDelivery, error) {
	return s.repo.TestWebhook(merchantID, webhookID)
}

func (s *merchantService) CreateAPIKey(merchantID string, req models.CreateAPIKeyRequest) (models.APIKey, error) {
	if req.Name == "" {
		return models.APIKey{}, errors.New("api key name is required")
	}
	if len(req.Scopes) == 0 {
		return models.APIKey{}, errors.New("api key scopes are required")
	}
	for _, scope := range req.Scopes {
		if !supportedAPIKeyScopes[scope] {
			return models.APIKey{}, errors.New("unsupported api key scope")
		}
	}
	return s.repo.CreateAPIKey(merchantID, req)
}

func (s *merchantService) ListAPIKeys(merchantID string) ([]models.APIKey, error) {
	return s.repo.ListAPIKeys(merchantID)
}

func (s *merchantService) RevokeAPIKey(merchantID string, keyID string) error {
	return s.repo.RevokeAPIKey(merchantID, keyID)
}

var supportedWebhookEvents = map[string]bool{
	"invoice.created":         true,
	"invoice.paid":            true,
	"invoice.expired":         true,
	"transaction.created":     true,
	"transaction.submitted":   true,
	"transaction.validated":   true,
	"transaction.broadcasted": true,
	"transaction.confirmed":   true,
	"transaction.failed":      true,
	"transaction.cancelled":   true,
}

var supportedAPIKeyScopes = map[string]bool{
	"invoice:read":     true,
	"invoice:write":    true,
	"transaction:read": true,
	"webhook:write":    true,
}

func validateWebhookRequest(req models.CreateWebhookRequest) error {
	parsed, err := url.ParseRequestURI(req.URL)
	if err != nil || parsed.Host == "" {
		return errors.New("invalid webhook URL")
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return errors.New("webhook URL must use http or https")
	}
	for _, eventType := range req.EventTypes {
		if !supportedWebhookEvents[eventType] {
			return errors.New("unsupported webhook event type")
		}
	}
	return nil
}
