# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Essential Commands

### Full-Stack Development
```bash
# Complete setup (from project root)
make setup              # Database, dependencies, env files
make dev                # Start both backend + frontend servers
make test               # Run all tests (backend + frontend)
make build              # Build both for production
make clean              # Remove all build artifacts

# Quick status check
make status             # Check if servers and database are running
make health             # Health check API endpoints
```

### Backend Commands (Go)
```bash
cd backend

# Development
make run                # Start backend server (localhost:8080)
make watch              # Auto-reload on changes (requires air)
make dev                # Clean, build, and test workflow

# Testing
make test               # Run all tests
make test-unit          # Unit tests only (./internal/...)
make test-e2e           # E2E integration tests only
make test-coverage      # Generate coverage.html report

# Single test execution
go test -v -run TestAuthHandler_Register ./internal/handlers/
go test -v ./internal/handlers/auth_test.go

# Database migrations
make migrate-up         # Run pending migrations (or just start server)
make migrate-create NAME=add_column_name  # Create new migration
make migrate-version    # Show current version

# Build
make build              # Build to bin/server
make prod               # Production build (CGO_ENABLED=0)
```

### Frontend Commands (Vue.js)
```bash
cd frontend

# Development
npm run dev             # Start frontend (localhost:5173)
npm run build           # Build for production
npm run preview         # Preview production build

# Testing
npm test                # Run unit tests (vitest)
npm test -- --watch     # Watch mode
npm run test:coverage   # Generate coverage report
npm run test:e2e        # Run Playwright E2E tests
npm run test:e2e:ui     # Playwright UI mode

# Code quality
npm run lint            # ESLint with auto-fix
```

## Architecture Overview

### Backend Structure (Go + Gin)

**Request Flow Pattern:**
```
HTTP Request → Gin Router → Auth Middleware → Handler → Database
                                ↓
                         Set userID in context
```

**Critical Pattern - Ownership Validation:**
All handlers follow this security pattern:
```go
// 1. Extract userID from middleware context
userID, _ := c.Get("userID")

// 2. Verify user owns the resource before any operation
var exists bool
db.QueryRow(ctx,
    `SELECT EXISTS(
        SELECT 1 FROM resource r
        JOIN parent p ON r.parent_id = p.id
        WHERE r.id = $1 AND p.user_id = $2
    )`, resourceID, userID).Scan(&exists)

if !exists {
    c.JSON(404, gin.H{"error": "Not found"})
    return
}

// 3. Proceed with operation
```

**Handler Layer Responsibilities:**
- `handlers/auth.go`: Registration, login, JWT generation, Google OAuth flow
- `handlers/projects.go`: Project CRUD with user ownership
- `handlers/boards.go`: Board management within projects
- `handlers/tasks.go`: Task CRUD, move between boards, history tracking
- `handlers/comments.go`, `handlers/attachments.go`, `handlers/labels.go`: Sub-resource management

**Transaction Pattern:**
All write operations that affect multiple tables use transactions:
```go
tx, _ := db.Pool.Begin(ctx)
defer tx.Rollback(ctx)

// Multiple operations...
tx.QueryRow(...) // Main operation
tx.Exec(...)     // Related operations
tx.Exec(...)     // History/audit logging

tx.Commit(ctx)   // Commit only if all succeed
```

### Frontend Structure (Vue 3 + Pinia)

**State Management Pattern:**
```
User Action → Component → Pinia Store → API Client → Backend
                              ↓              ↓
                         Update State   Add JWT Token
                              ↓
                         Reactive UI Update
```

**Pinia Store Responsibilities:**
- `stores/auth.js`: User session, token management, login/logout
- `stores/projects.js`: Projects list, current project, CRUD
- `stores/boards.js`: Boards array, board CRUD, position updates
- `stores/tasks.js`: Tasks by board, task operations, comments, attachments
- `stores/labels.js`: Project labels, label CRUD

