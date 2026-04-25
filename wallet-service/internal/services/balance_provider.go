package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"
)

type BalanceProvider interface {
	GetBalance(ctx context.Context, address string) (string, error)
}

func NewBalanceProviderFromEnv() BalanceProvider {
	rpcURL := os.Getenv("WALLET_BALANCE_RPC_URL")
	if rpcURL == "" {
		rpcURL = os.Getenv("BLOCKCHAIN_RPC_URL")
	}
	if rpcURL == "" {
		return mockBalanceProvider{}
	}
	return &rpcBalanceProvider{
		url:    rpcURL,
		client: &http.Client{Timeout: walletEnvDuration("WALLET_BALANCE_RPC_TIMEOUT_MS", 10*time.Second)},
	}
}

type mockBalanceProvider struct{}

func (mockBalanceProvider) GetBalance(_ context.Context, address string) (string, error) {
	return "1000000000000000000", nil
}

type rpcBalanceProvider struct {
	url    string
	client *http.Client
}

func (p *rpcBalanceProvider) GetBalance(ctx context.Context, address string) (string, error) {
	body, err := json.Marshal(walletJSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "eth_getBalance",
		Params:  []string{address, "latest"},
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var rpcResp walletJSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return "", err
	}
	if rpcResp.Error != nil {
		return "", errors.New(rpcResp.Error.Message)
	}
	if rpcResp.Result == "" {
		return "", errors.New("empty RPC balance result")
	}
	value := new(big.Int)
	if _, ok := value.SetString(stripHexPrefix(rpcResp.Result), 16); !ok {
		return "", errors.New("invalid RPC balance result")
	}
	return value.String(), nil
}

type walletJSONRPCRequest struct {
	JSONRPC string   `json:"jsonrpc"`
	ID      int      `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

type walletJSONRPCResponse struct {
	Result string              `json:"result"`
	Error  *walletJSONRPCError `json:"error"`
}

type walletJSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func stripHexPrefix(value string) string {
	if len(value) >= 2 && value[:2] == "0x" {
		return value[2:]
	}
	return value
}

func walletEnvDuration(name string, fallback time.Duration) time.Duration {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	ms, err := strconv.Atoi(value)
	if err != nil || ms <= 0 {
		return fallback
	}
	return time.Duration(ms) * time.Millisecond
}
