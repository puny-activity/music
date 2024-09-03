package genrerepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/pkg/queryer"
	"github.com/puny-activity/music/pkg/util"
)

type createAllEntity struct {
	ID      uuid.UUID  `db:"id"`
	Name    string     `db:"name"`
	CoverID *uuid.UUID `db:"cover_id"`
}

func (r *Repository) CreateAll(ctx context.Context, genres []genre.Genre) error {
	return r.createAll(ctx, r.db, genres)
}

func (r *Repository) CreateAllTx(ctx context.Context, tx *sqlx.Tx, genres []genre.Genre) error {
	return r.createAll(ctx, tx, genres)
}

func (r *Repository) createAll(ctx context.Context, queryer queryer.Queryer, genres []genre.Genre) error {
	if len(genres) == 0 {
		return nil
	}

	query := `
INSERT INTO genres(id, name, cover_id) 
VALUES (:id, :name, :cover_id)
`

	genresRepo := make([]createAllEntity, len(genres))
	for i, genreItem := range genres {
		var coverID *uuid.UUID
		if genreItem.Cover != nil {
			coverID = util.ToPointer(uuid.UUID(*genreItem.Cover.ID))
		}
		genresRepo[i] = createAllEntity{
			ID:      uuid.UUID(*genreItem.ID),
			Name:    genreItem.Name,
			CoverID: coverID,
		}
	}

	_, err := queryer.NamedExecContext(ctx, query, genresRepo)
	if err != nil {
		return err
	}

	return nil
}