**API Client Pattern:**
The Axios client in `api/client.js` handles:
```javascript
// Request interceptor - automatic JWT injection
config.headers.Authorization = `Bearer ${authStore.token}`

// Response interceptor - automatic error handling
if (error.response.status === 401) {
    authStore.logout()
    router.push('/login')
}
```

**Optimistic UI Updates:**
Tasks use optimistic updates for better UX:
```javascript
// Add to UI immediately
tasks.value.push(newTask)

// Then make API call
const result = await apiClient.post(...)

// Rollback on failure
if (!result.success) {
    tasks.value = tasks.value.filter(t => t.id !== newTask.id)
}
```

### Drag-and-Drop Architecture

The Kanban board uses **Atlassian's pragmatic-drag-and-drop** library (~4.7kB core) for high-performance, accessible drag-and-drop functionality.

**Key Components:**
- [TaskCard.vue](frontend/src/components/board/TaskCard.vue) - Draggable task cards
- [BoardColumn.vue](frontend/src/components/board/BoardColumn.vue) - Drop target columns
- [ProjectView.vue](frontend/src/views/ProjectView.vue) - Board container

**TaskCard.vue - Draggable Pattern:**
```vue
<script setup>
import { draggable } from '@atlaskit/pragmatic-drag-and-drop/element/adapter'
import { setCustomNativeDragPreview } from '@atlaskit/pragmatic-drag-and-drop/element/set-custom-native-drag-preview'

const taskCardRef = ref(null)
const isDragging = ref(false)

onMounted(() => {
  cleanupDraggable = draggable({
    element: taskCardRef.value.$el,
    getInitialData: () => ({
      type: 'task',
      taskId: props.task.id,
      boardId: props.task.board_id,
      task: props.task
    }),
    onGenerateDragPreview: ({ nativeSetDragImage }) => {
      setCustomNativeDragPreview({
        nativeSetDragImage,
        render: ({ container }) => {
          // Custom preview with rotation and shadow
          const preview = taskCardRef.value.$el.cloneNode(true)
          preview.style.transform = 'rotate(3deg)'
          preview.style.boxShadow = '0 12px 48px rgba(0, 0, 0, 0.25)'
          container.appendChild(preview)
        }
      })
    },
    onDragStart: () => isDragging.value = true,
    onDrop: () => isDragging.value = false
  })
})
</script>
```

**BoardColumn.vue - Drop Target Pattern:**
```vue
<script setup>
import { dropTargetForElements } from '@atlaskit/pragmatic-drag-and-drop/element/adapter'
import { autoScrollForElements } from '@atlaskit/pragmatic-drag-and-drop-auto-scroll/element'
import { attachClosestEdge, extractClosestEdge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/closest-edge'

const tasksContainerRef = ref(null)
const dropIndicatorIndex = ref(null)

onMounted(() => {
  cleanupDropTarget = dropTargetForElements({
    element: tasksContainerRef.value,
    canDrop: ({ source }) => source.data.type === 'task',
    getData: ({ input }) => attachClosestEdge(
      { boardId: props.board.id },
      { element: tasksContainerRef.value, input, allowedEdges: ['top', 'bottom'] }
    ),
    onDragEnter: () => isDropTarget.value = true,
    onDrag: ({ location }) => {
      // Show drop indicator at appropriate position
      const closestEdge = extractClosestEdge(location.current.dropTargets[0]?.data)
      dropIndicatorIndex.value = closestEdge === 'top' ? 0 : tasks.value.length
    },
    onDrop: async ({ source, location }) => {
      const taskId = source.data.taskId
      const sourceBoardId = source.data.boardId
      const targetBoardId = props.board.id

      if (sourceBoardId === targetBoardId) {
        // Reorder within same board
        const oldIndex = tasks.value.findIndex(t => t.id === taskId)
        const newIndex = dropIndicatorIndex.value
        tasks.value.splice(oldIndex, 1)
        tasks.value.splice(newIndex, 0, source.data.task)
        await tasksStore.updateTask(taskId, { position: newIndex })
      } else {
        // Move to different board
        await tasksStore.moveTask(taskId, targetBoardId, dropIndicatorIndex.value)
      }
    }
  })

  // Enable auto-scroll when dragging near edges
  cleanupAutoScroll = autoScrollForElements({
    element: tasksContainerRef.value,
    canScroll: () => true
  })
})
```

