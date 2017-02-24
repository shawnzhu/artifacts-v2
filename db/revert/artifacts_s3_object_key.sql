-- Revert artifacts:artifacts_s3_object_key from pg

BEGIN;

ALTER TABLE artifacts_v2.artifacts
ADD s3_object_key VARCHAR(1024);

COMMIT;
