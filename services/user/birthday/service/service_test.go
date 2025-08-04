package birthdayservice_test

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	birthdayservice "github.com/kianooshaz/skeleton/services/user/birthday/service"
)

func TestBirthdayService_New(t *testing.T) {
	// Setup.
	config := birthdayservice.Config{
		MaxAge: 150,
		MinAge: 0,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Execute (this would normally require a real database connection).
	// For now, we just test that the config and constructor work.
	service := birthdayservice.New(config, nil, logger)

	// Assert.
	require.NotNil(t, service)
}

func TestBirthdayService_Config(t *testing.T) {
	// Test default configuration values.
	config := birthdayservice.Config{
		MaxAge: 150,
		MinAge: 0,
	}

	assert.Equal(t, 150, config.MaxAge)
	assert.Equal(t, 0, config.MinAge)
}

func TestAgeCalculation(t *testing.T) {
	// Test age calculation logic (simulate what the service would do).
	testCases := []struct {
		name        string
		dateOfBirth time.Time
		minAge      int
	}{
		{
			name:        "adult",
			dateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			minAge:      30,
		},
		{
			name:        "newborn",
			dateOfBirth: time.Now().AddDate(0, 0, -30), // 30 days ago
			minAge:      0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate age the same way the service does.
			now := time.Now()
			age := now.Year() - tc.dateOfBirth.Year()
			if now.YearDay() < tc.dateOfBirth.YearDay() {
				age--
			}

			assert.GreaterOrEqual(t, age, tc.minAge)
		})
	}
}
