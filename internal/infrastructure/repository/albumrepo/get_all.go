package albumrepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
)

type getAllEntity struct {
	ID    uuid.UUID `db:"id"`
	Title string    `db:"title"`
}

func (r *Repository) GetAll(ctx context.Context) ([]album.Base, error) {
	return r.getAll(ctx, r.db)
}

func (r *Repository) GetAllTx(ctx context.Context, tx *sqlx.Tx) ([]album.Base, error) {
	return r.getAll(ctx, tx)
}

func (r *Repository) getAll(ctx context.Context, queryer queryer.Queryer) ([]album.Base, error) {
	query := `
SELECT a.id,
       a.title
FROM albums a
`

	albumsRepo := make([]getAllEntity, 0)
	err := queryer.SelectContext(ctx, &albumsRepo, query)
	if err != nil {
		return nil, err
	}

	albums := make([]album.Base, len(albumsRepo))
	for i := range albumsRepo {
		albums[i] = album.Base{
			ID:    util.ToPointer(album.NewID(albumsRepo[i].ID)),
			Title: albumsRepo[i].Title,
		}
	}

	return albums, nil
}
