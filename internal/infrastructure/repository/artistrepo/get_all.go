package artistrepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
)

type getAllEntity struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (r *Repository) GetAll(ctx context.Context) ([]artist.Base, error) {
	return r.getAll(ctx, r.db)
}

func (r *Repository) GetAllTx(ctx context.Context, tx *sqlx.Tx) ([]artist.Base, error) {
	return r.getAll(ctx, tx)
}

func (r *Repository) getAll(ctx context.Context, queryer queryer.Queryer) ([]artist.Base, error) {
	query := `
SELECT a.id,
       a.name
FROM artists a
`

	artistsRepo := make([]getAllEntity, 0)
	err := queryer.SelectContext(ctx, &artistsRepo, query)
	if err != nil {
		return nil, err
	}

	artists := make([]artist.Base, len(artistsRepo))
	for i := range artistsRepo {
		artists[i] = artist.Base{
			ID:   util.ToPointer(artist.NewID(artistsRepo[i].ID)),
			Name: artistsRepo[i].Name,
		}
	}

	return artists, nil
}
