DELETE FROM roles WHERE role_name IN ('compliance', 'service');

DROP TABLE IF EXISTS merchant_api_keys;
