package songrepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/werr"
)

func (r Repository) SetCoverByPath(ctx context.Context, path string, fileServiceInfo fileservice.ID, coverID cover.ID) error {
	return r.setCoverByPath(ctx, r.db, path, fileServiceInfo, coverID)
}

func (r Repository) SetCoverByPathTx(ctx context.Context, tx *sqlx.Tx, path string, fileServiceInfo fileservice.ID, coverID cover.ID) error {
	return r.setCoverByPath(ctx, tx, path, fileServiceInfo, coverID)
}

func (r Repository) setCoverByPath(ctx context.Context, queryer queryer.Queryer, path string, fileServiceInfo fileservice.ID, coverID cover.ID) error {
	query := `
UPDATE songs
SET cover_id = $3
WHERE id IN (SELECT s.id
             FROM songs s
                      JOIN files f ON s.file_id = f.id
             WHERE f.file_service_id = $1
               AND f.path = $2)
`

	_, err := queryer.ExecContext(ctx, query, fileServiceInfo.String(), path, coverID.String())
	if err != nil {
		return werr.WrapSE("failed to execute query", err)
	}

	return nil
}
