-- db schema

CREATE TABLE artifacts (
	artifact_id SERIAL PRIMARY KEY,
	build_id VARCHAR,
	path VARCHAR(512),
	s3_object_key VARCHAR(1024)
);
