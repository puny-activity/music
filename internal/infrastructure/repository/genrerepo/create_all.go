package genrerepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/pkg/queryer"
)

type createAllEntity struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (r *Repository) CreateAll(ctx context.Context, genres []genre.Base) error {
	return r.createAll(ctx, r.db, genres)
}

func (r *Repository) CreateAllTx(ctx context.Context, tx *sqlx.Tx, genres []genre.Base) error {
	return r.createAll(ctx, tx, genres)
}

func (r *Repository) createAll(ctx context.Context, queryer queryer.Queryer, genres []genre.Base) error {
	if len(genres) == 0 {
		return nil
	}

	query := `
INSERT INTO genres(id, name) 
VALUES (:id, :name)
`

	genresRepo := make([]createAllEntity, len(genres))
	for i, genreItem := range genres {
		genresRepo[i] = createAllEntity{
			ID:   uuid.UUID(*genreItem.ID),
			Name: genreItem.Name,
		}
	}

	_, err := queryer.NamedExecContext(ctx, query, genresRepo)
	if err != nil {
		return err
	}

	return nil
}
