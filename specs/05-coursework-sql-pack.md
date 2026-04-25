# 05. Coursework SQL Pack

## 12 SQL-функций

Каждый запрос должен быть реализован как PostgreSQL function и показан в записке с формулировкой, кодом и скриншотом результата.

1. `get_user_kyc_wallet_summary()` — статус KYC, последняя проверка, число кошельков по сетям.
2. `get_user_wallet_balances(p_user_id BIGINT)` — кошельки пользователя и балансы по активам.
3. `get_wallet_asset_distribution(p_wallet_id BIGINT)` — активы кошелька и доля каждого актива.
4. `get_transaction_card(p_tx_id BIGINT)` — карточка транзакции с пользователем, магазином, активом и длительностями.
5. `get_user_transaction_history(p_user_id, p_date_from, p_date_to)` — история пользователя по магазинам и статусам.
6. `get_store_transaction_history(p_store_id, p_date_from, p_date_to)` — история магазина по пользователям, активам, статусам.
7. `get_blockchain_asset_activity(p_chain_id, p_date_from, p_date_to)` — активность активов сети.
8. `get_rpc_nodes_activity(p_chain_id, p_date_from, p_date_to)` — RPC-ноды, средняя задержка, проверки.
9. `get_store_turnover(p_store_id, p_date_from, p_date_to)` — оборот магазина по активам и статусам.
10. `get_store_success_rate(p_store_id, p_date_from, p_date_to)` — доля успешных транзакций и среднее время.
11. `get_unverified_active_users(p_min_tx_count, p_min_amount_usdt, p_date_from, p_date_to)` — активные пользователи без KYC.
12. `get_failed_transactions_analytics(p_date_from, p_date_to)` — ошибки по сетям, активам, причинам.

## Процедуры бизнес-flow

- `create_payment_invoice`;
- `create_nfc_session`;
- `submit_signed_transaction`;
- `approve_kyc_application`;
- `reject_kyc_application`;
- `approve_merchant_verification`;
- `reject_merchant_verification`.

## Триггеры

Сделать минимум 7:

1. `trg_wallet_owner_check` — кошелек не может быть без владельца или с двумя владельцами.
2. `trg_balance_non_negative` — запрет отрицательного баланса.
3. `trg_balance_history` — запись истории изменения баланса.
4. `trg_transaction_status_flow` — контроль переходов статусов.
5. `trg_transaction_completed_at` — автозаполнение `completed_at`.
6. `trg_audit_payment_transactions` — аудит INSERT/UPDATE/DELETE.
7. `trg_risk_alert_large_unverified_payment` — risk-alert для крупного платежа без KYC.

## Роли PostgreSQL

- `depay_admin`;
- `depay_app`;
- `depay_compliance`;
- `depay_merchant_readonly`;
- `depay_user_readonly`.

## Views для интерфейса

- `vw_user_wallet_balances`;
- `vw_store_transactions`;
- `vw_failed_transactions`;
- `vw_rpc_node_status`;
- `vw_compliance_kyc_queue`.

## Скриншоты

Нужно подготовить:

- DBeaver-схему;
- заполненные таблицы;
- результаты всех функций;
- код функций;
- код триггеров;
- демонстрацию ошибки триггера;
- графики;
- веб-интерфейс таблиц и функций.
