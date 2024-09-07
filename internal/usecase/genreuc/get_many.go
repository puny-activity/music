package genreuc

import (
	"context"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/werr"
)

func (u *UseCase) GetMany(ctx context.Context, pgn pagination.CursorPagination) ([]genre.Genre, pagination.CursorPair, error) {
	if pgn.Limit < 1 || pgn.Limit > 100 {
		return nil, pagination.CursorPair{}, errs.InvalidLimitParameter
	}

	genres, cursorPair, err := u.genreRepository.GetMany(ctx, pgn)
	if err != nil {
		return nil, pagination.CursorPair{}, werr.WrapSE("failed to get genres", err)
	}

	return genres, cursorPair, nil
}
