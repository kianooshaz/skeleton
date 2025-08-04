# Documentation Index

This directory contains documentation for the skeleton project.

## Architecture Documentation

- [ADD_SERVICE_GUIDE.md](ADD_SERVICE_GUIDE.md) - Comprehensive guide for adding new services to the container package and Wire dependency injection system
- [ADD_SERVICE_QUICK_REFERENCE.md](ADD_SERVICE_QUICK_REFERENCE.md) - Quick checklist and templates for adding services
- [CLEAN_WIRE.md](CLEAN_WIRE.md) - Wire cleanup procedures
- [WIRE_CONTAINER.md](WIRE_CONTAINER.md) - Wire container documentation

## Service Development

When adding new services to this project, follow these steps:

1. **Quick Start**: Use the [Quick Reference Guide](ADD_SERVICE_QUICK_REFERENCE.md) for a step-by-step checklist
2. **Detailed Guide**: Refer to the [Complete Service Guide](ADD_SERVICE_GUIDE.md) for comprehensive instructions and examples
3. **Wire Cleanup**: Use [Clean Wire Guide](CLEAN_WIRE.md) if you need to reset or troubleshoot Wire dependencies

## Key Concepts

### Container Package

The `internal/container` package manages dependency injection using Google Wire. It provides:

- Configuration management through `AppConfig`
- Service lifecycle management (start/stop)
- Dependency graph definition through Wire sets
- Clean separation of concerns

### Service Structure

Services follow a consistent pattern:

```text
services/{domain}/{service}/
├── proto/           # Interfaces and models
├── service/         # Implementation and config
└── persistence/     # Database layer (if needed)
```

### Wire Integration

Each service requires:

1. Configuration struct with YAML tags
2. Constructor function compatible with Wire
3. Integration into Wire sets
4. Registration in the container

For detailed information, see the [Service Addition Guide](ADD_SERVICE_GUIDE.md).
