CREATE OR REPLACE FUNCTION get_user_kyc_wallet_summary()
RETURNS TABLE (
    user_id BIGINT,
    username VARCHAR,
    email VARCHAR,
    kyc_status kyc_status_enum,
    last_kyc_reviewed_at TIMESTAMPTZ,
    wallet_count BIGINT,
    chains TEXT
) AS $$
    SELECT
        u.user_id,
        u.username,
        u.email,
        u.kyc_status,
        max(ka.reviewed_at) AS last_kyc_reviewed_at,
        count(DISTINCT w.wallet_id) AS wallet_count,
        string_agg(DISTINCT b.chain_name, ', ' ORDER BY b.chain_name) AS chains
    FROM users u
    LEFT JOIN kyc_applications ka ON ka.user_id = u.user_id
    LEFT JOIN wallets w ON w.user_id = u.user_id
    LEFT JOIN blockchains b ON b.chain_id = w.chain_id
    GROUP BY u.user_id, u.username, u.email, u.kyc_status
    ORDER BY u.user_id;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_user_wallet_balances(p_user_id BIGINT)
RETURNS TABLE (
    wallet_id BIGINT,
    wallet_address VARCHAR,
    chain_name VARCHAR,
    asset_symbol VARCHAR,
    balance NUMERIC,
    updated_at TIMESTAMPTZ
) AS $$
    SELECT
        w.wallet_id,
        w.wallet_address,
        bc.chain_name,
        a.symbol,
        COALESCE(b.balance, 0),
        b.updated_at
    FROM wallets w
    JOIN blockchains bc ON bc.chain_id = w.chain_id
    LEFT JOIN balances b ON b.wallet_id = w.wallet_id
    LEFT JOIN assets a ON a.asset_id = b.asset_id
    WHERE w.user_id = p_user_id
    ORDER BY w.wallet_id, a.symbol;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_wallet_asset_distribution(p_wallet_id BIGINT)
RETURNS TABLE (
    wallet_id BIGINT,
    asset_symbol VARCHAR,
    balance NUMERIC,
    balance_share NUMERIC
) AS $$
    SELECT
        b.wallet_id,
        a.symbol,
        b.balance,
        round((b.balance / NULLIF(sum(b.balance) OVER (PARTITION BY b.wallet_id), 0)) * 100, 2) AS balance_share
    FROM balances b
    JOIN assets a ON a.asset_id = b.asset_id
    WHERE b.wallet_id = p_wallet_id
    ORDER BY 4 DESC NULLS LAST;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_transaction_card(p_tx_id BIGINT)
RETURNS TABLE (
    transaction_id BIGINT,
    username VARCHAR,
    store_name VARCHAR,
    chain_name VARCHAR,
    asset_symbol VARCHAR,
    amount NUMERIC,
    amount_in_usdt NUMERIC,
    status transaction_status_enum,
    created_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    duration_seconds NUMERIC
) AS $$
    SELECT
        pt.transaction_id,
        u.username,
        s.store_name,
        bc.chain_name,
        a.symbol,
        pt.amount,
        pt.amount_in_usdt,
        pt.status,
        pt.created_at,
        pt.completed_at,
        extract(epoch FROM COALESCE(pt.completed_at, now()) - pt.created_at) AS duration_seconds
    FROM payment_transactions pt
    JOIN users u ON u.user_id = pt.user_id
    JOIN stores s ON s.store_id = pt.store_id
    JOIN assets a ON a.asset_id = pt.asset_id
    JOIN blockchains bc ON bc.chain_id = a.chain_id
    WHERE pt.transaction_id = p_tx_id;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_user_transaction_history(p_user_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)
RETURNS TABLE (
    transaction_id BIGINT,
    store_name VARCHAR,
    asset_symbol VARCHAR,
    amount NUMERIC,
    amount_in_usdt NUMERIC,
    status transaction_status_enum,
    created_at TIMESTAMPTZ
) AS $$
    SELECT
        pt.transaction_id,
        s.store_name,
        a.symbol,
        pt.amount,
        pt.amount_in_usdt,
        pt.status,
        pt.created_at
    FROM payment_transactions pt
    JOIN stores s ON s.store_id = pt.store_id
    JOIN assets a ON a.asset_id = pt.asset_id
    WHERE pt.user_id = p_user_id
      AND pt.created_at >= p_date_from
      AND pt.created_at < p_date_to
    ORDER BY pt.created_at DESC;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_store_transaction_history(p_store_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)
