DROP VIEW IF EXISTS vw_compliance_kyc_queue;
DROP VIEW IF EXISTS vw_rpc_node_status;
DROP VIEW IF EXISTS vw_failed_transactions;
DROP VIEW IF EXISTS vw_store_transactions;
DROP VIEW IF EXISTS vw_user_wallet_balances;

DROP FUNCTION IF EXISTS get_failed_transactions_analytics(TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_unverified_active_users(BIGINT, NUMERIC, TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_store_success_rate(BIGINT, TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_store_turnover(BIGINT, TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_rpc_nodes_activity(BIGINT, TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_blockchain_asset_activity(BIGINT, TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_store_transaction_history(BIGINT, TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_user_transaction_history(BIGINT, TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_transaction_card(BIGINT);
DROP FUNCTION IF EXISTS get_wallet_asset_distribution(BIGINT);
DROP FUNCTION IF EXISTS get_user_wallet_balances(BIGINT);
DROP FUNCTION IF EXISTS get_user_kyc_wallet_summary();
