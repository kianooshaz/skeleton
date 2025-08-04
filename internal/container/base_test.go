package container

import (
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      BaseConfig
		expectValid bool
		description string
	}{
		{
			name: "valid_base_config",
			config: BaseConfig{
				ShutdownTimeout: 30 * time.Second,
				Logger: log.LoggerConfig{
					Environment: "development",
					Level:       "debug",
					AddSource:   true,
					Format:      "json",
				},
				Database: postgres.Config{
					Name:        "testdb",
					Host:        "localhost",
					Port:        5432,
					User:        "testuser",
					Password:    "testpass",
					SSLMode:     "disable",
					PingTimeout: 10 * time.Second,
				},
			},
			expectValid: true,
			description: "Valid base configuration with all required fields",
		},
		{
			name: "minimal_base_config",
			config: BaseConfig{
				ShutdownTimeout: 15 * time.Second,
				Logger: log.LoggerConfig{
					Environment: "production",
					Level:       "info",
					AddSource:   false,
					Format:      "text",
				},
				Database: postgres.Config{
					Name:        "proddb",
					Host:        "db.example.com",
					Port:        5432,
					User:        "produser",
					Password:    "prodpass",
					SSLMode:     "require",
					PingTimeout: 5 * time.Second,
				},
			},
			expectValid: true,
			description: "Minimal valid configuration for production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the configuration is structurally valid
			assert.NotEmpty(t, tt.config.Logger.Environment, "Logger environment should not be empty")
			assert.NotEmpty(t, tt.config.Logger.Level, "Logger level should not be empty")
			assert.NotEmpty(t, tt.config.Logger.Format, "Logger format should not be empty")
			assert.Greater(t, tt.config.ShutdownTimeout, time.Duration(0), "Shutdown timeout should be positive")

			// Test database config
			assert.NotEmpty(t, tt.config.Database.Name, "Database name should not be empty")
			assert.NotEmpty(t, tt.config.Database.Host, "Database host should not be empty")
			assert.Greater(t, tt.config.Database.Port, 0, "Database port should be positive")
			assert.NotEmpty(t, tt.config.Database.User, "Database user should not be empty")
		})
	}
}

func TestBaseContainer(t *testing.T) {
	// Create a test app config
	cfg := &AppConfig{
		ShutdownTimeout: 30 * time.Second,
		Logger: log.LoggerConfig{
			Environment: "development",
			Level:       "debug",
			AddSource:   true,
			Format:      "json",
		},
		Database: postgres.Config{
			Name:        "testdb",
			Host:        "localhost",
			Port:        5432,
			User:        "testuser",
			Password:    "testpass",
			SSLMode:     "disable",
			PingTimeout: 10 * time.Second,
		},
	}

	// Create logger
	logger := log.NewLogger(cfg.Logger)
	require.NotNil(t, logger)

	// Create base container (without database connection for testing)
	baseContainer := &BaseContainer{
		Config: cfg,
		DB:     nil, // Skip database for unit test
		Logger: logger,
	}

	// Test that base container implements IContainer interface
	var _ IContainer = baseContainer

	// Test configuration access
	assert.Equal(t, cfg, baseContainer.Config)
	assert.Equal(t, logger, baseContainer.Logger)

	// Test close (should not error even with nil DB)
	err := baseContainer.Close()
	assert.NoError(t, err)
}

func TestProvideAppConfig(t *testing.T) {
	// Set up test config file
	os.Setenv("CONFIG_PATH", "../../config.yml")
	defer os.Unsetenv("CONFIG_PATH")

	config, err := ProvideAppConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	// Test that all required config sections are present
	assert.NotEmpty(t, config.Logger.Environment, "Logger environment should be configured")
	assert.NotEmpty(t, config.Logger.Level, "Logger level should be configured")
	assert.NotEmpty(t, config.Logger.Format, "Logger format should be configured")

	assert.NotEmpty(t, config.RestServer.Address, "REST server address should be configured")

	assert.NotEmpty(t, config.Database.Name, "Database name should be configured")
	assert.NotEmpty(t, config.Database.Host, "Database host should be configured")
	assert.Greater(t, config.Database.Port, 0, "Database port should be configured")

	assert.Greater(t, config.ShutdownTimeout, time.Duration(0), "Shutdown timeout should be configured")
}

