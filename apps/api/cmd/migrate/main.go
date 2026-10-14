package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/satset/api/internal/platform/config"
	"github.com/satset/api/internal/platform/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(cfg.AppEnv)

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	command := args[0]
	migrationSource := "file://db/migrations"
	dbURL := fmt.Sprintf("mysql://%s", cfg.DSN())

	m, err := migrate.New(migrationSource, dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database migration engine")
	}
	defer m.Close()

	switch command {
	case "up":
		log.Info().Msg("Applying all pending migrations (up)...")
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal().Err(err).Msg("Migration failed")
		}
		log.Info().Msg("Migrations applied successfully or already up-to-date")

	case "down":
		log.Info().Msg("Rolling back one migration step (down)...")
		if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal().Err(err).Msg("Rollback failed")
		}
		log.Info().Msg("Migration rollback completed")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get current migration version")
		}
		log.Info().Uint("version", version).Bool("dirty", dirty).Msg("Current database migration version")

	case "force":
		if len(args) < 2 {
			log.Fatal().Msg("Version number required for force command: migrate force <version>")
		}
		v, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal().Err(err).Msg("Invalid version number")
		}
		log.Warn().Int("version", v).Msg("Forcing migration version...")
		if err := m.Force(v); err != nil {
			log.Fatal().Err(err).Msg("Failed to force migration version")
		}
		log.Info().Msg("Migration version forced successfully")

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run cmd/migrate/main.go <command>")
	fmt.Println("Commands:")
	fmt.Println("  up       Apply all pending migrations")
	fmt.Println("  down     Rollback 1 migration step")
	fmt.Println("  version  Show current migration version")
	fmt.Println("  force <version> Force set database migration version")
}
