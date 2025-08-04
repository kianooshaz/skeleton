# Adding Services to Container Package and Wire

This guide explains how to add new services to the container package and integrate them with the Wire dependency injection system.

## Overview

The container package uses Google Wire for dependency injection to manage service dependencies. When adding a new service, you need to:

1. Create the service implementation
2. Update container configuration
3. Modify the Wire setup
4. Update the container implementation
5. Regenerate Wire dependencies

## Step-by-Step Guide

### 1. Create Your Service

First, ensure your service follows the established patterns:

```go
// services/example/service/service.go
package service

import (
    "context"
    "github.com/kianooshaz/skeleton/services/example/proto"
)

type Config struct {
    // Service-specific configuration
    Timeout time.Duration `yaml:"timeout"`
    MaxRetries int `yaml:"max_retries"`
}

type Service struct {
    config Config
    // other dependencies
}

func New(config Config /* other dependencies */) proto.ExampleService {
    return &Service{
        config: config,
        // initialize other fields
    }
}

func (s *Service) Shutdown(ctx context.Context) {
    // cleanup logic
}
```

### 2. Update Container Configuration

Add your service configuration to `internal/container/config.go`:

```go
// AppConfig represents the root application configuration.
type AppConfig struct {
    ShutdownTimeout time.Duration          `yaml:"shutdown_timeout"`
    Logger          log.LoggerConfig       `yaml:"logger"`
    RestServer      rest.Config            `yaml:"rest_server"`
    Postgres        postgres.Config        `yaml:"postgres"`
    Password        passwordservice.Config `yaml:"password"`
    Username        usernameservice.Config `yaml:"username"`
    Audit           auditservice.Config    `yaml:"audit"`
    // Add your new service config
    Example         exampleservice.Config  `yaml:"example"`
}
```

### 3. Add Configuration Provider Function

In `internal/container/wire.go`, add a provider function for your service config:

```go
import (
    // ... existing imports
    exampleproto "github.com/kianooshaz/skeleton/services/example/proto"
    exampleservice "github.com/kianooshaz/skeleton/services/example/service"
)

// Add your config provider function
func ProvideExampleConfig(cfg *AppConfig) exampleservice.Config { return cfg.Example }
```

### 4. Update Wire Sets

Add your configuration provider to the `ConfigSet` and your service constructor to the `WebContainerSet`:

```go
var ConfigSet = wire.NewSet(
    config.LoadConfigWithDefaults,
    ProvideAppConfig,
    ProvidePasswordConfig,
    ProvideUsernameConfig,
    ProvideAuditConfig,
    ProvideRestConfig,
    ProvideLoggerConfig,
    ProvidePostgresConfig,
    // Add your config provider
    ProvideExampleConfig,
)

var WebContainerSet = wire.NewSet(
    ConfigSet,
    LoggerSet,
    DatabaseSet,
    userservice.New,
    orgservice.New,
    passwordservice.New,
    usernameservice.New,
    auditservice.New,
    // Add your service constructor
    exampleservice.New,
    rest.New,
    ProvideWebContainer,
)
```

### 5. Update Container Provider Function

Modify the `ProvideWebContainer` function to accept your new service:

```go
func ProvideWebContainer(
    cfg *AppConfig,
    logger *slog.Logger,
    db *sql.DB,
    webService protocol.WebService,
    userService userproto.UserService,
    orgService orgproto.OrganizationService,
    passwordService passwordproto.PasswordService,
    usernameService usernameproto.UsernameService,
    auditService auditproto.AuditService,
    // Add your service parameter
    exampleService exampleproto.ExampleService,
) Container {
    return &WebContainer{
        config:              cfg,
        logger:              logger,
        db:                  db,
        webService:          webService,
        userService:         userService,
        organizationService: orgService,
        passwordService:     passwordService,
        usernameService:     usernameService,
        auditService:        auditService,
        // Add your service field
        exampleService:      exampleService,
    }
}
```

### 6. Update WebContainer Struct

Add your service field to the `WebContainer` struct in `internal/container/web_container.go`:

```go
import (
    // ... existing imports
    exampleproto "github.com/kianooshaz/skeleton/services/example/proto"
)

// WebContainer holds all dependencies for the web application.
type WebContainer struct {
    config              *AppConfig
    logger              *slog.Logger
    db                  *sql.DB
    webService          protocol.WebService
    userService         userproto.UserService
    organizationService orgproto.OrganizationService
    passwordService     passwordproto.PasswordService
    usernameService     usernameproto.UsernameService
    auditService        auditproto.AuditService
    // Add your service field
    exampleService      exampleproto.ExampleService
}
```

