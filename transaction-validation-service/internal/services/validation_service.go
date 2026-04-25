package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"shared/validation"
	"transaction-validation-service/internal/models"
)

// ValidationService описывает логику валидации транзакций.
type ValidationService interface {
	Validate(req models.ValidationRequest) error
}

type validationService struct{}

func NewValidationService() ValidationService {
	return &validationService{}
}

func (s *validationService) Validate(req models.ValidationRequest) error {
	return validateStateless(req)
}

type postgresValidationService struct {
	db *sql.DB
}

func NewPostgresValidationService(db *sql.DB) ValidationService {
	return &postgresValidationService{db: db}
}

func (s *postgresValidationService) Validate(req models.ValidationRequest) error {
	if err := validateStateless(req); err != nil {
		return err
	}

	tx, err := s.loadTransaction(context.Background(), req.TransactionID)
	if err != nil {
		return err
	}

	if tx.Status == "confirmed" || tx.Status == "failed" || tx.Status == "cancelled" {
		return fmt.Errorf("transaction status %s is terminal", tx.Status)
	}
	if !equalDecimal(req.Amount, tx.Amount) {
		return errors.New("amount does not match transaction")
	}
	if !strings.EqualFold(req.Currency, tx.AssetSymbol) {
		return errors.New("currency does not match transaction asset")
	}
	if !strings.EqualFold(req.SenderAddress, tx.UserWalletAddress) {
		return errors.New("sender address does not match transaction wallet")
	}
	if !strings.EqualFold(req.RecipientAddress, tx.StoreWalletAddress) {
		return errors.New("recipient address does not match merchant wallet")
	}
	if tx.KYCStatus != "approved" {
		return errors.New("user KYC approval is required")
	}
	if tx.MerchantVerificationStatus != "approved" {
		return errors.New("merchant verification is required")
	}
	if tx.BlacklistedWallet {
		return errors.New("wallet is blacklisted")
	}
	if tx.CriticalRiskAlerts > 0 {
		return errors.New("open critical risk alert blocks transaction")
	}
	if !tx.HasSufficientBalance {
		return errors.New("insufficient funds")
	}
	if !tx.UserWalletOwned {
		return errors.New("user wallet ownership check failed")
	}
	if !tx.StoreWalletOwned {
		return errors.New("merchant wallet ownership check failed")
	}
	return nil
}

func validateStateless(req models.ValidationRequest) error {
	if strings.TrimSpace(req.Amount) == "0" {
		return errors.New("insufficient funds")
	}
	if err := validation.PositiveAmount(req.Amount); err != nil {
		return err
	}
	if err := validation.EVMAddress(req.SenderAddress); err != nil {
		return err
	}
	if err := validation.EVMAddress(req.RecipientAddress); err != nil {
		return err
	}
	if strings.TrimSpace(req.SignedData) == "" {
		return errors.New("signature is required")
	}
	return nil
}

type persistedTransaction struct {
	Amount                     string
	Status                     string
	AssetSymbol                string
	KYCStatus                  string
	MerchantVerificationStatus string
	UserWalletAddress          string
	StoreWalletAddress         string
	UserWalletOwned            bool
	StoreWalletOwned           bool
	HasSufficientBalance       bool
	BlacklistedWallet          bool
	CriticalRiskAlerts         int
}

func (s *postgresValidationService) loadTransaction(ctx context.Context, transactionID string) (persistedTransaction, error) {
	var tx persistedTransaction
	err := s.db.QueryRowContext(ctx, `
		SELECT
			pt.amount::text,
			pt.status::text,
			a.symbol,
			u.kyc_status::text,
			st.verification_status::text,
			uw.wallet_address,
			sw.wallet_address,
			uw.user_id = pt.user_id,
			sw.store_id = pt.store_id,
			COALESCE(b.balance >= pt.amount, false),
			EXISTS (
				SELECT 1
				FROM blacklisted_wallets bw
				WHERE bw.chain_id = uw.chain_id
				  AND lower(bw.wallet_address) IN (lower(uw.wallet_address), lower(sw.wallet_address))
			),
			(
				SELECT count(*)
				FROM risk_alerts ra
				WHERE ra.status IN ('open', 'in_review')
				  AND ra.risk_level IN ('high', 'critical')
				  AND (ra.user_id = pt.user_id OR ra.store_id = pt.store_id OR ra.transaction_id = pt.transaction_id)
			)::int
		FROM payment_transactions pt
		JOIN users u ON u.user_id = pt.user_id
		JOIN stores st ON st.store_id = pt.store_id
		JOIN assets a ON a.asset_id = pt.asset_id
		JOIN wallets uw ON uw.wallet_id = pt.user_wallet_id
		JOIN wallets sw ON sw.wallet_id = pt.store_wallet_id
		LEFT JOIN balances b ON b.wallet_id = pt.user_wallet_id AND b.asset_id = pt.asset_id
		WHERE pt.transaction_id::text = $1
		   OR pt.external_transaction_id = $1
	`, transactionID).Scan(
		&tx.Amount,
		&tx.Status,
		&tx.AssetSymbol,
		&tx.KYCStatus,
		&tx.MerchantVerificationStatus,
		&tx.UserWalletAddress,
		&tx.StoreWalletAddress,
		&tx.UserWalletOwned,
		&tx.StoreWalletOwned,
		&tx.HasSufficientBalance,
		&tx.BlacklistedWallet,
		&tx.CriticalRiskAlerts,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return persistedTransaction{}, errors.New("transaction not found")
		}
		return persistedTransaction{}, err
	}
	return tx, nil
}

func equalDecimal(left string, right string) bool {
	leftRat, ok := new(big.Rat).SetString(left)
	if !ok {
		return false
	}
	rightRat, ok := new(big.Rat).SetString(right)
	if !ok {
		return false
	}
	return leftRat.Cmp(rightRat) == 0
}
