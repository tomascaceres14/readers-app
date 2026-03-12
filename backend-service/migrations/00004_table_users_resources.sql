-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_resources (
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, resource_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_resources;
-- +goose StatementEnd
