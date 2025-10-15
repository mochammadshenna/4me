-- Drop attachments table and related objects
DROP INDEX IF EXISTS idx_attachments_uploaded_at;
DROP INDEX IF EXISTS idx_attachments_task_id;
DROP TABLE IF EXISTS attachments;
