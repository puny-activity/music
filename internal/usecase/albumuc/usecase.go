package albumuc

import (
	"context"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/txmanager"
	"github.com/rs/zerolog"
)

type UseCase struct {
	albumRepository albumRepository
	txManager       txmanager.Transactor
	log             *zerolog.Logger
}

func New(albumRepository albumRepository, txManager txmanager.Transactor, log *zerolog.Logger) *UseCase {
	return &UseCase{
		albumRepository: albumRepository,
		txManager:       txManager,
		log:             log,
	}
}

type albumRepository interface {
	GetMany(ctx context.Context, pgn pagination.CursorPagination) ([]album.Album, pagination.CursorPair, error)
}
