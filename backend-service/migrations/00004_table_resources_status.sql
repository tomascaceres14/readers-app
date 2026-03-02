-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resources_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO resources_status (name) VALUES ('PENDING'), ('OK'), ('FAILED');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resources_status;
-- +goose StatementEnd
