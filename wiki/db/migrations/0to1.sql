-- Migrate a new database to v1 format.

CREATE TABLE "users"
(
    id          TEXT(36) PRIMARY KEY,
    external_id TEXT UNIQUE,
    name        TEXT
);

CREATE TABLE "session_key"
(
    key TEXT
);

CREATE TABLE "version"
(
    version INT
);

INSERT INTO "version"(version)
VALUES (1);