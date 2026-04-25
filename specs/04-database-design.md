# 04. Database Design

## Группы таблиц

1. Auth/users.
2. Stores/merchants.
3. KYC/verification.
4. Blockchains/assets.
5. Wallets/balances.
6. Invoices/NFC sessions.
7. Payment transactions.
8. RPC/gas monitoring.
9. Risk/audit.

## Основные таблицы

### `users`

Поля:

- `user_id BIGSERIAL PK`;
- `username VARCHAR(80) UNIQUE NOT NULL`;
- `email VARCHAR(255) UNIQUE`;
- `phone_number VARCHAR(20) UNIQUE`;
- `password_hash VARCHAR(255) NOT NULL`;
- `last_name`, `first_name`, `middle_name`;
- `birth_date DATE`;
- `address TEXT`;
- `kyc_status kyc_status_enum DEFAULT 'not_submitted'`;
- `created_at`, `updated_at`.

Ограничения: email или phone обязательны; дата рождения меньше текущей даты.

### `roles`, `user_roles`

Роли: user, merchant, compliance_operator, admin, terminal.

### `stores`

Магазины:

- `store_id`;
- `owner_id`;
- `store_name`;
- `legal_name`;
- `contact_email`;
- `contact_phone`;
- `store_address`;
- `verification_status`;
- timestamps.

### `kyc_applications`, `kyc_documents`

История KYC-заявок и документов. Не хранить только один статус в `users`, нужна история проверок.

### `merchant_verification_applications`

История верификации магазинов.

### `blockchains`

Справочник сетей: Ethereum, BSC, Polygon, Solana mock и т.д.

### `assets`

Активы внутри сетей: ETH, USDT, USDC. Для токенов нужен `contract_address`.

### `wallets`

Кошельки пользователей и магазинов:

- `user_id` nullable;
- `store_id` nullable;
- `chain_id`;
- `wallet_address`;
- `public_key`;
- `is_store_wallet`.

Правило: кошелек принадлежит либо user, либо store, но не обоим.

### `balances`

Текущие остатки по активам. PK `(wallet_id, asset_id)`. Баланс неотрицательный.

### `balance_history`

История изменений баланса. Заполняется триггером.

### `payment_invoices`

Платежные счета магазина:

- магазин;
- пользователь optional;
- external_order_id;
- amount_usdt;
- status;
- expires_at;
- paid_at.

### `payment_invoice_items`

Позиции счета.

### `terminals`

NFC-терминалы магазина: serial number, secret hash, status, last seen.

### `nfc_sessions`

Временные платежные сессии терминала.

### `payment_transactions`

Платежи:

- invoice/session;
- user/store;
- asset;
- user wallet/store wallet;
- amount;
- amount_in_usdt;
- fees;
- blockchain_tx_hash;
- status;
- signatures;
- timestamps.

### `rpc_nodes`, `rpc_node_checks`

RPC-ноды и история проверок.

### `exchange_rates`

История курсов активов.

### `risk_alerts`, `blacklisted_wallets`

Антифрод.

### `audit_logs`

Аудит действий и изменений.

## Enum-типы

- `kyc_status_enum`;
- `verification_status_enum`;
- `transaction_status_enum`;
- `invoice_status_enum`;
- `nfc_session_status_enum`;
- `terminal_status_enum`;
- `rpc_node_status_enum`;
- `risk_level_enum`;
- `risk_alert_status_enum`;
- `document_type_enum`.

## Нормализация

Схема соответствует 3НФ:

- справочники вынесены отдельно;
- KYC-заявки отделены от пользователей;
- документы KYC вынесены отдельно;
- кошельки отделены от балансов;
- счета отделены от транзакций;
- история и аудит отделены от текущих данных.
