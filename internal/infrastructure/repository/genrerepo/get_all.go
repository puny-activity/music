package genrerepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
)

type getAllEntity struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (r *Repository) GetAll(ctx context.Context) ([]genre.Base, error) {
	return r.getAll(ctx, r.db)
}

func (r *Repository) GetAllTx(ctx context.Context, tx *sqlx.Tx) ([]genre.Base, error) {
	return r.getAll(ctx, tx)
}

func (r *Repository) getAll(ctx context.Context, queryer queryer.Queryer) ([]genre.Base, error) {
	query := `
SELECT a.id,
       a.name
FROM genres a
`

	genresRepo := make([]getAllEntity, 0)
	err := queryer.SelectContext(ctx, &genresRepo, query)
	if err != nil {
		return nil, err
	}

	genres := make([]genre.Base, len(genresRepo))
	for i := range genresRepo {
		genres[i] = genre.Base{
			ID:   util.ToPointer(genre.NewID(genresRepo[i].ID)),
			Name: genresRepo[i].Name,
		}
	}

	return genres, nil
}
