package songrepo

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/song"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/filters"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
	"github.com/puny-activity/music/pkg/werr"
	"strings"
	"time"
)

type getManyDTO struct {
	ID           string  `db:"id"`
	FileID       *string `db:"file_id"`
	Title        string  `db:"title"`
	Duration     int64   `db:"duration"`
	CoverID      *string `db:"cover_id"`
	CoverWidth   *int    `db:"cover_width"`
	CoverHeight  *int    `db:"cover_height"`
	CoverFileID  *string `db:"cover_file_id"`
	GenreID      string  `db:"genre_id"`
	GenreName    string  `db:"genre_name"`
	AlbumID      string  `db:"album_id"`
	AlbumTitle   string  `db:"album_title"`
	ArtistID     string  `db:"artist_id"`
	ArtistName   string  `db:"artist_name"`
	Year         *int    `db:"year"`
	Number       *int    `db:"number"`
	Comment      *string `db:"comment"`
	Channels     int     `db:"channels"`
	BitrateKbps  int     `db:"bitrate_kbps"`
	SampleRateHz int     `db:"sample_rate_hz"`
	MD5          string  `db:"md5"`
}

func (dto getManyDTO) ToSong() (song.Song, error) {
	songID, err := song.ParseID(dto.ID)
	if err != nil {
		return song.Song{}, werr.WrapSE("failed to parse song id", err)
	}

	var fileID *remotefile.ID = nil
	if dto.FileID != nil {
		fileIDLocal, err := remotefile.ParseID(*dto.FileID)
		if err != nil {
			return song.Song{}, werr.WrapSE("failed to parse file id", err)
		}
		fileID = &fileIDLocal
	}

	var songsCover *cover.Cover = nil
	if dto.CoverID != nil {
		songsCoverID, err := cover.ParseID(*dto.CoverID)
		if err != nil {
			return song.Song{}, werr.WrapSE("failed to parse cover id", err)
		}
		coverFileID, err := remotefile.ParseID(*dto.CoverFileID)
		if err != nil {
			return song.Song{}, werr.WrapSE("failed to parse cover file id", err)
		}
		songsCover = &cover.Cover{
			ID:     &songsCoverID,
			Width:  *dto.CoverWidth,
			Height: *dto.CoverHeight,
			FileID: coverFileID,
		}
	}

	genreID, err := genre.ParseID(dto.GenreID)
	if err != nil {
		return song.Song{}, werr.WrapSE("failed to parse genre id", err)
	}
	albumID, err := album.ParseID(dto.AlbumID)
	if err != nil {
		return song.Song{}, werr.WrapSE("failed to parse album id", err)
	}
	artistID, err := artist.ParseID(dto.ArtistID)
	if err != nil {
		return song.Song{}, werr.WrapSE("failed to parse album id", err)
	}

	return song.Song{
		ID:       &songID,
		FileID:   fileID,
		Title:    dto.Title,
		Duration: time.Duration(dto.Duration),
		Cover:    songsCover,
		Genre: genre.Base{
			ID:   &genreID,
			Name: dto.GenreName,
		},
		Album: album.Base{
			ID:    &albumID,
			Title: dto.AlbumTitle,
		},
		Artist: artist.Base{
			ID:   &artistID,
			Name: dto.ArtistName,
		},
		Year:         dto.Year,
		Number:       dto.Number,
		Comment:      dto.Comment,
		Channels:     dto.Channels,
		BitrateKbps:  dto.BitrateKbps,
		SampleRateHz: dto.SampleRateHz,
		MD5:          dto.MD5,
	}, nil
}

