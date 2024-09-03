package songrepo

import (
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/pkg/txmanager"
	"github.com/rs/zerolog"
)

type Repository struct {
	db        *sqlx.DB
	txManager txmanager.Transactor
	log       *zerolog.Logger
}

func New(db *sqlx.DB, txManager txmanager.Transactor, log *zerolog.Logger) *Repository {
	return &Repository{
		db:        db,
		txManager: txManager,
		log:       log,
	}
}
