package filerepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/pkg/queryer"
)

type createAllDTO struct {
	ID            string `db:"id"`
	Name          string `db:"name"`
	Path          string `db:"path"`
	FileServiceID string `db:"file_service_id"`
}

func (r Repository) CreateAll(ctx context.Context, fileServiceID fileservice.ID, filesToCreate []remotefile.File) error {
	return r.createAll(ctx, r.db, fileServiceID, filesToCreate)
}

func (r Repository) CreateAllTx(ctx context.Context, tx *sqlx.Tx, fileServiceID fileservice.ID, filesToCreate []remotefile.File) error {
	return r.createAll(ctx, tx, fileServiceID, filesToCreate)
}

func (r Repository) createAll(ctx context.Context, queryer queryer.Queryer, fileServiceID fileservice.ID, filesToCreate []remotefile.File) error {
	if len(filesToCreate) == 0 {
		return nil
	}

	query := `
INSERT INTO files (id, name, path, file_service_id)
VALUES (:id, :name, :path, :file_service_id)
`

	songsDTO := make([]createAllDTO, len(filesToCreate))
	for i, fileToCreate := range filesToCreate {

		songsDTO[i] = createAllDTO{
			ID:            fileToCreate.ID.String(),
			Name:          fileToCreate.Name,
			Path:          fileToCreate.Path,
			FileServiceID: fileServiceID.String(),
		}
	}

	_, err := queryer.NamedExecContext(ctx, query, songsDTO)
	if err != nil {
		return err
	}

	return nil
}