RETURNS TABLE (
    transaction_id BIGINT,
    username VARCHAR,
    asset_symbol VARCHAR,
    amount NUMERIC,
    amount_in_usdt NUMERIC,
    status transaction_status_enum,
    created_at TIMESTAMPTZ
) AS $$
    SELECT
        pt.transaction_id,
        u.username,
        a.symbol,
        pt.amount,
        pt.amount_in_usdt,
        pt.status,
        pt.created_at
    FROM payment_transactions pt
    JOIN users u ON u.user_id = pt.user_id
    JOIN assets a ON a.asset_id = pt.asset_id
    WHERE pt.store_id = p_store_id
      AND pt.created_at >= p_date_from
      AND pt.created_at < p_date_to
    ORDER BY pt.created_at DESC;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_blockchain_asset_activity(p_chain_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)
RETURNS TABLE (
    chain_name VARCHAR,
    asset_symbol VARCHAR,
    transaction_count BIGINT,
    total_amount NUMERIC,
    total_amount_usdt NUMERIC,
    confirmed_count BIGINT
) AS $$
    SELECT
        bc.chain_name,
        a.symbol,
        count(pt.transaction_id),
        COALESCE(sum(pt.amount), 0),
        COALESCE(sum(pt.amount_in_usdt), 0),
        count(*) FILTER (WHERE pt.status = 'confirmed')
    FROM assets a
    JOIN blockchains bc ON bc.chain_id = a.chain_id
    LEFT JOIN payment_transactions pt ON pt.asset_id = a.asset_id
        AND pt.created_at >= p_date_from
        AND pt.created_at < p_date_to
    WHERE a.chain_id = p_chain_id
    GROUP BY bc.chain_name, a.symbol
    ORDER BY 5 DESC;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_rpc_nodes_activity(p_chain_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)
RETURNS TABLE (
    rpc_node_id BIGINT,
    node_name VARCHAR,
    current_status rpc_node_status_enum,
    check_count BIGINT,
    avg_latency_ms NUMERIC,
    degraded_or_down_count BIGINT
) AS $$
    SELECT
        rn.rpc_node_id,
        rn.node_name,
        rn.status,
        count(rnc.rpc_node_check_id),
        round(avg(rnc.latency_ms), 2),
        count(*) FILTER (WHERE rnc.status IN ('degraded', 'down'))
    FROM rpc_nodes rn
    LEFT JOIN rpc_node_checks rnc ON rnc.rpc_node_id = rn.rpc_node_id
        AND rnc.checked_at >= p_date_from
        AND rnc.checked_at < p_date_to
    WHERE rn.chain_id = p_chain_id
    GROUP BY rn.rpc_node_id, rn.node_name, rn.status
    ORDER BY 5 DESC NULLS LAST;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_store_turnover(p_store_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)
RETURNS TABLE (
    store_id BIGINT,
    store_name VARCHAR,
    asset_symbol VARCHAR,
    status transaction_status_enum,
    transaction_count BIGINT,
    total_amount_usdt NUMERIC
) AS $$
    SELECT
        s.store_id,
        s.store_name,
        a.symbol,
        pt.status,
        count(pt.transaction_id),
        COALESCE(sum(pt.amount_in_usdt), 0)
    FROM stores s
    JOIN payment_transactions pt ON pt.store_id = s.store_id
    JOIN assets a ON a.asset_id = pt.asset_id
    WHERE s.store_id = p_store_id
      AND pt.created_at >= p_date_from
      AND pt.created_at < p_date_to
    GROUP BY s.store_id, s.store_name, a.symbol, pt.status
    ORDER BY 6 DESC;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_store_success_rate(p_store_id BIGINT, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)
RETURNS TABLE (
    store_id BIGINT,
    store_name VARCHAR,
    total_transactions BIGINT,
    confirmed_transactions BIGINT,
    success_rate NUMERIC,
    avg_duration_seconds NUMERIC
) AS $$
    SELECT
        s.store_id,
        s.store_name,
        count(pt.transaction_id),
        count(*) FILTER (WHERE pt.status = 'confirmed'),
        round((count(*) FILTER (WHERE pt.status = 'confirmed')::NUMERIC / NULLIF(count(pt.transaction_id), 0)) * 100, 2),
        round((avg(extract(epoch FROM pt.completed_at - pt.created_at)) FILTER (WHERE pt.completed_at IS NOT NULL))::NUMERIC, 2)
    FROM stores s
    LEFT JOIN payment_transactions pt ON pt.store_id = s.store_id
        AND pt.created_at >= p_date_from
        AND pt.created_at < p_date_to
    WHERE s.store_id = p_store_id
    GROUP BY s.store_id, s.store_name;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_unverified_active_users(p_min_tx_count BIGINT, p_min_amount_usdt NUMERIC, p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)
