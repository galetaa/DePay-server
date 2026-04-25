DROP TRIGGER IF EXISTS trg_risk_alert_large_unverified_payment ON payment_transactions;
DROP FUNCTION IF EXISTS create_large_unverified_payment_alert();

DROP TRIGGER IF EXISTS trg_audit_payment_transactions ON payment_transactions;
DROP FUNCTION IF EXISTS audit_payment_transactions();

DROP TRIGGER IF EXISTS trg_transaction_completed_at ON payment_transactions;
DROP FUNCTION IF EXISTS set_transaction_completed_at();

DROP TRIGGER IF EXISTS trg_transaction_status_flow ON payment_transactions;
DROP FUNCTION IF EXISTS validate_transaction_status_flow();

DROP TRIGGER IF EXISTS trg_balance_history ON balances;
DROP FUNCTION IF EXISTS write_balance_history();

DROP TRIGGER IF EXISTS trg_balance_non_negative ON balances;
DROP FUNCTION IF EXISTS validate_balance_non_negative();

DROP TRIGGER IF EXISTS trg_wallet_owner_check ON wallets;
DROP FUNCTION IF EXISTS validate_wallet_owner();
