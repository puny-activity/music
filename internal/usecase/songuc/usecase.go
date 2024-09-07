package songuc

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/internal/entity/song"
	"github.com/puny-activity/music/pkg/txmanager"
	"github.com/rs/zerolog"
)

type UseCase struct {
	songRepository songRepository
	fileRepository fileRepository
	txManager      txmanager.Transactor
	log            *zerolog.Logger
}

func New(songRepository songRepository, fileRepository fileRepository, txManager txmanager.Transactor,
	log *zerolog.Logger) *UseCase {
	return &UseCase{
		songRepository: songRepository,
		fileRepository: fileRepository,
		txManager:      txManager,
		log:            log,
	}
}

type songRepository interface {
	GetFileTx(ctx context.Context, tx *sqlx.Tx, songID song.ID) (remotefile.File, error)
}

type fileRepository interface {
	GetFileServiceTx(ctx context.Context, tx *sqlx.Tx, fileID remotefile.ID) (fileservice.FileService, error)
}
