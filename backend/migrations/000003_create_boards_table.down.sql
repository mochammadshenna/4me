-- Drop boards table and related objects
DROP TRIGGER IF EXISTS update_boards_updated_at ON boards;
DROP INDEX IF EXISTS idx_boards_position;
DROP INDEX IF EXISTS idx_boards_project_id;
DROP TABLE IF EXISTS boards;
