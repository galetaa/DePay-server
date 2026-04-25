# 09. Security and Compliance

## Принципы

1. Не хранить приватные ключи.
2. Пароли хранить только как hash.
3. Refresh tokens хранить в БД в виде hash.
4. Все важные изменения писать в audit.
5. Крупные платежи без KYC отправлять в risk alerts.
6. Доступ разграничить ролями.

## Auth

Access token claims:

```json
{
  "sub": "user_id",
  "roles": ["user"],
  "type": "access",
  "exp": 1234567890
}
```

Refresh token table:

- `token_id`;
- `user_id`;
- `token_hash`;
- `created_at`;
- `expires_at`;
- `revoked_at`;
- `ip_address`;
- `user_agent`.

## KYC workflow

1. User submits KYC.
2. Создается `kyc_applications`.
3. User status -> `pending`.
4. Compliance checks.
5. Approve -> user status `approved`.
6. Reject -> user status `rejected`, сохраняется причина.

## Merchant verification

Магазин не может создавать invoice и terminal, пока `verification_status != approved`.

## Risk rules

- `LARGE_UNVERIFIED_PAYMENT`;
- `BLACKLISTED_WALLET`;
- `MANY_FAILED_TRANSACTIONS`;
- `HIGH_AMOUNT_NEW_USER`;
- `RPC_FAILURE_SPIKE`.

## Audit

Аудировать:

- изменение KYC;
- изменение merchant verification;
- создание invoice;
- изменение transaction status;
- изменение balances;
- добавление blacklist wallet;
- изменение ролей.

## Секреты

- убрать реальные приватные ключи из git;
- добавить `.env.example`;
- `.env` не коммитить;
- hardcoded RabbitMQ/Redis credentials вынести в env.
