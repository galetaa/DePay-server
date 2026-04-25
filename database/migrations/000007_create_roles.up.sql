DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'depay_admin') THEN
        CREATE ROLE depay_admin;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'depay_app') THEN
        CREATE ROLE depay_app;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'depay_compliance') THEN
        CREATE ROLE depay_compliance;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'depay_merchant_readonly') THEN
        CREATE ROLE depay_merchant_readonly;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'depay_user_readonly') THEN
        CREATE ROLE depay_user_readonly;
    END IF;
END
$$;

GRANT USAGE ON SCHEMA public TO depay_admin, depay_app, depay_compliance, depay_merchant_readonly, depay_user_readonly;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO depay_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO depay_admin;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO depay_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO depay_app;
GRANT SELECT, UPDATE ON users, kyc_applications, merchant_verification_applications, risk_alerts, blacklisted_wallets TO depay_compliance;
GRANT SELECT ON vw_store_transactions, payment_invoices, payment_transactions TO depay_merchant_readonly;
GRANT SELECT ON vw_user_wallet_balances, payment_transactions TO depay_user_readonly;