**Visual States:**
```css
/* Dragging state - reduced opacity */
.task-card.is-dragging {
  opacity: 0.5;
  transform: scale(0.98);
}

/* Drop target highlight */
.tasks-container.drop-target-active {
  border: 2px dashed #1976D2;
  background: rgba(25, 118, 210, 0.08);
}

/* Drop indicator - animated line */
.drop-indicator {
  background: linear-gradient(90deg, #1976D2 0%, #2196F3 100%);
  height: 4px;
  border-radius: 2px;
  animation: pulse 1s ease-in-out infinite;
}
```

**Data Flow:**
```
1. User starts drag → TaskCard.onDragStart()
2. Mouse moves over column → BoardColumn.onDragEnter()
3. Calculate drop position → Show drop indicator
4. User releases → BoardColumn.onDrop()
5. Update Pinia store → API call → Database update
6. Optimistic UI update → Rollback on error
```

**Key Features:**
- ✅ Auto-scroll when dragging near edges
- ✅ Visual drop indicators between tasks
- ✅ Custom drag preview with rotation
- ✅ Drag state visual feedback (opacity, scale)
- ✅ Keyboard navigation support (future)
- ✅ Screen reader announcements (future)
- ✅ Touch device support
- ✅ Performance optimized (~4.7kB core)

**Testing:**
- Comprehensive Playwright E2E tests in [frontend/playwright/drag-drop.spec.js](frontend/playwright/drag-drop.spec.js)
- 20+ test scenarios covering drag between columns, reordering, visual states, performance
- Run with: `npm run test:e2e` (frontend directory)

**Dependencies:**
```json
{
  "@atlaskit/pragmatic-drag-and-drop": "latest",
  "@atlaskit/pragmatic-drag-and-drop-hitbox": "latest",
  "@atlaskit/pragmatic-drag-and-drop-auto-scroll": "latest"
}
```

### Authentication Flow

**JWT Token Lifecycle:**
1. User logs in → Backend generates JWT with `userID`, `username`, `email`, `exp`
2. Frontend stores in `localStorage.setItem('token', jwt)`
3. Axios interceptor adds `Authorization: Bearer {token}` to all requests
4. Auth middleware validates signature and expiration
5. Middleware injects `userID` into Gin context for handlers

**Google OAuth Flow:**
1. Frontend requests OAuth URL from `GET /api/auth/google`
2. User redirects to Google, approves permissions
3. Google redirects to `GET /api/auth/google/callback?code=xyz`
4. Backend exchanges code for user info, creates/finds user
5. Backend generates JWT and redirects to frontend with token
6. Frontend extracts token from URL and stores it

### Database Schema Patterns

**Key Relationships:**
```
users → projects → boards → tasks
                 → labels
                        ↓
                    task_labels (join table)

tasks → comments (user_id FK for attribution)
     → attachments (file_url points to Supabase)
     → task_history (complete audit trail)
```

**Automatic Migrations:**
- Migrations run automatically on server startup
- Located in `backend/internal/database/database.go`
- SQL embedded directly in code for reliability
- Idempotent: safe to run multiple times

## Development Guidelines

### Testing Approach

**Backend Testing:**
- Unit tests colocated with handlers: `auth_test.go` next to `auth.go`
- Use `httptest.NewRecorder()` for handler testing
- Clean database before each test with `TRUNCATE ... CASCADE`
- Test database: `4me_todos_test` (created with `make test-db-setup`)

**Frontend Testing:**
- Component tests in `__tests__/` directories
- Mock API client with `vi.mock('@/api/client')`
- Mock Pinia stores for isolated component testing
- Use `@vue/test-utils` for component mounting

**E2E Testing:**
- Playwright tests in `frontend/playwright.config.js`
- Test complete user workflows (register → create project → add tasks)
- Backend E2E tests in `backend/tests/e2e/`

