package artistuc

import (
	"context"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/txmanager"
	"github.com/rs/zerolog"
)

type UseCase struct {
	artistRepository artistRepository
	txManager        txmanager.Transactor
	log              *zerolog.Logger
}

func New(artistRepository artistRepository, txManager txmanager.Transactor, log *zerolog.Logger) *UseCase {
	return &UseCase{
		artistRepository: artistRepository,
		txManager:        txManager,
		log:              log,
	}
}

type artistRepository interface {
	GetMany(ctx context.Context, pgn pagination.CursorPagination) ([]artist.Artist, pagination.CursorPair, error)
}
