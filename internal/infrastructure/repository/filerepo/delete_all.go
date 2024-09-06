package filerepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
)

func (r Repository) DeleteAllTx(ctx context.Context, tx *sqlx.Tx, fileToDeleteIDs []remotefile.ID) error {
	if len(fileToDeleteIDs) == 0 {
		return nil
	}

	//TODO implement me
	panic("implement me")
}
