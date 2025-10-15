# 4me Todos Backend

Go backend API for the 4me Todos task management application.

## Tech Stack

- **Language**: Go 1.24+
- **Framework**: Gin
- **Database**: PostgreSQL
- **Authentication**: JWT + Google OAuth 2.0
- **File Storage**: Supabase

## Setup

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 14+ installed and running
- Supabase account (for file storage)
- Google Cloud Console project (for OAuth)

### Installation

1. Install dependencies:
```bash
go mod tidy
```

2. Create a PostgreSQL database:
```bash
createdb 4me_todos
```

3. Copy the environment file:
```bash
cp .env.example .env
```

4. Update `.env` with your credentials:
   - `DATABASE_URL`: PostgreSQL connection string
   - `JWT_SECRET`: Random secret key for JWT tokens
   - `GOOGLE_CLIENT_ID`: Google OAuth client ID
   - `GOOGLE_CLIENT_SECRET`: Google OAuth client secret
   - `SUPABASE_URL`: Your Supabase project URL
   - `SUPABASE_KEY`: Supabase anon/public key

### Running the Server

Development mode:
```bash
go run cmd/api/main.go
```

Build and run:
```bash
go build -o bin/server cmd/api/main.go
./bin/server
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login with credentials
- `GET /api/auth/google` - Get Google OAuth URL
- `GET /api/auth/google/callback` - Google OAuth callback
- `GET /api/auth/me` - Get current user (protected)

### Projects

- `POST /api/projects` - Create project
- `GET /api/projects` - List all projects
- `GET /api/projects/:id` - Get project details
- `PUT /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Delete project

### Boards

- `POST /api/projects/:id/boards` - Create board
- `GET /api/projects/:id/boards` - List project boards
- `PUT /api/boards/:id` - Update board
- `DELETE /api/boards/:id` - Delete board

### Tasks

- `POST /api/boards/:id/tasks` - Create task
- `GET /api/tasks/:id` - Get task details
- `PUT /api/tasks/:id` - Update task
- `PATCH /api/tasks/:id/move` - Move task to different board
- `DELETE /api/tasks/:id` - Delete task
- `GET /api/tasks/:id/history` - Get task history

### Labels

- `POST /api/projects/:id/labels` - Create label
- `GET /api/projects/:id/labels` - List project labels
- `PUT /api/labels/:id` - Update label
- `DELETE /api/labels/:id` - Delete label

### Comments

- `POST /api/tasks/:id/comments` - Add comment
- `GET /api/tasks/:id/comments` - List task comments
- `PUT /api/comments/:id` - Update comment
- `DELETE /api/comments/:id` - Delete comment

### Attachments

- `POST /api/tasks/:id/attachments` - Upload attachment
- `GET /api/tasks/:id/attachments` - List task attachments
- `DELETE /api/attachments/:id` - Delete attachment

## Database Schema

The application uses PostgreSQL with the following tables:

- `users` - User accounts
- `projects` - User projects
- `boards` - Kanban boards/columns
- `tasks` - Individual tasks
- `labels` - Task labels
- `task_labels` - Task-label relationships
- `comments` - Task comments
- `attachments` - File attachments
- `task_history` - Task change history

Migrations run automatically on server startup.

## Authentication

The API supports two authentication methods:

1. **Username/Password**: Traditional registration and login with bcrypt password hashing
2. **Google OAuth**: Sign in with Google account

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

## Development

### Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── database/           # Database connection and migrations
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # Middleware (auth, CORS)
│   ├── models/             # Data models and DTOs
│   └── utils/              # Utility functions (JWT, password)
├── migrations/             # Database migrations (if using migrate tool)
├── .env                    # Environment variables (not in git)
├── .env.example           # Example environment file
├── go.mod                 # Go module definition
└── README.md              # This file
```

## License

MIT

