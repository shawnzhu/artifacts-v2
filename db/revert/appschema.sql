-- Revert artifacts:appschema from pg

BEGIN;

DROP schema artifacts_v2;

COMMIT;
