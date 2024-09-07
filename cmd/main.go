package main

import (
	httpcontroller "github.com/puny-activity/music/api/http/controller"
	httpmiddleware "github.com/puny-activity/music/api/http/middleware"
	httprouter "github.com/puny-activity/music/api/http/router"
	"github.com/puny-activity/music/config"
	"github.com/puny-activity/music/internal/app"
	appconfig "github.com/puny-activity/music/internal/config"
	"github.com/puny-activity/music/pkg/chimux"
	"github.com/puny-activity/music/pkg/httpresp"
	"github.com/puny-activity/music/pkg/httpsrvr"
	"github.com/puny-activity/music/pkg/werr"
	"github.com/puny-activity/music/pkg/zerologger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		panic(werr.WrapSE("failed to parse config", err))
	}

	log, err := zerologger.NewLogger(cfg.Logger.Level)
	if err != nil {
		panic(werr.WrapSE("failed to create logger", err))
	}

	appConfig := appconfig.App{
		Database: appconfig.Database{
			Host:           cfg.App.Database.Host,
			Port:           cfg.App.Database.Port,
			Name:           cfg.App.Database.Name,
			User:           cfg.App.Database.User,
			Password:       cfg.App.Database.Password,
			MigrationsPath: cfg.App.Database.MigrationsPath,
		},
	}

	application := app.New(appConfig, log)

	chiMux := chimux.New()
	httpMiddleware := httpmiddleware.New()
	httpRespWriter := httpresp.NewWriter()
	httpWrapper := httprouter.NewWrapper(httpRespWriter, nil, log)
	controller := httpcontroller.New(application, httpRespWriter, log)
	httpRouter := httprouter.New(&cfg.API.HTTP, chiMux, httpMiddleware, httpWrapper, controller, log)
	httpRouter.Setup()

	httpServer := httpsrvr.New(
		chiMux,
		httpsrvr.Addr(cfg.API.HTTP.Host, cfg.API.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info().Str("signal", s.String()).Msg("interrupt")
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
	}
	err = application.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close application")
	}
}
