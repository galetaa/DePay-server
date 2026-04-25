BEGIN;

TRUNCATE TABLE
    merchant_webhook_deliveries,
    merchant_webhooks,
    audit_logs,
    risk_alerts,
    blacklisted_wallets,
    exchange_rates,
    rpc_node_checks,
    rpc_nodes,
    payment_transactions,
    nfc_sessions,
    terminals,
    payment_invoice_items,
    payment_invoices,
    balance_history,
    balances,
    wallets,
    assets,
    blockchains,
    merchant_verification_applications,
    kyc_documents,
    kyc_applications,
    stores,
    refresh_tokens,
    user_roles,
    roles,
    users
RESTART IDENTITY CASCADE;

INSERT INTO roles(role_name) VALUES
    ('user'),
    ('merchant'),
    ('compliance_operator'),
    ('admin'),
    ('terminal');

INSERT INTO users(username, email, phone_number, password_hash, last_name, first_name, birth_date, address, kyc_status, created_at) VALUES
    ('alice', 'alice@example.com', '+10000000001', crypt('password123', gen_salt('bf')), 'Stone', 'Alice', '1995-01-10', 'Berlin', 'approved', now() - interval '80 days'),
    ('boris', 'boris@example.com', '+10000000002', crypt('password123', gen_salt('bf')), 'Volkov', 'Boris', '1991-03-14', 'Moscow', 'pending', now() - interval '70 days'),
    ('cora', 'cora@example.com', '+10000000003', crypt('password123', gen_salt('bf')), 'Miller', 'Cora', '1997-07-07', 'Prague', 'not_submitted', now() - interval '60 days'),
    ('dmitry', 'dmitry@example.com', '+10000000004', crypt('password123', gen_salt('bf')), 'Ivanov', 'Dmitry', '1989-11-20', 'Saint Petersburg', 'approved', now() - interval '55 days'),
    ('eva', 'eva@example.com', '+10000000005', crypt('password123', gen_salt('bf')), 'Frost', 'Eva', '1993-05-25', 'Paris', 'rejected', now() - interval '44 days'),
    ('farid', 'farid@example.com', '+10000000006', crypt('password123', gen_salt('bf')), 'Khan', 'Farid', '1990-08-18', 'Dubai', 'approved', now() - interval '38 days'),
    ('gina', 'gina@example.com', '+10000000007', crypt('password123', gen_salt('bf')), 'Lopez', 'Gina', '1998-02-03', 'Madrid', 'pending', now() - interval '34 days'),
    ('haruto', 'haruto@example.com', '+10000000008', crypt('password123', gen_salt('bf')), 'Sato', 'Haruto', '1994-09-12', 'Tokyo', 'not_submitted', now() - interval '25 days'),
    ('irina', 'irina@example.com', '+10000000009', crypt('password123', gen_salt('bf')), 'Smirnova', 'Irina', '1996-12-01', 'Tbilisi', 'approved', now() - interval '18 days'),
    ('jon', 'jon@example.com', '+10000000010', crypt('password123', gen_salt('bf')), 'Reed', 'Jon', '1992-04-22', 'London', 'not_submitted', now() - interval '12 days');

INSERT INTO user_roles(user_id, role_id)
SELECT user_id, (SELECT role_id FROM roles WHERE role_name = 'user')
FROM users;

INSERT INTO user_roles(user_id, role_id)
SELECT user_id, (SELECT role_id FROM roles WHERE role_name = 'merchant')
FROM users
WHERE username IN ('alice', 'boris', 'dmitry', 'farid', 'irina');

INSERT INTO user_roles(user_id, role_id)
SELECT user_id, (SELECT role_id FROM roles WHERE role_name = 'admin')
FROM users
WHERE username = 'alice';

INSERT INTO user_roles(user_id, role_id)
SELECT user_id, (SELECT role_id FROM roles WHERE role_name = 'compliance_operator')
FROM users
WHERE username = 'dmitry';

INSERT INTO stores(owner_id, store_name, legal_name, contact_email, contact_phone, store_address, verification_status, created_at) VALUES
    (1, 'North Coffee', 'North Coffee LLC', 'merchant1@example.com', '+20000000001', 'Berlin, Mitte', 'approved', now() - interval '60 days'),
    (2, 'Pixel Market', 'Pixel Market GmbH', 'merchant2@example.com', '+20000000002', 'Prague, Old Town', 'pending', now() - interval '50 days'),
    (4, 'Metro Books', 'Metro Books Ltd', 'merchant3@example.com', '+20000000003', 'London, Soho', 'approved', now() - interval '40 days'),
    (6, 'Orbit Electronics', 'Orbit Electronics FZ', 'merchant4@example.com', '+20000000004', 'Dubai Marina', 'approved', now() - interval '35 days'),
    (9, 'Fresh Deli', 'Fresh Deli OOO', 'merchant5@example.com', '+20000000005', 'Tbilisi Center', 'rejected', now() - interval '30 days');

