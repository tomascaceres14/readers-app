-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
id TEXT PRIMARY KEY,
username varchar(255) NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
