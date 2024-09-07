package app

import (
	"context"
	"fmt"
	"github.com/puny-activity/music/internal/config"
	"github.com/puny-activity/music/internal/infrastructure/fileserviceclient"
	"github.com/puny-activity/music/internal/infrastructure/repository/albumrepo"
	"github.com/puny-activity/music/internal/infrastructure/repository/artistrepo"
	"github.com/puny-activity/music/internal/infrastructure/repository/coverrepo"
	"github.com/puny-activity/music/internal/infrastructure/repository/filerepo"
	"github.com/puny-activity/music/internal/infrastructure/repository/fileservicerepo"
	"github.com/puny-activity/music/internal/infrastructure/repository/genrerepo"
	"github.com/puny-activity/music/internal/infrastructure/repository/songrepo"
	"github.com/puny-activity/music/internal/usecase/albumuc"
	"github.com/puny-activity/music/internal/usecase/artistuc"
	"github.com/puny-activity/music/internal/usecase/coveruc"
	"github.com/puny-activity/music/internal/usecase/fileserviceuc"
	"github.com/puny-activity/music/internal/usecase/genreuc"
	"github.com/puny-activity/music/internal/usecase/songuc"
	"github.com/puny-activity/music/internal/usecase/updatesonguc"
	"github.com/puny-activity/music/pkg/postgres"
	"github.com/puny-activity/music/pkg/txmanager"
	"github.com/puny-activity/music/pkg/werr"
	"github.com/rs/zerolog"
)

type App struct {
	FileServiceUseCase           *fileserviceuc.UseCase
	UpdateSongUseCase            *updatesonguc.UseCase
	SongUseCase                  *songuc.UseCase
	CoverUseCase                 *coveruc.UseCase
	GenreUseCase                 *genreuc.UseCase
	AlbumUseCase                 *albumuc.UseCase
	ArtistUseCase                *artistuc.UseCase
	db                           *postgres.Postgres
	fileServiceClientsController fileserviceclient.Controller
	log                          *zerolog.Logger
}

func New(cfg config.App, log *zerolog.Logger) *App {
	db, err := postgres.New(cfg.Database.ConnectionString())
	if err != nil {
		panic(err)
	}
	err = db.RunMigrations(cfg.Database.MigrationsPath)
	if err != nil {
		panic(err)
	}

	txManager := txmanager.New(db.DB)

	fileServiceRepository := fileservicerepo.New(db.DB, txManager, log)
	fileRepository := filerepo.New(db.DB, txManager, log)
	genreRepository := genrerepo.New(db.DB, txManager, log)
	albumRepository := albumrepo.New(db.DB, txManager, log)
	artistRepository := artistrepo.New(db.DB, txManager, log)
	songRepository := songrepo.New(db.DB, txManager, log)
	coverRepository := coverrepo.New(db.DB, txManager, log)

	fileServiceClientsController := fileserviceclient.NewController(log)

	fileServiceUseCase := fileserviceuc.New(fileServiceRepository, fileServiceClientsController, txManager, log)
	updateSongUseCase := updatesonguc.New(fileServiceRepository, fileRepository, genreRepository, albumRepository, artistRepository,
		songRepository, coverRepository, fileServiceClientsController, txManager, log)
	songUseCase := songuc.New(songRepository, fileRepository, txManager, log)
	coverUseCase := coveruc.New(coverRepository, fileRepository, txManager, log)
	genreUseCase := genreuc.New(genreRepository, txManager, log)
	albumUseCase := albumuc.New(albumRepository, txManager, log)
	artistUseCase := artistuc.New(artistRepository, txManager, log)

	err = fileServiceUseCase.ReloadClients(context.Background())
	if err != nil {
		fmt.Println(werr.WrapSE("failed to reload clients", err))
	}

	return &App{
		FileServiceUseCase: fileServiceUseCase,
		UpdateSongUseCase:  updateSongUseCase,
		SongUseCase:        songUseCase,
		CoverUseCase:       coverUseCase,
		GenreUseCase:       genreUseCase,
		AlbumUseCase:       albumUseCase,
		ArtistUseCase:      artistUseCase,
		db:                 db,
		log:                log,
	}
}

func (a *App) Close() error {
	err := a.db.Close()
	if err != nil {
		a.log.Error().Err(err).Msg("failed to close database connection")
	}

	err = a.fileServiceClientsController.Reset()
	if err != nil {
		a.log.Error().Err(err).Msg("failed to close file service clients")
	}

	return nil
}
