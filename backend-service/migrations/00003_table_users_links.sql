-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_links (
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    link_id UUID REFERENCES links(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, link_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
