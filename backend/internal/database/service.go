package database

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var migratetionFiles = "internal/database/migrations"

func MigrateUP() {
	// initialize migratetion
	m, err := migrate.New(
		// migration files
		"file://"+migratetionFiles,
		// connection string
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		fmt.Println("failed to create migrate : ", err)
	}
	// if the database was already migrated, it will return ErrNoChange
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println("failed to migrate : ", err)
	}
	fmt.Println("migratrate success")
}

func MigrateDown() {
	// initialize migratetion
	m, err := migrate.New(
		// migration files
		"file://"+migratetionFiles,
		// connection string
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		fmt.Println("failed to create migrate : ", err)
	}
	// if the database was already migrated, it will return ErrNoChange
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		fmt.Println("failed to migrate : ", err)
	}
	fmt.Println("migrate success")
}
