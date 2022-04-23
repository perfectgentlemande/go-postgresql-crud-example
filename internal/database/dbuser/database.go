package dbuser

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	ConnStr string `yaml:"connstr"`
}

type Database struct {
	db *sqlx.DB
}

func NewDatabase(ctx context.Context, config *Config) (*Database, error) {
	db, err := sqlx.Connect("pgx", config.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("cannot create postgresql client: %w", err)
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot create driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("cannot create db instance: %w", err)
	}
	err = m.Up()
	if err != nil && !strings.Contains(err.Error(), "no change") {
		return nil, fmt.Errorf("cannot migrate db: %w", err)
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