INSERT INTO merchant_webhooks(store_id, url, event_types, secret_hash, is_active)
VALUES
    (1, 'https://webhooks.example.com/depay/north-coffee', ARRAY['transaction.created', 'transaction.broadcasted', 'transaction.confirmed', 'transaction.failed'], crypt('north-coffee-webhook-secret', gen_salt('bf')), true),
    (3, 'https://webhooks.example.com/depay/metro-books', ARRAY['transaction.confirmed', 'transaction.failed'], crypt('metro-books-webhook-secret', gen_salt('bf')), true);

INSERT INTO kyc_applications(user_id, status, reviewer_id, rejection_reason, submitted_at, reviewed_at)
SELECT
    u.user_id,
    CASE u.kyc_status
        WHEN 'approved' THEN 'approved'::kyc_status_enum
        WHEN 'rejected' THEN 'rejected'::kyc_status_enum
        WHEN 'pending' THEN 'pending'::kyc_status_enum
        ELSE 'pending'::kyc_status_enum
    END,
    CASE WHEN u.kyc_status IN ('approved', 'rejected') THEN 4 ELSE NULL END,
    CASE WHEN u.kyc_status = 'rejected' THEN 'Document photo is unreadable' ELSE NULL END,
    u.created_at + interval '2 days',
    CASE WHEN u.kyc_status IN ('approved', 'rejected') THEN u.created_at + interval '4 days' ELSE NULL END
FROM users u
WHERE u.kyc_status <> 'not_submitted'
UNION ALL
SELECT user_id, 'pending', NULL, NULL, created_at + interval '1 day', NULL
FROM users
WHERE username IN ('haruto', 'jon');

INSERT INTO kyc_documents(kyc_application_id, document_type, document_url)
SELECT kyc_application_id, 'passport', 'https://example.com/kyc/' || kyc_application_id || '/passport.png'
FROM kyc_applications;

INSERT INTO kyc_documents(kyc_application_id, document_type, document_url)
SELECT kyc_application_id, 'selfie', 'https://example.com/kyc/' || kyc_application_id || '/selfie.png'
FROM kyc_applications
WHERE kyc_application_id % 2 = 0;

INSERT INTO merchant_verification_applications(store_id, status, reviewer_id, rejection_reason, submitted_at, reviewed_at)
SELECT
    store_id,
    verification_status,
    CASE WHEN verification_status IN ('approved', 'rejected') THEN 4 ELSE NULL END,
    CASE WHEN verification_status = 'rejected' THEN 'Legal address mismatch' ELSE NULL END,
    created_at + interval '1 day',
    CASE WHEN verification_status IN ('approved', 'rejected') THEN created_at + interval '3 days' ELSE NULL END
FROM stores;

INSERT INTO blockchains(chain_name, native_symbol, chain_type, explorer_url) VALUES
    ('Ethereum', 'ETH', 'evm', 'https://etherscan.io'),
    ('BSC', 'BNB', 'evm', 'https://bscscan.com'),
    ('Polygon', 'MATIC', 'evm', 'https://polygonscan.com'),
    ('Solana Mock', 'SOL', 'solana', 'https://explorer.solana.com');

INSERT INTO assets(chain_id, symbol, asset_name, contract_address, decimals, is_native) VALUES
    (1, 'ETH', 'Ether', NULL, 18, TRUE),
    (1, 'USDT', 'Tether USD', '0xdac17f958d2ee523a2206206994597c13d831ec7', 6, FALSE),
    (2, 'BNB', 'BNB', NULL, 18, TRUE),
    (2, 'USDT', 'Tether USD', '0x55d398326f99059ff775485246999027b3197955', 18, FALSE),
    (3, 'MATIC', 'Matic', NULL, 18, TRUE),
    (3, 'USDC', 'USD Coin', '0x2791bca1f2de4661ed88a30c99a7a9449aa84174', 6, FALSE),
    (4, 'SOL', 'Solana', NULL, 9, TRUE),
    (4, 'USDC', 'USD Coin', 'mock-sol-usdc', 6, FALSE);

