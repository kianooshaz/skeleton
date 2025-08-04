# Quick Reference: Adding Services to Container

This is a quick checklist for adding new services to the container package.

## Checklist

### 1. Service Implementation

- [ ] Create service in `services/{domain}/{service}/`
- [ ] Implement proto interface
- [ ] Add `Config` struct with YAML tags
- [ ] Add `New()` constructor function
- [ ] Add `Shutdown(ctx context.Context)` method if needed

### 2. Container Configuration

- [ ] Add service config to `AppConfig` in `internal/container/config.go`
- [ ] Add config provider function in `internal/container/wire.go`
- [ ] Add config provider to `ConfigSet`

### 3. Wire Integration

- [ ] Import service packages in `internal/container/wire.go`
- [ ] Add service constructor to `WebContainerSet`
- [ ] Add service parameter to `ProvideWebContainer` function
- [ ] Add service field to `WebContainer` struct

### 4. Container Updates

- [ ] Add service field to `WebContainer` in `internal/container/web_container.go`
- [ ] Add import for service proto package
- [ ] Add shutdown logic to `Stop()` method if needed
- [ ] Add getter method if needed

### 5. Configuration Files

- [ ] Add service config section to `config.yaml`
- [ ] Add service config to test configs if needed

### 6. Wire Generation

- [ ] Run `go generate ./internal/container/` or `wire` command
- [ ] Verify `wire_gen.go` is updated correctly

## Template Code Snippets

### Config Provider Function

```go
func Provide{Service}Config(cfg *AppConfig) {service}service.Config { 
    return cfg.{Service} 
}
```

### WebContainer Field

```go
{service}Service {service}proto.{Service}Service
```

### Shutdown Logic

```go
if c.{service}Service != nil {
    c.logger.Info("Shutting down {service} service")
    c.{service}Service.Shutdown(ctx)
}
```

### Service Constructor Pattern

```go
func New(config Config, deps ...Dependencies) proto.ServiceInterface {
    return &Service{
        config: config,
        // other fields
    }
}
```

## Common Patterns

### Service Config Structure

```go
type Config struct {
    Timeout     time.Duration `yaml:"timeout"`
    MaxRetries  int          `yaml:"max_retries"`
    BatchSize   int          `yaml:"batch_size"`
}
```

### YAML Configuration

```yaml
app:
  {service}:
    timeout: 30s
    max_retries: 3
    batch_size: 100
```

## Validation Commands

```bash
# Check wire setup
cd internal/container && wire check

# Generate wire dependencies  
cd internal/container && wire

# View dependency graph
cd internal/container && wire show
```

## See Also

- [ADD_SERVICE_GUIDE.md](ADD_SERVICE_GUIDE.md) - Detailed step-by-step guide
- [CLEAN_WIRE.md](CLEAN_WIRE.md) - Wire cleanup procedures
