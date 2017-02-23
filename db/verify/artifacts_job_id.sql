-- Verify artifacts:artifacts_job_id on pg

BEGIN;

SELECT job_id
    FROM artifacts_v2.artifacts
    WHERE FALSE;

ROLLBACK;
