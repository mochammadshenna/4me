-- Drop labels table and related objects
DROP INDEX IF EXISTS idx_labels_project_name;
DROP INDEX IF EXISTS idx_labels_project_id;
DROP TABLE IF EXISTS labels;
