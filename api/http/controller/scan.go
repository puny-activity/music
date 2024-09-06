package httpcontroller

import (
	httpapi "github.com/puny-activity/music/api/http"
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/werr"
	"net/http"
)

func (c *Controller) Scan(w http.ResponseWriter, r *http.Request) error {
	version := r.Header.Get(httpapi.APIVersion)
	switch version {
	case "1":
		return c.scanV1(w, r)
	default:
		return errs.InvalidAPIVersion
	}
}

func (c *Controller) scanV1(w http.ResponseWriter, r *http.Request) error {
	err := c.app.UpdateSongUseCase.Update(r.Context())
	if err != nil {
		return werr.WrapSE("failed to get genres", err)
	}

	return c.writer.Write(w, http.StatusOK, nil)
}
