-- Revert artifacts:artifacts_job_id from pg

BEGIN;

ALTER TABLE artifacts_v2.artifacts
DROP COLUMN job_id;

COMMIT;
