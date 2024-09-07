package coverrepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/werr"
)

type getFileDTO struct {
	ID   string `db:"id"`
	Name string `db:"name"`
	Path string `db:"path"`
}

func (r *Repository) GetFile(ctx context.Context, coverID cover.ID) (remotefile.File, error) {
	return r.getFile(ctx, r.db, coverID)
}

func (r *Repository) GetFileTx(ctx context.Context, tx *sqlx.Tx, coverID cover.ID) (remotefile.File, error) {
	return r.getFile(ctx, tx, coverID)
}

func (r *Repository) getFile(ctx context.Context, queryer queryer.Queryer, coverID cover.ID) (remotefile.File, error) {
	query := `
SELECT f.id,
       f.name,
       f.path
FROM covers c
         JOIN files f ON c.file_id = f.id
WHERE c.id = $1
`

	fileDTO := &getFileDTO{}
	err := queryer.GetContext(ctx, fileDTO, query, coverID.String())
	if err != nil {
		return remotefile.File{}, werr.WrapSE("failed to execute query", err)
	}

	fileID, err := remotefile.ParseID(fileDTO.ID)
	if err != nil {
		return remotefile.File{}, werr.WrapSE("failed to parse file id", err)
	}

	return remotefile.File{
		ID:   fileID,
		Name: fileDTO.Name,
		Path: fileDTO.Path,
	}, nil
}
