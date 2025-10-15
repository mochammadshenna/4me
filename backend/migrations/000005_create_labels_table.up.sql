-- Create labels table
CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    color VARCHAR(7) NOT NULL DEFAULT '#6b7280',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for labels table
CREATE INDEX idx_labels_project_id ON labels(project_id);

-- Create unique constraint for label names per project
CREATE UNIQUE INDEX idx_labels_project_name ON labels(project_id, name);
