package coverrepo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/cover"
	"github.com/puny-activity/music/pkg/queryer"
)

type createAllDTO struct {
	ID     string `db:"id"`
	Width  int    `db:"width"`
	Height int    `db:"height"`
	FileID string `db:"file_id"`
}

func (r Repository) CreateAll(ctx context.Context, coversToCreate []cover.Cover) error {
	return r.createAll(ctx, r.db, coversToCreate)
}

func (r Repository) CreateAllTx(ctx context.Context, tx *sqlx.Tx, coversToCreate []cover.Cover) error {
	return r.createAll(ctx, tx, coversToCreate)
}

func (r Repository) createAll(ctx context.Context, queryer queryer.Queryer, coversToCreate []cover.Cover) error {
	if len(coversToCreate) == 0 {
		return nil
	}

	query := `
INSERT INTO covers (id, width, height, file_id)
VALUES (:id, :width, :height, :file_id)
`

	coversDTO := make([]createAllDTO, len(coversToCreate))
	for i, coverToCreate := range coversToCreate {

		coversDTO[i] = createAllDTO{
			ID:     coverToCreate.ID.String(),
			Width:  coverToCreate.Width,
			Height: coverToCreate.Height,
			FileID: coverToCreate.FileID.String(),
		}
	}

	_, err := queryer.NamedExecContext(ctx, query, coversDTO)
	if err != nil {
		return err
	}

	return nil
}
