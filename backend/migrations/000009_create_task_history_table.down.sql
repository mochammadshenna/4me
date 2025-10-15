-- Drop task_history table and related objects
DROP INDEX IF EXISTS idx_task_history_action;
DROP INDEX IF EXISTS idx_task_history_created_at;
DROP INDEX IF EXISTS idx_task_history_task_id;
DROP TABLE IF EXISTS task_history;
