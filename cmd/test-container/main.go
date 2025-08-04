package main

import (
	"fmt"
	"log/slog"

	"github.com/kianooshaz/skeleton/internal/container"
)

func main() {
	fmt.Println("=== Testing Wire Container ===")

	// Test that we can create the container (this will test all our wire setup)
	fmt.Println("Creating web container...")
	c, err := container.NewWebContainer()
	if err != nil {
		slog.Error("Failed to create container", "error", err)
		// Don't exit, show that container creation worked
		fmt.Printf("Container creation failed (expected if DB not running): %v\n", err)
		return
	}

	fmt.Println("✅ Container created successfully!")
	fmt.Printf("✅ Logger configured: %v\n", c.Logger != nil)
	fmt.Printf("✅ Config loaded: %v\n", c.Config != nil)
	fmt.Printf("✅ WebService created: %v\n", c.WebService != nil)
	fmt.Printf("✅ Database connection: %v\n", c.DB != nil)

	if c.Config != nil {
		fmt.Printf("   - Shutdown timeout: %v\n", c.Config.ShutdownTimeout)
		fmt.Printf("   - Server address: %s\n", c.Config.RestServer.Address)
		fmt.Printf("   - Log level: %s\n", c.Config.Logger.Level)
		fmt.Printf("   - DB name: %s\n", c.Config.Postgres.Name)
	}

	// Cleanup
	if err := c.Close(); err != nil {
		fmt.Printf("Close error: %v\n", err)
	}

	fmt.Println("=== Wire Container Test Complete ===")
}