var a = `
SELECT s.id             AS id,
       s.file_id        AS file_id,
       s.title          AS title,
       s.duration_ns    AS duration,
       c.id             AS cover_id,
       c.width          AS cover_width,
       c.height         AS cover_height,
       c.file_id        AS cover_file_id,
       g.id             AS genre_id,
       g.name           AS genre_name,
       al.id            AS album_id,
       al.title         AS album_title,
       ar.id            AS artist_id,
       ar.name          AS artist_name,
       s.year           AS year,
       s.number         AS number,
       s.comment        AS comment,
       s.channels       AS channels,
       s.bitrate_kbps   AS bitrate_kbps,
       s.sample_rate_hz AS sample_rate_hz,
       s.md5            AS md5
FROM songs s
         LEFT JOIN covers c ON s.cover_id = c.id
         JOIN genres g ON s.genre_id = g.id
         JOIN albums al ON s.album_id = al.id
         JOIN artists ar ON s.artist_id = ar.id
`

const (
	getManyPaginationNumber     = "s.number"
	getManyPaginationTitle      = "s.title"
	getManyPaginationYear       = "s.year"
	getManyPaginationDuration   = "s.duration_ns"
	getManyPaginationBitrate    = "s.bitrate_kbps"
	getManyPaginationSampleRate = "s.sample_rate_hz"
	getManyPaginationGenreName  = "g.name"
	getManyPaginationAlbumTitle = "al.title"
	getManyPaginationArtistName = "ar.name"
	getManyFilterGenre          = "g.id"
	getManyFilterAlbum          = "al.id"
	getManyFilterArtist         = "ar.id"
)

