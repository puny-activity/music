package httpcontroller

import (
	httpapi "github.com/puny-activity/music/api/http"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/werr"
	"net/http"
)

func (c *Controller) GetGenres(w http.ResponseWriter, r *http.Request) error {
	version := r.Header.Get(httpapi.APIVersion)
	switch version {
	case "1":
		return c.getGenresV1(w, r)
	default:
		return errs.InvalidAPIVersion
	}
}

type getGenresV1ResponseGenreItemCoverItem struct {
	ID string `json:"id"`
}

type getGenresV1ResponseGenreItem struct {
	ID        string                                  `json:"id"`
	Name      string                                  `json:"name"`
	SongCount int                                     `json:"song_count"`
	Covers    []getGenresV1ResponseGenreItemCoverItem `json:"covers"`
}

type getGenresV1Response struct {
	Genres     []getGenresV1ResponseGenreItem `json:"genres"`
	Pagination Pagination                     `json:"pagination"`
}

func (c *Controller) getGenresV1(w http.ResponseWriter, r *http.Request) error {
	pgn, err := ExtractCursorPagination(r)
	if err != nil {
		return werr.WrapSE("failed to extract cursor pagination", err)
	}

	if pgn.Parameters != nil {
		parametersAfter := make([]pagination.Parameter, 0)
		parametersBefore := pgn.Parameters
		for _, parameterBefore := range parametersBefore {
			switch parameterBefore.FieldName {
			case "name":
				parametersAfter = append(parametersAfter, pagination.Parameter{
					FieldName: genre.PaginationName,
					SortOrder: parameterBefore.SortOrder,
				})
			case "song_count":
				parametersAfter = append(parametersAfter, pagination.Parameter{
					FieldName: genre.PaginationSongCount,
					SortOrder: parameterBefore.SortOrder,
				})
			default:
				return errs.UnknownSortParameter
			}
		}
		pgn.Parameters = parametersAfter
	}

	genres, cursorPair, err := c.app.GenreUseCase.GetMany(r.Context(), pgn)
	if err != nil {
		return werr.WrapSE("failed to get genres", err)
	}

	genresResponse := make([]getGenresV1ResponseGenreItem, len(genres))
	for i, genre := range genres {
		covers := make([]getGenresV1ResponseGenreItemCoverItem, len(genre.CoversIDs))
		for j, coverID := range genre.CoversIDs {
			covers[j] = getGenresV1ResponseGenreItemCoverItem{
				ID: coverID.String(),
			}
		}
		genresResponse[i] = getGenresV1ResponseGenreItem{
			ID:        genre.ID.String(),
			Name:      genre.Name,
			SongCount: genre.SongCount,
			Covers:    covers,
		}
	}

	var prevCursor *string = nil
	if cursorPair.PrevCursor != nil {
		prevCursorLocal, err := cursorPair.PrevCursor.Encode()
		if err != nil {
			return werr.WrapSE("failed to get prev cursor", err)
		}
		prevCursor = &prevCursorLocal
	}

	var nextCursor *string = nil
	if cursorPair.NextCursor != nil {
		nextCursorLocal, err := cursorPair.NextCursor.Encode()
		if err != nil {
			return werr.WrapSE("failed to get next cursor", err)
		}
		nextCursor = &nextCursorLocal
	}

	response := getGenresV1Response{
		Genres: genresResponse,
		Pagination: Pagination{
			PrevCursor: prevCursor,
			NextCursor: nextCursor,
		},
	}

	return c.writer.Write(w, http.StatusOK, response)
}
