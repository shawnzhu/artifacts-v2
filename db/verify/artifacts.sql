-- Verify artifacts:artifacts on pg

BEGIN;

SELECT build_id, path
    FROM artifacts_v2.artifacts
    WHERE FALSE;

ROLLBACK;
