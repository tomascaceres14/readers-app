-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
username varchar(255) NOT NULL,
email varchar(255) NOT NULL,
password TEXT NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
