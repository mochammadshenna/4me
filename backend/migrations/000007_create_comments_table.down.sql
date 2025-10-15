-- Drop comments table and related objects
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP INDEX IF EXISTS idx_comments_created_at;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_task_id;
DROP TABLE IF EXISTS comments;
