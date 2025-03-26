package database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationService struct {
	dbURL string
	files string
}

func NewMigrationService(dbURL string, files string) *MigrationService {
	return &MigrationService{
		dbURL: dbURL,
		files: files,
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
		fmt.Println("failed to create migrate : ", err)
	}
	// if the database was already migrated, it will return ErrNoChange
	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("failed to migrate : ", err)
	}
	fmt.Println("migration success")
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
		fmt.Println("failed to create migrate : ", err)
	}
	// if the database was already migrated, it will return ErrNoChange
	if err := migration.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("failed to migrate : ", err)
	}
	fmt.Println("migrate success")
}
