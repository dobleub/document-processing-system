# Document Processing API

[< Back to README.md](/README.md#api-overview)

Go API service for:

- Running asynchronous document processing jobs
- Getting process status and results
- Streaming process updates over WebSocket
- Summarizing content with Gemini via HTTP and MCP

## Overview

This service is implemented with:

- `gin` for HTTP routing and middleware
- `gorilla/websocket` for real-time status streaming
- `swaggo` for OpenAPI/Swagger docs
- `google.golang.org/genai` for summarization
- in-memory process state (`sync.Map`) shared across handlers

Main entrypoint: `apps/api/main.go`

## Base URL

- Local: `http://localhost:8080`

## Authentication

All routes require a Bearer token except:

- `GET /health`
- `GET /swagger/*any`

Header format:

```http
Authorization: Bearer <API_AUTH_TOKEN>
```

The token is validated against `API_AUTH_TOKEN` from environment variables.

## Environment Variables

Configured in `apps/api/config/config.go`.

Required for normal operation:

- `API_AUTH_TOKEN` use it as `Bearer` token, it's for basic Auth.
- `GEMINI_API_KEY` use for summarize file content in genai client.
- `MONGODB_URI` allows mongodb connections to persists data.
- `MONGODB_DB` mongodb database name to store process info.

## Run Locally

Before run the API locally ensure:

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
- Docker and Docker Compose for local deployment and testing
    ```bash
    docker --version
    # Docker version 24.0.5, build 0aa7e65
    docker-compose --version
    # docker-compose version 2.17.2, build 8a1c60b
    ```

Running the API from `apps/api`:

```bash
go mod tidy # check go version >= 1.25

go run main.go
```

Running the API from repository root:

```bash
npx nx tidy api # check go version installed >= 125

npx nx dev api
```

The API starts on port `8080` when not running in Lambda.

## Swagger Docs

- UI: `GET /swagger/index.html`

## Architecture Summary

- `main.go` initializes config, MongoDB, MCP client, router, and middleware.
- Middleware injects objects into request context for handlers, enable basic auth, handle CORS and injects logger/config/db/state/MCP client into request context.
- Handlers read from context and operate on shared process state.
- Process execution starts in background goroutines and updates shared state.
- WebSocket endpoint polls state every second and emits updates on change.

## Data Models

Core process types are in:

- `apps/api/src/processDomain/interfaces/operations.go`
- `apps/api/src/processDomain/interfaces/response.go`
- `apps/api/src/processDomain/interfaces/results.go`

Processing Status Flags:

- `PENDING`
- `RUNNING`
- `PAUSED` # not implemented yet
- `COMPLETED`
- `FAILED`
- `STOPPED`

## Endpoint Specifications

### Health

#### `GET /health`

Returns service liveness.

Response `200`:

```text
OK
```

### Process Endpoints

#### `POST /process/start`

Starts a new asynchronous processing job over files in `targetFiles` directory.

Request body: currently ignored by implementation.

Response `200`:

```json
{
  "message": "Process Started",
  "id": "0f60f7b7-8d43-4d87-bf9d-84703ed2f67a"
}
```

Errors:

- `405` method not allowed
- `500` internal server error

#### `POST /process/stop/{id}`

Stops a running process by ID.

Path params:

- `id` (string, required)

Response `200` (`OperationStatus`):

```json
{
  "process_id": "0f60f7b7-8d43-4d87-bf9d-84703ed2f67a",
  "status": "STOPPED",
  "progress": {
    "total_files": 6,
    "processed_files": 3,
    "percentage": 50
  },
  "results": {
    "total_words": 12345,
    "total_lines": 678,
    "most_frequent_words": ["the", "and"],
    "files_processed": ["a.txt", "b.txt"],
    "files_to_process": ["c.txt"]
  },
  "started_at": "2026-03-10T00:00:00Z",
  "estimated_completion": "12s",
  "completed_at": "2026-03-10T00:00:12Z"
}
```

Errors:

- `400` missing process id
- `404` process not found

#### `GET /process/status/{id}`

Returns current status for one process.

Path params:

- `id` (string, required)

Response `200`: `OperationStatus` (same shape as above)

Errors:

- `400` missing process id
- `404` process not found

#### `GET /process/list`

Lists all tracked processes in review format.

Response `200`:

```json
{
  "processes": [
    {
      "id": "0f60f7b7-8d43-4d87-bf9d-84703ed2f67a",
      "status": "RUNNING",
      "error": "",
      "started_at": "2026-03-10T00:00:00Z",
      "estimated_completion": "8s",
      "files_processed": ["Don_Quijote.txt"],
      "files_to_process": ["Mobi_Dick.txt", "Peter_Pan.txt"],
      "completed_at": ""
    }
  ]
}
```

Notes:

- `files_processed` and `files_to_process` are base file names.

#### `GET /process/results/{id}`

Returns per-file analysis details for a process.

Path params:

- `id` (string, required)

Response `200`:

```json
{
  "progress": {
    "total_files": 6,
    "processed_files": 6,
    "percentage": 100
  },
  "analysis": {
    "process_id": "0f60f7b7-8d43-4d87-bf9d-84703ed2f67a",
    "status": "COMPLETED",
    "analysis": [
      {
        "file_name": "Don_Quijote.txt",
        "total_words": 123,
        "total_lines": 20,
        "most_frequent_words": ["de", "la"],
        "total_characters": 999,
        "summary": "..."
      }
    ]
  }
}
```

Errors:

- `400` missing process id
- `401` unauthorized
- `404` process not found

### Summarizer Endpoints

#### `POST /summarizer/summarize`

Summarizes either:

- multipart file form field `file`, or
- raw request body

Size limit: 32 MB for either input mode.

Request options:

- `multipart/form-data` with `file`
- `text/plain` or JSON/raw body content

Response `200`:

```json
{
  "message": "Process Completed",
  "duration": "1.234s",
  "summary": "Summarized content here..."
}
```

Errors:

- `400` file/body exceeds 32 MB
- `500` summarize failure or file read failure

#### `POST /summarizer/mcp`

Streamable MCP endpoint backed by the registered summarizer tool.

Used by MCP clients, not typical REST clients.

### WebSocket Endpoint

#### `GET /ws/status`

WebSocket stream of all process statuses.

Connection requirements:

- Include `Authorization: Bearer <API_AUTH_TOKEN>` in upgrade request.

Behavior:

- Sends an initial snapshot immediately after connection
- Checks process state every second
- Sends update only when snapshot changes per process:
  - process status changed
  - `files_processed` length changed
  - completion timestamp changed
  - process created or removed
- Sends ping frames every 30 seconds
- Uses synchronized writes to avoid concurrent write panics

Message format (`text` frame, JSON array of `OperationReview`):

```json
[
  {
    "id": "0f60f7b7-8d43-4d87-bf9d-84703ed2f67a",
    "status": "RUNNING",
    "error": "",
    "started_at": "2026-03-10T00:00:00Z",
    "estimated_completion": "8s",
    "files_processed": ["Don_Quijote.txt"],
    "files_to_process": ["Mobi_Dick.txt", "Peter_Pan.txt"],
    "completed_at": ""
  }
]
```

Client example:

```bash
wscat -c ws://localhost:8080/ws/status -H "Authorization: Bearer <API_AUTH_TOKEN>"
```

## cURL Examples

```bash
TOKEN="your-token"

# Start process
curl -X POST http://localhost:8080/process/start \
  -H "Authorization: Bearer $TOKEN"

# List processes
curl http://localhost:8080/process/list \
  -H "Authorization: Bearer $TOKEN"

# Get status by id
curl http://localhost:8080/process/status/<id> \
  -H "Authorization: Bearer $TOKEN"

# Stop process
curl -X POST http://localhost:8080/process/stop/<id> \
  -H "Authorization: Bearer $TOKEN"

# Get analysis results
curl http://localhost:8080/process/results/<id> \
  -H "Authorization: Bearer $TOKEN"

# Summarize file
curl -X POST http://localhost:8080/summarizer/summarize \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/absolute/path/to/file.txt"

# Summarize raw body
curl -X POST http://localhost:8080/summarizer/summarize \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: text/plain" \
  --data "Long text to summarize"
```

## Important Implementation Notes

- Process state is in-memory (`sync.Map`). Restarting the service clears all process state.
- Start handler processes files from `targetFiles` relative to process working directory.
- CORS is currently open (`Access-Control-Allow-Origin: *`).
- WebSocket origin check currently allows all origins and should be restricted in production.
