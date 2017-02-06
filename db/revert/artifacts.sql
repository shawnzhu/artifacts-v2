-- Revert artifacts:artifacts from pg

BEGIN;

DROP TABLE artifacts_v2.artifacts;

COMMIT;
