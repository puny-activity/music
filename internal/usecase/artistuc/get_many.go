package artistuc

import (
	"context"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/werr"
)

func (u *UseCase) GetMany(ctx context.Context, pgn pagination.CursorPagination) ([]artist.Artist, pagination.CursorPair, error) {
	artists, cursorPair, err := u.artistRepository.GetMany(ctx, pgn)
	if err != nil {
		return nil, pagination.CursorPair{}, werr.WrapSE("failed to get artists", err)
	}

	return artists, cursorPair, nil
}