INSERT INTO wallets(user_id, chain_id, wallet_address, public_key, is_store_wallet, wallet_label)
SELECT
    ((g - 1) % 10) + 1,
    ((g - 1) % 4) + 1,
    '0x' || lpad(to_hex(g), 40, '0'),
    'pub-user-' || g,
    FALSE,
    'User wallet ' || g
FROM generate_series(1, 15) AS g;

INSERT INTO wallets(store_id, chain_id, wallet_address, public_key, is_store_wallet, wallet_label)
SELECT
    g,
    ((g - 1) % 4) + 1,
    '0x' || lpad(to_hex(100 + g), 40, '0'),
    'pub-store-' || g,
    TRUE,
    'Store wallet ' || g
FROM generate_series(1, 5) AS g;

INSERT INTO balances(wallet_id, asset_id, balance)
SELECT
    w.wallet_id,
    a.asset_id,
    (10 + ((w.wallet_id * a.asset_id) % 75))::NUMERIC(30,8)
FROM wallets w
JOIN assets a ON a.chain_id = w.chain_id;

INSERT INTO payment_invoices(store_id, user_id, external_order_id, amount_usdt, status, expires_at, paid_at, created_at)
SELECT
    ((g - 1) % 5) + 1,
    ((g - 1) % 10) + 1,
    'ORD-' || lpad(g::TEXT, 4, '0'),
    (25 + g * 12.5)::NUMERIC(30,8),
    CASE
        WHEN g % 6 = 0 THEN 'expired'::invoice_status_enum
        WHEN g % 5 = 0 THEN 'cancelled'::invoice_status_enum
        WHEN g % 3 = 0 THEN 'paid'::invoice_status_enum
        ELSE 'issued'::invoice_status_enum
    END,
    now() + interval '14 days' - (g || ' hours')::INTERVAL,
    CASE WHEN g % 3 = 0 THEN now() - (g || ' hours')::INTERVAL ELSE NULL END,
    now() - (g || ' days')::INTERVAL
FROM generate_series(1, 20) AS g;

INSERT INTO payment_invoice_items(invoice_id, item_name, quantity, unit_price_usdt)
SELECT invoice_id, 'Coursework item ' || invoice_id, (invoice_id % 3) + 1, round((amount_usdt / ((invoice_id % 3) + 1))::NUMERIC, 8)
FROM payment_invoices;

INSERT INTO terminals(store_id, serial_number, secret_hash, status, last_seen_at)
SELECT store_id, 'TERM-' || lpad(store_id::TEXT, 4, '0'), crypt('terminal-' || store_id, gen_salt('bf')), 'active', now() - (store_id || ' minutes')::INTERVAL
FROM stores;

INSERT INTO nfc_sessions(terminal_id, invoice_id, status, session_token_hash, expires_at, completed_at, created_at)
SELECT
    ((g - 1) % 5) + 1,
    g,
    CASE WHEN g % 4 = 0 THEN 'completed'::nfc_session_status_enum ELSE 'active'::nfc_session_status_enum END,
    encode(digest('session-' || g, 'sha256'), 'hex'),
    now() + interval '30 minutes',
    CASE WHEN g % 4 = 0 THEN now() - (g || ' minutes')::INTERVAL ELSE NULL END,
    now() - (g || ' hours')::INTERVAL
FROM generate_series(1, 8) AS g;

