package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"transaction-core-service/internal/models"
)

type BroadcastResult struct {
	TxHash string
}

type BlockchainBroadcaster interface {
	Broadcast(ctx context.Context, tx models.Transaction) (BroadcastResult, error)
}

func NewBlockchainBroadcasterFromEnv() BlockchainBroadcaster {
	rpcURL := os.Getenv("BLOCKCHAIN_RPC_URL")
	if rpcURL == "" {
		return mockBlockchainBroadcaster{}
	}
	return &rpcBlockchainBroadcaster{
		url:    rpcURL,
		client: &http.Client{Timeout: envDuration("BLOCKCHAIN_RPC_TIMEOUT_MS", 10*time.Second)},
	}
}

type mockBlockchainBroadcaster struct{}

func (mockBlockchainBroadcaster) Broadcast(_ context.Context, tx models.Transaction) (BroadcastResult, error) {
	seed := tx.TransactionID + ":" + tx.StoreID + ":" + tx.Amount
	sum := sha256.Sum256([]byte(seed))
	return BroadcastResult{TxHash: "0x" + hex.EncodeToString(sum[:])}, nil
}

type rpcBlockchainBroadcaster struct {
	url    string
	client *http.Client
}

func (b *rpcBlockchainBroadcaster) Broadcast(ctx context.Context, tx models.Transaction) (BroadcastResult, error) {
	if tx.SignedTransaction == "" {
		return BroadcastResult{}, errors.New("signed_transaction is required for RPC broadcast")
	}

	reqBody := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "eth_sendRawTransaction",
		Params:  []string{tx.SignedTransaction},
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return BroadcastResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.url, bytes.NewReader(payload))
	if err != nil {
		return BroadcastResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return BroadcastResult{}, err
	}
	defer resp.Body.Close()

	var rpcResp jsonRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return BroadcastResult{}, err
	}
	if rpcResp.Error != nil {
		return BroadcastResult{}, fmt.Errorf("blockchain rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}
	if rpcResp.Result == "" {
		return BroadcastResult{}, errors.New("blockchain rpc returned empty transaction hash")
	}
	return BroadcastResult{TxHash: rpcResp.Result}, nil
}

type jsonRPCRequest struct {
	JSONRPC string   `json:"jsonrpc"`
	ID      int      `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

type jsonRPCResponse struct {
	Result string        `json:"result"`
	Error  *jsonRPCError `json:"error"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func envDuration(name string, fallback time.Duration) time.Duration {
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
