package main

import (
	"github.com/puny-activity/music/config"
	"github.com/puny-activity/music/internal/app"
	appconfig "github.com/puny-activity/music/internal/config"
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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info().Str("signal", s.String()).Msg("interrupt")
	}

	err = application.Close()
	if err != nil {
		panic(err)
	}
}