INSERT INTO payment_transactions(
    external_transaction_id,
    invoice_id,
    nfc_session_id,
    user_id,
    store_id,
    asset_id,
    user_wallet_id,
    store_wallet_id,
    amount,
    amount_in_usdt,
    network_fee,
    service_fee,
    blockchain_tx_hash,
    status,
    failure_reason,
    signed_payload,
    created_at,
    submitted_at,
    validated_at,
    broadcasted_at,
    completed_at
)
SELECT
    'seed-tx-' || g,
    ((g - 1) % 20) + 1,
    CASE WHEN g <= 8 THEN g ELSE NULL END,
    ((g - 1) % 10) + 1,
    ((g - 1) % 5) + 1,
    ((g - 1) % 8) + 1,
    ((g - 1) % 15) + 1,
    15 + (((g - 1) % 5) + 1),
    (0.01 * g)::NUMERIC(30,8),
    CASE WHEN g IN (7, 14, 21, 28, 35) THEN (1200 + g)::NUMERIC(30,8) ELSE (25 + g * 8)::NUMERIC(30,8) END,
    (0.0002 * g)::NUMERIC(30,8),
    (0.15 + g * 0.01)::NUMERIC(30,8),
    CASE WHEN g % 5 IN (0, 1) THEN '0x' || lpad(to_hex(1000 + g), 64, '0') ELSE NULL END,
    CASE
        WHEN g % 10 = 0 THEN 'failed'::transaction_status_enum
        WHEN g % 9 = 0 THEN 'cancelled'::transaction_status_enum
        WHEN g % 5 = 0 THEN 'confirmed'::transaction_status_enum
        WHEN g % 4 = 0 THEN 'broadcasted'::transaction_status_enum
        WHEN g % 3 = 0 THEN 'validated'::transaction_status_enum
        WHEN g % 2 = 0 THEN 'submitted'::transaction_status_enum
        ELSE 'created'::transaction_status_enum
    END,
    CASE WHEN g % 10 = 0 THEN 'insufficient funds' WHEN g % 9 = 0 THEN 'user cancelled' ELSE NULL END,
    jsonb_build_object('signature', 'sig-' || g, 'nonce', g),
    now() - (g || ' hours')::INTERVAL,
    CASE WHEN g % 2 = 0 THEN now() - (g || ' hours')::INTERVAL + interval '2 minutes' ELSE NULL END,
    CASE WHEN g % 3 = 0 OR g % 4 = 0 OR g % 5 = 0 THEN now() - (g || ' hours')::INTERVAL + interval '4 minutes' ELSE NULL END,
    CASE WHEN g % 4 = 0 OR g % 5 = 0 THEN now() - (g || ' hours')::INTERVAL + interval '6 minutes' ELSE NULL END,
    CASE WHEN g % 5 = 0 OR g % 9 = 0 OR g % 10 = 0 THEN now() - (g || ' hours')::INTERVAL + interval '9 minutes' ELSE NULL END
FROM generate_series(1, 80) AS g;

INSERT INTO rpc_nodes(chain_id, node_name, rpc_url, status, is_primary)
SELECT
    c.chain_id,
    c.chain_name || ' RPC ' || n,
    'https://rpc.example.com/' || lower(replace(c.chain_name, ' ', '-')) || '/' || n,
    CASE WHEN n = 1 THEN 'healthy'::rpc_node_status_enum ELSE 'degraded'::rpc_node_status_enum END,
    n = 1
FROM blockchains c
CROSS JOIN generate_series(1, 2) AS n;

INSERT INTO rpc_node_checks(rpc_node_id, status, latency_ms, block_height, error_message, checked_at)
SELECT
    rn.rpc_node_id,
    CASE WHEN g % 11 = 0 THEN 'down'::rpc_node_status_enum WHEN g % 5 = 0 THEN 'degraded'::rpc_node_status_enum ELSE 'healthy'::rpc_node_status_enum END,
    CASE WHEN g % 11 = 0 THEN NULL ELSE 80 + ((rn.rpc_node_id * g) % 340) END,
    1000000 + rn.rpc_node_id * 1000 + g,
    CASE WHEN g % 11 = 0 THEN 'timeout' ELSE NULL END,
    now() - (g || ' hours')::INTERVAL
FROM rpc_nodes rn
CROSS JOIN generate_series(1, 12) AS g;

INSERT INTO exchange_rates(asset_id, rate_to_usdt, captured_at)
SELECT
    a.asset_id,
    CASE a.symbol
        WHEN 'ETH' THEN 3100
        WHEN 'BNB' THEN 620
        WHEN 'MATIC' THEN 0.9
        WHEN 'SOL' THEN 145
        ELSE 1
    END + (g * 0.01),
    now() - (g || ' days')::INTERVAL
FROM assets a
CROSS JOIN generate_series(1, 5) AS g;

INSERT INTO blacklisted_wallets(chain_id, wallet_address, reason, created_by)
VALUES
    (1, '0x0000000000000000000000000000000000000bad', 'Known phishing wallet', 4),
    (3, '0x0000000000000000000000000000000000000cab', 'Chargeback fraud pattern', 4);

INSERT INTO risk_alerts(user_id, store_id, transaction_id, alert_type, risk_level, details)
SELECT user_id, store_id, transaction_id, 'MANY_FAILED_TRANSACTIONS', 'medium', jsonb_build_object('source', 'seed')
FROM payment_transactions
WHERE status = 'failed'
LIMIT 3;

COMMIT;
