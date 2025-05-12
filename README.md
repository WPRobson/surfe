# Surfe API

## Prerequisites

- Go 1.24.2 or higher
- Git
- Docker (optional)
- Docker Compose (optional)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/surfe.git
cd surfe
```

2. Install dependencies:
```bash
go mod download
```

3. Install Swagger CLI tool (for API documentation):
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Running the Application

### Option 1: Local Development

1. Start the server:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8000`

### Option 2: Using Docker

1. Build the Docker image:
```bash
docker build -t surfe-api .
```

2. Run the container:
```bash
docker run -p 8000:8000 surfe-api
```

The server will be available at `http://localhost:8000`

### Option 3: Using Docker Compose

1. Start the application:
```bash
docker-compose up
```

To run in detached mode:
```bash
docker-compose up -d
```

2. Stop the application:
```bash
docker-compose down
```

The server will be available at `http://localhost:8000`

## Access the Swagger documentation:
```
http://localhost:8000/swagger/index.html
```

## API Documentation

### Users

#### Get User by ID
```http
GET /api/v1/users/{id}
```
Returns user details by their ID.

#### Get User Action Count
```http
GET /api/v1/users/{id}/actions/count
```
Returns the total number of actions performed by a user.

### Actions

#### Get Next Action Probabilities
```http
GET /api/v1/actions/{type}/next
```
Returns probabilities of next actions based on current action type.

#### Get Referral Index
```http
GET /api/v1/actions/referral
```
Returns the referral index showing how many users each user has referred.

## Project Structure

```
surfe/
├── cmd/
│   └── api/
│       └── main.go         # Application entry point
├── internal/
│   ├── handlers/          # HTTP request handlers
│   ├── models/            # Data models
│   ├── repository/        # Data access layer
│   └── services/          # Business logic
├── docs/                  # Swagger documentation
└── README.md
```

## Testing

Run the test suite:
```bash
go test ./...
```

## API Response Examples

### Get User Response
```json
{
    "id": 1,
    "name": "John Doe",
    "createdAt": "2024-03-11T20:00:00Z"
}
```

### Get Action Count Response
```json
{
    "count": 42
}
```

### Get Next Action Probabilities Response
```json
{
	"ADD_TO_CRM": 0.70,
	"REFER_USER": 0.20,
	"VIEW_CONVERSATION": 0.10
}
```

### Get Referral Index Response
```json
{
	"1": 3, 
	"2": 0,
	"3": 7
}
```