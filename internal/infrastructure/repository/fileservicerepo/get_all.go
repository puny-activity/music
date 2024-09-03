package fileservicerepo

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
	"github.com/puny-activity/music/pkg/werr"
)

type getAllEntity struct {
	ID        uuid.UUID `db:"id"`
	Address   string    `db:"address"`
	ScannedAt *string   `db:"scanned_at"`
}

func (r *Repository) GetAll(ctx context.Context) ([]fileservice.FileService, error) {
	return r.getAll(ctx, r.db)
}

func (r *Repository) GetAllTx(ctx context.Context, tx *sqlx.Tx) ([]fileservice.FileService, error) {
	return r.getAll(ctx, tx)
}

func (r *Repository) getAll(ctx context.Context, queryer queryer.Queryer) ([]fileservice.FileService, error) {
	query := `
SELECT fs.id,
       fs.address,
       fs.scanned_at
FROM file_services fs
`

	fileServicesRepo := make([]getAllEntity, 0)
	err := queryer.SelectContext(ctx, &fileServicesRepo, query)
	if err != nil {
		return nil, err
	}

	fileServices := make([]fileservice.FileService, len(fileServicesRepo))
	for i := range fileServicesRepo {
		var scannedAt *carbon.Carbon = nil
		if fileServicesRepo[i].ScannedAt != nil {
			scannedAt = util.ToPointer(carbon.Parse(*fileServicesRepo[i].ScannedAt))
			if scannedAt.Error != nil {
				return nil, werr.WrapSE("failed to parse scanned at time", scannedAt.Error)
			}
		}

		fileServices[i] = fileservice.FileService{
			ID:        util.ToPointer(fileservice.NewID(fileServicesRepo[i].ID)),
			Address:   fileServicesRepo[i].Address,
			ScannedAt: scannedAt,
		}
	}

	return fileServices, nil
}
