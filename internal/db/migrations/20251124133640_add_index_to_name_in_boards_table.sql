-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS IDX_boards_name ON boards (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS IDX_boards_name;
-- +goose StatementEnd
