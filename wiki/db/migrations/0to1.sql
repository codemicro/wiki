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
    id         TEXT PRIMARY KEY,
    title      TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    content    TEXT      NOT NULL
);

INSERT INTO "pages"("id", "title", "created_at", "updated_at", "content")
VALUES ('index', 'system:index', '2000-01-01 00:00:00.00+00:00', '2000-01-01 00:00:00.00+00:00',
        '# Wiki

Welcome to the Wiki!');

CREATE TABLE "page_tag_mapping"
(
    page_id TEXT NOT NULL,
    tag_id  TEXT NOT NULL
);

CREATE TABLE "version"
(
    version INT
);

INSERT INTO "version"(version)
VALUES (1);