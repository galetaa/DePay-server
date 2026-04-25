\set ON_ERROR_STOP on

BEGIN;

DO $$
BEGIN
    INSERT INTO wallets(chain_id, wallet_address, is_store_wallet)
    VALUES (1, '0x0000000000000000000000000000000000000aaa', FALSE);
    RAISE EXCEPTION 'trg_wallet_owner_check did not reject ownerless wallet';
EXCEPTION WHEN OTHERS THEN
    IF SQLERRM NOT LIKE '%wallet must belong to exactly one owner%' THEN
        RAISE;
    END IF;
END;
$$;

DO $$
BEGIN
    UPDATE balances
    SET balance = -1
    WHERE wallet_id = 1 AND asset_id = 1;
    RAISE EXCEPTION 'trg_balance_non_negative did not reject negative balance';
EXCEPTION WHEN OTHERS THEN
    IF SQLERRM NOT LIKE '%balance cannot be negative%' THEN
        RAISE;
    END IF;
END;
$$;

DO $$
DECLARE
    before_count BIGINT;
    after_count BIGINT;
BEGIN
    SELECT count(*) INTO before_count FROM balance_history WHERE wallet_id = 1 AND asset_id = 1;

    UPDATE balances
    SET balance = balance + 1
    WHERE wallet_id = 1 AND asset_id = 1;

    SELECT count(*) INTO after_count FROM balance_history WHERE wallet_id = 1 AND asset_id = 1;

    IF after_count <= before_count THEN
        RAISE EXCEPTION 'trg_balance_history did not create history row';
    END IF;
END;
$$;

DO $$
DECLARE
    tx_id BIGINT;
BEGIN
    INSERT INTO payment_transactions(
        user_id, store_id, asset_id, user_wallet_id, store_wallet_id, amount, amount_in_usdt, status
    )
    VALUES (1, 1, 1, 1, 16, 1, 10, 'created')
    RETURNING transaction_id INTO tx_id;

    UPDATE payment_transactions
    SET status = 'confirmed'
    WHERE transaction_id = tx_id;

    RAISE EXCEPTION 'trg_transaction_status_flow did not reject invalid transition';
EXCEPTION WHEN OTHERS THEN
    IF SQLERRM NOT LIKE '%invalid transaction status transition%' THEN
        RAISE;
    END IF;
END;
$$;

DO $$
DECLARE
    tx_id BIGINT;
    audit_count BIGINT;
BEGIN
    INSERT INTO payment_transactions(
        user_id, store_id, asset_id, user_wallet_id, store_wallet_id, amount, amount_in_usdt, status
    )
    VALUES (3, 1, 1, 3, 16, 1, 2000, 'created')
    RETURNING transaction_id INTO tx_id;

    UPDATE payment_transactions SET status = 'submitted' WHERE transaction_id = tx_id;
    UPDATE payment_transactions SET status = 'validated' WHERE transaction_id = tx_id;

    SELECT count(*) INTO audit_count
    FROM audit_logs
    WHERE entity_type = 'payment_transactions' AND entity_id = tx_id::TEXT;

    IF audit_count < 3 THEN
        RAISE EXCEPTION 'trg_audit_payment_transactions did not write audit rows';
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM risk_alerts
        WHERE transaction_id = tx_id
          AND alert_type = 'LARGE_UNVERIFIED_PAYMENT'
    ) THEN
        RAISE EXCEPTION 'trg_risk_alert_large_unverified_payment did not create alert';
    END IF;
END;
$$;

ROLLBACK;