### Environment Configuration

**Backend `.env` (required):**
```bash
DATABASE_URL=postgres://user:pass@localhost:5432/4me_todos?sslmode=disable
JWT_SECRET=random-secret-string-change-in-production
GOOGLE_CLIENT_ID=your-google-oauth-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-oauth-client-secret
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-anon-key
PORT=8080
```

**Frontend `.env` (required):**
```bash
VITE_API_URL=http://localhost:8080/api
VITE_GOOGLE_CLIENT_ID=your-google-oauth-client-id.apps.googleusercontent.com
```

### Common Development Tasks

**Adding a New API Endpoint:**
1. Define model in `internal/models/models.go`
2. Create handler function in appropriate handler file
3. Add route in `cmd/api/main.go` (use `protected` group for auth)
4. Add ownership validation check if accessing user resources
5. Write unit tests in corresponding `_test.go` file

**Adding a New Database Table:**
1. Add migration SQL to `internal/database/database.go` in `RunMigrations()`
2. Define Go struct in `internal/models/models.go`
3. Server restart will run migration automatically

**Adding a New Frontend Feature:**
1. Create Pinia store action in appropriate store file
2. Add API call to `api/client.js` or directly in store
3. Create/update Vue components
4. Add route in `router/index.js` if needed
5. Update UI to call store action

### Key Technical Notes

**Ownership Security:**
Never trust client-provided IDs for authorization. Always validate through database joins:
```sql
-- ✅ Correct: Validates user owns parent resource
SELECT * FROM tasks t
JOIN boards b ON t.board_id = b.id
JOIN projects p ON b.project_id = p.id
WHERE t.id = $1 AND p.user_id = $2

-- ❌ Wrong: Only checks task exists, not ownership
SELECT * FROM tasks WHERE id = $1
```

**Migration Philosophy:**
- Migrations are NOT rolled back in production
- Use `migrate-create` to add new migrations
- Always test migrations against copy of production data
- Server handles migration state automatically

**Supabase Integration:**
- File uploads handled in `handlers/attachments.go`
- Requires Supabase Storage bucket named `4me-attachments`
- Files stored with format: `{user_id}/{task_id}/{filename}`
- Public access or RLS policies must be configured in Supabase dashboard

**Frontend Path Aliases:**
```javascript
// Configured in vite.config.js
import Component from '@/components/Component.vue'  // @ = src/
```

**State Persistence:**
- Auth token stored in `localStorage` (key: `token`)
- User object stored in Pinia auth store
- No other state persisted across sessions

### Debugging Tips

**Backend Issues:**
```bash
# Check database connection
psql -U postgres -d 4me_todos -c "SELECT 1"

# Verbose logging
go run cmd/api/main.go  # Gin logs all requests

# Check if server is running
curl http://localhost:8080/api/auth/me
```

**Frontend Issues:**
```bash
# Check API connection
console.log(import.meta.env.VITE_API_URL)

# View Pinia store state in browser
// Use Vue DevTools extension

# Check auth token
localStorage.getItem('token')
```

**Database Issues:**
```bash
# Reset database completely
make db-reset  # WARNING: Deletes all data

# Check migrations ran
psql -U postgres -d 4me_todos -c "\d"  # List all tables
```

## Project Context

This is a personal task management application inspired by Jira/Trello. Key design decisions:

- **Monorepo structure**: Backend and frontend in same repo for simplicity
- **JWT authentication**: Stateless auth for easy horizontal scaling
- **Automatic migrations**: No separate migration tool needed for development
- **Optimistic UI**: Tasks update immediately in UI before server confirmation
- **Ownership model**: All resources are scoped to the creating user
- **History tracking**: Complete audit trail in `task_history` table

For comprehensive architecture details, see [ARCHITECTURE.md](./ARCHITECTURE.md).
For deployment instructions, see [DEPLOYMENT.md](./DEPLOYMENT.md).
For testing strategies, see [TESTING.md](./TESTING.md).
