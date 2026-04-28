# Security Notes

DePay is a portfolio/demo payment platform, not a production payment processor. The local stack is intentionally easy to run, while production-like controls are documented and partially enforced.

## Authentication

- Access tokens are JWTs signed with `JWT_SECRET`.
- User access tokens are short-lived.
- User refresh tokens are opaque random strings.
- Refresh tokens are stored only as hashes in `refresh_tokens`.
- Refresh rotates the refresh token and revokes the old hash.
- Logout revokes the supplied refresh token.

## Roles

Supported roles are:

- `user`
- `merchant`
- `compliance`
- `admin`
- `service`

The legacy seed role `compliance_operator` is accepted as `compliance` by shared role middleware.

## Endpoint Policy

- Admin table/function endpoints are gated by `ProductionRoleGate("admin")` when `APP_ENV=production` or `REQUIRE_AUTH=true`.
- Demo endpoints are blocked in production unless `DEMO_ENDPOINTS_ENABLED=true`.
- Merchant repository queries are scoped by `store_id`.
- User profile/KYC routes resolve the user from JWT claims.

## Merchant API Keys

Merchant API keys are for server-to-server integrations.

- The raw secret is returned only on creation.
- The database stores `secret_hash`, `key_prefix`, scopes, `revoked_at`, and `last_used_at`.
- Supported scopes: `invoice:read`, `invoice:write`, `transaction:read`, `webhook:write`.
- Revoked keys remain auditable but cannot be used by future API-key auth.

## Private Key Policy

DePay does not store private keys, seed phrases, or raw wallet secrets. It may store public wallet addresses, public keys, signed transaction payloads, transaction hashes, and signature metadata.

## Production Caveats

- Local provider defaults are mock/demo implementations.
- Dev JWT and Vault/Kong files are local examples and must be replaced with managed secrets before real deployment.
- Demo endpoints and permissive local admin routes must stay disabled in production.
