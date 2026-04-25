REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM depay_admin;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM depay_admin;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM depay_app;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM depay_app;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM depay_compliance;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM depay_merchant_readonly;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM depay_user_readonly;
REVOKE USAGE ON SCHEMA public FROM depay_admin, depay_app, depay_compliance, depay_merchant_readonly, depay_user_readonly;

DROP ROLE IF EXISTS depay_user_readonly;
DROP ROLE IF EXISTS depay_merchant_readonly;
DROP ROLE IF EXISTS depay_compliance;
DROP ROLE IF EXISTS depay_app;
DROP ROLE IF EXISTS depay_admin;
