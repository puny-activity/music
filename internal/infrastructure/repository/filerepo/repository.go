package filerepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/pkg/txmanager"
	"github.com/rs/zerolog"
)

type Repository struct {
	db        *sqlx.DB
	txManager txmanager.Transactor
	log       *zerolog.Logger
}

func (r Repository) DeleteTx(ctx context.Context, tx *sqlx.Tx, fileID remotefile.ID) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) GetAllTx(ctx context.Context, tx *sqlx.Tx) ([]remotefile.File, error) {
	//TODO implement me
	panic("implement me")
}

func New(db *sqlx.DB, txManager txmanager.Transactor, log *zerolog.Logger) *Repository {
	return &Repository{
		db:        db,
		txManager: txManager,
		log:       log,
	}
}
