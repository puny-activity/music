package filerepo

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/werr"
)

type getFileServiceDTO struct {
	ID          string  `db:"id"`
	HTTPAddress string  `db:"http_address"`
	GRPCAddress string  `db:"grpc_address"`
	ScannedAt   *string `db:"scanned_at"`
}

func (r *Repository) GetFileService(ctx context.Context, fileID remotefile.ID) (fileservice.FileService, error) {
	return r.getFileService(ctx, r.db, fileID)
}

func (r *Repository) GetFileServiceTx(ctx context.Context, tx *sqlx.Tx, fileID remotefile.ID) (fileservice.FileService, error) {
	return r.getFileService(ctx, tx, fileID)
}

func (r *Repository) getFileService(ctx context.Context, queryer queryer.Queryer, fileID remotefile.ID) (fileservice.FileService, error) {
	query := `
SELECT fs.id,
       fs.http_address,
       fs.grpc_address,
       fs.scanned_at
FROM files f
         JOIN file_services fs ON f.file_service_id = fs.id
WHERE f.id = $1
`

	fileServiceDTO := &getFileServiceDTO{}
	err := queryer.GetContext(ctx, fileServiceDTO, query, fileID.String())
	if err != nil {
		return fileservice.FileService{}, werr.WrapSE("failed to execute query", err)
	}

	fileServiceID, err := fileservice.ParseID(fileServiceDTO.ID)
	if err != nil {
		return fileservice.FileService{}, werr.WrapSE("failed to parse file id", err)
	}

	var scannedAt *carbon.Carbon = nil
	if fileServiceDTO.ScannedAt != nil {
		scannedAtLocal := carbon.Parse(*fileServiceDTO.ScannedAt)
		if scannedAtLocal.Error != nil {
			return fileservice.FileService{}, werr.WrapSE("failed to parse scanned at", err)
		}
	}

	return fileservice.FileService{
		ID:          &fileServiceID,
		HTTPAddress: fileServiceDTO.HTTPAddress,
		GRPCAddress: fileServiceDTO.GRPCAddress,
		ScannedAt:   scannedAt,
	}, nil
}
