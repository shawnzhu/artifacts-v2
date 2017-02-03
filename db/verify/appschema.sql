-- Verify artifacts:appschema on pg

BEGIN;

SELECT 1/COUNT(*) FROM information_schema.schemata WHERE schema_name = 'artifacts_v2';
SELECT pg_catalog.has_schema_privilege('artifacts_v2', 'usage');

ROLLBACK;
