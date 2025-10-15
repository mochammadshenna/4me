# 4me Todos Frontend

Vue.js 3 frontend for the 4me Todos personal task management application.

## Tech Stack

- **Framework**: Vue.js 3 + Vite
- **UI Libraries**: Tailwind CSS + Vuetify 3
- **State Management**: Pinia
- **Routing**: Vue Router
- **HTTP Client**: Axios
- **Date Handling**: date-fns
- **Drag & Drop**: vuedraggable

## Features

- 🔐 Authentication (Username/Password + Google OAuth)
- 📊 Project Management
- 📋 Kanban Boards with Drag & Drop
- ✅ Task Management with:
  - Priorities
  - Labels
  - Due Dates
  - Comments
  - File Attachments
  - Activity History
- 🎨 Modern, Responsive UI
- 🌙 Light/Dark Mode Support (via Vuetify)

## Setup

### Prerequisites

- Node.js 18+ and npm

### Installation

1. Install dependencies:

```bash
npm install
```

2. Create environment file:

```bash
cp .env.example .env
```

3. Update `.env` with your configuration:

```
VITE_API_URL=http://localhost:8080/api
VITE_GOOGLE_CLIENT_ID=your-google-oauth-client-id
```

### Development

Run the development server:

```bash
npm run dev
```

The app will be available at `http://localhost:5173`

### Build

Build for production:

```bash
npm run build
```

Preview production build:

```bash
npm run preview
```

## Project Structure

```
src/
├── api/                  # API client and interceptors
├── assets/              # Static assets
├── components/          # Vue components
│   ├── auth/           # Authentication components
│   ├── board/          # Kanban board components
│   ├── task/           # Task-related components
│   ├── common/         # Reusable components
│   └── layout/         # Layout components
├── plugins/            # Vue plugins (Vuetify)
├── router/             # Vue Router configuration
├── stores/             # Pinia stores
│   ├── auth.js        # Authentication state
│   ├── projects.js    # Projects state
│   ├── boards.js      # Boards state
│   ├── tasks.js       # Tasks state
│   └── labels.js      # Labels state
├── views/              # Page components
│   ├── LoginView.vue
│   ├── RegisterView.vue
│   ├── DashboardView.vue
│   └── ProjectView.vue
├── App.vue             # Root component
├── main.js             # Application entry point
└── style.css           # Global styles (Tailwind)
```

## Key Components

### Authentication

- **LoginView**: Username/password login with Google OAuth option
- **RegisterView**: New user registration
- **AuthCallbackView**: Handles Google OAuth callback

### Dashboard

- **DashboardView**: Shows all user projects
- Create, edit, delete projects
- Project color customization

### Project Board

- **ProjectView**: Main Kanban board view
- **BoardColumn**: Individual board column with tasks
- **TaskCard**: Task card with drag-and-drop
- **TaskDetailDialog**: Complete task details with:
  - Title & Description editing
  - Priority & Due date management
  - Labels assignment
  - Comments section
  - File attachments
  - Activity history
- **LabelsDialog**: Label management interface

## State Management

The app uses Pinia for state management with separate stores for:

- **auth**: User authentication and session
- **projects**: Project CRUD operations
- **boards**: Board/column management
- **tasks**: Task operations, comments, attachments
- **labels**: Label management

## Deployment

### Vercel (Recommended)

1. Push code to GitHub
2. Import project in Vercel
3. Set environment variables:
   - `VITE_API_URL`
   - `VITE_GOOGLE_CLIENT_ID`
4. Deploy

Or use the Vercel CLI:

```bash
npm install -g vercel
vercel
```

### Other Platforms

Build the project and deploy the `dist` folder to any static hosting service:

- Netlify
- GitHub Pages
- Firebase Hosting
- AWS S3 + CloudFront

## API Integration

The frontend communicates with the Go backend API. The API client is configured with:

- Automatic JWT token injection
- Request/response interceptors
- Error handling
- Automatic redirect to login on 401 errors

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `VITE_API_URL` | Backend API URL | Yes |
| `VITE_GOOGLE_CLIENT_ID` | Google OAuth Client ID | Yes |

## License

MIT
