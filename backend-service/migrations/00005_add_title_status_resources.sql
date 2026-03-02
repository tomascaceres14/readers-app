-- +goose Up
-- +goose StatementBegin
ALTER TABLE resources ADD COLUMN IF NOT EXISTS title TEXT;
ALTER TABLE resources ADD COLUMN IF NOT EXISTS status_id UUID REFERENCES resources_status(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE resources DROP COLUMN IF EXISTS title;
ALTER TABLE resources DROP COLUMN IF EXISTS status_id;
-- +goose StatementEnd
