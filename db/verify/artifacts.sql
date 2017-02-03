-- Verify artifacts:artifacts on pg

BEGIN;

SELECT build_id, path, s3_object_key
    FROM artifacts_v2.artifacts
    WHERE FALSE;

ROLLBACK;
