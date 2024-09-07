package albumuc

import (
	"context"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/werr"
)

func (u *UseCase) GetMany(ctx context.Context, pgn pagination.CursorPagination) ([]album.Album, pagination.CursorPair, error) {
	if pgn.Limit < 1 || pgn.Limit > 100 {
		return nil, pagination.CursorPair{}, errs.InvalidLimitParameter
	}

	albums, cursorPair, err := u.albumRepository.GetMany(ctx, pgn)
	if err != nil {
		return nil, pagination.CursorPair{}, werr.WrapSE("failed to get albums", err)
	}

	return albums, cursorPair, nil
}
