-- Deploy artifacts:artifacts to pg
-- requires: appschema

BEGIN;

CREATE TABLE artifacts_v2.artifacts (
	artifact_id    SERIAL     PRIMARY KEY,
	build_id       VARCHAR,
	path           VARCHAR(512),
	s3_object_key  VARCHAR(1024)
);

COMMIT;
