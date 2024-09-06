package httpcontroller

import (
	"github.com/puny-activity/music/internal/app"
	"github.com/puny-activity/music/pkg/httpresp"
	"github.com/rs/zerolog"
)

type Controller struct {
	app    *app.App
	writer *httpresp.Writer
	log    *zerolog.Logger
}

func New(app *app.App, writer *httpresp.Writer, log *zerolog.Logger) *Controller {
	return &Controller{
		app:    app,
		writer: writer,
		log:    log,
	}
}
