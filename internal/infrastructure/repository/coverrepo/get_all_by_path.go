package coverrepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
	"github.com/puny-activity/music/pkg/werr"
)

type getAllDTO struct {
	ID     string `db:"id"`
	Width  int    `db:"width"`
	Height int    `db:"height"`
	FileID string `db:"file_id"`
}

func (r *Repository) GetAllByPath(ctx context.Context, fileServiceID fileservice.ID, path string) ([]cover.Cover, error) {
	return r.getAllByPath(ctx, r.db, fileServiceID, path)
}

func (r *Repository) GetAllByPathTx(ctx context.Context, tx *sqlx.Tx, fileServiceID fileservice.ID, path string) ([]cover.Cover, error) {
	return r.getAllByPath(ctx, tx, fileServiceID, path)
}

func (r *Repository) getAllByPath(ctx context.Context, queryer queryer.Queryer, fileServiceID fileservice.ID, path string) ([]cover.Cover, error) {
	query := `
SELECT c.id,
       c.width,
       c.height,
       c.file_id
FROM covers c
         JOIN files f ON f.id = c.file_id
WHERE f.file_service_id = $1
  AND f.path = $2
`

	coversDTO := make([]getAllDTO, 0)
	err := queryer.SelectContext(ctx, &coversDTO, query, fileServiceID.String(), path)
	if err != nil {
		return nil, err
	}

	covers := make([]cover.Cover, len(coversDTO))
	for i := range coversDTO {
		coverID, err := cover.ParseID(coversDTO[i].ID)
		if err != nil {
			return nil, werr.WrapSE("failed to parse cover id", err)
		}

		fileID, err := remotefile.ParseID(coversDTO[i].FileID)
		if err != nil {
			return nil, werr.WrapSE("failed to parse file id", err)
		}

		covers[i] = cover.Cover{
			ID:     util.ToPointer(coverID),
			Width:  coversDTO[i].Width,
			Height: coversDTO[i].Height,
			FileID: fileID,
		}
	}

	return covers, nil
}
