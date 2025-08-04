# Clean Wire Implementation - No Custom Providers Needed

This document shows how to implement Google Wire dependency injection in a clean way by using service constructors directly, avoiding the need for custom provider functions.

## Key Principles

1. **Use service constructors directly** - `userservice.New`, `passwordservice.New`, etc.
2. **Simple config extraction functions** - Just extract config sections from `AppConfig`
3. **Use `wire.Bind`** to bind concrete types to interfaces
4. **Embed service configs** directly in `AppConfig`

## Clean AppConfig Structure

Instead of creating separate config types, embed the service configs directly:

```go
// AppConfig represents the root application configuration.
type AppConfig struct {
    ShutdownTimeout time.Duration          `yaml:"shutdown_timeout"`
    Logger          log.LoggerConfig       `yaml:"logger"`
    RestServer      rest.Config            `yaml:"rest_server"`
    Postgres        postgres.Config        `yaml:"postgres"`
    Password        passwordservice.Config `yaml:"password"`    // Embedded directly
    Username        usernameservice.Config `yaml:"username"`    // Embedded directly
    Audit           auditservice.Config    `yaml:"audit"`       // Embedded directly
}
```

## Simple Config Extraction Functions

These are not custom providers - they're just simple field extractors:

```go
// Simple config extraction functions - Wire can use these automatically
func ProvidePasswordConfig(cfg *AppConfig) passwordservice.Config { return cfg.Password }
func ProvideUsernameConfig(cfg *AppConfig) usernameservice.Config { return cfg.Username }
func ProvideAuditConfig(cfg *AppConfig) auditservice.Config       { return cfg.Audit }
func ProvideRestConfig(cfg *AppConfig) rest.Config               { return cfg.RestServer }
func ProvideLoggerConfig(cfg *AppConfig) log.LoggerConfig         { return cfg.Logger }
func ProvidePostgresConfig(cfg *AppConfig) postgres.Config        { return cfg.Postgres }
```

## Wire Sets Using Service Constructors Directly

```go
var ServiceSet = wire.NewSet(
    // Use service constructors directly - Wire will inject the right configs
    userservice.New,        // func New(db *sql.DB, logger *slog.Logger) *Service
    orgservice.New,         // func New(db *sql.DB, logger *slog.Logger) *Service
    passwordservice.New,    // func New(cfg Config, db *sql.DB, logger *slog.Logger) *Service
    usernameservice.New,    // func New(cfg Config, db *sql.DB, logger *slog.Logger) UsernameService
    auditservice.New,       // func New(cfg Config, db *sql.DB, logger *slog.Logger) AuditService
    
    // Bind concrete types to interfaces (only where needed)
    wire.Bind(new(userproto.UserService), new(*userservice.Service)),
    wire.Bind(new(orgproto.OrganizationService), new(*orgservice.Service)),
    wire.Bind(new(passwordproto.PasswordService), new(*passwordservice.Service)),
    // usernameservice.New and auditservice.New already return interfaces
    
    ProvideRestServices,  // Only aggregation function needed
)
```

## How Wire Resolves Dependencies Automatically

1. **userservice.New** needs `*sql.DB` and `*slog.Logger` → Wire provides them
2. **passwordservice.New** needs `passwordservice.Config`, `*sql.DB`, `*slog.Logger` → Wire calls `ProvidePasswordConfig(cfg)` to get the config
3. **Wire automatically matches** function parameters to available providers

## Benefits of This Approach

1. **No Custom Provider Functions** - Use service constructors directly
2. **Clean Configuration** - Service configs embedded in `AppConfig`
3. **Type Safety** - Wire validates at compile time
4. **Easy to Understand** - Clear dependency flow
5. **Easy to Test** - Each service can be tested independently
6. **Easy to Add Services** - Just add constructor + config + binding if needed

## When You Need Custom Functions

You only need custom functions for:

1. **Config extraction** - Simple field access: `func(cfg *AppConfig) ServiceConfig { return cfg.ServiceConfig }`
2. **Aggregation** - Combining multiple services: `ProvideRestServices`
3. **Type conversion** - When service constructor signatures don't match exactly

## Example: Adding a New Service

To add a new service, you only need:

1. **Add config to AppConfig**:

   ```go
   type AppConfig struct {
       // ... existing fields
       NewService newservice.Config `yaml:"new_service"`
   }
   ```

2. **Add config extractor**:

   ```go
   func ProvideNewServiceConfig(cfg *AppConfig) newservice.Config { return cfg.NewService }
   ```

3. **Add to ServiceSet**:

   ```go
   var ServiceSet = wire.NewSet(
       // ... existing services
       newservice.New,  // Use constructor directly
       wire.Bind(new(newproto.NewService), new(*newservice.Service)), // If needed
       // ...
   )
   ```

4. **Add to REST Services** (if needed):

   ```go
   func ProvideRestServices(
       // ... existing services
       newService newproto.NewService,
   ) rest.Services {
       return rest.Services{
           // ... existing services
           NewService: newService,
       }
   }
   ```

That's it! No custom provider functions needed.

## Wire Dependency Resolution Flow

```
config.yaml → koanf.Koanf → AppConfig → Service Configs → Service Constructors → Services → REST Services → WebService → WebContainer
```

Wire automatically figures out this flow and calls the right functions in the right order.

## Summary

- **Use service constructors directly** instead of creating provider functions
- **Embed service configs** in `AppConfig` instead of creating separate config types
- **Use simple extraction functions** for config fields
- **Use `wire.Bind`** when concrete types need to implement interfaces
- **Only create custom functions** for aggregation (like `ProvideRestServices`)

This approach is much cleaner and follows the principle of using existing constructors rather than wrapping them unnecessarily.
