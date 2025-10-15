# 4me Todos - Personal Task Management Application

A full-stack task management application with Jira-like Kanban boards, built for personal use.

![Tech Stack](https://img.shields.io/badge/Frontend-Vue.js_3-4FC08D)
![Tech Stack](https://img.shields.io/badge/Backend-Go-00ADD8)
![Tech Stack](https://img.shields.io/badge/Database-PostgreSQL-336791)
![Tech Stack](https://img.shields.io/badge/Storage-Supabase-3ECF8E)

## Features

### 🔐 Authentication

- Username/password authentication with JWT tokens
- Google OAuth 2.0 integration
- Secure password hashing with bcrypt

### 📊 Project Management

- Create and manage multiple projects
- Color-coded projects
- Project descriptions

### 📋 Kanban Boards

- Drag-and-drop task cards between columns
- Create custom board columns
- Reorder tasks within columns
- Visual board organization

### ✅ Advanced Task Management

- **Task Details**
  - Title and description
  - Priority levels (Low, Medium, High, Urgent)
  - Due dates with visual indicators
  - Custom labels with colors

- **Collaboration Features**
  - Comments with edit/delete
  - User avatars and timestamps
  - @mentions ready (user display)

- **Attachments**
  - File upload to Supabase storage
  - Support for documents, images, PDFs
  - File preview and download

- **Activity History**
  - Complete task change log
  - User attribution
  - Timestamps for all actions

### 🎨 Modern UI/UX

- Clean, professional interface
- Responsive design (desktop, tablet, mobile)
- Smooth animations and transitions
- Loading states and skeletons
- Error handling with user-friendly messages

## Tech Stack

### Frontend

- **Framework**: Vue.js 3 with Composition API
- **Build Tool**: Vite
- **UI Libraries**:
  - Tailwind CSS for utility-first styling
  - Vuetify 3 for Material Design components
- **State Management**: Pinia
- **Routing**: Vue Router
- **HTTP Client**: Axios with interceptors
- **Drag & Drop**: vuedraggable
- **Date Handling**: date-fns
- **Icons**: Material Design Icons

### Backend

- **Language**: Go 1.24+
- **Framework**: Gin
- **Database**: PostgreSQL 14+
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **OAuth**: Google OAuth 2.0
- **Database Driver**: pgx/v5

### Storage & Deployment

- **File Storage**: Supabase Storage
- **Frontend Hosting**: Vercel
- **Backend Hosting**: Railway/Render/DigitalOcean (or local)
- **Database Hosting**: Supabase PostgreSQL

## Project Structure

```
4me/
├── backend/                 # Go backend API
│   ├── cmd/api/            # Application entry point
│   ├── internal/
│   │   ├── config/         # Configuration management
│   │   ├── database/       # Database connection & migrations
│   │   ├── handlers/       # HTTP request handlers
│   │   ├── middleware/     # Auth, CORS middleware
│   │   ├── models/         # Data models
│   │   └── utils/          # JWT, password utilities
│   ├── go.mod              # Go dependencies
│   └── README.md           # Backend documentation
│
├── frontend/               # Vue.js frontend
│   ├── src/
│   │   ├── api/           # API client
│   │   ├── components/    # Vue components
│   │   ├── plugins/       # Vuetify config
│   │   ├── router/        # Vue Router
│   │   ├── stores/        # Pinia stores
│   │   ├── views/         # Page components
│   │   └── main.js        # App entry point
│   ├── package.json       # NPM dependencies
│   ├── vercel.json        # Vercel config
│   └── README.md          # Frontend documentation
│
└── README.md              # This file
```

## Getting Started

### Prerequisites

- **Node.js** 18+ and npm
- **Go** 1.24+
- **PostgreSQL** 14+
- **Supabase** account (for file storage)
- **Google Cloud Console** project (for OAuth)

### 1. Clone the Repository

```bash
git clone <your-repo-url>
cd 4me
```

### 2. Backend Setup

```bash
cd backend

# Install dependencies
go mod tidy

# Create PostgreSQL database
createdb 4me_todos

# Copy environment file
cp .env.example .env

# Edit .env with your credentials
# DATABASE_URL, JWT_SECRET, GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, SUPABASE_URL, SUPABASE_KEY

# Run the server (migrations run automatically)
go run cmd/api/main.go
```

Backend will start on `http://localhost:8080`

### 3. Create Demo User (Optional)

For testing purposes, create a demo user:

```bash
# From project root
make seed-demo
```

This creates:
- **Username**: `demo`
- **Password**: `password`
- **Email**: `demo@example.com`

You can now login with these credentials at the frontend.

### 4. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Copy environment file
cp .env.example .env

# Edit .env with your configuration
# VITE_API_URL, VITE_GOOGLE_CLIENT_ID

# Run development server
npm run dev
```

Frontend will start on `http://localhost:5173`

### 4. Supabase Setup

1. Create a Supabase project at <https://supabase.com>
2. Create a storage bucket named `4me-attachments`
3. Set the bucket to public or configure appropriate policies
4. Copy your project URL and anon key to the backend `.env` file

### 5. Google OAuth Setup

1. Go to Google Cloud Console
2. Create OAuth 2.0 credentials
3. Add authorized redirect URI: `http://localhost:8080/api/auth/google/callback`
4. Copy Client ID and Secret to backend `.env`
5. Copy Client ID to frontend `.env`

## API Documentation

### Authentication Endpoints

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login with credentials
- `GET /api/auth/google` - Get Google OAuth URL
- `GET /api/auth/google/callback` - Handle Google OAuth callback
- `GET /api/auth/me` - Get current user (protected)

### Projects

- `POST /api/projects` - Create project
- `GET /api/projects` - List all projects
- `GET /api/projects/:id` - Get project details
- `PUT /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Delete project

### Boards

- `POST /api/projects/:id/boards` - Create board
- `GET /api/projects/:id/boards` - List boards
- `PUT /api/boards/:id` - Update board
- `DELETE /api/boards/:id` - Delete board

### Tasks

- `POST /api/boards/:id/tasks` - Create task
- `GET /api/tasks/:id` - Get task details
- `PUT /api/tasks/:id` - Update task
- `PATCH /api/tasks/:id/move` - Move task between boards
- `DELETE /api/tasks/:id` - Delete task
- `GET /api/tasks/:id/history` - Get task history

### Labels, Comments, Attachments

See detailed API documentation in `/backend/README.md`

## Database Schema

The application uses the following main tables:

- `users` - User accounts and authentication
- `projects` - User projects
- `boards` - Kanban board columns
- `tasks` - Individual tasks
- `labels` - Color-coded task labels
- `task_labels` - Task-label relationships
- `comments` - Task comments
- `attachments` - File attachments metadata
- `task_history` - Complete audit trail

## Deployment

### Frontend (Vercel)

1. Push code to GitHub
2. Import project in Vercel
3. Configure environment variables
4. Deploy automatically on push

Or use CLI:

```bash
cd frontend
npm install -g vercel
vercel
```

### Backend (Railway/Render)

1. Create new project
2. Connect GitHub repository
3. Set environment variables
4. Deploy

### Database (Supabase)

Use Supabase's managed PostgreSQL database with automatic backups.

## Development

### Running Tests

Backend:

```bash
cd backend
go test ./...
```

Frontend:

```bash
cd frontend
npm run test
```

### Building for Production

Backend:

```bash
cd backend
go build -o bin/server cmd/api/main.go
```

Frontend:

```bash
cd frontend
npm run build
```

## Security Considerations

- JWT tokens with configurable expiration
- Password hashing with bcrypt
- CORS protection
- SQL injection prevention with parameterized queries
- Input validation on both frontend and backend
- Protected routes with authentication middleware
- Secure file upload handling

## Contributing

This is a personal project, but suggestions and feedback are welcome!

## License

MIT License - feel free to use this project for your own purposes.

## Acknowledgments

- Inspired by Jira, Trello, and modern project management tools
- Built with modern web technologies and best practices
- Designed for simplicity and ease of use

---

**Author**: Mochammad Shenna Wardana  
**Email**: <shenawardana@gmail.com>  
**Version**: 1.0.0
