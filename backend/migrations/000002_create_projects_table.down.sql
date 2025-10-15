-- Drop projects table and related objects
DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;
DROP INDEX IF EXISTS idx_projects_created_at;
DROP INDEX IF EXISTS idx_projects_user_id;
DROP TABLE IF EXISTS projects;
