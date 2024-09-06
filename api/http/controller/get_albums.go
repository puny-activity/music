package httpcontroller

import (
	httpapi "github.com/puny-activity/music/api/http"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/werr"
	"net/http"
)

func (c *Controller) GetAlbums(w http.ResponseWriter, r *http.Request) error {
	version := r.Header.Get(httpapi.APIVersion)
	switch version {
	case "1":
		return c.getAlbumsV1(w, r)
	default:
		return errs.InvalidAPIVersion
	}
}

type getAlbumsV1ResponseAlbumItemCoverItem struct {
	ID string `json:"id"`
}

type getAlbumsV1ResponseAlbumItem struct {
	ID        string                                  `json:"id"`
	Title     string                                  `json:"title"`
	SongCount int                                     `json:"song_count"`
	Covers    []getAlbumsV1ResponseAlbumItemCoverItem `json:"covers"`
}

type getAlbumsV1Response struct {
	Albums     []getAlbumsV1ResponseAlbumItem `json:"albums"`
	Pagination Pagination                     `json:"pagination"`
}

func (c *Controller) getAlbumsV1(w http.ResponseWriter, r *http.Request) error {
	pgn, err := ExtractCursorPagination(r)
	if err != nil {
		return werr.WrapSE("failed to extract cursor pagination", err)
	}

	if pgn.Parameters != nil {
		parametersAfter := make([]pagination.Parameter, 0)
		parametersBefore := pgn.Parameters
		for _, parameterBefore := range parametersBefore {
			switch parameterBefore.FieldName {
			case "title":
				parametersAfter = append(parametersAfter, pagination.Parameter{
					FieldName: album.PaginationTitle,
					SortOrder: parameterBefore.SortOrder,
				})
			default:
				return errs.UnknownSortParameter
			}
		}
		pgn.Parameters = parametersAfter
	}

	albums, cursorPair, err := c.app.AlbumUseCase.GetMany(r.Context(), pgn)
	if err != nil {
		return werr.WrapSE("failed to get albums", err)
	}

	albumsResponse := make([]getAlbumsV1ResponseAlbumItem, len(albums))
	for i, album := range albums {
		covers := make([]getAlbumsV1ResponseAlbumItemCoverItem, len(album.CoversIDs))
		for j, coverID := range album.CoversIDs {
			covers[j] = getAlbumsV1ResponseAlbumItemCoverItem{
				ID: coverID.String(),
			}
		}
		albumsResponse[i] = getAlbumsV1ResponseAlbumItem{
			ID:        album.ID.String(),
			Title:     album.Title,
			SongCount: album.SongCount,
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

	response := getAlbumsV1Response{
		Albums: albumsResponse,
		Pagination: Pagination{
			PrevCursor: prevCursor,
			NextCursor: nextCursor,
		},
	}

	return c.writer.Write(w, http.StatusOK, response)
}
