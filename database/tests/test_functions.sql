\set ON_ERROR_STOP on

BEGIN;

CREATE TEMP TABLE function_test_results (
    test_name TEXT PRIMARY KEY,
    rows_returned BIGINT NOT NULL
) ON COMMIT DROP;

INSERT INTO function_test_results
SELECT 'get_user_kyc_wallet_summary', count(*) FROM get_user_kyc_wallet_summary();

INSERT INTO function_test_results
SELECT 'get_user_wallet_balances', count(*) FROM get_user_wallet_balances(1);

INSERT INTO function_test_results
SELECT 'get_wallet_asset_distribution', count(*) FROM get_wallet_asset_distribution(1);

INSERT INTO function_test_results
SELECT 'get_transaction_card', count(*) FROM get_transaction_card(1);

INSERT INTO function_test_results
SELECT 'get_user_transaction_history', count(*) FROM get_user_transaction_history(1, now() - interval '30 days', now() + interval '1 day');

INSERT INTO function_test_results
SELECT 'get_store_transaction_history', count(*) FROM get_store_transaction_history(1, now() - interval '30 days', now() + interval '1 day');

INSERT INTO function_test_results
SELECT 'get_blockchain_asset_activity', count(*) FROM get_blockchain_asset_activity(1, now() - interval '30 days', now() + interval '1 day');

INSERT INTO function_test_results
SELECT 'get_rpc_nodes_activity', count(*) FROM get_rpc_nodes_activity(1, now() - interval '30 days', now() + interval '1 day');

INSERT INTO function_test_results
SELECT 'get_store_turnover', count(*) FROM get_store_turnover(1, now() - interval '30 days', now() + interval '1 day');

INSERT INTO function_test_results
SELECT 'get_store_success_rate', count(*) FROM get_store_success_rate(1, now() - interval '30 days', now() + interval '1 day');

INSERT INTO function_test_results
SELECT 'get_unverified_active_users', count(*) FROM get_unverified_active_users(1, 1, now() - interval '30 days', now() + interval '1 day');

INSERT INTO function_test_results
SELECT 'get_failed_transactions_analytics', count(*) FROM get_failed_transactions_analytics(now() - interval '30 days', now() + interval '1 day');

DO $$
DECLARE
    failed_tests TEXT;
BEGIN
    SELECT string_agg(test_name, ', ' ORDER BY test_name)
    INTO failed_tests
    FROM function_test_results
    WHERE rows_returned = 0;

    IF failed_tests IS NOT NULL THEN
        RAISE EXCEPTION 'SQL functions returned no rows: %', failed_tests;
    END IF;
END;
$$;

SELECT * FROM function_test_results ORDER BY test_name;

ROLLBACK;
