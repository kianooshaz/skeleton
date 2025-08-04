package container

import (
	"os"
	"strings"
	"testing"

	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvideLogger(t *testing.T) {
	// Create a test config
	cfg := &AppConfig{
		Logger: log.LoggerConfig{
			Environment: "development",
			Level:       "debug",
			AddSource:   true,
			Format:      "json",
		},
	}

	logger := ProvideLogger(cfg)
	require.NotNil(t, logger)

	// Test that we can log without panic
	logger.Info("Test log message")
}

func TestProvideLoggerWithDifferentConfigs(t *testing.T) {
	tests := []struct {
		name   string
		config log.LoggerConfig
	}{
		{
			name: "production_info_text",
			config: log.LoggerConfig{
				Environment: "production",
				Level:       "info",
				AddSource:   false,
				Format:      "text",
			},
		},
		{
			name: "staging_warn_json",
			config: log.LoggerConfig{
				Environment: "staging",
				Level:       "warn",
				AddSource:   false,
				Format:      "json",
			},
		},
		{
			name: "development_debug_json",
			config: log.LoggerConfig{
				Environment: "development",
				Level:       "debug",
				AddSource:   true,
				Format:      "json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AppConfig{
				Logger: tt.config,
			}

			logger := ProvideLogger(cfg)
			require.NotNil(t, logger)

			// Test that we can log without panic
			logger.Info("Test log message for", "config", tt.name)
		})
	}
}

func TestProvideAppConfigWithValidFile(t *testing.T) {
	// Set config path for testing
	os.Setenv("CONFIG_PATH", "../../config.yml")
	defer os.Unsetenv("CONFIG_PATH")

	config, err := ProvideAppConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	// Test that config values are loaded correctly
	assert.NotEmpty(t, config.RestServer.Address, "REST server address should not be empty")
	assert.NotEmpty(t, config.Database.Host, "Database host should not be empty")
	assert.Greater(t, config.Database.Port, 0, "Database port should be greater than 0")
}

func TestNewWebContainer(t *testing.T) {
	// Skip this test if database is not available (expected in CI/testing environments)
	// Set config path for testing
	os.Setenv("CONFIG_PATH", "../../config.yml")
	defer os.Unsetenv("CONFIG_PATH")

	container, err := NewWebContainer()
	if err != nil {
		// If it's a database connection error, that's expected in test environment
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "dial tcp") {
			t.Skip("Skipping test due to database connection issue (expected in test environment)")
			return
		}
		// For other errors, fail the test
		require.NoError(t, err)
	}
	require.NotNil(t, container)

	// Test graceful shutdown
	err = container.Close()
	assert.NoError(t, err)
}
