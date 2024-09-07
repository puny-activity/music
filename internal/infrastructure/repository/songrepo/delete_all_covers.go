package songrepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/werr"
)

func (r *Repository) DeleteAllCovers(ctx context.Context, fileServiceID fileservice.ID) error {
	return r.deleteAllCovers(ctx, r.db, fileServiceID)
}

func (r *Repository) DeleteAllCoversTx(ctx context.Context, tx *sqlx.Tx, fileServiceID fileservice.ID) error {
	return r.deleteAllCovers(ctx, tx, fileServiceID)
}

func (r *Repository) deleteAllCovers(ctx context.Context, queryer queryer.Queryer, fileServiceID fileservice.ID) error {
	query := `
UPDATE songs s
SET cover_id = NULL
FROM files f
WHERE s.file_id = f.id
  AND f.file_service_id = $1
`

	_, err := queryer.ExecContext(ctx, query, fileServiceID.String())
	if err != nil {
		return werr.WrapSE("failed to execute query", err)
	}

	return nil
}
