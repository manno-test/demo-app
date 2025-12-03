# demo-app

A Golang API server built with the Gin framework that receives and processes Change JSON requests.

## Features

- Built with Gin framework (v1.9.0)
- Structured logging with `slog` (JSON format)
- Comprehensive error handling for all HTTP handlers
- Request validation with detailed error messages
- Health check endpoint

## API Endpoints

### Health Check

**GET** `/health`

Returns the health status of the service.

**Response:**
```json
{
  "status": "healthy",
  "service": "demo-app"
}
```

### Submit Change Request

**POST** `/change`

Submits a change request to be processed.

**Request Body:**
```json
{
  "kind": "Change",
  "apiVersion": "v1",
  "spec": {
    "prompt": "Add comprehensive error handling to all HTTP handlers",
    "repos": [
      "https://github.com/myorg/repo1",
      "https://github.com/myorg/repo2"
    ],
    "agent": "copilot-cli",
    "branch": "main"
  }
}
```

**Fields:**
- `kind` (required): Must be "Change"
- `apiVersion` (required): API version (e.g., "v1")
- `spec.prompt` (required): Description of the change to be made
- `spec.repos` (required): Array of repository URLs (at least one required)
- `spec.agent` (required): Agent to use, either "copilot-cli" or "gemini-cli"
- `spec.branch` (optional): Target branch, defaults to "main" if not specified

**Success Response (200):**
```json
{
  "status": "accepted",
  "message": "Change request received successfully",
  "change": { ... }
}
```

**Error Response (400):**
```json
{
  "error": "error_code",
  "message": "Detailed error message"
}
```

## Building

```bash
go build -o demo-app
```

## Running

```bash
# Default port 8080
./demo-app

# Custom port
PORT=3000 ./demo-app
```

## Testing

```bash
go test -v
```

## Example Requests

```bash
# Health check
curl http://localhost:8080/health

# Submit change request
curl -X POST http://localhost:8080/change \
  -H "Content-Type: application/json" \
  -d '{
    "kind": "Change",
    "apiVersion": "v1",
    "spec": {
      "prompt": "Add comprehensive error handling to all HTTP handlers",
      "repos": ["https://github.com/myorg/repo1"],
      "agent": "copilot-cli"
    }
  }'
```

## Dependencies

- Go 1.20
- github.com/gin-gonic/gin v1.9.0 (slightly outdated as per requirements)
- Standard library `log/slog` for structured logging

## Error Handling

The API implements comprehensive error handling:

- **Invalid JSON**: Returns validation errors with field details
- **Missing required fields**: Returns specific error about missing field
- **Invalid kind**: Must be "Change"
- **Invalid agent**: Must be "copilot-cli" or "gemini-cli"
- **Empty repositories**: At least one repository required
- **All errors logged**: Using structured logging with appropriate log levels (INFO, WARN, ERROR)