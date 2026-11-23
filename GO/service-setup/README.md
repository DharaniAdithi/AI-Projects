# Project Setup Script - Instructions

## Prerequisites

- [Go](https://go.dev/dl/) installed and available in PATH
- **macOS/Linux**: Bash shell
- **Windows**: PowerShell 5.1+ (pre-installed) or [PowerShell 7+](https://github.com/PowerShell/PowerShell)

## Files

| File | Platform |
|------|----------|
| `project_setup.sh` | macOS / Linux |
| `project_setup.ps1` | Windows |

## How to Run

### macOS / Linux

```bash
bash project_setup.sh
```

### Windows

```powershell
powershell -ExecutionPolicy Bypass -File project_setup.ps1
```

Or with PowerShell 7+:

```powershell
pwsh -ExecutionPolicy Bypass -File project_setup.ps1
```

## What It Does

The script prompts for three inputs:

1. **Root project name** – used for `go mod init` (e.g. `library-management-service`)
2. **Service name** – the folder name under `services/` (e.g. `library-service`)
3. **Entity name** – used to name handler, model, repository, service, and test files (e.g. `book`)

It then creates the following project structure:

```
services/<service_name>/
├── cmd/server/main.go
├── deployments/
│   ├── Dockerfile
│   └── docker-compose.yml
├── docs/
│   ├── openapi.yaml
│   └── swaggerui/
│       ├── index.html
│       └── scripts/
├── envs/.env.dev
├── internal/
│   ├── config/config.go
│   ├── dtos/
│   │   ├── response_dto.go
│   │   └── <entity>_dto.go
│   ├── errorcodes/error_codes.go
│   ├── handlers/
│   │   ├── routes.go
│   │   └── <entity>_handler.go
│   ├── loggers/logger.go
│   ├── models/<entity>.go
│   ├── repository/<entity>_repository.go
│   ├── services/<entity>_service.go
│   └── utils/utils.go
├── pkg/
│   ├── database/
│   │   ├── database.go
│   │   └── migration.go
│   └── server/server.go
├── tests/
│   ├── integration/
│   └── unit/
│       ├── config/
│       │   ├── config_test.go
│       │   └── mocks/config_mock.go
│       ├── database/database_test.go
│       ├── handlers/
│       │   ├── <entity>_handler_test.go
│       │   └── routes_test.go
│       ├── pkg/server/server_test.go
│       ├── repository/
│       │   ├── <entity>_repository_test.go
│       │   └── mock/<entity>_repository_mock.go
│       ├── server/server_test.go
│       ├── services/
│       │   ├── <entity>_service_test.go
│       │   └── mock/<entity>_service_mock.go
│       └── utils/utils_test.go
├── go.mod
├── go.sum
└── README.md
```

## Go Dependencies Installed

- `github.com/gorilla/mux` – HTTP router
- `go.uber.org/zap` – Structured logging
- `github.com/joho/godotenv` – .env file loader
- `github.com/stretchr/testify/assert` – Test assertions
- `github.com/lib/pq` – PostgreSQL driver

## Example

```
$ bash project_setup.sh
1. Create service
2. Exit
Enter your choice: 1
Enter the root project name: library-management-service
Enter the name of the new service: library-service
Enter the entity name (e.g. user, book, order): book
...
Service 'library-service' created at /Users/you/Public/services/library-service
```

This creates files like `book_handler.go`, `book_repository.go`, `book_service.go`, etc.
