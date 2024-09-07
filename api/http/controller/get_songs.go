package httpcontroller

import (
	httpapi "github.com/puny-activity/music/api/http"
	"github.com/puny-activity/music/internal/entity/song"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/filters"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/util"
	"github.com/puny-activity/music/pkg/werr"
	"net/http"
)

func (c *Controller) GetSongs(w http.ResponseWriter, r *http.Request) error {
	version := r.Header.Get(httpapi.APIVersion)
	switch version {
	case "1":
		return c.getSongsV1(w, r)
	default:
		return errs.InvalidAPIVersion
	}
}

type getSongsV1ResponseSongItemGenre struct {
	ID   string `json:"id"`
	Name string `json:"title"`
}

type getSongsV1ResponseSongItemAlbum struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type getSongsV1ResponseSongItemArtist struct {
	ID   string `json:"id"`
	Name string `json:"title"`
}

type getSongsV1ResponseSongItem struct {
	ID           string                           `json:"id"`
	Title        string                           `json:"title"`
	Duration     int64                            `json:"duration"`
	CoverID      *string                          `json:"cover_id,omitempty"`
	Genre        getSongsV1ResponseSongItemGenre  `json:"genre"`
	Album        getSongsV1ResponseSongItemAlbum  `json:"album"`
	Artists      getSongsV1ResponseSongItemArtist `json:"artist"`
	Year         *int                             `json:"year,omitempty"`
	Number       *int                             `json:"number,omitempty"`
	Comment      *string                          `json:"comment,omitempty"`
	Channels     int                              `json:"channels"`
	BitrateKbps  int                              `json:"bitrate_kbps"`
	SampleRateHz int                              `json:"sample_rate_hz"`
}

type getSongsV1Response struct {
	Songs      []getSongsV1ResponseSongItem `json:"songs"`
	Pagination Pagination                   `json:"pagination"`
}

func (c *Controller) getSongsV1(w http.ResponseWriter, r *http.Request) error {
	pgn, err := ExtractCursorPagination(r)
	if err != nil {
		return werr.WrapSE("failed to extract cursor pagination", err)
	}

	search := r.URL.Query().Get("search")

	fltr := filters.NewFilter()
	genreIDsStr := r.URL.Query()["genre_id"]
	if len(genreIDsStr) > 0 {
		andFltr := filters.NewAndFilter()
		for _, genreIDStr := range genreIDsStr {
			genreID, err := genre.ParseID(genreIDStr)
			if err != nil {
				return werr.WrapES(errs.InvalidGenreParameter, err.Error())
			}
			orFltr := filters.NewOrFilter(song.FilterGenre, genreID.String())
			andFltr = andFltr.Add(orFltr)
		}
		fltr = fltr.Add(andFltr)
	}
	albumIDsStr := r.URL.Query()["album_id"]
	if len(albumIDsStr) > 0 {
		andFltr := filters.NewAndFilter()
		for _, albumIDStr := range albumIDsStr {
			albumID, err := album.ParseID(albumIDStr)
			if err != nil {
				return werr.WrapES(errs.InvalidAlbumParameter, err.Error())
			}
			orFltr := filters.NewOrFilter(song.FilterAlbum, albumID.String())
			andFltr = andFltr.Add(orFltr)
		}
		fltr = fltr.Add(andFltr)
	}
	artistIDsStr := r.URL.Query()["artist_id"]
	if len(artistIDsStr) > 0 {
		andFltr := filters.NewAndFilter()
		for _, artistIDStr := range artistIDsStr {
			artistID, err := artist.ParseID(artistIDStr)
			if err != nil {
				return werr.WrapES(errs.InvalidArtistParameter, err.Error())
			}
			orFltr := filters.NewOrFilter(song.FilterArtist, artistID.String())
			andFltr = andFltr.Add(orFltr)
		}
		fltr = fltr.Add(andFltr)
	}

	if pgn.Parameters != nil {
		parametersAfter := make([]pagination.Parameter, 0)
		parametersBefore := pgn.Parameters
		for _, parameterBefore := range parametersBefore {
			switch parameterBefore.FieldName {
			case "number":
				parametersAfter = append(parametersAfter, pagination.Parameter{
					FieldName: song.PaginationNumber,
					SortOrder: parameterBefore.SortOrder,
				})
			case "title":
				parametersAfter = append(parametersAfter, pagination.Parameter{
					FieldName: song.PaginationTitle,
					SortOrder: parameterBefore.SortOrder,
				})
			case "year":
				parametersAfter = append(parametersAfter, pagination.Parameter{
					FieldName: song.PaginationYear,
					SortOrder: parameterBefore.SortOrder,
				})
			case "duration":
				parametersAfter = append(parametersAfter, pagination.Parameter{
					FieldName: song.PaginationDuration,
					SortOrder: parameterBefore.SortOrder,
				})
			case "bitrate":
				parametersAfter = append(parametersAfter, pagination.Parameter{
					FieldName: song.PaginationBitrate,
					SortOrder: parameterBefore.SortOrder,
				})
			case "sample_rate":
				parametersAfter = append(parametersAfter, pagination.Parameter{
					FieldName: song.PaginationSampleRate,
					SortOrder: parameterBefore.SortOrder,
				})
			case "genre_name":
				parametersAfter = append(parametersAfter, pagination.Parameter{
					FieldName: genre.PaginationName,
					SortOrder: parameterBefore.SortOrder,
				})
			case "album_title":
				parametersAfter = append(parametersAfter, pagination.Parameter{
					FieldName: album.PaginationTitle,
					SortOrder: parameterBefore.SortOrder,
				})
			case "artist_name":
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

	songs, cursorPair, err := c.app.SongUseCase.Get(r.Context(), search, fltr, pgn)
	if err != nil {
		return err
	}

	songsResponse := make([]getSongsV1ResponseSongItem, len(songs))
	for i, song := range songs {
		var coverID *string = nil
		if song.Cover != nil {
			coverID = util.ToPointer(song.Cover.ID.String())
		}
		songsResponse[i] = getSongsV1ResponseSongItem{
			ID:       song.ID.String(),
			Title:    song.Title,
			Duration: song.Duration.Nanoseconds(),
			CoverID:  coverID,
			Genre: getSongsV1ResponseSongItemGenre{
				ID:   song.Genre.ID.String(),
				Name: song.Genre.Name,
			},
			Album: getSongsV1ResponseSongItemAlbum{
				ID:    song.Album.ID.String(),
				Title: song.Album.Title,
			},
			Artists: getSongsV1ResponseSongItemArtist{
				ID:   song.Artist.ID.String(),
				Name: song.Artist.Name,
			},
			Year:         song.Year,
			Number:       song.Number,
			Comment:      song.Comment,
			Channels:     song.Channels,
			BitrateKbps:  song.BitrateKbps,
			SampleRateHz: song.SampleRateHz,
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

	response := getSongsV1Response{
		Songs: songsResponse,
		Pagination: Pagination{
			PrevCursor: prevCursor,
			NextCursor: nextCursor,
		},
	}

	return c.writer.Write(w, http.StatusOK, response)
}
