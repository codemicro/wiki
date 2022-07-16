-- Migrate a new database to v1 format.

CREATE TABLE "version"
(
    version INT
);

INSERT INTO "version"(version)
VALUES (1);