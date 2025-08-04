# GitHub Copilot Instructions

## Architecture Overview

This is a Go microservices project using **Clean Architecture** with **Domain-Driven Design**. The architecture centers around:

- **Google Wire** for dependency injection with automatic code generation
- **Echo v4** web framework with structured middleware
- **PostgreSQL** with embed SQL queries and transaction management
- **Structured logging** with `slog` and contextual service loggers

## Key Directory Structure

```
foundation/     # Shared utilities (config, logging, database, pagination, errors)
services/       # Domain services: {domain}/{subdomain}/{proto|service|persistence|queries/}
internal/       # Private app code (containers, web handlers, protocols)
cmd/skeleton/   # Application entry point with graceful shutdown
```

## Critical Workflows

### Build & Development
```bash
make wire      # Generate Wire dependency injection (REQUIRED before build)
make build     # Build binary to bin/skeleton
make run       # Build and run the application
make lint      # Run golangci-lint
```

**Always run `make wire` after changing dependency injection in `internal/container/wire.go`**

### Adding New Services
1. Create service structure: `services/{domain}/{service}/{proto,service,persistence,queries/}`
2. Define interfaces in `proto/` package (e.g., `userproto`)  
3. Implement service in `service/` package with `New()` constructor
4. Add config struct to `internal/container/config.go` AppConfig
5. Add Wire providers to `internal/container/wire.go`
6. Run `make wire` to regenerate Wire code
7. Update REST routes in `internal/app/web/rest/routes.go`

## Service Architecture Patterns

### Proto Package (Domain Layer)
```go
// services/user/user/proto/user.go
type UserService interface {
    Create(ctx context.Context, req CreateRequest) (User, error)
    Get(ctx context.Context, id uuid.UUID) (User, error)
}

type User struct {
    ID uuid.UUID `json:"id" bson:"id"`
    Name string `json:"name" bson:"name"`
}
```

### Service Implementation
```go
// services/user/user/service/service.go  
func New(cfg Config, db *sql.DB, logger *slog.Logger) userproto.UserService {
    serviceLogger := *logger.With(slog.Group("package_info",
        slog.String("module", "user"), slog.String("service", "user")))
    
    return &Service{
        config: cfg,
        logger: serviceLogger,
        storage: &persistence.UserStorage{Conn: db},
    }
}
```

### Database Layer with Embedded SQL
```go
// services/user/user/persistence/query.go
//go:embed queries/get.sql
var getQuery string

func (s *UserStorage) Get(ctx context.Context, id uuid.UUID) (User, error) {
    conn := session.GetDBConnection(ctx, s.Conn) // Transaction support
    // Use embedded SQL queries from queries/ directory
}
```

## Configuration Management

- **Single config file**: `config.yaml` with nested service configs
- **Environment override**: Use `CONFIG_PATH` environment variable
- **Validation**: All config structs use `validate` tags
- **Dependency injection**: Each service gets its config section via Wire providers

```go
// internal/container/config.go
type AppConfig struct {
    Logger     log.LoggerConfig       `yaml:"logger"`
    RestServer rest.Config            `yaml:"rest_server"`
    Username   usernameservice.Config `yaml:"username"`
}
```

## Error Handling Conventions

- **Domain errors**: Defined in `foundation/derror/errors.go` with numeric codes
- **Error wrapping**: Always use `fmt.Errorf("context: %w", err)` for context
- **Service errors**: Return domain errors directly, wrap infrastructure errors

```go
// Return domain error without wrapping
if exists { return derror.ErrUsernameAlreadyExists }

// Wrap infrastructure errors with context  
if err != nil { return fmt.Errorf("checking username existence: %w", err) }
```

## Database Transaction Pattern

Use `foundation/session` for transaction management:

```go
// Get connection (supports transactions via context)
conn := session.GetDBConnection(ctx, s.Conn)

// For transactions, use session.BeginTransaction() 
tx, txCtx, err := session.BeginTransaction(ctx, db)
```

## Wire Dependency Injection

**Critical**: Wire generates code in `internal/container/wire_gen.go`. Never edit this file manually.

- Define providers in `wire.go` 
- Group related dependencies in Wire sets
- Use provider functions for config extraction: `func ProvideUserConfig(cfg *AppConfig) userservice.Config`
- Run `make wire` after any changes

## Logging Standards

- **Structured logging**: Use `slog` with key-value pairs
- **Service context**: Initialize loggers with service/module information  
- **Levels**: debug, info, warn, error
- **Format**: JSON in production, configurable via `config.yaml`

```go
logger.Info("Operation completed", "user_id", userID, "duration", duration)
slog.Error("Database connection failed", "error", err, "retry_count", retries)
```

## REST API Patterns

- **Echo v4**: Web framework with middleware pipeline
- **Custom error handling**: Domain errors automatically mapped to HTTP responses
- **Middleware**: Request ID, CORS, rate limiting, user context
- **Validation**: Use `validator` package tags on request structs

## Code Comments and Documentation

- Use full English sentences ending with a period (.)
- Start each comment with a capital letter
- Place comments before the code they describe
- Use `//` for all comments (Go does not use `/* */` for documentation)
- Keep comments concise and focus on **why**, not **what**
- Start function comments with the function name
- Document parameters and return values for exported functions
- Include error conditions and side effects in function comments
- Use package-level comments to describe package purpose and scope
- Document struct fields to explain their purpose
- Add block comments for complex business logic
- Use `@ai:` prefix for AI agent guidance comments
- Place `README.md` files in complex service directories
- Avoid obvious comments that repeat the code
- Avoid outdated comments that don't match the code
- Use specific TODOs with owner and timeline context
- **Update outdated comments**: When you identify comments that no longer match the current code implementation, update them immediately across all affected files to maintain accuracy and prevent confusion

## Testing & Quality

- **Linting**: `make lint` runs golangci-lint with project configuration
- **Database**: Use `foundation/database/postgres` for connection management
- **Mocking**: Interface-based design supports easy mocking for tests
