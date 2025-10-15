# Quick Start Guide - 4me Todos

Get the application running in 5 minutes!

## Prerequisites Checklist

- [ ] Go 1.24+ installed (`go version`)
- [ ] Node.js 18+ installed (`node --version`)
- [ ] PostgreSQL 14+ installed and running (`pg_isready`)
- [ ] Git installed

## Step 1: Clone and Setup (2 minutes)

```bash
# Clone the repository
git clone <your-repo-url>
cd 4me

# Check your current location
pwd  # Should show: .../4me
```

## Step 2: Backend Setup (2 minutes)

```bash
# Navigate to backend
cd backend

# Install Go dependencies
go mod download

# Create PostgreSQL database
createdb 4me_todos

# Copy environment file
cp .env.example .env

# Edit .env file with your settings
# Minimum required:
#   DATABASE_URL=postgres://localhost:5432/4me_todos?sslmode=disable
#   JWT_SECRET=your-secret-key-change-this
```

Edit `backend/.env`:

```env
DATABASE_URL=postgres://localhost:5432/4me_todos?sslmode=disable
JWT_SECRET=my-super-secret-jwt-key-change-this-in-production
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
SUPABASE_URL=
SUPABASE_KEY=
FRONTEND_URL=http://localhost:5173
PORT=8080
```

```bash
# Build and run the backend
go run cmd/api/main.go

# You should see:
# Database connection established successfully
# Database migrations completed successfully
# Server starting on port 8080...
```

**Leave this terminal running!**

## Step 3: Frontend Setup (1 minute)

Open a **new terminal**:

```bash
# Navigate to frontend (from project root)
cd frontend

# Install dependencies
npm install

# Copy environment file
cp .env.example .env

# Edit .env with your backend URL
echo "VITE_API_URL=http://localhost:8080/api" > .env
echo "VITE_GOOGLE_CLIENT_ID=" >> .env
```

```bash
# Start the development server
npm run dev

# You should see:
# VITE v6.0.3  ready in XXX ms
# ➜  Local:   http://localhost:5173/
```

## Step 4: Open and Test

1. Open your browser to: <http://localhost:5173>
2. Click **"Sign up"**
3. Register a new account:
   - Username: `testuser`
   - Email: `test@example.com`
   - Password: `password123`
4. You should be automatically logged in and see the Dashboard!

## Verification Checklist

- [ ] Backend running on <http://localhost:8080>
- [ ] Frontend running on <http://localhost:5173>
- [ ] Can register a new user
- [ ] Can create a project
- [ ] Can create boards
- [ ] Can create tasks
- [ ] Can drag tasks between boards

## Common Issues

### Backend won't start

**Error**: `Failed to connect to database`

```bash
# Check if PostgreSQL is running
pg_isready

# If not running, start it:
# macOS (Homebrew):
brew services start postgresql

# Linux:
sudo service postgresql start

# Create the database again
createdb 4me_todos
```

**Error**: `Port 8080 already in use`

```bash
# Find and kill the process
lsof -ti:8080 | xargs kill -9

# Or change PORT in backend/.env
PORT=8081
```

### Frontend won't start

**Error**: `Cannot find module`

```bash
# Delete and reinstall
rm -rf node_modules package-lock.json
npm install
```

**Error**: `API calls fail with CORS error`

```bash
# Ensure FRONTEND_URL in backend/.env matches your frontend URL
FRONTEND_URL=http://localhost:5173
```

## Testing the Application

### Quick Manual Test Flow

1. **Register** → Create account
2. **Create Project** → Click "New Project"
3. **Add Boards** → Create "To Do", "In Progress", "Done"
4. **Create Tasks** → Add tasks to boards
5. **Drag & Drop** → Move tasks between boards
6. **Task Details** → Click a task to see full details
7. **Add Labels** → Click label icon, create colored labels
8. **Add Comments** → Open task, add a comment

### Run Automated Tests

```bash
# Backend tests (new terminal)
cd backend
make test

# Frontend tests (new terminal)
cd frontend
npm test
```

## Development Workflow

### Terminal Layout (3 terminals)

**Terminal 1 - Backend**:

```bash
cd backend
go run cmd/api/main.go
```

**Terminal 2 - Frontend**:

```bash
cd frontend
npm run dev
```

**Terminal 3 - Commands**:

```bash
# Run tests, build, git commands, etc.
```

## Next Steps

### Enable Google OAuth (Optional)

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create OAuth 2.0 credentials
3. Add to `backend/.env`:

   ```env
   GOOGLE_CLIENT_ID=your-client-id
   GOOGLE_CLIENT_SECRET=your-client-secret
   ```

4. Add to `frontend/.env`:

   ```env
   VITE_GOOGLE_CLIENT_ID=your-client-id
   ```

5. Restart both servers

### Enable File Uploads (Optional)

1. Create [Supabase](https://supabase.com) account
2. Create project and storage bucket named `4me-attachments`
3. Add to `backend/.env`:

   ```env
   SUPABASE_URL=https://your-project.supabase.co
   SUPABASE_KEY=your-anon-key
   ```

4. Restart backend

## Useful Commands

### Backend

```bash
# Build
make build

# Run tests
make test

# Test coverage
make test-coverage

# Format code
make fmt

# Lint
make lint
```

### Frontend

```bash
# Development
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Run tests
npm test

# Test coverage
npm run test:coverage
```

### Database

```bash
# Connect to database
psql 4me_todos

# View tables
\dt

# View users
SELECT * FROM users;

# View projects
SELECT * FROM projects;

# Reset database
DROP DATABASE 4me_todos;
CREATE DATABASE 4me_todos;
# Then restart backend (migrations run automatically)
```

## Documentation

- [README.md](./README.md) - Project overview and features
- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture with diagrams
- [TESTING.md](./TESTING.md) - Testing guide
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Deployment instructions
- [Backend README](./backend/README.md) - Backend API documentation
- [Frontend README](./frontend/README.md) - Frontend documentation

## Support

If you encounter issues:

1. Check the [Common Issues](#common-issues) section above
2. Review the error messages carefully
3. Check that all prerequisites are installed
4. Ensure both backend and frontend are running
5. Check the browser console for frontend errors
6. Check the terminal output for backend errors

## Success! 🎉

You should now have:

- ✅ Backend API running on port 8080
- ✅ Frontend app running on port 5173
- ✅ PostgreSQL database with migrations
- ✅ Ability to create projects, boards, and tasks
- ✅ Working authentication system
- ✅ Drag-and-drop kanban boards

Happy task managing!
