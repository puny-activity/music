package httpcontroller

import (
	httpapi "github.com/puny-activity/music/api/http"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/werr"
	"net/http"
)

func (c *Controller) GetArtists(w http.ResponseWriter, r *http.Request) error {
	version := r.Header.Get(httpapi.APIVersion)
	switch version {
	case "1":
		return c.getArtistsV1(w, r)
	default:
		return errs.InvalidAPIVersion
	}
}

type getArtistsV1ResponseArtistItemCoverItem struct {
	ID string `json:"id"`
}

type getArtistsV1ResponseArtistItem struct {
	ID        string                                    `json:"id"`
	Name      string                                    `json:"name"`
	SongCount int                                       `json:"song_count"`
	Covers    []getArtistsV1ResponseArtistItemCoverItem `json:"covers"`
}

type getArtistsV1Response struct {
	Artists    []getArtistsV1ResponseArtistItem `json:"artists"`
	Pagination Pagination                       `json:"pagination"`
}

func (c *Controller) getArtistsV1(w http.ResponseWriter, r *http.Request) error {
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
					FieldName: artist.PaginationName,
					SortOrder: parameterBefore.SortOrder,
				})
			default:
				return errs.UnknownSortParameter
			}
		}
		pgn.Parameters = parametersAfter
	}

	artists, cursorPair, err := c.app.ArtistUseCase.GetMany(r.Context(), pgn)
	if err != nil {
		return werr.WrapSE("failed to get artists", err)
	}

	artistsResponse := make([]getArtistsV1ResponseArtistItem, len(artists))
	for i, artist := range artists {
		covers := make([]getArtistsV1ResponseArtistItemCoverItem, len(artist.CoversIDs))
		for j, coverID := range artist.CoversIDs {
			covers[j] = getArtistsV1ResponseArtistItemCoverItem{
				ID: coverID.String(),
			}
		}
		artistsResponse[i] = getArtistsV1ResponseArtistItem{
			ID:        artist.ID.String(),
			Name:      artist.Name,
			SongCount: artist.SongCount,
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

	response := getArtistsV1Response{
		Artists: artistsResponse,
		Pagination: Pagination{
			PrevCursor: prevCursor,
			NextCursor: nextCursor,
		},
	}

	return c.writer.Write(w, http.StatusOK, response)
}
