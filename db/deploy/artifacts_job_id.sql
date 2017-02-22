-- Deploy artifacts:artifacts_job_id to pg
-- requires: artifacts
-- requires: appschema

BEGIN;

ALTER TABLE artifacts_v2.artifacts
ADD job_id VARCHAR;

COMMIT;
