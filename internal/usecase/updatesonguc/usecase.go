package updatesonguc

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/internal/infrastructure/fileserviceclient"
	"github.com/puny-activity/music/pkg/txmanager"
	"github.com/rs/zerolog"
)

type UseCase struct {
	fileServiceRepository fileServiceRepository
	fileRepository        fileRepository
	genreRepository       genreRepository
	albumRepository       albumRepository
	artistRepository      artistRepository
	songRepository        songRepository
	fileServiceController *fileserviceclient.Controller
	txManager             txmanager.Transactor
	log                   *zerolog.Logger
}

func New(fileServiceRepository fileServiceRepository, fileRepository fileRepository, genreRepository genreRepository,
	albumRepository albumRepository, artistRepository artistRepository, songRepository songRepository,
	fileServiceController *fileserviceclient.Controller, txManager txmanager.Transactor,
	log *zerolog.Logger) *UseCase {
	return &UseCase{
		fileServiceRepository: fileServiceRepository,
		fileRepository:        fileRepository,
		genreRepository:       genreRepository,
		albumRepository:       albumRepository,
		artistRepository:      artistRepository,
		songRepository:        songRepository,
		fileServiceController: fileServiceController,
		txManager:             txManager,
		log:                   log,
	}
}

type fileServiceRepository interface {
	GetAll(ctx context.Context) ([]fileservice.FileService, error)
}

type fileRepository interface {
	DeleteTx(ctx context.Context, tx *sqlx.Tx, fileID remotefile.ID) error
	GetAllTx(ctx context.Context, tx *sqlx.Tx) ([]remotefile.File, error)
}

type genreRepository interface {
	GetAllTx(ctx context.Context, tx *sqlx.Tx) ([]genre.Genre, error)
	CreateAllTx(ctx context.Context, tx *sqlx.Tx, genresToCreate []genre.Genre) error
}

type albumRepository interface {
	GetAllTx(ctx context.Context, tx *sqlx.Tx) ([]album.Album, error)
	CreateAllTx(ctx context.Context, tx *sqlx.Tx, albumsToCreate []album.Album) error
}

type artistRepository interface {
	GetAllTx(ctx context.Context, tx *sqlx.Tx) ([]artist.Base, error)
	CreateAllTx(ctx context.Context, tx *sqlx.Tx, artistsToCreate []artist.Base) error
}

type songRepository interface {
}