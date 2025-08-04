# Birthday Service

The Birthday Service manages user birthday information within the skeleton application. It follows the Clean Architecture pattern with domain-driven design principles.

## Features

- Create, read, update, and delete birthday records
- Automatic age calculation based on date of birth
- Age validation with configurable min/max bounds
- One birthday record per user (enforced by unique constraint)
- Paginated listing with advanced filtering options
- Structured logging and error handling

## Architecture

```
birthday/
├── proto/                  # Domain layer (interfaces and models)
│   ├── birthday.go        # Service interface and request/response models
│   └── birthdayid.go      # Birthday ID type with UUID implementation
├── service/               # Business logic layer
│   ├── service.go         # Service constructor and configuration
│   └── business_logic.go  # Core business logic implementation
├── persistence/           # Data access layer
│   ├── query.go          # Database queries and operations
│   ├── order.go          # SQL ordering logic
│   └── queries/          # Embedded SQL files
│       ├── create.sql
│       ├── get.sql
│       ├── get_by_user_id.sql
│       ├── update.sql
│       ├── delete.sql
│       ├── list.sql
│       ├── count.sql
│       └── exists_by_user_id.sql
├── schema.sql            # Database schema definition
└── README.md            # This file
```

## Configuration

The birthday service is configured in `config.yaml`:

```yaml
birthday:
  max_age: 150    # Maximum allowed age (default: 150)
  min_age: 0      # Minimum allowed age (default: 0)
```

## Database Schema

The service uses a single `birthdays` table with the following structure:

```sql
CREATE TABLE birthdays (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    date_of_birth DATE NOT NULL,
    age INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

Run the `schema.sql` file to create the table and indexes.

## API Operations

### Create Birthday

- **Purpose**: Create a new birthday record for a user
- **Validation**:
  - Ensures user doesn't already have a birthday record
  - Validates age is within configured bounds
  - Automatically calculates age from date of birth

### Get Birthday

- **Purpose**: Retrieve a birthday record by ID
- **Returns**: Birthday record with calculated age

### Get Birthday by User ID

- **Purpose**: Retrieve a birthday record for a specific user
- **Returns**: Birthday record for the user

### Update Birthday

- **Purpose**: Update an existing birthday record
- **Validation**:
  - Recalculates age based on new date of birth
  - Validates new age is within bounds

### Delete Birthday

- **Purpose**: Remove a birthday record
- **Effect**: Permanently deletes the record

### List Birthdays

- **Purpose**: Retrieve paginated list of birthdays with optional filters
- **Filters**:
  - `user_id`: Filter by specific user
  - `min_age`/`max_age`: Filter by age range
  - `birth_month`: Filter by birth month (1-12)
- **Sorting**: Supports sorting by any field (id, user_id, date_of_birth, age, created_at, updated_at)

## Business Rules

1. **One Birthday Per User**: Each user can have only one birthday record
2. **Age Validation**: Age must be within configured min/max bounds
3. **Automatic Age Calculation**: Age is automatically calculated and updated
4. **Date Validation**: Date of birth cannot be in the future
5. **Immutable User Association**: Once created, birthday records cannot be transferred between users

## Error Handling

The service uses domain-specific errors from `foundation/derror`:

- `ErrUserAlreadyExists`: When trying to create a duplicate birthday for a user
- Standard database errors are wrapped with context for better debugging

## Usage in Wire Container

The birthday service is automatically wired into the dependency injection container:

```go
// Configuration
ProvideBirthdayConfig(cfg *AppConfig) birthdayservice.Config

// Service creation
birthdayservice.New(cfg Config, db *sql.DB, logger *slog.Logger) birthdayproto.BirthdayService
```

## Testing

The service follows interface-based design for easy mocking and testing:

- `persister` interface allows database layer mocking
- Business logic is separated for unit testing
- All operations include structured logging for debugging

## Dependencies

- **Foundation packages**: config, database/postgres, session, pagination, order, derror, log
- **External**: database/sql, slog, google/uuid
- **User service**: Uses UserID type from user proto package
