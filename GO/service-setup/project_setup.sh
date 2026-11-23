#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"
SERVICES_DIR="$PROJECT_ROOT/services"

ensure_dir() { mkdir -p "$1"; }
ensure_file() { mkdir -p "$(dirname "$1")" && [ -f "$1" ] || touch "$1"; }

create_service() {
    read -rp "Enter the root project name: " project_name
    read -rp "Enter the name of the new service: " service_name
    read -rp "Enter the entity name (e.g. user, book, order): " entity
    local svc="$SERVICES_DIR/$service_name"

    # cmd
    ensure_file "$svc/cmd/server/main.go"

    # deployments
    ensure_file "$svc/deployments/Dockerfile"
    ensure_file "$svc/deployments/docker-compose.yml"

    # docs
    ensure_file "$svc/docs/openapi.yaml"
    ensure_file "$svc/docs/swaggerui/index.html"
    ensure_dir "$svc/docs/swaggerui/scripts"

    # envs
    ensure_file "$svc/envs/.env.dev"

    # internal
    ensure_file "$svc/internal/config/config.go"
    ensure_file "$svc/internal/dtos/response_dto.go"
    ensure_file "$svc/internal/dtos/${entity}_dto.go"
    ensure_file "$svc/internal/errorcodes/error_codes.go"
    ensure_file "$svc/internal/handlers/routes.go"
    ensure_file "$svc/internal/handlers/${entity}_handler.go"
    ensure_file "$svc/internal/loggers/logger.go"
    ensure_file "$svc/internal/models/${entity}.go"
    ensure_file "$svc/internal/repository/${entity}_repository.go"
    ensure_file "$svc/internal/services/${entity}_service.go"
    ensure_file "$svc/internal/utils/utils.go"

    # pkg
    ensure_file "$svc/pkg/database/database.go"
    ensure_file "$svc/pkg/database/migration.go"
    ensure_file "$svc/pkg/server/server.go"

    # tests
    ensure_file "$svc/tests/integration/.gitkeep"
    ensure_file "$svc/tests/unit/config/config_test.go"
    ensure_file "$svc/tests/unit/config/mocks/config_mock.go"
    ensure_file "$svc/tests/unit/database/database_test.go"
    ensure_file "$svc/tests/unit/handlers/${entity}_handler_test.go"
    ensure_file "$svc/tests/unit/handlers/routes_test.go"
    ensure_file "$svc/tests/unit/pkg/server/server_test.go"
    ensure_file "$svc/tests/unit/repository/${entity}_repository_test.go"
    ensure_file "$svc/tests/unit/repository/mock/${entity}_repository_mock.go"
    ensure_file "$svc/tests/unit/services/${entity}_service_test.go"
    ensure_file "$svc/tests/unit/services/mock/${entity}_service_mock.go"
    ensure_file "$svc/tests/unit/server/server_test.go"
    ensure_file "$svc/tests/unit/utils/utils_test.go"

    # go module
    cd "$svc" || exit
    go mod init "$project_name/$service_name"
    go get github.com/gorilla/mux
    go get go.uber.org/zap
    go get github.com/joho/godotenv
    go get github.com/stretchr/testify/assert
    go get github.com/lib/pq
    go mod tidy
    cd "$PROJECT_ROOT" || exit

    ensure_file "$svc/README.md"
    echo "Service '$service_name' created at $svc"
}

show_menu() {
    echo "1. Create service"
    echo "2. Exit"
    read -rp "Enter your choice: " choice
    case $choice in
        1) create_service ;;
        2) exit 0 ;;
        *) echo "Invalid choice." ;;
    esac
    show_menu
}

show_menu
