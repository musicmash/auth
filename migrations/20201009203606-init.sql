-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "users" (
    "name"       VARCHAR(255) PRIMARY KEY,
    "created_at" TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "photo"      VARCHAR(255) NULL
);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS "sessions" (
    "id"         SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "user_name"  VARCHAR(255) NOT NULL,
    "value"      VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_name) REFERENCES "users" (name) ON UPDATE RESTRICT ON DELETE RESTRICT
);
-- +migrate StatementEnd
CREATE INDEX CONCURRENTLY "idx_sessions_value" ON "sessions" USING HASH (value);

-- +migrate Down
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "users";