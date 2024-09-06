package genrerepo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
	"github.com/puny-activity/music/pkg/werr"
	"strings"
)

type getManyDTO struct {
	GenreID   string `db:"genre_id"`
	GenreName string `db:"genre_name"`
	SongCount int    `db:"song_count"`
	Covers    []byte `db:"covers"`
}

func (dto getManyDTO) ToGenre() (genre.Genre, error) {
	genreID, err := genre.ParseID(dto.GenreID)
	if err != nil {
		return genre.Genre{}, werr.WrapSE("failed to parse genre id", err)
	}

	coverIDsStr := make([]string, 0)
	err = json.Unmarshal(dto.Covers, &coverIDsStr)
	if err != nil {
		return genre.Genre{}, werr.WrapSE("failed to parse covers", err)
	}
	coverIDs := make([]cover.ID, len(coverIDsStr))
	for i := range coverIDsStr {
		coverIDs[i], err = cover.ParseID(coverIDsStr[i])
		if err != nil {
			return genre.Genre{}, werr.WrapSE("failed to parse cover id", err)
		}
	}

	return genre.Genre{
		Base: genre.Base{
			ID:   &genreID,
			Name: dto.GenreName,
		},
		SongCount: dto.SongCount,
		CoversIDs: coverIDs,
	}, nil
}

const (
	getManyPaginationName      = "ra.genre_name"
	getManyPaginationSongCount = "ra.song_count"
)

func getManyPgnParameterConvert(before []pagination.Parameter) ([]pagination.Parameter, error) {
	defaultPaginators := []pagination.Parameter{
		{
			FieldName: getManyPaginationName,
			SortOrder: pagination.ASC,
		},
		{
			FieldName: getManyPaginationSongCount,
			SortOrder: pagination.DESC,
		},
	}
	appliedPaginators := make(map[string]struct{})

	after := make([]pagination.Parameter, 0)
	for _, param := range before {
		switch param.FieldName {
		case genre.PaginationName:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationName,
				SortOrder: param.SortOrder,
			})
		case genre.PaginationSongCount:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationSongCount,
				SortOrder: param.SortOrder,
			})
		default:
			return nil, errs.UnknownSortParameter
		}
		appliedPaginators[after[len(after)-1].FieldName] = struct{}{}
	}

	for _, defaultPaginator := range defaultPaginators {
		_, ok := appliedPaginators[defaultPaginator.FieldName]
		if !ok {
			after = append(after, defaultPaginator)
		}
	}

	return after, nil
}

func (r *Repository) GetMany(ctx context.Context, pgn pagination.CursorPagination) ([]genre.Genre, pagination.CursorPair, error) {
	return r.getMany(ctx, r.db, pgn)
}

func (r *Repository) GetManyTx(ctx context.Context, tx *sqlx.Tx, pgn pagination.CursorPagination) ([]genre.Genre, pagination.CursorPair, error) {
	return r.getMany(ctx, tx, pgn)
}

func (r *Repository) getMany(ctx context.Context, queryer queryer.Queryer, pgn pagination.CursorPagination) ([]genre.Genre, pagination.CursorPair, error) {
	if pgn.Limit < 1 {
		return nil, pagination.CursorPair{}, errs.InvalidLimitParameter
	}
	if pgn.Cursor == nil {
		var err error
		pgn.Parameters, err = getManyPgnParameterConvert(pgn.Parameters)
		if err != nil {
			return nil, pagination.CursorPair{}, werr.WrapSE("failed to convert parameters", err)
		}
	}

	query, args, err := getManyBuildQuery(pgn)
	if err != nil {
		return nil, pagination.CursorPair{}, werr.WrapSE("failed to build query", err)
	}

	var dtoList []getManyDTO
	err = queryer.SelectContext(ctx, &dtoList, query, args...)
	if err != nil {
		return nil, pagination.CursorPair{}, werr.WrapSE("failed to get genre dtos", err)
	}
	if len(dtoList) == 0 {
		return make([]genre.Genre, 0), pagination.CursorPair{}, nil
	}

	var genres []genre.Genre
	for i, dto := range dtoList {
		if (len(dtoList) > pgn.Limit) && (i == len(dtoList)-1) {
			continue
		}
		oneGenre, err := dto.ToGenre()
		if err != nil {
			return nil, pagination.CursorPair{}, err
		}
		genres = append(genres, oneGenre)
	}
	if pgn.Cursor != nil && pgn.Cursor.Type() == pagination.CursorPrev {
		genres = util.Reverse(genres)
	}

	var prevCursor *pagination.Cursor = nil
	var nextCursor *pagination.Cursor = nil

	if pgn.Parameters != nil {
		if len(dtoList) == pgn.Limit+1 {
			lastGenre := genres[len(genres)-1]
			nextCursor = util.ToPointer(pagination.NewCursor(pagination.CursorNext))
			nextCursor = util.ToPointer(getManyIncludeParameters(*nextCursor, pgn.Parameters, lastGenre))
		}
	} else if pgn.Cursor != nil {
		if pgn.Cursor.Type() == pagination.CursorNext {
			if len(genres) > 0 {
				firstGenre := genres[0]
				prevCursor = util.ToPointer(pagination.NewCursor(pagination.CursorPrev))
				prevCursor = util.ToPointer(getManyIncludeCursor(*prevCursor, pgn.Cursor, firstGenre))
			}
			if len(dtoList) == pgn.Limit+1 {
				lastGenre := genres[len(genres)-1]
				nextCursor = util.ToPointer(pagination.NewCursor(pagination.CursorNext))
				nextCursor = util.ToPointer(getManyIncludeCursor(*nextCursor, pgn.Cursor, lastGenre))
			}
		} else if pgn.Cursor.Type() == pagination.CursorPrev {
			if len(dtoList) == pgn.Limit+1 {
				firstGenre := genres[0]
				prevCursor = util.ToPointer(pagination.NewCursor(pagination.CursorPrev))
				prevCursor = util.ToPointer(getManyIncludeCursor(*prevCursor, pgn.Cursor, firstGenre))
			}
			if len(genres) > 0 {
				lastGenre := genres[len(genres)-1]
				nextCursor = util.ToPointer(pagination.NewCursor(pagination.CursorNext))
				nextCursor = util.ToPointer(getManyIncludeCursor(*nextCursor, pgn.Cursor, lastGenre))
			}
		}
	}

	cursorPair := pagination.CursorPair{
		PrevCursor: prevCursor,
		NextCursor: nextCursor,
	}

	return genres, cursorPair, nil
}

