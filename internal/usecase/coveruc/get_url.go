package coveruc

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/pkg/werr"
)

func (u *UseCase) GetURL(ctx context.Context, coverID cover.ID) (string, error) {
	url := ""

	err := u.txManager.Transaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		coverFile, err := u.coverRepository.GetFileTx(ctx, tx, coverID)
		if err != nil {
			return werr.WrapSE("failed to get file", err)
		}

		fileService, err := u.fileRepository.GetFileServiceTx(ctx, tx, coverFile.ID)
		if err != nil {
			return werr.WrapSE("failed to get file service", err)
		}

		url = fmt.Sprintf("%s/stream/%s", fileService.HTTPAddress, coverFile.ID)
		return nil
	})
	if err != nil {
		return "", err
	}

	return url, nil
}