func getManyPgnParameterConvert(before []pagination.Parameter) ([]pagination.Parameter, error) {
	defaultPaginators := []pagination.Parameter{
		{
			FieldName: getManyPaginationNumber,
			SortOrder: pagination.ASC,
		},
		{
			FieldName: getManyPaginationTitle,
			SortOrder: pagination.ASC,
		},
		{
			FieldName: getManyPaginationGenreName,
			SortOrder: pagination.ASC,
		},
		{
			FieldName: getManyPaginationAlbumTitle,
			SortOrder: pagination.ASC,
		},
		{
			FieldName: getManyPaginationArtistName,
			SortOrder: pagination.ASC,
		},
		{
			FieldName: getManyPaginationYear,
			SortOrder: pagination.DESC,
		},
		{
			FieldName: getManyPaginationDuration,
			SortOrder: pagination.DESC,
		},
		{
			FieldName: getManyPaginationBitrate,
			SortOrder: pagination.DESC,
		},
		{
			FieldName: getManyPaginationSampleRate,
			SortOrder: pagination.DESC,
		},
	}
	appliedPaginators := make(map[string]struct{})

	after := make([]pagination.Parameter, 0)
	for _, param := range before {
		switch param.FieldName {
		case song.PaginationNumber:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationNumber,
				SortOrder: param.SortOrder,
			})
		case song.PaginationTitle:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationTitle,
				SortOrder: param.SortOrder,
			})
		case genre.PaginationName:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationGenreName,
				SortOrder: param.SortOrder,
			})
		case album.PaginationTitle:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationAlbumTitle,
				SortOrder: param.SortOrder,
			})
		case artist.PaginationName:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationArtistName,
				SortOrder: param.SortOrder,
			})
		case song.PaginationYear:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationYear,
				SortOrder: param.SortOrder,
			})
		case song.PaginationDuration:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationDuration,
				SortOrder: param.SortOrder,
			})
		case song.PaginationBitrate:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationBitrate,
				SortOrder: param.SortOrder,
			})
		case song.PaginationSampleRate:
			after = append(after, pagination.Parameter{
				FieldName: getManyPaginationSampleRate,
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

func getManyFltrParameterConvert(before filters.Filter) (filters.Filter, error) {
	after := filters.NewFilter()
	for _, andFilter := range before.AndFilters {
		afterAndFilter := filters.NewAndFilter()
		for _, orFilter := range andFilter.OrFilter {
			var newOrFilter filters.OrFilter
			switch orFilter.Key {
			case song.FilterGenre:
				newOrFilter = filters.NewOrFilter(getManyFilterGenre, orFilter.Value)
			case song.FilterAlbum:
				newOrFilter = filters.NewOrFilter(getManyFilterAlbum, orFilter.Value)
			case song.FilterArtist:
				newOrFilter = filters.NewOrFilter(getManyFilterArtist, orFilter.Value)
			default:
				return filters.Filter{}, fmt.Errorf("unknown filter key: %s", orFilter.Key)
			}
			afterAndFilter = afterAndFilter.Add(newOrFilter)
		}
		after = after.Add(afterAndFilter)
	}
	return after, nil
}

func (r *Repository) GetMany(ctx context.Context, search string, fltr filters.Filter, pgn pagination.CursorPagination) ([]song.Song, pagination.CursorPair, error) {
	return r.getMany(ctx, r.db, search, fltr, pgn)
}

func (r *Repository) GetManyTx(ctx context.Context, tx *sqlx.Tx, search string, fltr filters.Filter, pgn pagination.CursorPagination) ([]song.Song, pagination.CursorPair, error) {
	return r.getMany(ctx, tx, search, fltr, pgn)
}

func (r *Repository) getMany(ctx context.Context, queryer queryer.Queryer, search string, fltr filters.Filter, pgn pagination.CursorPagination) ([]song.Song, pagination.CursorPair, error) {
	fltr, err := getManyFltrParameterConvert(fltr)
	if err != nil {
		return nil, pagination.CursorPair{}, werr.WrapSE("failed to convert filter parameters", err)
	}

	if pgn.Cursor == nil {
		var err error
		pgn.Parameters, err = getManyPgnParameterConvert(pgn.Parameters)
		if err != nil {
			return nil, pagination.CursorPair{}, werr.WrapSE("failed to convert parameters", err)
		}
	}

	query, args, err := getManyBuildQuery(search, fltr, pgn)
	if err != nil {
		return nil, pagination.CursorPair{}, werr.WrapSE("failed to build query", err)
	}

	var dtoList []getManyDTO
	err = queryer.SelectContext(ctx, &dtoList, query, args...)
	if err != nil {
		return nil, pagination.CursorPair{}, werr.WrapSE("failed to get song dtos", err)
	}
	if len(dtoList) == 0 {
		return make([]song.Song, 0), pagination.CursorPair{}, nil
	}

	var songs []song.Song
	for i, dto := range dtoList {
		if (len(dtoList) > pgn.Limit) && (i == len(dtoList)-1) {
			continue
		}
		oneAlbum, err := dto.ToSong()
		if err != nil {
			return nil, pagination.CursorPair{}, err
		}
		songs = append(songs, oneAlbum)
	}
	if pgn.Cursor != nil && pgn.Cursor.Type() == pagination.CursorPrev {
		songs = util.Reverse(songs)
	}

	var prevCursor *pagination.Cursor = nil
	var nextCursor *pagination.Cursor = nil

	if pgn.Parameters != nil {
		if len(dtoList) == pgn.Limit+1 {
			lastAlbum := songs[len(songs)-1]
			nextCursor = util.ToPointer(pagination.NewCursor(pagination.CursorNext))
			nextCursor = util.ToPointer(getManyIncludeParameters(*nextCursor, pgn.Parameters, lastAlbum))
		}
	} else if pgn.Cursor != nil {
		if pgn.Cursor.Type() == pagination.CursorNext {
			if len(songs) > 0 {
				firstAlbum := songs[0]
				prevCursor = util.ToPointer(pagination.NewCursor(pagination.CursorPrev))
				prevCursor = util.ToPointer(getManyIncludeCursor(*prevCursor, pgn.Cursor, firstAlbum))
			}
			if len(dtoList) == pgn.Limit+1 {
				lastAlbum := songs[len(songs)-1]
				nextCursor = util.ToPointer(pagination.NewCursor(pagination.CursorNext))
				nextCursor = util.ToPointer(getManyIncludeCursor(*nextCursor, pgn.Cursor, lastAlbum))
			}
		} else if pgn.Cursor.Type() == pagination.CursorPrev {
			if len(dtoList) == pgn.Limit+1 {
				firstAlbum := songs[0]
				prevCursor = util.ToPointer(pagination.NewCursor(pagination.CursorPrev))
				prevCursor = util.ToPointer(getManyIncludeCursor(*prevCursor, pgn.Cursor, firstAlbum))
			}
			if len(songs) > 0 {
				lastAlbum := songs[len(songs)-1]
				nextCursor = util.ToPointer(pagination.NewCursor(pagination.CursorNext))
				nextCursor = util.ToPointer(getManyIncludeCursor(*nextCursor, pgn.Cursor, lastAlbum))
			}
		}
	}

	cursorPair := pagination.CursorPair{
		PrevCursor: prevCursor,
		NextCursor: nextCursor,
	}

	return songs, cursorPair, nil
}

func getManyIncludeParameters(cursor pagination.Cursor, parameters []pagination.Parameter, targetSong song.Song) pagination.Cursor {
	for _, parameter := range parameters {
		var value any
		switch parameter.FieldName {
		case getManyPaginationNumber:
			value = targetSong.Number
		case getManyPaginationTitle:
			value = targetSong.Title
		case getManyPaginationYear:
			value = targetSong.Year
		case getManyPaginationDuration:
			value = targetSong.Duration
		case getManyPaginationBitrate:
			value = targetSong.BitrateKbps
		case getManyPaginationSampleRate:
			value = targetSong.SampleRateHz
		case getManyPaginationGenreName:
			value = targetSong.Genre.Name
		case getManyPaginationAlbumTitle:
			value = targetSong.Album.Title
		case getManyPaginationArtistName:
			value = targetSong.Artist.Name
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

func getManyIncludeCursor(cursor pagination.Cursor, cursor2 *pagination.Cursor, targetSong song.Song) pagination.Cursor {
	for _, parameter := range cursor2.GetAll() {
		var value any
		switch parameter.FieldName {
		case getManyPaginationNumber:
			value = targetSong.Number
		case getManyPaginationTitle:
			value = targetSong.Title
		case getManyPaginationYear:
			value = targetSong.Year
		case getManyPaginationDuration:
			value = targetSong.Duration
		case getManyPaginationBitrate:
			value = targetSong.BitrateKbps
		case getManyPaginationSampleRate:
			value = targetSong.SampleRateHz
		case getManyPaginationGenreName:
			value = targetSong.Genre.Name
		case getManyPaginationAlbumTitle:
			value = targetSong.Album.Title
		case getManyPaginationArtistName:
			value = targetSong.Artist.Name
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

func getManyBuildQuery(search string, fltr filters.Filter, pgn pagination.CursorPagination) (string, []any, error) {
	var queryBuilder strings.Builder
	var args []any

	queryBuilder.WriteString(`
SELECT s.id             AS id,
       s.file_id        AS file_id,
       s.title          AS title,
       s.duration_ns    AS duration,
       c.id             AS cover_id,
       c.width          AS cover_width,
       c.height         AS cover_height,
       c.file_id        AS cover_file_id,
       g.id             AS genre_id,
       g.name           AS genre_name,
       al.id            AS album_id,
       al.title         AS album_title,
       ar.id            AS artist_id,
       ar.name          AS artist_name,
       s.year           AS year,
       s.number         AS number,
       s.comment        AS comment,
       s.channels       AS channels,
       s.bitrate_kbps   AS bitrate_kbps,
       s.sample_rate_hz AS sample_rate_hz,
       s.md5            AS md5
FROM songs s
         LEFT JOIN covers c ON s.cover_id = c.id
         JOIN genres g ON s.genre_id = g.id
         JOIN albums al ON s.album_id = al.id
         JOIN artists ar ON s.artist_id = ar.id
`)
	queryBuilder.WriteString(fmt.Sprintf(" WHERE (s.title LIKE '%%%s%%') ", search))

	if !fltr.IsEmpty() {
		queryBuilder.WriteString(fmt.Sprintf(" AND %s ", fltr.BuildWhere()))
	}

	if pgn.Cursor != nil {
		whereClause, err := pgn.Cursor.BuildWhere()
		if err != nil {
			return "", nil, werr.WrapSE("failed to build where query", err)
		}
		if whereClause != "" {
			queryBuilder.WriteString(fmt.Sprintf(" AND %s ", whereClause))
		}
	}

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
