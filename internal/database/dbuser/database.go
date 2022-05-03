package dbuser

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	ConnStr string `yaml:"connstr"`
}

type Database struct {
	db *sqlx.DB
}

func NewDatabase(ctx context.Context, config *Config, migrationsContent fs.FS) (*Database, error) {
	sourceInstance, err := iofs.New(migrationsContent, "migrations")
	if err != nil {
		log.Fatal("cannot create source instance", err)
	}
	db, err := sqlx.Connect("pgx", config.ConnStr)
	if err != nil {
		log.Fatal("cannot open postgres:", err)
	}
	targetInstance, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal("cannot target instance:", err)
	}
	m, err := migrate.NewWithInstance("iofs", sourceInstance, "postgres", targetInstance)
	if err != nil {
		return &Database{}, fmt.Errorf("cannot create migrate object: %w", err)
	}
	err = m.Up()
	if err != nil && !strings.Contains(err.Error(), "no change") {
		return &Database{}, fmt.Errorf("cannot migrate db: %w", err)
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *Database) Close(ctx context.Context) error {
	return d.db.Close()
}
