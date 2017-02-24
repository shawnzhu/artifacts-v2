-- Verify artifacts:artifacts_s3_object_key on pg

BEGIN;

SELECT
    (SELECT count(*)
        FROM information_schema.columns
        WHERE table_schema = 'artifacts_v2' AND table_name = 'artifacts' AND column_name = 's3_object_key'
    ) = 0;

ROLLBACK;
