package albumrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
	"github.com/puny-activity/music/pkg/werr"
	"strings"
)

type getManyDTO struct {
	AlbumID    string `db:"album_id"`
	AlbumTitle string `db:"album_title"`
	SongCount  int    `db:"song_count"`
	Covers     []byte `db:"covers"`
}

func (dto getManyDTO) ToAlbum() (album.Album, error) {
	albumID, err := album.ParseID(dto.AlbumID)
	if err != nil {
		return album.Album{}, werr.WrapSE("failed to parse album id", err)
	}

	coverIDsStr := make([]string, 0)
	err = json.Unmarshal(dto.Covers, &coverIDsStr)
	if err != nil {
		return album.Album{}, werr.WrapSE("failed to parse covers", err)
	}
	coverIDs := make([]cover.ID, len(coverIDsStr))
	for i := range coverIDsStr {
		coverIDs[i], err = cover.ParseID(coverIDsStr[i])
		if err != nil {
			return album.Album{}, werr.WrapSE("failed to parse cover id", err)
		}
	}

	return album.Album{
		Base: album.Base{
			ID:    &albumID,
			Title: dto.AlbumTitle,
		},
		SongCount: dto.SongCount,
		CoversIDs: coverIDs,
	}, nil
}

const (
	getManyPaginationTitle     = "ra.album_title"
	getManyPaginationSongCount = "ra.song_count"
)

func getManyPgnParameterConvert(before []pagination.Parameter) ([]pagination.Parameter, error) {
	defaultPaginators := []pagination.Parameter{
		{
			FieldName: getManyPaginationTitle,
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
		case album.PaginationTitle:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationTitle,
				SortOrder: param.SortOrder,
			})
		case album.PaginationSongCount:
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

func (r *Repository) GetMany(ctx context.Context, pgn pagination.CursorPagination) ([]album.Album, pagination.CursorPair, error) {
	return r.getMany(ctx, r.db, pgn)
}

func (r *Repository) GetManyTx(ctx context.Context, tx *sqlx.Tx, pgn pagination.CursorPagination) ([]album.Album, pagination.CursorPair, error) {
	return r.getMany(ctx, tx, pgn)
}

func (r *Repository) getMany(ctx context.Context, queryer queryer.Queryer, pgn pagination.CursorPagination) ([]album.Album, pagination.CursorPair, error) {
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
		return nil, pagination.CursorPair{}, werr.WrapSE("failed to get album dtos", err)
	}
	if len(dtoList) == 0 {
		return make([]album.Album, 0), pagination.CursorPair{}, nil
	}

	var albums []album.Album
	for i, dto := range dtoList {
		if (len(dtoList) > pgn.Limit) && (i == len(dtoList)-1) {
			continue
		}
		oneAlbum, err := dto.ToAlbum()
		if err != nil {
			return nil, pagination.CursorPair{}, err
		}
		albums = append(albums, oneAlbum)
	}
	if pgn.Cursor != nil && pgn.Cursor.Type() == pagination.CursorPrev {
		albums = util.Reverse(albums)
	}

	var prevCursor *pagination.Cursor = nil
	var nextCursor *pagination.Cursor = nil

	if pgn.Parameters != nil {
		if len(dtoList) == pgn.Limit+1 {
			lastAlbum := albums[len(albums)-1]
			nextCursor = util.ToPointer(pagination.NewCursor(pagination.CursorNext))
			nextCursor = util.ToPointer(getManyIncludeParameters(*nextCursor, pgn.Parameters, lastAlbum))
		}
	} else if pgn.Cursor != nil {
		if pgn.Cursor.Type() == pagination.CursorNext {
			if len(albums) > 0 {
				firstAlbum := albums[0]
				prevCursor = util.ToPointer(pagination.NewCursor(pagination.CursorPrev))
				prevCursor = util.ToPointer(getManyIncludeCursor(*prevCursor, pgn.Cursor, firstAlbum))
			}
			if len(dtoList) == pgn.Limit+1 {
				lastAlbum := albums[len(albums)-1]
				nextCursor = util.ToPointer(pagination.NewCursor(pagination.CursorNext))
				nextCursor = util.ToPointer(getManyIncludeCursor(*nextCursor, pgn.Cursor, lastAlbum))
			}
		} else if pgn.Cursor.Type() == pagination.CursorPrev {
			if len(dtoList) == pgn.Limit+1 {
				firstAlbum := albums[0]
				prevCursor = util.ToPointer(pagination.NewCursor(pagination.CursorPrev))
				prevCursor = util.ToPointer(getManyIncludeCursor(*prevCursor, pgn.Cursor, firstAlbum))
			}
			if len(albums) > 0 {
				lastAlbum := albums[len(albums)-1]
				nextCursor = util.ToPointer(pagination.NewCursor(pagination.CursorNext))
				nextCursor = util.ToPointer(getManyIncludeCursor(*nextCursor, pgn.Cursor, lastAlbum))
			}
		}
	}

	cursorPair := pagination.CursorPair{
		PrevCursor: prevCursor,
		NextCursor: nextCursor,
	}

	return albums, cursorPair, nil
}

func getManyIncludeParameters(cursor pagination.Cursor, parameters []pagination.Parameter, targetAlbum album.Album) pagination.Cursor {
	for _, parameter := range parameters {
		var value any
		switch parameter.FieldName {
		case getManyPaginationTitle:
			value = targetAlbum.Title
		case getManyPaginationSongCount:
			value = targetAlbum.SongCount
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

func getManyIncludeCursor(cursor pagination.Cursor, cursor2 *pagination.Cursor, targetAlbum album.Album) pagination.Cursor {
	for _, parameter := range cursor2.GetAll() {
		var value any
		switch parameter.FieldName {
		case getManyPaginationTitle:
			value = targetAlbum.Title
		case getManyPaginationSongCount:
			value = targetAlbum.SongCount
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
WITH resp_covers AS (SELECT a.id        album_id,
                            c.id        cover_id,
                            COUNT(c.id) song_count
                     FROM albums a
                              JOIN songs s ON a.id = s.album_id
                              JOIN covers c ON s.cover_id = c.id
                     GROUP BY a.id, c.id
                     ORDER BY album_id, song_count DESC),
     resp_albums AS (SELECT a.id        AS album_id,
                            a.title     AS album_title,
                            COUNT(s.id) AS song_count
                     FROM albums a
                              JOIN songs s ON a.id = s.album_id
                     GROUP BY a.id, a.title)
SELECT ra.album_id,
       ra.album_title,
       ra.song_count,
       TO_JSON(ARRAY_REMOVE(ARRAY_AGG(rc.cover_id), NULL)) AS covers
FROM resp_albums ra
         LEFT JOIN resp_covers rc ON ra.album_id = rc.album_id
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

	queryBuilder.WriteString(" GROUP BY ra.album_id, ra.album_title, ra.song_count ")

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
