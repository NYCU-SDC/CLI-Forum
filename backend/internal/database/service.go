package database

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

type MigrationService struct {
	logger *zap.Logger
	dbURL  string
	files  string
}

func NewMigrationService(logger *zap.Logger, dbURL string, files string) *MigrationService {
	return &MigrationService{
		logger: logger,
		dbURL:  dbURL,
		files:  files,
	}
}

func (m MigrationService) Up() {
	// initialize migration
	migration, err := migrate.New(
		// migration files
		"file://"+m.files,
		// connection string
		m.dbURL,
	)
	if err != nil {
		m.logger.Fatal("failed to create migration", zap.String("migration file", m.files), zap.String("DB URL", m.dbURL), zap.Error(err))
	}
	// if the database was already migrated, it will return ErrNoChange
	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		m.logger.Fatal("failed to migrate", zap.String("migration file", m.files), zap.String("DB URL", m.dbURL), zap.Error(err))
	}
	m.logger.Info("migrated successfully", zap.String("migration file", m.files))
}

func (m MigrationService) Down() {
	// initialize migration
	migration, err := migrate.New(
		// migration files
		"file://"+m.files,
		// connection string
		m.dbURL,
	)
	if err != nil {
		m.logger.Fatal("failed to create migration", zap.String("migration file", m.files), zap.String("DB URL", m.dbURL), zap.Error(err))
	}
	// if the database was already migrated, it will return ErrNoChange
	if err := migration.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		m.logger.Fatal("failed to migrate", zap.String("migration file", m.files), zap.String("DB URL", m.dbURL), zap.Error(err))
	}
	m.logger.Info("migrated successfully", zap.String("migration file", m.files))
}
