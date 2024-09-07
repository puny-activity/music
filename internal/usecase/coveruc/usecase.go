package coveruc

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/pkg/txmanager"
	"github.com/rs/zerolog"
)

type UseCase struct {
	coverRepository coverRepository
	fileRepository  fileRepository
	txManager       txmanager.Transactor
	log             *zerolog.Logger
}

func New(coverRepository coverRepository, fileRepository fileRepository, txManager txmanager.Transactor,
	log *zerolog.Logger) *UseCase {
	return &UseCase{
		coverRepository: coverRepository,
		fileRepository:  fileRepository,
		txManager:       txManager,
		log:             log,
	}
}

type coverRepository interface {
	GetFileTx(ctx context.Context, tx *sqlx.Tx, coverID cover.ID) (remotefile.File, error)
}

type fileRepository interface {
	GetFileServiceTx(ctx context.Context, tx *sqlx.Tx, fileID remotefile.ID) (fileservice.FileService, error)
}
