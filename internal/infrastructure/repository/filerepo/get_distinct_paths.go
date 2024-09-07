package filerepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/pkg/queryer"
)

func (r Repository) GetDistinctPaths(ctx context.Context, fileServiceID fileservice.ID) ([]string, error) {
	return r.getDistinctPaths(ctx, r.db, fileServiceID)
}

func (r Repository) GetDistinctPathsTx(ctx context.Context, tx *sqlx.Tx, fileServiceID fileservice.ID) ([]string, error) {
	return r.getDistinctPaths(ctx, tx, fileServiceID)
}

func (r Repository) getDistinctPaths(ctx context.Context, queryer queryer.Queryer, fileServiceID fileservice.ID) ([]string, error) {
	query := `
SELECT DISTINCT f.path
FROM files f
WHERE f.file_service_id = $1
`

	paths := make([]string, 0)
	err := queryer.SelectContext(ctx, &paths, query, fileServiceID.String())
	if err != nil {
		return nil, err
	}

	return paths, nil
}
