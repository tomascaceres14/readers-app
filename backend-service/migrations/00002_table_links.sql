-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links (
id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
url TEXT NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS links;
-- +goose StatementEnd
