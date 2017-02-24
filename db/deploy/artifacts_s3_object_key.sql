-- Deploy artifacts:artifacts_s3_object_key to pg
-- requires: artifacts
-- requires: appschema

BEGIN;

ALTER TABLE artifacts_v2.artifacts
DROP COLUMN s3_object_key;

COMMIT;