func getManyIncludeParameters(cursor pagination.Cursor, parameters []pagination.Parameter, targetGenre genre.Genre) pagination.Cursor {
	for _, parameter := range parameters {
		var value any
		switch parameter.FieldName {
		case getManyPaginationName:
			value = targetGenre.Name
		case getManyPaginationSongCount:
			value = targetGenre.SongCount
		default:
			continue
		}
		cursor = cursor.Add(pagination.NewCursorParameter(pagination.Parameter{
			FieldName: parameter.FieldName,
			SortOrder: parameter.SortOrder,
		}, value))
	}
	return cursor
}

func getManyIncludeCursor(cursor pagination.Cursor, cursor2 *pagination.Cursor, targetGenre genre.Genre) pagination.Cursor {
	for _, parameter := range cursor2.GetAll() {
		var value any
		switch parameter.FieldName {
		case getManyPaginationName:
			value = targetGenre.Name
		case getManyPaginationSongCount:
			value = targetGenre.SongCount
		default:
			continue
		}
		cursor = cursor.Add(pagination.NewCursorParameter(pagination.Parameter{
			FieldName: parameter.FieldName,
			SortOrder: parameter.SortOrder,
		}, value))
	}
	return cursor
}

func getManyBuildQuery(pgn pagination.CursorPagination) (string, []any, error) {
	var queryBuilder strings.Builder
	var args []any

	queryBuilder.WriteString(`
WITH resp_covers AS (SELECT a.id        genre_id,
                            c.id        cover_id,
                            COUNT(c.id) song_count
                     FROM genres a
                              JOIN songs s ON a.id = s.genre_id
                              JOIN covers c ON s.cover_id = c.id
                     GROUP BY a.id, c.id
                     ORDER BY genre_id, song_count DESC),
     resp_genres AS (SELECT a.id        AS genre_id,
                            a.name      AS genre_name,
                            COUNT(s.id) AS song_count
                     FROM genres a
                              JOIN songs s ON a.id = s.genre_id
                     GROUP BY a.id, a.name)
SELECT ra.genre_id,
       ra.genre_name,
       ra.song_count,
       TO_JSON(ARRAY_REMOVE(ARRAY_AGG(rc.cover_id), NULL)) AS covers
FROM resp_genres ra
         LEFT JOIN resp_covers rc ON ra.genre_id = rc.genre_id
`)

	if pgn.Cursor != nil {
		whereClause, err := pgn.Cursor.BuildWhere()
		if err != nil {
			return "", nil, werr.WrapSE("failed to build where query", err)
		}
		if whereClause != "" {
			queryBuilder.WriteString(" WHERE ")
			queryBuilder.WriteString(whereClause)
		}
	}

	queryBuilder.WriteString(" GROUP BY ra.genre_id, ra.genre_name, ra.song_count ")

	orderClauses := make([]string, 0)
	if pgn.Parameters != nil {
		for _, param := range pgn.Parameters {
			orderClause := fmt.Sprintf("%s %s", param.FieldName, param.SortOrder.SQLString())
			orderClauses = append(orderClauses, orderClause)
		}
	} else if pgn.Cursor != nil {
		for _, param := range pgn.Cursor.GetAll() {
			orderClause := ""
			if pgn.Cursor.Type() == pagination.CursorPrev {
				orderClause = fmt.Sprintf("%s %s", param.FieldName, param.SortOrder.Invert().SQLString())
			} else if pgn.Cursor.Type() == pagination.CursorNext {
				orderClause = fmt.Sprintf("%s %s", param.FieldName, param.SortOrder.SQLString())
			}
			orderClauses = append(orderClauses, orderClause)
		}
	}
	queryBuilder.WriteString(" ORDER BY ")
	queryBuilder.WriteString(strings.Join(orderClauses, ", "))

	queryBuilder.WriteString(fmt.Sprintf(" LIMIT %d ", pgn.Limit+1))

	return queryBuilder.String(), args, nil
}
