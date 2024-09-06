package genreuc

import (
	"context"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/txmanager"
	"github.com/rs/zerolog"
)

type UseCase struct {
	genreRepository genreRepository
	txManager       txmanager.Transactor
	log             *zerolog.Logger
}

func New(genreRepository genreRepository, txManager txmanager.Transactor, log *zerolog.Logger) *UseCase {
	return &UseCase{
		genreRepository: genreRepository,
		txManager:       txManager,
		log:             log,
	}
}

type genreRepository interface {
	GetMany(ctx context.Context, pgn pagination.CursorPagination) ([]genre.Genre, pagination.CursorPair, error)
}
