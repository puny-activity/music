package songuc

import (
	"context"
	"github.com/puny-activity/music/internal/entity/song"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/filters"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/werr"
)

func (u *UseCase) Get(ctx context.Context, search string, fltr filters.Filter, pgn pagination.CursorPagination) ([]song.Song, pagination.CursorPair, error) {
	if pgn.Limit < 1 || pgn.Limit > 100 {
		return nil, pagination.CursorPair{}, errs.InvalidLimitParameter
	}

	songs, cursorPair, err := u.songRepository.GetMany(ctx, search, fltr, pgn)
	if err != nil {
		return nil, pagination.CursorPair{}, werr.WrapSE("failed to get songs", err)
	}

	return songs, cursorPair, nil
}
