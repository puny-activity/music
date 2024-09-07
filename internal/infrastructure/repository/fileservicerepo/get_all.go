package fileservicerepo

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
	"github.com/puny-activity/music/pkg/werr"
)

type getAllEntity struct {
	ID          string  `db:"id"`
	HTTPAddress string  `db:"http_address"`
	GRPCAddress string  `db:"grpc_address"`
	ScannedAt   *string `db:"scanned_at"`
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
       fs.http_address,
       fs.grpc_address,
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

		fileServiceID, err := fileservice.ParseID(fileServicesRepo[i].ID)
		if err != nil {
			return nil, werr.WrapSE("failed to parse file service id", err)
		}

		fileServices[i] = fileservice.FileService{
			ID:          &fileServiceID,
			GRPCAddress: fileServicesRepo[i].GRPCAddress,
			ScannedAt:   scannedAt,
		}
	}

	return fileServices, nil
}
