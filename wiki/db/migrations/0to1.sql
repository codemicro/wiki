-- Migrate a new database to v1 format.

CREATE TABLE "users"
(
    id          TEXT PRIMARY KEY,
    external_id TEXT UNIQUE,
    name        TEXT,
    email       TEXT NOT NULL UNIQUE
);

CREATE TABLE "session_key"
(
    key TEXT
);

CREATE TABLE "tags"
(
    id   TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE "pages"
(
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    content TEXT NOT NULL
);

CREATE TABLE "version"
(
    version INT
);

INSERT INTO "version"(version)
VALUES (1);