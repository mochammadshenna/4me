# 4me Todos - Architecture Documentation

This document provides a comprehensive overview of the 4me Todos application architecture, including data flow, component interactions, and detailed explanations for developers.

## Table of Contents

1. [System Overview](#system-overview)
2. [Architecture Diagram](#architecture-diagram)
3. [Request Flow Control](#request-flow-control)
4. [Data Lineage Examples](#data-lineage-examples)
5. [Component Structure](#component-structure)
6. [Database Schema](#database-schema)
7. [Authentication Flow](#authentication-flow)
8. [API Design Patterns](#api-design-patterns)

---

## System Overview

4me Todos is a full-stack task management application built with:

- **Frontend**: Vue.js 3 + Vite (SPA)
- **Backend**: Go + Gin Framework (REST API)
- **Database**: PostgreSQL
- **Storage**: Supabase Storage
- **Authentication**: JWT + Google OAuth 2.0

### High-Level Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        Browser[Web Browser]
        VueApp[Vue.js Application]
    end
    
    subgraph "Application Layer"
        API[Go API Server<br/>Gin Framework]
        Auth[Auth Middleware<br/>JWT Validation]
        Handlers[Request Handlers]
    end
    
    subgraph "Data Layer"
        DB[(PostgreSQL<br/>Database)]
        Supabase[Supabase Storage<br/>File Attachments]
    end
    
    subgraph "External Services"
        Google[Google OAuth 2.0]
    end
    
    Browser --> VueApp
    VueApp -->|HTTP/REST| API
    API --> Auth
    Auth --> Handlers
    Handlers --> DB
    Handlers --> Supabase
    API -->|OAuth Flow| Google
    
    style Browser fill:#e1f5ff
    style VueApp fill:#4fc3f7
    style API fill:#81c784
    style DB fill:#ffd54f
    style Supabase fill:#ff8a65
    style Google fill:#f48fb1
```

---

## Request Flow Control

### 1. Complete Request Lifecycle

This diagram shows the complete path of a request from the browser to the database and back.

```mermaid
sequenceDiagram
    participant Browser
    participant VueRouter as Vue Router
    participant PiniaStore as Pinia Store
    participant AxiosClient as Axios Client
    participant GinRouter as Gin Router
    participant AuthMiddleware as Auth Middleware
    participant Handler
    participant Database as PostgreSQL
    
    Browser->>VueRouter: Navigate to /projects/1
    VueRouter->>VueRouter: Check Auth Guard
    VueRouter->>PiniaStore: Verify Token Exists
    
    alt Token Valid
        VueRouter->>Browser: Render ProjectView
        Browser->>PiniaStore: fetchProject(1)
        PiniaStore->>AxiosClient: GET /api/projects/1
        
        Note over AxiosClient: Request Interceptor<br/>Adds Authorization Header
        
        AxiosClient->>GinRouter: HTTP GET /api/projects/1<br/>Header: Bearer {token}
        
        GinRouter->>AuthMiddleware: Validate Request
        AuthMiddleware->>AuthMiddleware: Extract JWT Token
        AuthMiddleware->>AuthMiddleware: Validate Token Signature
        AuthMiddleware->>AuthMiddleware: Check Expiration
        
        alt Token Valid
            AuthMiddleware->>Handler: Set userID in Context<br/>Continue Request
            Handler->>Handler: Parse Project ID from URL
            Handler->>Database: SELECT * FROM projects<br/>WHERE id = $1 AND user_id = $2
            Database-->>Handler: Return Project Data
            Handler-->>GinRouter: JSON Response (200 OK)
        else Token Invalid
            AuthMiddleware-->>GinRouter: JSON Error (401 Unauthorized)
        end
        
        GinRouter-->>AxiosClient: HTTP Response
        
        Note over AxiosClient: Response Interceptor<br/>Handle Errors
        
        AxiosClient-->>PiniaStore: Return Data/Error
        PiniaStore-->>Browser: Update UI State
    else No Token
        VueRouter->>Browser: Redirect to /login
    end
```

### 2. Detailed Handler Flow (Task Creation Example)

```mermaid
flowchart TD
    Start([HTTP POST /api/boards/:id/tasks]) --> ParseRequest[Parse Request Body]
    
    ParseRequest --> ValidateInput{Valid Input?}
    ValidateInput -->|No| ReturnError400[Return 400 Bad Request]
    ValidateInput -->|Yes| ExtractUserID[Extract userID from Context]
    
    ExtractUserID --> ParseBoardID[Parse Board ID from URL]
    ParseBoardID --> CheckOwnership[Check Board Ownership]
    
    CheckOwnership --> OwnershipQuery[(Query: SELECT EXISTS<br/>boards + projects<br/>WHERE board.id = $1<br/>AND project.user_id = $2)]
    
    OwnershipQuery --> IsOwner{User Owns<br/>Board?}
    IsOwner -->|No| ReturnError404[Return 404 Not Found]
    IsOwner -->|Yes| BeginTransaction[Begin Database Transaction]
    
    BeginTransaction --> InsertTask[(INSERT INTO tasks<br/>board_id, title, description,<br/>priority, due_date)]
    
    InsertTask --> TaskCreated{Task<br/>Created?}
    TaskCreated -->|Error| RollbackTx[Rollback Transaction]
    RollbackTx --> ReturnError500[Return 500 Internal Error]
    
    TaskCreated -->|Success| HasLabels{Has<br/>Labels?}
    HasLabels -->|Yes| InsertLabels[(INSERT INTO task_labels<br/>FOR EACH label_id)]
    HasLabels -->|No| CreateHistory
    
    InsertLabels --> CreateHistory[(INSERT INTO task_history<br/>action: 'created'<br/>changes_json)]
    
    CreateHistory --> CommitTx[Commit Transaction]
    CommitTx --> ReturnSuccess[Return 201 Created<br/>with Task JSON]
    
    ReturnSuccess --> End([Response Sent])
    ReturnError400 --> End
    ReturnError404 --> End
    ReturnError500 --> End
    
    style Start fill:#4caf50
    style End fill:#f44336
    style BeginTransaction fill:#ff9800
    style CommitTx fill:#ff9800
    style InsertTask fill:#2196f3
    style InsertLabels fill:#2196f3
    style CreateHistory fill:#2196f3
```

---

## Data Lineage Examples

### Example 1: Task Creation - Data Transformation Journey

Let's trace how task data flows through the system when creating a new task.

```mermaid
graph LR
    subgraph "Frontend - User Input"
        A1[User fills form:<br/>title: 'Fix bug'<br/>priority: 'high'<br/>labels: 1,3]
    end
    
    subgraph "Frontend - Vue Component"
        A2[TaskForm.vue<br/>Form Data Object]
        A3[Validation:<br/>title.length > 0<br/>priority in list]
    end
    
    subgraph "Frontend - Store"
        B1[tasksStore.createTask<br/>boardId: 5<br/>taskData: object]
        B2[API Client Request:<br/>POST /api/boards/5/tasks<br/>Body: JSON]
    end
    
    subgraph "Backend - Network Layer"
        C1[Gin Router<br/>Route Match:<br/>POST /boards/:id/tasks]
        C2[Auth Middleware<br/>Extract userID: 42<br/>from JWT token]
    end
    
    subgraph "Backend - Handler"
        D1[TaskHandler.Create<br/>c.Param id: '5'<br/>c.Get userID: 42]
        D2[Parse JSON Body:<br/>CreateTaskRequest struct]
        D3[Set defaults:<br/>priority: 'high'<br/>status: 'todo'<br/>position: 0]
    end
    
    subgraph "Backend - Database Layer"
        E1[Ownership Check:<br/>SELECT EXISTS<br/>board.id=5 AND<br/>project.user_id=42]
        E2[Transaction Start]
        E3[INSERT INTO tasks:<br/>board_id: 5<br/>title: 'Fix bug'<br/>status: 'todo'<br/>priority: 'high']
        E4[RETURNING:<br/>id: 123<br/>created_at: timestamp]
    end
    
    subgraph "Backend - Related Data"
        F1[INSERT task_labels:<br/>task_id: 123<br/>label_id: 1]
        F2[INSERT task_labels:<br/>task_id: 123<br/>label_id: 3]
        F3[INSERT task_history:<br/>task_id: 123<br/>user_id: 42<br/>action: 'created']
    end
    
    subgraph "Backend - Response"
        G1[Transaction Commit]
        G2[Build Response:<br/>models.Task JSON<br/>status: 201]
    end
    
    subgraph "Frontend - Update"
        H1[Store receives:<br/>Task object<br/>id: 123]
        H2[Add to tasks array<br/>tasks.value.push]
        H3[UI re-renders:<br/>New task card appears<br/>in board column]
    end
    
    A1 --> A2 --> A3
    A3 --> B1 --> B2
    B2 --> C1 --> C2
    C2 --> D1 --> D2 --> D3
    D3 --> E1 --> E2 --> E3 --> E4
    E4 --> F1 --> F2 --> F3
    F1 & F2 & F3 --> G1 --> G2
    G2 --> H1 --> H2 --> H3
    
    style A1 fill:#e3f2fd
    style E3 fill:#fff9c4
    style F1 fill:#fff9c4
    style F2 fill:#fff9c4
    style F3 fill:#fff9c4
    style H3 fill:#c8e6c9
```

### Example 2: Authentication Token Lifecycle

```mermaid
stateDiagram-v2
    [*] --> UserInput: User enters credentials
    
    UserInput --> FrontendValidation: username + password
    FrontendValidation --> APIRequest: POST /api/auth/login
    
    APIRequest --> BackendValidation: JSON body received
    BackendValidation --> DatabaseQuery: Query user by username
    
    DatabaseQuery --> PasswordCheck: User found
    PasswordCheck --> TokenGeneration: Password matches (bcrypt)
    
    TokenGeneration --> JWTCreation: Generate JWT with:<br/>- userID<br/>- username<br/>- email<br/>- expiry: 24h
    
    JWTCreation --> SignToken: Sign with JWT_SECRET
    SignToken --> ResponseSent: Return token + user data
    
    ResponseSent --> StorageLocal: localStorage.setItem('token')
    StorageLocal --> PiniaStore: authStore.setAuthData()
    
    PiniaStore --> AxiosInterceptor: Axios request interceptor
    AxiosInterceptor --> RequestHeaders: Add header:<br/>Authorization: Bearer {token}
    
    RequestHeaders --> BackendMiddleware: Every protected request
    BackendMiddleware --> TokenValidation: Validate signature<br/>Check expiration
    
    TokenValidation --> SetContext: Valid: Set userID in context
    TokenValidation --> Unauthorized: Invalid: Return 401
    
    SetContext --> HandlerExecution: Handler processes request
    Unauthorized --> FrontendLogout: Clear token, redirect
    
    FrontendLogout --> [*]
    HandlerExecution --> [*]: Request complete
    
    note right of JWTCreation
        JWT Payload:
        {
          "user_id": 42,
          "username": "john",
          "email": "john@example.com",
          "exp": 1234567890,
          "iat": 1234481490
        }
    end note
```

---

## Component Structure

### Backend Architecture (Go)

```mermaid
graph TB
    subgraph "Entry Point"
        Main[cmd/api/main.go<br/>Application Bootstrap]
    end
    
    subgraph "Configuration"
        Config[internal/config<br/>Environment Variables<br/>Database URL<br/>JWT Secret<br/>OAuth Credentials]
    end
    
    subgraph "Database Layer"
        DB[internal/database<br/>Connection Pool<br/>Migrations<br/>Query Execution]
    end
    
    subgraph "Middleware"
        CORS[middleware/cors.go<br/>CORS Headers<br/>Preflight Requests]
        Auth[middleware/auth.go<br/>JWT Validation<br/>Context Injection]
    end
    
    subgraph "Handlers - Business Logic"
        AuthH[handlers/auth.go<br/>Register<br/>Login<br/>Google OAuth<br/>Me]
        ProjectsH[handlers/projects.go<br/>CRUD Operations<br/>Ownership Validation]
        BoardsH[handlers/boards.go<br/>Board Management<br/>Position Ordering]
        TasksH[handlers/tasks.go<br/>Task CRUD<br/>Move Between Boards<br/>History Tracking]
        LabelsH[handlers/labels.go<br/>Label Management<br/>Color Assignment]
        CommentsH[handlers/comments.go<br/>Comment CRUD<br/>User Association]
        AttachmentsH[handlers/attachments.go<br/>File Upload<br/>Supabase Integration]
    end
    
    subgraph "Models"
        Models[internal/models<br/>Data Structures<br/>Request DTOs<br/>Response DTOs<br/>Validation Tags]
    end
    
    subgraph "Utilities"
        JWT[utils/jwt.go<br/>Token Generation<br/>Token Validation<br/>Claims Extraction]
        Password[utils/password.go<br/>Hash Password<br/>Compare Password<br/>bcrypt]
    end
    
    Main --> Config
    Main --> DB
    Main --> CORS
    Main --> Auth
    Main --> AuthH
    Main --> ProjectsH
    Main --> BoardsH
    Main --> TasksH
    Main --> LabelsH
    Main --> CommentsH
    Main --> AttachmentsH
    
    AuthH --> Models
    ProjectsH --> Models
    BoardsH --> Models
    TasksH --> Models
    LabelsH --> Models
    CommentsH --> Models
    AttachmentsH --> Models
    
    AuthH --> DB
    ProjectsH --> DB
    BoardsH --> DB
    TasksH --> DB
    LabelsH --> DB
    CommentsH --> DB
    AttachmentsH --> DB
    
    AuthH --> JWT
    AuthH --> Password
    Auth --> JWT
    
    AttachmentsH --> Config
    
    style Main fill:#4caf50
    style DB fill:#ff9800
    style Models fill:#2196f3
```

### Frontend Architecture (Vue.js)

```mermaid
graph TB
    subgraph "Application Entry"
        MainJS[main.js<br/>Vue App Initialization<br/>Plugin Registration]
        AppVue[App.vue<br/>Root Component<br/>Router View]
    end
    
    subgraph "Routing"
        Router[router/index.js<br/>Route Definitions<br/>Auth Guards<br/>Navigation Guards]
    end
    
    subgraph "State Management - Pinia"
        AuthStore[stores/auth.js<br/>User Session<br/>Token Management<br/>Login/Logout<br/>Google OAuth]
        ProjectsStore[stores/projects.js<br/>Projects List<br/>Current Project<br/>CRUD Operations]
        BoardsStore[stores/boards.js<br/>Boards Array<br/>Board CRUD<br/>Position Updates]
        TasksStore[stores/tasks.js<br/>Tasks by Board<br/>Task CRUD<br/>Move Task<br/>Comments/Attachments]
        LabelsStore[stores/labels.js<br/>Project Labels<br/>Label CRUD]
    end
    
    subgraph "API Layer"
        APIClient[api/client.js<br/>Axios Instance<br/>Request Interceptor<br/>Response Interceptor<br/>Error Handling]
    end
    
    subgraph "Views - Pages"
        LoginView[LoginView.vue<br/>Login Form<br/>OAuth Button]
        RegisterView[RegisterView.vue<br/>Registration Form]
        DashboardView[DashboardView.vue<br/>Projects Grid<br/>Create Project]
        ProjectView[ProjectView.vue<br/>Kanban Board<br/>Boards/Tasks Display]
        AuthCallback[AuthCallbackView.vue<br/>OAuth Redirect Handler]
    end
    
    subgraph "Components - Board"
        BoardColumn[BoardColumn.vue<br/>Column Header<br/>Task List<br/>Drag & Drop Container]
        TaskCard[TaskCard.vue<br/>Task Display<br/>Labels/Priority<br/>Metadata Icons]
    end
    
    subgraph "Components - Task"
        TaskDetail[TaskDetailDialog.vue<br/>Full Task Editor<br/>Comments Section<br/>Attachments Upload<br/>History Timeline]
        LabelsDialog[LabelsDialog.vue<br/>Label Management<br/>Color Picker<br/>CRUD Interface]
    end
    
    subgraph "Plugins"
        Vuetify[plugins/vuetify.js<br/>UI Component Library<br/>Theme Configuration<br/>Material Design Icons]
    end
    
    MainJS --> AppVue
    MainJS --> Router
    MainJS --> AuthStore
    MainJS --> ProjectsStore
    MainJS --> BoardsStore
    MainJS --> TasksStore
    MainJS --> LabelsStore
    MainJS --> Vuetify
    
    AppVue --> LoginView
    AppVue --> RegisterView
    AppVue --> DashboardView
    AppVue --> ProjectView
    AppVue --> AuthCallback
    
    Router --> LoginView
    Router --> DashboardView
    Router --> ProjectView
    
    DashboardView --> ProjectsStore
    ProjectView --> BoardsStore
    ProjectView --> TasksStore
    ProjectView --> LabelsStore
    ProjectView --> BoardColumn
    
    BoardColumn --> TaskCard
    BoardColumn --> TasksStore
    
    TaskCard --> TaskDetail
    ProjectView --> TaskDetail
    ProjectView --> LabelsDialog
    
    TaskDetail --> TasksStore
    LabelsDialog --> LabelsStore
    
    AuthStore --> APIClient
    ProjectsStore --> APIClient
    BoardsStore --> APIClient
    TasksStore --> APIClient
    LabelsStore --> APIClient
    
    style MainJS fill:#4caf50
    style AuthStore fill:#ff9800
    style ProjectsStore fill:#ff9800
    style BoardsStore fill:#ff9800
    style TasksStore fill:#ff9800
    style APIClient fill:#2196f3
```

---

## Database Schema

### Entity Relationship Diagram

```mermaid
erDiagram
    USERS ||--o{ PROJECTS : creates
    USERS ||--o{ COMMENTS : writes
    USERS ||--o{ TASK_HISTORY : records
    
    PROJECTS ||--o{ BOARDS : contains
    PROJECTS ||--o{ LABELS : defines
    
    BOARDS ||--o{ TASKS : organizes
    
    TASKS ||--o{ COMMENTS : has
    TASKS ||--o{ ATTACHMENTS : includes
    TASKS ||--o{ TASK_HISTORY : tracks
    TASKS ||--o{ TASK_LABELS : tagged_with
    
    LABELS ||--o{ TASK_LABELS : applied_to
    
    USERS {
        int id PK
        string username UK
        string email UK
        string password_hash
        string google_id UK
        string avatar_url
        timestamp created_at
        timestamp updated_at
    }
    
    PROJECTS {
        int id PK
        int user_id FK
        string name
        text description
        string color
        timestamp created_at
        timestamp updated_at
    }
    
    BOARDS {
        int id PK
        int project_id FK
        string name
        int position
        timestamp created_at
        timestamp updated_at
    }
    
    TASKS {
        int id PK
        int board_id FK
        string title
        text description
        string status
        string priority
        int assignee_id FK
        timestamp due_date
        int position
        timestamp created_at
        timestamp updated_at
    }
    
    LABELS {
        int id PK
        int project_id FK
        string name
        string color
        timestamp created_at
    }
    
    TASK_LABELS {
        int task_id FK
        int label_id FK
    }
    
    COMMENTS {
        int id PK
        int task_id FK
        int user_id FK
        text content
        timestamp created_at
        timestamp updated_at
    }
    
    ATTACHMENTS {
        int id PK
        int task_id FK
        string filename
        string file_url
        string file_type
        bigint size
        timestamp uploaded_at
    }
    
    TASK_HISTORY {
        int id PK
        int task_id FK
        int user_id FK
        string action
        jsonb changes_json
        timestamp created_at
    }
```

### Table Indexes and Constraints

```sql
-- Critical Indexes for Performance

-- Users
CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;

-- Projects
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_created_at ON projects(created_at DESC);

-- Boards
CREATE INDEX idx_boards_project_id ON boards(project_id);
CREATE INDEX idx_boards_position ON boards(project_id, position);

-- Tasks
CREATE INDEX idx_tasks_board_id ON tasks(board_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id) WHERE assignee_id IS NOT NULL;
CREATE INDEX idx_tasks_due_date ON tasks(due_date) WHERE due_date IS NOT NULL;
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_position ON tasks(board_id, position);

-- Task Labels (Composite for JOIN optimization)
CREATE INDEX idx_task_labels_task_id ON task_labels(task_id);
CREATE INDEX idx_task_labels_label_id ON task_labels(label_id);

-- Comments
CREATE INDEX idx_comments_task_id ON comments(task_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);

-- Attachments
CREATE INDEX idx_attachments_task_id ON attachments(task_id);

-- Task History
CREATE INDEX idx_task_history_task_id ON task_history(task_id);
CREATE INDEX idx_task_history_created_at ON task_history(created_at DESC);
```

---

## Authentication Flow

### Complete OAuth 2.0 Flow (Google)

```mermaid
sequenceDiagram
    participant User
    participant Browser
    participant VueApp as Vue App
    participant GoAPI as Go API
    participant Google as Google OAuth
    participant Database as PostgreSQL
    
    User->>Browser: Click "Sign in with Google"
    Browser->>VueApp: authStore.googleLogin()
    VueApp->>GoAPI: GET /api/auth/google
    
    GoAPI->>GoAPI: Generate OAuth URL<br/>with client_id, redirect_uri, scopes
    GoAPI-->>VueApp: Return {url: "https://accounts.google.com/..."}
    
    VueApp->>Browser: window.location.href = url
    Browser->>Google: Redirect to Google OAuth
    
    User->>Google: Authenticate & Approve Permissions
    
    Google->>Browser: Redirect with auth code<br/>to redirect_uri
    Browser->>GoAPI: GET /api/auth/google/callback?code=xyz
    
    GoAPI->>Google: Exchange code for access_token
    Google-->>GoAPI: Return access_token
    
    GoAPI->>Google: GET /oauth2/v2/userinfo<br/>Header: Bearer {access_token}
    Google-->>GoAPI: Return user profile<br/>{id, email, name, picture}
    
    GoAPI->>Database: SELECT user WHERE google_id = ?
    
    alt User Exists
        Database-->>GoAPI: Return existing user
    else New User
        GoAPI->>Database: INSERT INTO users<br/>(username, email, google_id, avatar_url)
        Database-->>GoAPI: Return new user with id
    end
    
    GoAPI->>GoAPI: Generate JWT token<br/>Sign with JWT_SECRET
    GoAPI->>GoAPI: Generate refresh_token
    
    GoAPI->>Browser: Redirect to frontend<br/>/auth/callback?token=jwt&refresh_token=rt
    
    Browser->>VueApp: AuthCallbackView mounted
    VueApp->>VueApp: Extract tokens from URL
    VueApp->>VueApp: localStorage.setItem('token', jwt)
    VueApp->>VueApp: authStore.handleGoogleCallback()
    VueApp->>GoAPI: GET /api/auth/me<br/>Header: Bearer {jwt}
    
    GoAPI->>Database: SELECT user WHERE id = ?
    Database-->>GoAPI: Return user details
    GoAPI-->>VueApp: Return user object
    
    VueApp->>VueApp: authStore.user = userData
    VueApp->>Browser: router.push('/')
    Browser->>User: Show Dashboard
```

### JWT Token Structure

```javascript
// JWT Header
{
  "alg": "HS256",
  "typ": "JWT"
}

// JWT Payload
{
  "user_id": 42,
  "username": "john_doe",
  "email": "john@example.com",
  "exp": 1735689600,  // Expiration (24 hours)
  "iat": 1735603200   // Issued at
}

// JWT Signature
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  JWT_SECRET
)
```

---

## API Design Patterns

### RESTful Resource Hierarchy

```mermaid
graph LR
    subgraph "Resource Hierarchy"
        Users[Users]
        Projects[Projects]
        Boards[Boards]
        Tasks[Tasks]
        Labels[Labels]
        Comments[Comments]
        Attachments[Attachments]
        History[Task History]
    end
    
    Users -->|owns| Projects
    Projects -->|contains| Boards
    Projects -->|defines| Labels
    Boards -->|organizes| Tasks
    Tasks -->|has| Comments
    Tasks -->|includes| Attachments
    Tasks -->|tracks| History
    Tasks -->|tagged_with| Labels
    
    style Users fill:#e1bee7
    style Projects fill:#c5e1a5
    style Boards fill:#90caf9
    style Tasks fill:#ffcc80
```

### API Endpoint Patterns

```
# Authentication
POST   /api/auth/register          # Create new user account
POST   /api/auth/login             # Authenticate user
GET    /api/auth/google            # Get OAuth URL
GET    /api/auth/google/callback   # Handle OAuth callback
GET    /api/auth/me                # Get current user [Protected]

# Projects (All Protected)
POST   /api/projects               # Create project
GET    /api/projects               # List user's projects
GET    /api/projects/:id           # Get project by ID
PUT    /api/projects/:id           # Update project
DELETE /api/projects/:id           # Delete project

# Boards (Nested under Projects)
POST   /api/projects/:id/boards    # Create board in project
GET    /api/projects/:id/boards    # List project's boards
PUT    /api/boards/:id             # Update board
DELETE /api/boards/:id             # Delete board

# Tasks (Nested under Boards)
POST   /api/boards/:id/tasks       # Create task in board
GET    /api/tasks/:id              # Get task details
PUT    /api/tasks/:id              # Update task
PATCH  /api/tasks/:id/move         # Move task to another board
DELETE /api/tasks/:id              # Delete task

# Task Sub-resources
GET    /api/tasks/:id/history      # Get task change history
POST   /api/tasks/:id/comments     # Add comment to task
GET    /api/tasks/:id/comments     # List task comments
POST   /api/tasks/:id/attachments  # Upload file to task
GET    /api/tasks/:id/attachments  # List task attachments

# Labels (Project-scoped)
POST   /api/projects/:id/labels    # Create label in project
GET    /api/projects/:id/labels    # List project labels
PUT    /api/labels/:id             # Update label
DELETE /api/labels/:id             # Delete label

# Comments (Direct access)
PUT    /api/comments/:id           # Update comment
DELETE /api/comments/:id           # Delete comment

# Attachments (Direct access)
DELETE /api/attachments/:id        # Delete attachment
```

### Error Response Format

```json
{
  "error": "Descriptive error message",
  "code": "ERROR_CODE",
  "details": {
    "field": "Additional context"
  }
}
```

### Success Response Format

```json
{
  "id": 123,
  "field1": "value1",
  "field2": "value2",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

---

## Step-by-Step: Creating a Task

Let's walk through the complete process of creating a task, from user interaction to database storage.

### Step 1: User Interaction (Frontend)

**Location**: `frontend/src/components/board/BoardColumn.vue`

```javascript
// User clicks "Add Task" button and types title
const newTaskTitle = ref('')

async function handleAddTask() {
  // Step 1a: Validate input locally
  if (!newTaskTitle.value.trim()) return
  
  // Step 1b: Call Pinia store method
  addingTask.value = true
  const result = await tasksStore.createTask(props.board.id, {
    title: newTaskTitle.value,
    priority: 'medium',
  })
  
  // Step 1c: Handle response
  if (result.success) {
    tasks.value.push(result.data)  // Optimistic UI update
    newTaskTitle.value = ''
  }
  addingTask.value = false
}
```

### Step 2: Store Layer (State Management)

**Location**: `frontend/src/stores/tasks.js`

```javascript
async function createTask(boardId, taskData) {
  try {
    // Step 2a: Make API call via Axios client
    const response = await apiClient.post(`/boards/${boardId}/tasks`, taskData)
    
    // Step 2b: Update local state
    if (!tasks.value[boardId]) {
      tasks.value[boardId] = []
    }
    tasks.value[boardId].push(response.data)
    
    // Step 2c: Return success
    return { success: true, data: response.data }
  } catch (error) {
    return { 
      success: false, 
      error: error.response?.data?.error || 'Failed to create task' 
    }
  }
}
```

### Step 3: HTTP Request (API Client)

**Location**: `frontend/src/api/client.js`

```javascript
// Axios request interceptor adds authentication
apiClient.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.token) {
    // Step 3a: Add JWT token to request header
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

// Request sent:
// POST /api/boards/5/tasks
// Headers: { Authorization: 'Bearer eyJhbGc...' }
// Body: { title: 'Fix bug', priority: 'medium' }
```

### Step 4: Backend Routing (Gin)

**Location**: `backend/cmd/api/main.go`

```go
// Step 4a: Router matches route pattern
protected.POST("/boards/:id/tasks", taskHandler.Create)

// Step 4b: Auth middleware validates JWT
protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
```

### Step 5: Authentication Middleware

**Location**: `backend/internal/middleware/auth.go`

```go
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
  return func(c *gin.Context) {
    // Step 5a: Extract token from Authorization header
    authHeader := c.GetHeader("Authorization")
    tokenParts := strings.Split(authHeader, " ")
    token := tokenParts[1]
    
    // Step 5b: Validate JWT signature and expiration
    claims, err := utils.ValidateToken(token, jwtSecret)
    if err != nil {
      c.JSON(401, gin.H{"error": "Invalid token"})
      c.Abort()
      return
    }
    
    // Step 5c: Set user context for handler
    c.Set("userID", claims.UserID)  // userID = 42
    c.Next()
  }
}
```

### Step 6: Handler Business Logic

**Location**: `backend/internal/handlers/tasks.go`

```go
func (h *TaskHandler) Create(c *gin.Context) {
  // Step 6a: Extract user ID from context
  userID, _ := c.Get("userID")  // 42
  
  // Step 6b: Parse board ID from URL parameter
  boardID, err := strconv.Atoi(c.Param("id"))  // 5
  
  // Step 6c: Verify user owns the board's project
  var exists bool
  h.db.Pool.QueryRow(ctx,
    `SELECT EXISTS(
      SELECT 1 FROM boards b 
      JOIN projects p ON b.project_id = p.id 
      WHERE b.id = $1 AND p.user_id = $2
    )`, boardID, userID).Scan(&exists)
  
  if !exists {
    c.JSON(404, gin.H{"error": "Board not found"})
    return
  }
  
  // Step 6d: Parse and validate request body
  var req models.CreateTaskRequest
  c.ShouldBindJSON(&req)
  // req.Title = "Fix bug"
  // req.Priority = "medium"
  
  // Step 6e: Set defaults
  if req.Priority == "" {
    req.Priority = "medium"
  }
  
  // Step 6f: Begin database transaction
  tx, _ := h.db.Pool.Begin(ctx)
  defer tx.Rollback(ctx)
  
  // Step 6g: Insert task into database
  var task models.Task
  tx.QueryRow(ctx,
    `INSERT INTO tasks (board_id, title, priority, status) 
     VALUES ($1, $2, $3, 'todo') 
     RETURNING id, board_id, title, description, status, priority, 
               position, created_at, updated_at`,
    boardID, req.Title, req.Priority).
    Scan(&task.ID, &task.BoardID, &task.Title, &task.Description, 
          &task.Status, &task.Priority, &task.Position, 
          &task.CreatedAt, &task.UpdatedAt)
  // task.ID = 123 (auto-generated)
  
  // Step 6h: Create history entry
  changes := map[string]interface{}{
    "action": "created",
    "title":  req.Title,
  }
  changesJSON, _ := json.Marshal(changes)
  tx.Exec(ctx,
    "INSERT INTO task_history (task_id, user_id, action, changes_json) VALUES ($1, $2, $3, $4)",
    task.ID, userID, "created", changesJSON)
  
  // Step 6i: Commit transaction
  tx.Commit(ctx)
  
  // Step 6j: Return success response
  c.JSON(201, task)
}
```

### Step 7: Database Execution

```sql
-- Step 7a: Ownership verification query
SELECT EXISTS(
  SELECT 1 FROM boards b 
  JOIN projects p ON b.project_id = p.id 
  WHERE b.id = 5 AND p.user_id = 42
); -- Result: true

-- Step 7b: Task insertion
INSERT INTO tasks (board_id, title, priority, status) 
VALUES (5, 'Fix bug', 'medium', 'todo') 
RETURNING id, board_id, title, description, status, priority, 
          position, created_at, updated_at;
-- Returns: id=123, board_id=5, title='Fix bug', ...

-- Step 7c: History insertion
INSERT INTO task_history (task_id, user_id, action, changes_json) 
VALUES (123, 42, 'created', '{"action":"created","title":"Fix bug"}');
```

### Step 8: Response Journey Back

```mermaid
graph RL
    DB[(Database)] -->|Task Record| Handler
    Handler -->|JSON Serialization| Gin
    Gin -->|HTTP 201 Response| Network
    Network -->|Axios Response| APIClient
    APIClient -->|Response Interceptor| Store
    Store -->|Reactive Update| Component
    Component -->|DOM Update| Browser
```

**Final Response**:
```json
{
  "id": 123,
  "board_id": 5,
  "title": "Fix bug",
  "description": null,
  "status": "todo",
  "priority": "medium",
  "assignee_id": null,
  "due_date": null,
  "position": 0,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

---

## Performance Considerations

### Query Optimization

```mermaid
graph TD
    A[Request arrives] --> B{Is data cached?}
    B -->|Yes| C[Return from Pinia Store]
    B -->|No| D[Query Database]
    D --> E{Use Index?}
    E -->|Yes| F[Fast Index Scan]
    E -->|No| G[Slow Table Scan]
    F --> H[Return Results]
    G --> H
    H --> I[Cache in Store]
    I --> C
```

### Key Optimizations

1. **Database Indexes**: All foreign keys and frequently queried columns are indexed
2. **Frontend Caching**: Pinia stores cache API responses
3. **Optimistic Updates**: UI updates immediately before server confirmation
4. **Pagination**: Large lists are paginated (ready for implementation)
5. **Connection Pooling**: Database connection pool reuses connections
6. **Transaction Batching**: Related operations grouped in transactions

---

## Security Architecture

```mermaid
flowchart TD
    Request[Incoming Request] --> HTTPS{HTTPS?}
    HTTPS -->|No| Reject1[Reject - Use HTTPS]
    HTTPS -->|Yes| CORS{Valid Origin?}
    CORS -->|No| Reject2[Reject - CORS Error]
    CORS -->|Yes| Auth{Has Token?}
    Auth -->|No| Reject3[401 Unauthorized]
    Auth -->|Yes| ValidateJWT{Valid JWT?}
    ValidateJWT -->|No| Reject4[401 Invalid Token]
    ValidateJWT -->|Yes| ExtractUser[Extract User ID]
    ExtractUser --> CheckOwnership{User Owns<br/>Resource?}
    CheckOwnership -->|No| Reject5[404 Not Found]
    CheckOwnership -->|Yes| Sanitize[Sanitize Input]
    Sanitize --> ParamQuery[Parameterized Query]
    ParamQuery --> Execute[Execute Safely]
    Execute --> Response[Return Response]
    
    style Request fill:#4caf50
    style Response fill:#4caf50
    style Reject1 fill:#f44336
    style Reject2 fill:#f44336
    style Reject3 fill:#f44336
    style Reject4 fill:#f44336
    style Reject5 fill:#f44336
```

---

## Deployment Architecture

```mermaid
graph TB
    subgraph "Client Devices"
        Mobile[Mobile Browser]
        Desktop[Desktop Browser]
    end
    
    subgraph "CDN / Edge Network"
        Vercel[Vercel Edge Network<br/>Static Assets<br/>Vue.js App]
    end
    
    subgraph "Application Servers"
        API1[Go API Server 1<br/>Railway/Render]
        API2[Go API Server 2<br/>Load Balanced]
    end
    
    subgraph "Database Layer"
        Primary[(Supabase PostgreSQL<br/>Primary Database)]
        Replica[(Read Replica<br/>Optional)]
    end
    
    subgraph "Storage"
        SupaStorage[Supabase Storage<br/>File Attachments]
    end
    
    subgraph "External Services"
        GoogleAuth[Google OAuth 2.0]
    end
    
    Mobile --> Vercel
    Desktop --> Vercel
    Vercel --> API1
    Vercel --> API2
    API1 --> Primary
    API2 --> Primary
    API1 --> Replica
    API2 --> Replica
    API1 --> SupaStorage
    API2 --> SupaStorage
    API1 --> GoogleAuth
    API2 --> GoogleAuth
    
    style Vercel fill:#4caf50
    style API1 fill:#2196f3
    style API2 fill:#2196f3
    style Primary fill:#ff9800
    style SupaStorage fill:#e91e63
```

---

## Conclusion

This architecture provides:

- **Scalability**: Stateless API servers can scale horizontally
- **Security**: Multi-layer authentication and authorization
- **Performance**: Indexed queries, caching, optimistic updates
- **Maintainability**: Clear separation of concerns, modular design
- **Developer Experience**: Type-safe models, comprehensive error handling

For questions or contributions, refer to the main [README.md](./README.md) and [DEPLOYMENT.md](./DEPLOYMENT.md) files.

