package songuc

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song"
	"github.com/puny-activity/music/pkg/werr"
)

func (u *UseCase) GetURL(ctx context.Context, songID song.ID) (string, error) {
	url := ""

	err := u.txManager.Transaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		songFile, err := u.songRepository.GetFileTx(ctx, tx, songID)
		if err != nil {
			return werr.WrapSE("failed to get file", err)
		}

		fileService, err := u.fileRepository.GetFileServiceTx(ctx, tx, songFile.ID)
		if err != nil {
			return werr.WrapSE("failed to get file service", err)
		}

		url = fmt.Sprintf("%s/stream/%s", fileService.HTTPAddress, songFile.ID)
		return nil
	})
	if err != nil {
		return "", err
	}

	return url, nil
}
