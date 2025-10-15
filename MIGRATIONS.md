# Database Migrations Guide

This guide covers the database migration system for the 4me Todos application.

## Overview

The application uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management. Migrations are stored as SQL files in the `backend/migrations/` directory.

## Migration Files Structure

```
backend/migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_projects_table.up.sql
├── 000002_create_projects_table.down.sql
├── 000003_create_boards_table.up.sql
├── 000003_create_boards_table.down.sql
├── 000004_create_tasks_table.up.sql
├── 000004_create_tasks_table.down.sql
├── 000005_create_labels_table.up.sql
├── 000005_create_labels_table.down.sql
├── 000006_create_task_labels_table.up.sql
├── 000006_create_task_labels_table.down.sql
├── 000007_create_comments_table.up.sql
├── 000007_create_comments_table.down.sql
├── 000008_create_attachments_table.up.sql
├── 000008_create_attachments_table.down.sql
├── 000009_create_task_history_table.up.sql
└── 000009_create_task_history_table.down.sql
```

### Naming Convention

- `{version}_{description}.up.sql` - Migration to apply changes
- `{version}_{description}.down.sql` - Migration to rollback changes
- Version numbers are sequential: `000001`, `000002`, etc.
- Descriptions are descriptive: `create_users_table`, `add_index_to_tasks`, etc.

## Available Migrations

### 000001 - Users Table

**UP**: Creates users table with authentication fields
**DOWN**: Drops users table and related indexes

### 000002 - Projects Table

**UP**: Creates projects table linked to users
**DOWN**: Drops projects table

### 000003 - Boards Table

**UP**: Creates boards table linked to projects
**DOWN**: Drops boards table

### 000004 - Tasks Table

**UP**: Creates tasks table linked to boards
**DOWN**: Drops tasks table

### 000005 - Labels Table

**UP**: Creates labels table for task categorization
**DOWN**: Drops labels table

### 000006 - Task Labels Junction

**UP**: Creates many-to-many relationship between tasks and labels
**DOWN**: Drops task_labels junction table

### 000007 - Comments Table

**UP**: Creates comments table for task discussions
**DOWN**: Drops comments table

### 000008 - Attachments Table

**UP**: Creates attachments table for file uploads
**DOWN**: Drops attachments table

### 000009 - Task History Table

**UP**: Creates task history table for audit trail
**DOWN**: Drops task history table

## Installation

### Install Migration CLI Tool

```bash
# Install golang-migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Or using the Makefile
cd backend
make migrate-cli
```

### Verify Installation

```bash
migrate -version
```

## Usage

### Running Migrations

#### Option 1: Using Makefiles (Recommended)

```bash
# From project root
make migrate-up

# From backend directory
cd backend
make migrate-up

# From database directory
cd database
make migrate-up
```

#### Option 2: Using CLI directly

```bash
# Run all pending migrations
migrate -path backend/migrations -database "postgres://user:pass@localhost:5432/4me_todos?sslmode=disable" up

# Run specific number of migrations
migrate -path backend/migrations -database "$DATABASE_URL" up 2

# Force migration to specific version (use with caution)
migrate -path backend/migrations -database "$DATABASE_URL" force 5
```

### Rolling Back Migrations

```bash
# Rollback last migration
make migrate-down

# Or using CLI
migrate -path backend/migrations -database "$DATABASE_URL" down 1

# Rollback multiple migrations
migrate -path backend/migrations -database "$DATABASE_URL" down 3
```

### Checking Migration Status

```bash
# Check current version
make migrate-version

# Or using CLI
migrate -path backend/migrations -database "$DATABASE_URL" version
```

## Creating New Migrations

### Using Makefile (Recommended)

```bash
cd backend
make migrate-create NAME=add_user_preferences_table
```

This creates:

- `000010_add_user_preferences_table.up.sql`
- `000010_add_user_preferences_table.down.sql`

### Using CLI directly

```bash
migrate create -ext sql -dir backend/migrations -seq add_user_preferences_table
```

### Manual Creation

Create files manually following the naming convention:

```bash
touch backend/migrations/000010_add_user_preferences_table.up.sql
touch backend/migrations/000010_add_user_preferences_table.down.sql
```

## Writing Migration Files

### UP Migration Example

```sql
-- 000010_add_user_preferences_table.up.sql
CREATE TABLE user_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    theme VARCHAR(20) DEFAULT 'light',
    notifications_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);

CREATE TRIGGER update_user_preferences_updated_at 
    BEFORE UPDATE ON user_preferences 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

### DOWN Migration Example

```sql
-- 000010_add_user_preferences_table.down.sql
DROP TRIGGER IF EXISTS update_user_preferences_updated_at ON user_preferences;
DROP INDEX IF EXISTS idx_user_preferences_user_id;
DROP TABLE IF EXISTS user_preferences;
```

## Best Practices

### 1. Always Create Both UP and DOWN Migrations

- UP migrations apply changes
- DOWN migrations must completely reverse the UP migration
- Test both directions before deploying

### 2. Use Transactions for Complex Migrations

```sql
BEGIN;

-- Add column
ALTER TABLE tasks ADD COLUMN estimated_hours INTEGER;

-- Update existing data
UPDATE tasks SET estimated_hours = 1 WHERE estimated_hours IS NULL;

-- Add constraint
ALTER TABLE tasks ALTER COLUMN estimated_hours SET NOT NULL;

