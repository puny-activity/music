package postgres

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/puny-activity/music/pkg/werr"
)

type Postgres struct {
	*sqlx.DB
}

func New(connectionString string) (*Postgres, error) {
	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to postgres: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &Postgres{
		DB: db,
	}, nil
}

func (p *Postgres) RunMigrations(migrationPath string) error {
	goose.SetTableName("migrations")

	err := goose.Up(p.DB.DB, migrationPath)
	if err != nil {
		return werr.WrapSE("failed to run migrations", err)
	}

	return nil
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}
