CREATE TABLE merchant_api_keys (
    api_key_id BIGSERIAL PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES stores(store_id) ON DELETE CASCADE,
    key_name VARCHAR(120) NOT NULL,
    key_prefix VARCHAR(32) NOT NULL,
    secret_hash VARCHAR(255) NOT NULL,
    scopes TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    last_used_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (store_id, key_name),
    UNIQUE (key_prefix),
    CONSTRAINT chk_merchant_api_keys_scopes CHECK (
        scopes <@ ARRAY['invoice:read', 'invoice:write', 'transaction:read', 'webhook:write']::TEXT[]
    )
);

CREATE INDEX idx_merchant_api_keys_store_id ON merchant_api_keys(store_id);
CREATE INDEX idx_merchant_api_keys_active ON merchant_api_keys(store_id, revoked_at);

GRANT ALL PRIVILEGES ON merchant_api_keys TO depay_admin;
GRANT ALL PRIVILEGES ON merchant_api_keys_api_key_id_seq TO depay_admin;
GRANT SELECT, INSERT, UPDATE ON merchant_api_keys TO depay_app;
GRANT USAGE, SELECT ON merchant_api_keys_api_key_id_seq TO depay_app;
GRANT SELECT ON merchant_api_keys TO depay_merchant_readonly;

INSERT INTO roles(role_name)
VALUES ('compliance'), ('service')
ON CONFLICT DO NOTHING;
