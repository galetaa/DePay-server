CREATE OR REPLACE FUNCTION validate_wallet_owner()
RETURNS TRIGGER AS $$
BEGIN
    IF (NEW.user_id IS NULL AND NEW.store_id IS NULL) OR (NEW.user_id IS NOT NULL AND NEW.store_id IS NOT NULL) THEN
        RAISE EXCEPTION 'wallet must belong to exactly one owner';
    END IF;

    IF NEW.is_store_wallet AND NEW.store_id IS NULL THEN
        RAISE EXCEPTION 'store wallet must have store_id';
    END IF;

    IF NOT NEW.is_store_wallet AND NEW.user_id IS NULL THEN
        RAISE EXCEPTION 'user wallet must have user_id';
    END IF;

    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_wallet_owner_check
BEFORE INSERT OR UPDATE ON wallets
FOR EACH ROW EXECUTE FUNCTION validate_wallet_owner();

CREATE OR REPLACE FUNCTION validate_balance_non_negative()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.balance < 0 THEN
        RAISE EXCEPTION 'balance cannot be negative';
    END IF;

    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_balance_non_negative
BEFORE INSERT OR UPDATE ON balances
FOR EACH ROW EXECUTE FUNCTION validate_balance_non_negative();

CREATE OR REPLACE FUNCTION write_balance_history()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO balance_history(wallet_id, asset_id, old_balance, new_balance, reason)
    VALUES (
        NEW.wallet_id,
        NEW.asset_id,
        CASE WHEN TG_OP = 'UPDATE' THEN OLD.balance ELSE NULL END,
        NEW.balance,
        lower(TG_OP)
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_balance_history
AFTER INSERT OR UPDATE ON balances
FOR EACH ROW EXECUTE FUNCTION write_balance_history();

CREATE OR REPLACE FUNCTION validate_transaction_status_flow()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status = NEW.status THEN
        RETURN NEW;
    END IF;

    IF OLD.status IN ('confirmed', 'failed', 'cancelled') THEN
        RAISE EXCEPTION 'transaction status % is terminal', OLD.status;
    END IF;

    IF NOT (
        (OLD.status = 'created' AND NEW.status IN ('submitted', 'cancelled', 'failed')) OR
        (OLD.status = 'submitted' AND NEW.status IN ('validated', 'cancelled', 'failed')) OR
        (OLD.status = 'validated' AND NEW.status IN ('broadcasted', 'failed')) OR
        (OLD.status = 'broadcasted' AND NEW.status IN ('confirmed', 'failed'))
    ) THEN
        RAISE EXCEPTION 'invalid transaction status transition from % to %', OLD.status, NEW.status;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_transaction_status_flow
BEFORE UPDATE OF status ON payment_transactions
FOR EACH ROW EXECUTE FUNCTION validate_transaction_status_flow();

CREATE OR REPLACE FUNCTION set_transaction_completed_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status IN ('confirmed', 'failed', 'cancelled') AND NEW.completed_at IS NULL THEN
        NEW.completed_at = now();
    END IF;

    IF NEW.status = 'submitted' AND NEW.submitted_at IS NULL THEN
        NEW.submitted_at = now();
    END IF;

    IF NEW.status = 'validated' AND NEW.validated_at IS NULL THEN
        NEW.validated_at = now();
    END IF;

    IF NEW.status = 'broadcasted' AND NEW.broadcasted_at IS NULL THEN
        NEW.broadcasted_at = now();
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_transaction_completed_at
BEFORE UPDATE OF status ON payment_transactions
FOR EACH ROW EXECUTE FUNCTION set_transaction_completed_at();

CREATE OR REPLACE FUNCTION audit_payment_transactions()
RETURNS TRIGGER AS $$
DECLARE
    entity_id_text TEXT;
BEGIN
    entity_id_text = COALESCE(NEW.transaction_id, OLD.transaction_id)::TEXT;

    INSERT INTO audit_logs(action, entity_type, entity_id, old_values, new_values)
    VALUES (
        TG_OP,
        'payment_transactions',
        entity_id_text,
        CASE WHEN TG_OP IN ('UPDATE', 'DELETE') THEN to_jsonb(OLD) ELSE NULL END,
        CASE WHEN TG_OP IN ('INSERT', 'UPDATE') THEN to_jsonb(NEW) ELSE NULL END
    );

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_audit_payment_transactions
AFTER INSERT OR UPDATE OR DELETE ON payment_transactions
FOR EACH ROW EXECUTE FUNCTION audit_payment_transactions();

CREATE OR REPLACE FUNCTION create_large_unverified_payment_alert()
RETURNS TRIGGER AS $$
DECLARE
    current_kyc_status kyc_status_enum;
BEGIN
    SELECT u.kyc_status INTO current_kyc_status
    FROM users u
    WHERE u.user_id = NEW.user_id;

    IF NEW.amount_in_usdt >= 1000
        AND COALESCE(current_kyc_status, 'not_submitted') <> 'approved'
        AND NEW.status IN ('validated', 'broadcasted', 'confirmed')
        AND NOT EXISTS (
            SELECT 1
            FROM risk_alerts r
            WHERE r.transaction_id = NEW.transaction_id
              AND r.alert_type = 'LARGE_UNVERIFIED_PAYMENT'
        )
    THEN
        INSERT INTO risk_alerts(user_id, store_id, transaction_id, alert_type, risk_level, details)
        VALUES (
            NEW.user_id,
            NEW.store_id,
            NEW.transaction_id,
            'LARGE_UNVERIFIED_PAYMENT',
            'high',
            jsonb_build_object('amount_in_usdt', NEW.amount_in_usdt, 'kyc_status', current_kyc_status)
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_risk_alert_large_unverified_payment
AFTER INSERT OR UPDATE OF status, amount_in_usdt ON payment_transactions
FOR EACH ROW EXECUTE FUNCTION create_large_unverified_payment_alert();
