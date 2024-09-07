package songrepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/song"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/werr"
)

type getFileDTO struct {
	ID   string `db:"id"`
	Name string `db:"name"`
	Path string `db:"path"`
}

func (r *Repository) GetFile(ctx context.Context, songID song.ID) (remotefile.File, error) {
	return r.getFile(ctx, r.db, songID)
}

func (r *Repository) GetFileTx(ctx context.Context, tx *sqlx.Tx, songID song.ID) (remotefile.File, error) {
	return r.getFile(ctx, tx, songID)
}

func (r *Repository) getFile(ctx context.Context, queryer queryer.Queryer, songID song.ID) (remotefile.File, error) {
	query := `
SELECT f.id,
       f.name,
       f.path
FROM songs s
         JOIN files f ON s.file_id = f.id
WHERE s.id = $1
`

	fileDTO := &getFileDTO{}
	err := queryer.GetContext(ctx, fileDTO, query, songID.String())
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
