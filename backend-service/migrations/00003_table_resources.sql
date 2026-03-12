-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resources (
id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
url TEXT NOT NULL,
title TEXT,
excerpt TEXT,
language VARCHAR(10),
status_id UUID REFERENCES resources_status(id),
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resources;
-- +goose StatementEnd