### 7. Update Container Shutdown Logic

If your service needs cleanup, add it to the `Stop()` method:

```go
func (c *WebContainer) Stop() error {
    c.logger.Info("Starting graceful shutdown of web container")

    // Create a timeout context for graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), c.config.ShutdownTimeout)
    defer cancel()

    // Stop services in reverse order of dependency
    if c.webService != nil {
        c.logger.Info("Shutting down web service")
        if err := c.webService.Shutdown(ctx); err != nil {
            c.webService.Close()
            c.logger.Error("Failed to shutdown web service", "error", err)
        }
    }

    // Add your service shutdown logic
    if c.exampleService != nil {
        c.logger.Info("Shutting down example service")
        c.exampleService.Shutdown(ctx)
    }

    if c.auditService != nil {
        c.logger.Info("Shutting down audit service")
        c.auditService.Shutdown(ctx)
    }

    if c.db != nil {
        c.logger.Info("Closing database connection")
        if err := c.db.Close(); err != nil {
            c.logger.Error("Failed to close database connection", "error", err)
        }
    }

    return nil
}
```

### 8. Update Configuration Files

Add your service configuration to the YAML files:

```yaml
# config.yaml
app:
  shutdown_timeout: 30s
  logger:
    level: info
  # ... other configs
  example:
    timeout: 5s
    max_retries: 3
```

### 9. Regenerate Wire Dependencies

Run Wire to regenerate the dependency injection code:

```bash
make clean-wire  # if you have this target, otherwise:
go generate ./internal/container/
```

Or manually run:

```bash
cd internal/container
wire
```

### 10. Add Service Accessor (Optional)

If other parts of your application need to access the service, add getter methods to the container:

```go
func (c *WebContainer) ExampleService() exampleproto.ExampleService {
    return c.exampleService
}
```

## Example: Adding a Notification Service

Here's a complete example of adding a notification service:

### 1. Service Structure

```text
services/
  notification/
    notification/
      proto/
        notification.go
        models.go
      service/
        service.go
        config.go
      persistence/
        query.go
        queries/
          create_notification.sql
```

### 2. Configuration Update

```go
// internal/container/config.go
type AppConfig struct {
    // ... existing fields
    Notification notificationservice.Config `yaml:"notification"`
}
```

### 3. Wire Updates

```go
// internal/container/wire.go
import (
    notificationproto "github.com/kianooshaz/skeleton/services/notification/notification/proto"
    notificationservice "github.com/kianooshaz/skeleton/services/notification/notification/service"
)

func ProvideNotificationConfig(cfg *AppConfig) notificationservice.Config {
    return cfg.Notification
}

var ConfigSet = wire.NewSet(
    // ... existing providers
    ProvideNotificationConfig,
)

var WebContainerSet = wire.NewSet(
    // ... existing services
    notificationservice.New,
    ProvideWebContainer,
)
```

## Best Practices

1. **Consistent Naming**: Follow the established naming patterns for services and configurations
2. **Error Handling**: Always handle errors in service initialization and shutdown
3. **Configuration**: Use struct tags for YAML configuration mapping
4. **Testing**: Write unit tests for your service and integration tests for the container
5. **Documentation**: Update this guide when adding new patterns or requirements
6. **Dependencies**: Be mindful of circular dependencies between services

## Troubleshooting

### Common Issues

1. **Wire Generation Fails**:
   - Check import paths are correct
   - Ensure all required parameters are provided to constructors
   - Verify no circular dependencies exist

2. **Configuration Not Loading**:
   - Check YAML tags match configuration file structure
   - Verify configuration provider function is in ConfigSet

3. **Service Not Starting**:
   - Check service is added to WebContainerSet
   - Verify all dependencies are satisfied
   - Check for runtime errors in logs

### Debugging Wire Issues

Use `wire check` to validate your wire setup:

```bash
cd internal/container
wire check
```

Use `wire show` to see the dependency graph:

```bash
cd internal/container
wire show
```

## Related Files

- `internal/container/wire.go` - Wire dependency injection setup
- `internal/container/config.go` - Application configuration
- `internal/container/web_container.go` - Container implementation
- `internal/container/proto.go` - Container interface
- `docs/CLEAN_WIRE.md` - Wire cleanup procedures

## References

- [Google Wire Documentation](https://github.com/google/wire)
- [Dependency Injection Patterns](https://martinfowler.com/articles/injection.html)
