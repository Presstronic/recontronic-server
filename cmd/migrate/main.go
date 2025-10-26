package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/presstronic/recontronic-server/internal/config"
)

const (
	migrationsPath = "file://migrations"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Build database connection string
	dbURL := buildDatabaseURL(cfg.Database)

	// Create migrator instance
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}
	defer m.Close()

	// Execute command
	command := os.Args[1]
	switch command {
	case "up":
		handleUp(m)
	case "down":
		handleDown(m)
	case "version":
		handleVersion(m)
	case "force":
		handleForce(m)
	case "drop":
		handleDrop(m)
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func buildDatabaseURL(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
}

func handleUp(m *migrate.Migrate) {
	fmt.Println("Running migrations up...")
	err := m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("✓ No migrations to run (already up to date)")
			return
		}
		log.Fatalf("Failed to run migrations up: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Fatalf("Failed to get version: %v", err)
	}

	if dirty {
		fmt.Printf("⚠ Migrations completed but database is in dirty state at version %d\n", version)
		fmt.Println("  Run 'migrate force <version>' to fix")
	} else {
		fmt.Printf("✓ Successfully migrated to version %d\n", version)
	}
}

func handleDown(m *migrate.Migrate) {
	// Check if a step count was provided
	steps := 1
	if len(os.Args) > 2 {
		var err error
		steps, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid step count: %v", err)
		}
	}

	fmt.Printf("Rolling back %d migration(s)...\n", steps)
	err := m.Steps(-steps)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("✓ No migrations to roll back")
			return
		}
		log.Fatalf("Failed to roll back migrations: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("✓ Successfully rolled back all migrations")
			return
		}
		log.Fatalf("Failed to get version: %v", err)
	}

	if dirty {
		fmt.Printf("⚠ Rollback completed but database is in dirty state at version %d\n", version)
		fmt.Println("  Run 'migrate force <version>' to fix")
	} else {
		fmt.Printf("✓ Successfully rolled back to version %d\n", version)
	}
}

func handleVersion(m *migrate.Migrate) {
	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("No migrations have been run yet")
			return
		}
		log.Fatalf("Failed to get version: %v", err)
	}

	if dirty {
		fmt.Printf("Current version: %d (dirty)\n", version)
		fmt.Println("\n⚠ Warning: Database is in dirty state!")
		fmt.Println("This usually happens when a migration fails partway through.")
		fmt.Println("You can fix this by running: migrate force <version>")
	} else {
		fmt.Printf("Current version: %d\n", version)
	}
}

func handleForce(m *migrate.Migrate) {
	if len(os.Args) < 3 {
		fmt.Println("Error: version number required")
		fmt.Println("Usage: migrate force <version>")
		os.Exit(1)
	}

	version, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid version number: %v", err)
	}

	fmt.Printf("Forcing database version to %d...\n", version)
	err = m.Force(version)
	if err != nil {
		log.Fatalf("Failed to force version: %v", err)
	}

	fmt.Printf("✓ Successfully forced version to %d\n", version)
	fmt.Println("\nNote: This only updates the migration version number.")
	fmt.Println("It does NOT run any migrations. Use 'migrate up' to apply pending migrations.")
}

func handleDrop(m *migrate.Migrate) {
	fmt.Println("⚠  WARNING: This will DROP all tables in the database!")
	fmt.Println("Type 'yes' to confirm: ")

	var confirmation string
	fmt.Scanln(&confirmation)

	if confirmation != "yes" {
		fmt.Println("Aborted.")
		return
	}

	fmt.Println("Dropping all tables...")
	err := m.Drop()
	if err != nil {
		log.Fatalf("Failed to drop database: %v", err)
	}

	fmt.Println("✓ Successfully dropped all tables")
}

func printUsage() {
	fmt.Println("Database Migration Tool")
	fmt.Println("\nUsage:")
	fmt.Println("  migrate <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  up              Apply all pending migrations")
	fmt.Println("  down [N]        Rollback N migrations (default: 1)")
	fmt.Println("  version         Show current migration version")
	fmt.Println("  force <version> Force set migration version (fixes dirty state)")
	fmt.Println("  drop            Drop all tables (requires confirmation)")
	fmt.Println("\nExamples:")
	fmt.Println("  migrate up                  # Run all pending migrations")
	fmt.Println("  migrate down                # Rollback last migration")
	fmt.Println("  migrate down 3              # Rollback last 3 migrations")
	fmt.Println("  migrate version             # Show current version")
	fmt.Println("  migrate force 20251026040100  # Force version (emergency only)")
	fmt.Println("\nConfiguration:")
	fmt.Println("  Database connection is configured via config file or environment variables")
	fmt.Println("  See internal/config/config.go for details")
}
