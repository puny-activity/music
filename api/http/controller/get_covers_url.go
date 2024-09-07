package httpcontroller

import (
	"github.com/go-chi/chi/v5"
	httpapi "github.com/puny-activity/music/api/http"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/internal/errs"
	"net/http"
)

type getCoversURLV1Response struct {
	URL string `json:"url"`
}

func (c *Controller) GetCoversURL(w http.ResponseWriter, r *http.Request) error {
	version := r.Header.Get(httpapi.APIVersion)
	switch version {
	case "1":
		return c.getCoversURLV1(w, r)
	default:
		return errs.InvalidAPIVersion
	}
}

func (c *Controller) getCoversURLV1(w http.ResponseWriter, r *http.Request) error {
	coverIDStr := chi.URLParam(r, "cover_id")
	if coverIDStr == "" {
		return errs.InvalidCoverParameter
	}
	coverID, err := cover.ParseID(coverIDStr)
	if err != nil {
		return errs.InvalidCoverParameter
	}

	url, err := c.app.CoverUseCase.GetURL(r.Context(), coverID)
	if err != nil {
		return err
	}

	return c.writer.Write(w, http.StatusOK, &getCoversURLV1Response{
		URL: url,
	})
}