RETURNS TABLE (
    user_id BIGINT,
    username VARCHAR,
    kyc_status kyc_status_enum,
    transaction_count BIGINT,
    total_amount_usdt NUMERIC
) AS $$
    SELECT
        u.user_id,
        u.username,
        u.kyc_status,
        count(pt.transaction_id),
        COALESCE(sum(pt.amount_in_usdt), 0)
    FROM users u
    JOIN payment_transactions pt ON pt.user_id = u.user_id
    WHERE u.kyc_status <> 'approved'
      AND pt.created_at >= p_date_from
      AND pt.created_at < p_date_to
    GROUP BY u.user_id, u.username, u.kyc_status
    HAVING count(pt.transaction_id) >= p_min_tx_count
       AND COALESCE(sum(pt.amount_in_usdt), 0) >= p_min_amount_usdt
    ORDER BY 5 DESC;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_failed_transactions_analytics(p_date_from TIMESTAMPTZ, p_date_to TIMESTAMPTZ)
RETURNS TABLE (
    chain_name VARCHAR,
    asset_symbol VARCHAR,
    failure_reason TEXT,
    failed_count BIGINT,
    total_failed_amount_usdt NUMERIC
) AS $$
    SELECT
        bc.chain_name,
        a.symbol,
        COALESCE(pt.failure_reason, 'unknown'),
        count(pt.transaction_id),
        COALESCE(sum(pt.amount_in_usdt), 0)
    FROM payment_transactions pt
    JOIN assets a ON a.asset_id = pt.asset_id
    JOIN blockchains bc ON bc.chain_id = a.chain_id
    WHERE pt.status = 'failed'
      AND pt.created_at >= p_date_from
      AND pt.created_at < p_date_to
    GROUP BY bc.chain_name, a.symbol, COALESCE(pt.failure_reason, 'unknown')
    ORDER BY 4 DESC;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE VIEW vw_user_wallet_balances AS
SELECT
    u.user_id,
    u.username,
    w.wallet_id,
    w.wallet_address,
    bc.chain_name,
    a.symbol AS asset_symbol,
    b.balance,
    b.updated_at
FROM users u
JOIN wallets w ON w.user_id = u.user_id
JOIN blockchains bc ON bc.chain_id = w.chain_id
LEFT JOIN balances b ON b.wallet_id = w.wallet_id
LEFT JOIN assets a ON a.asset_id = b.asset_id;

CREATE OR REPLACE VIEW vw_store_transactions AS
SELECT
    s.store_id,
    s.store_name,
    pt.transaction_id,
    u.username,
    a.symbol AS asset_symbol,
    pt.amount_in_usdt,
    pt.status,
    pt.created_at
FROM stores s
JOIN payment_transactions pt ON pt.store_id = s.store_id
JOIN users u ON u.user_id = pt.user_id
JOIN assets a ON a.asset_id = pt.asset_id;

CREATE OR REPLACE VIEW vw_failed_transactions AS
SELECT
    pt.transaction_id,
    bc.chain_name,
    a.symbol AS asset_symbol,
    pt.failure_reason,
    pt.amount_in_usdt,
    pt.created_at
FROM payment_transactions pt
JOIN assets a ON a.asset_id = pt.asset_id
JOIN blockchains bc ON bc.chain_id = a.chain_id
WHERE pt.status = 'failed';

CREATE OR REPLACE VIEW vw_rpc_node_status AS
SELECT
    rn.rpc_node_id,
    bc.chain_name,
    rn.node_name,
    rn.status,
    max(rnc.checked_at) AS last_checked_at,
    round(avg(rnc.latency_ms), 2) AS avg_latency_ms
FROM rpc_nodes rn
JOIN blockchains bc ON bc.chain_id = rn.chain_id
LEFT JOIN rpc_node_checks rnc ON rnc.rpc_node_id = rn.rpc_node_id
GROUP BY rn.rpc_node_id, bc.chain_name, rn.node_name, rn.status;

CREATE OR REPLACE VIEW vw_compliance_kyc_queue AS
SELECT
    ka.kyc_application_id,
    u.user_id,
    u.username,
    u.email,
    ka.status,
    ka.submitted_at,
    ka.reviewed_at
FROM kyc_applications ka
JOIN users u ON u.user_id = ka.user_id
WHERE ka.status = 'pending';
