-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resources (
id uuid PRIMARY KEY DEFAULT gen_random_uuid() ,
url TEXT NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resources;
-- +goose StatementEnd
