CREATE TABLE users (
    user_id BIGSERIAL PRIMARY KEY,
    username VARCHAR(80) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    last_name VARCHAR(120),
    first_name VARCHAR(120),
    middle_name VARCHAR(120),
    birth_date DATE,
    address TEXT,
    kyc_status kyc_status_enum NOT NULL DEFAULT 'not_submitted',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_users_contact CHECK (email IS NOT NULL OR phone_number IS NOT NULL),
    CONSTRAINT chk_users_birth_date CHECK (birth_date IS NULL OR birth_date < CURRENT_DATE)
);

CREATE TABLE roles (
    role_id BIGSERIAL PRIMARY KEY,
    role_name VARCHAR(60) UNIQUE NOT NULL
);

CREATE TABLE user_roles (
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE refresh_tokens (
    token_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    ip_address INET,
    user_agent TEXT
);

CREATE TABLE stores (
    store_id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    store_name VARCHAR(160) NOT NULL,
    legal_name VARCHAR(255),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(20),
    store_address TEXT,
    verification_status verification_status_enum NOT NULL DEFAULT 'not_submitted',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE kyc_applications (
    kyc_application_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status kyc_status_enum NOT NULL DEFAULT 'pending',
    reviewer_id BIGINT REFERENCES users(user_id) ON DELETE SET NULL,
    rejection_reason TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    reviewed_at TIMESTAMPTZ
);

CREATE TABLE kyc_documents (
    kyc_document_id BIGSERIAL PRIMARY KEY,
    kyc_application_id BIGINT NOT NULL REFERENCES kyc_applications(kyc_application_id) ON DELETE CASCADE,
    document_type document_type_enum NOT NULL,
    document_url TEXT NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE merchant_verification_applications (
    merchant_verification_application_id BIGSERIAL PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES stores(store_id) ON DELETE CASCADE,
    status verification_status_enum NOT NULL DEFAULT 'pending',
    reviewer_id BIGINT REFERENCES users(user_id) ON DELETE SET NULL,
    rejection_reason TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    reviewed_at TIMESTAMPTZ
);

CREATE TABLE blockchains (
    chain_id BIGSERIAL PRIMARY KEY,
    chain_name VARCHAR(80) UNIQUE NOT NULL,
    native_symbol VARCHAR(20) NOT NULL,
    chain_type VARCHAR(40) NOT NULL DEFAULT 'evm',
    explorer_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE assets (
    asset_id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES blockchains(chain_id) ON DELETE CASCADE,
    symbol VARCHAR(20) NOT NULL,
    asset_name VARCHAR(120) NOT NULL,
    contract_address VARCHAR(120),
    decimals INT NOT NULL DEFAULT 18,
    is_native BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (chain_id, symbol)
);

CREATE TABLE wallets (
    wallet_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(user_id) ON DELETE CASCADE,
    store_id BIGINT REFERENCES stores(store_id) ON DELETE CASCADE,
    chain_id BIGINT NOT NULL REFERENCES blockchains(chain_id) ON DELETE RESTRICT,
    wallet_address VARCHAR(120) NOT NULL,
    public_key TEXT,
    is_store_wallet BOOLEAN NOT NULL DEFAULT FALSE,
    wallet_label VARCHAR(120),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (chain_id, wallet_address),
    CONSTRAINT chk_wallet_single_owner CHECK ((user_id IS NULL) <> (store_id IS NULL))
);

CREATE TABLE balances (
    wallet_id BIGINT NOT NULL REFERENCES wallets(wallet_id) ON DELETE CASCADE,
    asset_id BIGINT NOT NULL REFERENCES assets(asset_id) ON DELETE CASCADE,
    balance NUMERIC(30,8) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (wallet_id, asset_id),
    CONSTRAINT chk_balances_non_negative CHECK (balance >= 0)
);

CREATE TABLE balance_history (
    balance_history_id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallets(wallet_id) ON DELETE CASCADE,
    asset_id BIGINT NOT NULL REFERENCES assets(asset_id) ON DELETE CASCADE,
    old_balance NUMERIC(30,8),
    new_balance NUMERIC(30,8) NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    reason TEXT
);

CREATE TABLE payment_invoices (
    invoice_id BIGSERIAL PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES stores(store_id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(user_id) ON DELETE SET NULL,
    external_order_id VARCHAR(120),
    amount_usdt NUMERIC(30,8) NOT NULL,
    status invoice_status_enum NOT NULL DEFAULT 'issued',
    expires_at TIMESTAMPTZ NOT NULL,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_payment_invoices_amount CHECK (amount_usdt > 0),
    UNIQUE (store_id, external_order_id)
);

CREATE TABLE payment_invoice_items (
    invoice_item_id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES payment_invoices(invoice_id) ON DELETE CASCADE,
    item_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price_usdt NUMERIC(30,8) NOT NULL,
    CONSTRAINT chk_invoice_items_quantity CHECK (quantity > 0),
    CONSTRAINT chk_invoice_items_price CHECK (unit_price_usdt >= 0)
);

CREATE TABLE terminals (
    terminal_id BIGSERIAL PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES stores(store_id) ON DELETE CASCADE,
    serial_number VARCHAR(120) UNIQUE NOT NULL,
    secret_hash VARCHAR(255) NOT NULL,
    status terminal_status_enum NOT NULL DEFAULT 'active',
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE nfc_sessions (
    nfc_session_id BIGSERIAL PRIMARY KEY,
    terminal_id BIGINT NOT NULL REFERENCES terminals(terminal_id) ON DELETE CASCADE,
    invoice_id BIGINT REFERENCES payment_invoices(invoice_id) ON DELETE SET NULL,
    status nfc_session_status_enum NOT NULL DEFAULT 'created',
    session_token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payment_transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    external_transaction_id VARCHAR(120) UNIQUE,
    invoice_id BIGINT REFERENCES payment_invoices(invoice_id) ON DELETE SET NULL,
    nfc_session_id BIGINT REFERENCES nfc_sessions(nfc_session_id) ON DELETE SET NULL,
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    store_id BIGINT NOT NULL REFERENCES stores(store_id) ON DELETE RESTRICT,
    asset_id BIGINT NOT NULL REFERENCES assets(asset_id) ON DELETE RESTRICT,
    user_wallet_id BIGINT NOT NULL REFERENCES wallets(wallet_id) ON DELETE RESTRICT,
    store_wallet_id BIGINT NOT NULL REFERENCES wallets(wallet_id) ON DELETE RESTRICT,
    amount NUMERIC(30,8) NOT NULL,
    amount_in_usdt NUMERIC(30,8) NOT NULL,
    network_fee NUMERIC(30,8) NOT NULL DEFAULT 0,
    service_fee NUMERIC(30,8) NOT NULL DEFAULT 0,
    blockchain_tx_hash VARCHAR(160),
    status transaction_status_enum NOT NULL DEFAULT 'created',
    failure_reason TEXT,
    signed_payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    submitted_at TIMESTAMPTZ,
    validated_at TIMESTAMPTZ,
    broadcasted_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    CONSTRAINT chk_payment_transactions_amount CHECK (amount > 0),
    CONSTRAINT chk_payment_transactions_amount_usdt CHECK (amount_in_usdt > 0),
    CONSTRAINT chk_payment_transactions_fees CHECK (network_fee >= 0 AND service_fee >= 0)
);

CREATE TABLE rpc_nodes (
    rpc_node_id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES blockchains(chain_id) ON DELETE CASCADE,
    node_name VARCHAR(120) NOT NULL,
    rpc_url TEXT NOT NULL,
    status rpc_node_status_enum NOT NULL DEFAULT 'healthy',
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE rpc_node_checks (
    rpc_node_check_id BIGSERIAL PRIMARY KEY,
    rpc_node_id BIGINT NOT NULL REFERENCES rpc_nodes(rpc_node_id) ON DELETE CASCADE,
    status rpc_node_status_enum NOT NULL,
    latency_ms INT,
    block_height BIGINT,
    error_message TEXT,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_rpc_node_checks_latency CHECK (latency_ms IS NULL OR latency_ms >= 0)
);

CREATE TABLE exchange_rates (
    exchange_rate_id BIGSERIAL PRIMARY KEY,
    asset_id BIGINT NOT NULL REFERENCES assets(asset_id) ON DELETE CASCADE,
    rate_to_usdt NUMERIC(30,8) NOT NULL,
    captured_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_exchange_rates_positive CHECK (rate_to_usdt > 0)
);

CREATE TABLE blacklisted_wallets (
    blacklisted_wallet_id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL REFERENCES blockchains(chain_id) ON DELETE CASCADE,
    wallet_address VARCHAR(120) NOT NULL,
    reason TEXT NOT NULL,
    created_by BIGINT REFERENCES users(user_id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (chain_id, wallet_address)
);

CREATE TABLE risk_alerts (
    risk_alert_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(user_id) ON DELETE SET NULL,
    store_id BIGINT REFERENCES stores(store_id) ON DELETE SET NULL,
    transaction_id BIGINT REFERENCES payment_transactions(transaction_id) ON DELETE SET NULL,
    alert_type VARCHAR(80) NOT NULL,
    risk_level risk_level_enum NOT NULL,
    status risk_alert_status_enum NOT NULL DEFAULT 'open',
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE audit_logs (
    audit_log_id BIGSERIAL PRIMARY KEY,
    actor_user_id BIGINT REFERENCES users(user_id) ON DELETE SET NULL,
    action VARCHAR(120) NOT NULL,
    entity_type VARCHAR(120) NOT NULL,
    entity_id TEXT,
    old_values JSONB,
    new_values JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