func TestProvideLoggerFromBase(t *testing.T) {
	tests := []struct {
		name        string
		config      log.LoggerConfig
		expectPanic bool
		description string
	}{
		{
			name: "development_json_logger",
			config: log.LoggerConfig{
				Environment: "development",
				Level:       "debug",
				AddSource:   true,
				Format:      "json",
			},
			expectPanic: false,
			description: "Development logger with JSON format and debug level",
		},
		{
			name: "production_text_logger",
			config: log.LoggerConfig{
				Environment: "production",
				Level:       "info",
				AddSource:   false,
				Format:      "text",
			},
			expectPanic: false,
			description: "Production logger with text format and info level",
		},
		{
			name: "staging_logger",
			config: log.LoggerConfig{
				Environment: "staging",
				Level:       "warn",
				AddSource:   true,
				Format:      "json",
			},
			expectPanic: false,
			description: "Staging logger with warning level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AppConfig{
				Logger: tt.config,
			}

			var logger *slog.Logger
			if tt.expectPanic {
				assert.Panics(t, func() {
					logger = ProvideLogger(cfg)
				}, "Should panic with invalid config")
			} else {
				assert.NotPanics(t, func() {
					logger = ProvideLogger(cfg)
				}, "Should not panic with valid config")

				require.NotNil(t, logger, "Logger should not be nil")

				// Test that logger can be used without panic
				assert.NotPanics(t, func() {
					logger.Info("Test message", "test", tt.name)
				}, "Logger should be usable")
			}
		})
	}
}

func TestProvideDatabase(t *testing.T) {
	cfg := &AppConfig{
		Database: postgres.Config{
			Name:        "testdb",
			Host:        "nonexistent-host", // Use nonexistent host to avoid real DB dependency
			Port:        5432,
			User:        "testuser",
			Password:    "testpass",
			SSLMode:     "disable",
			PingTimeout: 1 * time.Second, // Short timeout for testing
		},
	}

	// This will fail to connect, but should not panic during creation
	db, err := ProvideDatabase(cfg)

	// We expect an error due to connection failure, but the function should handle it gracefully
	if err != nil {
		// Connection error is expected with nonexistent host
		assert.Contains(t, err.Error(), "connection refused", "Should get connection error")
		assert.Nil(t, db, "DB should be nil on connection error")
	} else {
		// If somehow it succeeds (shouldn't happen with nonexistent host)
		require.NotNil(t, db, "DB should not be nil if no error")
		// Clean up
		err = db.Close()
		assert.NoError(t, err, "Should be able to close DB")
	}
}

func TestProvideBaseContainer(t *testing.T) {
	cfg := &AppConfig{
		ShutdownTimeout: 30 * time.Second,
		Logger: log.LoggerConfig{
			Environment: "development",
			Level:       "debug",
			AddSource:   true,
			Format:      "json",
		},
		Database: postgres.Config{
			Name:        "testdb",
			Host:        "localhost",
			Port:        5432,
			User:        "testuser",
			Password:    "testpass",
			SSLMode:     "disable",
			PingTimeout: 10 * time.Second,
		},
	}

	logger := log.NewLogger(cfg.Logger)
	require.NotNil(t, logger)

	// Test with nil DB (common in unit tests)
	var db *sql.DB = nil

	baseContainer := ProvideBaseContainer(cfg, db, logger)
	require.NotNil(t, baseContainer)

	assert.Equal(t, cfg, baseContainer.Config)
	assert.Equal(t, db, baseContainer.DB)
	assert.Equal(t, logger, baseContainer.Logger)

	// Test that it implements IContainer
	var _ IContainer = baseContainer

	// Test close
	err := baseContainer.Close()
	assert.NoError(t, err, "Close should succeed even with nil DB")
}