COMMIT;
```

### 3. Handle Data Migration Carefully

```sql
-- Add new column with default
ALTER TABLE users ADD COLUMN timezone VARCHAR(50) DEFAULT 'UTC';

-- Migrate existing data
UPDATE users SET timezone = 'America/New_York' WHERE email LIKE '%@company.com';

-- Make column required (after data migration)
ALTER TABLE users ALTER COLUMN timezone SET NOT NULL;
```

### 4. Index Management

```sql
-- UP: Create index
CREATE INDEX CONCURRENTLY idx_tasks_due_date ON tasks(due_date) WHERE due_date IS NOT NULL;

-- DOWN: Drop index
DROP INDEX IF EXISTS idx_tasks_due_date;
```

### 5. Foreign Key Constraints

```sql
-- UP: Add foreign key
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_assignee 
    FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL;

-- DOWN: Drop foreign key
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_tasks_assignee;
```

## Environment Variables

Set the `DATABASE_URL` environment variable:

```bash
# Development
export DATABASE_URL="postgres://localhost:5432/4me_todos?sslmode=disable"

# Production
export DATABASE_URL="postgres://user:pass@host:5432/4me_todos?sslmode=require"
```

Or create a `.env` file:

```env
DATABASE_URL=postgres://localhost:5432/4me_todos?sslmode=disable
```

## Troubleshooting

### Migration Failed

1. **Check migration status**:

   ```bash
   make migrate-version
   ```

2. **Force to previous version**:

   ```bash
   make migrate-force VERSION=5
   ```

3. **Fix the migration file** and try again

### Database Connection Issues

1. **Check database is running**:

   ```bash
   pg_isready
   ```

2. **Verify connection string**:

   ```bash
   psql "$DATABASE_URL"
   ```

3. **Check environment variables**:

   ```bash
   echo $DATABASE_URL
   ```

### Migration Already Applied

If you get "no change" errors:

```bash
# Check current version
migrate -path backend/migrations -database "$DATABASE_URL" version

# Force to specific version if needed
migrate -path backend/migrations -database "$DATABASE_URL" force 9
```

## Development Workflow

### 1. Create New Migration

```bash
cd backend
make migrate-create NAME=add_task_priority_index
```

### 2. Write Migration Files

Edit the generated `.up.sql` and `.down.sql` files.

### 3. Test Migration

```bash
# Apply migration
make migrate-up

# Verify changes
make migrate-version
psql "$DATABASE_URL" -c "\d tasks"

# Test rollback
make migrate-down
make migrate-up
```

### 4. Commit to Version Control

```bash
git add backend/migrations/
git commit -m "Add migration for task priority index"
```

## Production Deployment

### 1. Backup Database

```bash
make backup
```

### 2. Run Migrations

```bash
# Check current version
make migrate-version

# Apply migrations
make migrate-up

# Verify
make migrate-version
```

### 3. Monitor Application

Check application logs for any issues after migration.

## Migration History

The migration tool creates a `schema_migrations` table to track applied migrations:

```sql
SELECT * FROM schema_migrations;
```

This table contains:

- `version` - Migration version number
- `dirty` - Whether migration failed (1) or succeeded (0)

## Rollback Strategies

### 1. Planned Rollback

```bash
# Rollback last migration
make migrate-down

# Rollback multiple migrations
migrate -path backend/migrations -database "$DATABASE_URL" down 3
```

### 2. Emergency Rollback

```bash
# Force to previous stable version
make migrate-force VERSION=8

# Fix issues and reapply
make migrate-up
```

### 3. Database Reset (Development Only)

```bash
# WARNING: This deletes all data
make db-reset
```

## Common Patterns

### Adding a Column

```sql
-- UP
ALTER TABLE tasks ADD COLUMN estimated_hours INTEGER;

-- DOWN
ALTER TABLE tasks DROP COLUMN estimated_hours;
```

### Renaming a Column

```sql
-- UP
ALTER TABLE tasks RENAME COLUMN description TO content;

-- DOWN
ALTER TABLE tasks RENAME COLUMN content TO description;
```

### Adding an Index

```sql
-- UP
CREATE INDEX CONCURRENTLY idx_tasks_priority ON tasks(priority);

-- DOWN
DROP INDEX IF EXISTS idx_tasks_priority;
```

### Modifying Constraints

```sql
-- UP
ALTER TABLE tasks DROP CONSTRAINT tasks_status_check;
ALTER TABLE tasks ADD CONSTRAINT tasks_status_check 
    CHECK (status IN ('todo', 'in_progress', 'done', 'blocked'));

-- DOWN
ALTER TABLE tasks DROP CONSTRAINT tasks_status_check;
ALTER TABLE tasks ADD CONSTRAINT tasks_status_check 
    CHECK (status IN ('todo', 'in_progress', 'done'));
```

## Integration with CI/CD

### GitHub Actions Example

```yaml
- name: Run Migrations
  run: |
    make migrate-up
  env:
    DATABASE_URL: ${{ secrets.DATABASE_URL }}
```

### Pre-deployment Checklist

- [ ] All migrations tested in development
- [ ] Database backup created
- [ ] Migration rollback plan ready
- [ ] Application compatibility verified
- [ ] Monitoring in place

## Resources

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Database Migration Best Practices](https://www.prisma.io/dataguide/types/relational/migrating-your-database)

## Support

For migration-related issues:

1. Check this guide first
2. Review migration files for syntax errors
3. Test migrations in development environment
4. Check database logs for detailed error messages
5. Consult the golang-migrate documentation
