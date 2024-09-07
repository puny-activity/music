package httpcontroller

import (
	"github.com/go-chi/chi/v5"
	httpapi "github.com/puny-activity/music/api/http"
	"github.com/puny-activity/music/internal/entity/song"
	"github.com/puny-activity/music/internal/errs"
	"net/http"
)

func (c *Controller) GetSongsURL(w http.ResponseWriter, r *http.Request) error {
	version := r.Header.Get(httpapi.APIVersion)
	switch version {
	case "1":
		return c.getSongsURLV1(w, r)
	default:
		return errs.InvalidAPIVersion
	}
}

type getSongsURLV1Response struct {
	URL string `json:"url"`
}

func (c *Controller) getSongsURLV1(w http.ResponseWriter, r *http.Request) error {
	songIDStr := chi.URLParam(r, "song_id")
	if songIDStr == "" {
		return errs.InvalidSongParameter
	}
	songID, err := song.ParseID(songIDStr)
	if err != nil {
		return errs.InvalidSongParameter
	}

	url, err := c.app.SongUseCase.GetURL(r.Context(), songID)
	if err != nil {
		return err
	}

	return c.writer.Write(w, http.StatusOK, &getSongsURLV1Response{
		URL: url,
	})
}
