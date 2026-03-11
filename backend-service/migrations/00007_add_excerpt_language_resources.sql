-- +goose Up
-- +goose StatementBegin
ALTER TABLE resources ADD COLUMN IF NOT EXISTS excerpt TEXT;
ALTER TABLE resources ADD COLUMN IF NOT EXISTS language VARCHAR(10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE resources DROP COLUMN IF EXISTS excerpt;
ALTER TABLE resources DROP COLUMN IF EXISTS language;
-- +goose StatementEnd
