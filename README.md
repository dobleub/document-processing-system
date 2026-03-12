# Document Processing System

A modern, full-stack document processing platform built with **Next.js**, **Go**, **MongoDB**, and **Google Gemini AI**. This monorepo provides a comprehensive solution for uploading, processing, and summarizing documents with real-time status updates. Powered by Nx.

<a alt="Nx logo" href="https://nx.dev" target="_blank" rel="noreferrer"><img src="https://raw.githubusercontent.com/nrwl/nx/master/images/nx-logo.png" width="45"></a>

---

## Table of Contents

- [Project Overview](#project-overview)
- [System Architecture](#system-architecture)
- [Application Structure](#application-structure)
- [Dashboard Overview](#dashboard-overview)
- [API Overview](#api-overview)
- [Prerequisites](#prerequisites)
- [Local Development Setup](#local-development-setup)
- [Running Services](#running-services)
- [Environment Configuration](#environment-configuration)
- [Available Commands](#available-commands)
- [API Overview](#api-overview)
- [> API Endpoints](/apps/api/README.md#api-overview)

---

## Project Overview

Document Processing System is a develop-ready monorepo built with **Nx** that enables users to:

- 📤 **Upload documents** for processing (text files, PDFs, etc.)
- 🔄 **Track processing status** in real-time via WebSocket
- 🤖 **Summarize content** using Google Gemini AI, could be disabled via a toggle comment if there is not enough tokens in Gemini in __apps/api/src/processDomain/lib/fileProcessing.go:312-313__
    - It could be implementing via another container with a basic summarization model if you want to avoid using Gemini API.
- 📊 **Manage processes** through an intuitive web dashboard
- 💾 **Persist data** using MongoDB for process history and results
- 🔐 **Secure endpoints** with Bearer token authentication

---

## System Architecture

```mermaid
graph TB
    subgraph Client["Client Layer"]
        UI[Web Dashboard<br/>Next.js + React]
    end
    
    subgraph Server["Application Layer"]
        API[REST API<br/>Go + Gin]
        StateManager[State Manager<br/>In-memory + MongoDB]
        WS["WebSocket Server<br/>Real-time Updates"]
    end

    subgraph Datasource["DataSources"]
        subgraph Services["External Services"]
            GEMINI["Google Gemini AI<br/>Content Summarization"]
            MCP["Model Context Protocol<br/>AI Integration"]
        end
        
        subgraph Data["Data Layer"]
            DB[(MongoDB<br/>Process Storage)]
        end
    end
    
    UI -->|HTTP Requests| API
    UI -->|WebSocket| WS
    API --> GEMINI
    API --> MCP
    API --> StateManager
    WS --> StateManager
    StateManager --> DB
    
    classDef frontend fill:#61dafb,stroke:#333,color:#000
    classDef backend fill:#00add8,stroke:#333,color:#fff
    classDef state fill:#0089d8,stroke:#333,color:#fff
    classDef external fill:#ff9900,stroke:#333,color:#000
    classDef database fill:#13aa52,stroke:#333,color:#fff
    
    class UI frontend
    class API,WS backend
    class StateManager state
    class GEMINI,MCP external
    class DB database
```

### The system consists of three main layers:

1. **Client Layer**: The web dashboard built with Next.js and React that allows users to interact with the system. It's responsible for sending API requests, receiving real-time updates via WebSocket, and displaying process status and results. 

    * This framework is not essential and can be replaced with any frontend technology but I chose Next.js because of my developer experience.

2. **Application Layer**: The Go backend service that handles API requests, manages processing state, and communicates with external services. The API is built with the Gin framework for high performance and simplicity. It includes endpoints for starting/stopping processes, retrieving status/results, and a WebSocket server for real-time updates.

    * I chose Go for its performance and concurrency capabilities, which are ideal for handling multiple document processing tasks simultaneously. I <3 Go.

3. **Data Layer**: MongoDB for persisting process data and results. This allows for historical tracking and retrieval of processing outcomes.

    * MongoDB was chosen for its flexibility in handling unstructured data, which is common in document processing scenarios and for its ease of integration with Go.

---

## Application Structure

This project is organized as a monorepo with the following structure:

```
document-processing-system/
├── apps/
│   ├── api/                        # Go backend service
│   │   ├── main.go
│   │   ├── config/                 # Configuration management
│   │   ├── interfaces/             # General interfaces & contracts
│   │   ├── middlewares/            # Auth & request processing
│   │   ├── lib/                    # Database & utilities
│   │   ├── docs/                   # Swagger API documentation
│   │   ├── targetFiles/            # Target files for processing
│   │   └── src/                    # Source code organized by domain
│   │       ├── processDomain/      # Source files for processing files content and managing processes states
│   │       ├── summarizeDomain/    # Source files for summarization files content
│   │       └── websocketDomain/    # Source files for WebSocket communication and real-time updates
│   │
│   └── dashboard/                  # Next.js web application
│       ├── src/
│       │   ├── app/                # Next.js app directory
│       │   ├── components/         # Reusable React components
│       │   ├── context/            # React context providers
│       │   ├── hooks/              # Custom React hooks
│       │   └── lib/                # API client & utilities
│       ├── public/                 # Static assets
│       └── tailwind.config.js      # Tailwind CSS configuration
│
├── db/                             # Database configuration
├── docker compose.yml              # Docker services orchestration
├── package.json                    # Root dependencies
└── nx.json                         # Nx workspace configuration
```

### Environment Variables

The project uses environment variables for configuration. Each app inside `apps/**` has its own `.env` file for local development, and these can be overridden in production environments.

- `apps/api/.env` for API configuration (authentication, database, AI keys)
    ```env
    API_AUTH_TOKEN=your-secure-token-here
    GEMINI_API_KEY=your-gemini-api-key-here
    MONGODB_URI=mongodb://localhost:27017 # Override with internal Docker network hostname in Docker Compose
    MONGODB_DB=document_processing
    ```
- `apps/dashboard/.env` for dashboard configuration (API base URL, WebSocket URL)
    ```env
    NEXT_PUBLIC_API_URL=http://dps-api:8080 # Use internal Docker network hostname for API in Docker Compose, override with API_URL_FRONT in production
    NEXT_PUBLIC_API_URL_FRONT=http://localhost:8080
    NEXT_PUBLIC_API_AUTH_TOKEN=your-secure-token-here
    ```

---

## Dashboard Overview

The **Dashboard** is a responsive admin interface built with **Next.js 16**, **React 19**, and **TypeScript**. It provides a very basic user experience for start/stop document processing workflows.

### Technology Stack

- **Framework**: Next.js 16.1.6
- **UI Library**: React 19.2.4
- **Styling**: Tailwind CSS 4.2.1
- **Components**: Custom components built on TailAdmin template
- **State Management**: React Context API
- **Real-time Updates**: React Use WebSocket
- **Language**: TypeScript

### Key Features

- 🎨 **Responsive Design**: Fully responsive from mobile to desktop
- 📊 **Data Visualization**: Charts, calendars, and data tables
- 🔄 **Real-time Updates**: WebSocket integration for live process status
- 🎯 **Process Management**: Upload, track, and manage document processing tasks
- 📈 **Process Results**: View summarization results and processing history
- 🌓 **Dark Mode**: Built-in theme switching
- 📱 **Mobile Optimized**: Mobile-first responsive design

### Dashboard Pages

- **Process List**: View all document processing jobs with status indicators

### Dashboard Components

Located in `apps/dashboard/src/components/`:

- **Tables**: Display process lists status
- **Header & Sidebar**: Navigation and branding

---

## API Overview

[> API Endpoints Documentation](/apps/api/README.md#api-overview)

The **API** is a high-performance backend service built with **Go** and the **Gin** framework. It handles document processing, AI summarization, and real-time WebSocket communications.

### Technology Stack

- **Language**: Go (Golang) +v1.25
- **Framework**: Gin for HTTP routing
- **WebSocket**: Gorilla WebSocket for real-time streaming
- **Database**: MongoDB with driver integration
- **AI**: Google Genai client for Gemini summarization
- **Documentation**: Swagger/OpenAPI via SwagGo
- **IPC**: Model Context Protocol (MCP) for AI integration

### API Features

#### 1. **Document Processing**
- Start async processing jobs for documents
- Monitor processing status in real-time
- Retrieve processing results and logs
- Stop/cancel running processes

#### 2. **Status & Results**
- HTTP endpoints for polling status
- WebSocket streaming for real-time updates
- Process history via MongoDB persistence
- Soft delete and archival support

#### 3. **AI Summarization**
- Integration with Google Gemini for content summarization
- MCP (Model Context Protocol) support
- Configurable summarization parameters
- Error handling and fallback strategies

#### 4. **Authentication**
- Bearer token authentication on all protected routes
- Environment-based token configuration
- Swagger UI available at `GET /swagger/`

#### 5. **Documentation**
- Auto-generated Swagger docs
- Clear API contracts
- Example requests/responses

### Core Endpoints

**Base URL**: `http://localhost:8080`

```
POST   /process/start           - Start a new document processing job
POST   /process/stop/:id        - Stop a processing job
GET    /process/status/:id      - Get current status of a job
GET    /process/results/:id     - Get processing results
GET    /process/list            - List all processes
WS     /ws/process/:id          - WebSocket for real-time status updates
GET    /health                  - Health check (no auth required)
GET    /swagger/*               - Swagger API documentation (no auth required, for production use, secure this endpoint)
```

---

## Prerequisites

Before setting up the project locally, ensure you have:

- **Node.js**: 20.x or later (recommended: 22.x+)
- **Go**: 1.21 or later
- **Docker**: Latest stable version
- **Docker Compose**: Latest stable version
- **MongoDB**: (via Docker) or local installation

### Required API Keys

- **Google Gemini API Key**: Get from [Google AI Studio](https://aistudio.google.com/app/apikey)
- **API Auth Token**: Generate a secure token for Bearer authentication

---

## Local Development Setup

### Before run the API locally ensure:

- Node.js >= v20.x
    ```bash
    node -v
    # v20.10.0
    ```
- go version >= 1.25
    ```bash
    go version
    # go version go1.26.1 linux/amd64
    ```
- Install gin server with hot-reload
    ```bash
    go install github.com/codegangsta/gin@latest
    ```
- Install Nx CLI globally
    ```bash
    npm install -g nx
    ```
- Docker and Docker Compose if you want to run the services in containers
    ```bash
    docker --version
    # Docker version 24.0.5, build 0aa7e65
    docker compose --version
    # docker compose version 2.17.2, build 8a1c60b
    ```

### Step 1: Clone the Repository

```bash
git clone <repository-url>
cd document-processing-system
```

### Step 2: Install Dependencies

```bash
# Install Node.js dependencies (for Nx monorepo management and dashboard)
npm install
```

### Step 3: Set Up Environment Variables

Create environment files for each service:

#### Dashboard Environment (`.env` in `apps/dashboard/`)

```env
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_API_URL_FRONT=http://localhost:8080
NEXT_PUBLIC_API_AUTH_TOKEN=your-secure-token-here       # Use the same token as API_AUTH_TOKEN in the API .env file, ensure this token is the same as the one used in the API for authentication

# Optional: Override in production
```

#### API Environment (`.env` in `apps/api/`)

```env
# Required: Authentication
API_AUTH_TOKEN=your-secure-token-here               # Use a secure token for Bearer authentication

# Required: Google Gemini
GEMINI_API_KEY=your-gemini-api-key-here

# Required: MongoDB
MONGODB_URI=mongodb://admin:your-secure-password-here@localhost:27017/?authMechanism=SCRAM-SHA-1&authSource=admin # Use the same credentials as the MongoDB container in Docker Compose, override with internal hostname in Docker Compose
MONGODB_DB=document_processing

# Optional: Server configuration
GIN_MODE=debug  # Use 'release' for production
```

#### Database Environment (`.env` in `db/`)

```env
# MongoDB configuration (if using local MongoDB)
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=your-secure-password-here
MONGO_INITDB_DATABASE=admin
```

### Step 4: Start Services with Docker Compose

```bash
# Start all services (Dashboard, API, MongoDB)
docker compose up -d --build

# View logs
docker compose logs -f

# View logs for a specific service
docker logs -f dps-app
docker logs -f dps-api

# Stop all services
docker compose down

# Stop one service
docker stop dps-app
docker stop dps-api
```

**Services will be available at:**
- 📊 Dashboard: `http://localhost:3030`
- 🔌 API: `http://localhost:8080`
- 🗄️ MongoDB: `mongodb://localhost:27017` (or check `docker compose.yml`)
- 📖 Swagger Docs: `http://localhost:8080/swagger/index.html`

### Step 5: Verify Setup

Check that all services are running:

```bash
# Check containers
docker compose ps

# Expected output:
# dps-app    running on 0.0.0.0:3030->3000/tcp
# dps-api    running on 0.0.0.0:8080->8080/tcp
# dps-db     running on 0.0.0.0:27018->27017/tcp
```

---

## Running Services

### Option 1: Docker Compose (Recommended)

```bash
# Start all services
docker compose up -d

# Start with live logs
docker compose up

# Stop all services
docker compose down

# Remove volumes (clear database)
docker compose down -v
```

### Option 2: Manual Development

#### Start MongoDB

```bash
docker run -d \
  --name dps-db \
  -p 27018:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=password123 \
  mongo:8-noble
```

#### Start API (Go)

```bash
cd apps/api
go run main.go
```

#### Start Dashboard (Next.js)

```bash
npm run dev -- apps/dashboard
# Or using Nx
npx nx dev dashboard
```

---

## Environment Configuration

### API Configuration (`apps/api/config/config.go`)

The API reads configuration from environment variables at startup:

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `API_AUTH_TOKEN` | ✅ Yes | Bearer token for authentication | `Bearer abc123xyz` |
| `GEMINI_API_KEY` | ✅ Yes | Google Gemini AI API key | `AIzaSyD...` |
| `MONGODB_URI` | ✅ Yes | MongoDB connection string | `mongodb://root:pass@localhost:27018` |
| `MONGODB_DB` | ✅ Yes | MongoDB database name | `document_processing` |
| `PORT` | ❌ No | Server port (default: 8080) | `8080` |
| `GIN_MODE` | ❌ No | Gin mode (debug/release) | `debug` |

### Dashboard Configuration

The dashboard uses environment variables with the `NEXT_PUBLIC_` prefix:

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `NEXT_PUBLIC_API_BASE_URL` | ✅ Yes | API base URL | `http://localhost:8080` |
| `NEXT_PUBLIC_WS_URL` | ✅ Yes | WebSocket URL | `ws://localhost:8080` |

---

## Available Commands

### Dashboard Commands

```bash
# Development server (with hot reload)
npx nx dev dashboard

# Production build
npx nx build dashboard

# Unit tests, not implemented yet because of time constraints but can be added in the future
npx nx test dashboard
```

### API Commands

```bash
# Development server (with hot reload)
npx nx dev api

# Production build
npx nx build api

# Install dependencies (if needed)
npx nx tidy api

# Recreate swagger docs
npx nx swagger api
# After running this command, the Swagger UI will be available at http://localhost:8080/swagger/index.html
# be careful with this command because you need to comment next lines in docs/docs.go
var SwaggerInfo = &swag.Spec{
	...
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	// LeftDelim:        "{{",
	// RightDelim:       "}}",
}

# Prettier formatting
npx nx prettier api

# Unit tests, not implemented yet because of time constraints but can be added in the future
npx nx test api
```

### Workspace Commands

```bash
# View project graph
npx nx graph

# Show project details
npx nx show project dashboard
npx nx show project api

# Run all tests
npm run test
```

---

## Screenshots

![Document Processing System](./docs/screencapture-dps-running.png)
<center><i>Document Processing System - Real-time Process Monitor</i></center><br />

This is a simple monitor dashboard that displays the first 10 statuses of document processing in real-time using WebSocket. You can see the list of processes, their IDs, and their current status. You can also start a new process using the **Start New Process** button and Stop/Cancel a running process using the **Stop** button.

<br />

![Document Processing System - Swagger Docs](./docs/screencapture-dps-swagger.png)
<center><i>Document Processing System - Swagger UI</i></center><br />

This is the Swagger UI for the API documentation. It provides a user-friendly interface to explore and test the API endpoints. You can see all available endpoints, their request/response formats, and even execute test requests directly from the UI.

---

## Project Statistics

- 📦 **Monorepo**: Single repository with 2 applications, could be add more as needed (e.g., worker service, CLI tool)
- 🎯 **Domains**: Process handling, Summarization, WebSocket communication
- 📁 **Structure**: Feature-based directory organization
- ⚙️ **Build System**: Nx workspace (NxGo + Next.js)
- 🧪 **Testing**: Jest for frontend + Go tests for backend, this will be done in next iterations

---

## Support & Documentation

Learn more:

- [Next.js Documentation](https://nextjs.org/docs)
- [Go Documentation](https://golang.org/doc)
- [Nx Workspace Docs](https://nx.dev)
- [MongoDB Docs](https://docs.mongodb.com)
- [Google Gemini API](https://ai.google.dev)
- [Gin Framework](https://gin-gonic.com)

---

## License

MIT License - See LICENSE file for details
