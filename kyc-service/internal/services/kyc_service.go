package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"kyc-service/internal/models"
)

// KYCService описывает бизнес-логику для обработки KYC запроса
type KYCService interface {
	ProcessKYC(req models.KYCRequest) (models.KYCResponse, error)
}

type KYCProvider interface {
	Process(ctx context.Context, req models.KYCRequest) (models.KYCResponse, error)
}

type kycService struct {
	provider KYCProvider
}

// NewKYCService создаёт новый экземпляр KYCService
func NewKYCService() KYCService {
	if providerURL := os.Getenv("KYC_PROVIDER_URL"); providerURL != "" {
		return &kycService{provider: &httpKYCProvider{
			url:    providerURL,
			client: &http.Client{Timeout: providerTimeout()},
		}}
	}
	return &kycService{provider: mockKYCProvider{delay: mockDelay()}}
}

func (s *kycService) ProcessKYC(req models.KYCRequest) (models.KYCResponse, error) {
	return s.provider.Process(context.Background(), req)
}

type mockKYCProvider struct {
	delay time.Duration
}

func (p mockKYCProvider) Process(ctx context.Context, req models.KYCRequest) (models.KYCResponse, error) {
	if p.delay > 0 {
		timer := time.NewTimer(p.delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return models.KYCResponse{}, ctx.Err()
		case <-timer.C:
		}
	}
	resp := models.KYCResponse{
		UserID:    req.UserID,
		KYCStatus: "verified",
		Message:   "KYC verification successful",
	}
	return resp, nil
}

type httpKYCProvider struct {
	url    string
	client *http.Client
}

func (p *httpKYCProvider) Process(ctx context.Context, req models.KYCRequest) (models.KYCResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return models.KYCResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.url, bytes.NewReader(payload))
	if err != nil {
		return models.KYCResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if apiKey := os.Getenv("KYC_PROVIDER_API_KEY"); apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return models.KYCResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return models.KYCResponse{}, fmt.Errorf("kyc provider returned status %d", resp.StatusCode)
	}

	var kycResp models.KYCResponse
	if err := json.NewDecoder(resp.Body).Decode(&kycResp); err != nil {
		return models.KYCResponse{}, err
	}
	if kycResp.UserID == "" {
		kycResp.UserID = req.UserID
	}
	return kycResp, nil
}

func providerTimeout() time.Duration {
	timeoutMS, err := strconv.Atoi(os.Getenv("KYC_PROVIDER_TIMEOUT_MS"))
	if err != nil || timeoutMS <= 0 {
		return 5 * time.Second
	}
	return time.Duration(timeoutMS) * time.Millisecond
}

func mockDelay() time.Duration {
	delayMS, err := strconv.Atoi(os.Getenv("KYC_MOCK_DELAY_MS"))
	if err != nil || delayMS <= 0 {
		return 0
	}
	return time.Duration(delayMS) * time.Millisecond
}
