package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/modylegi/service/internal/domain/repository"
)

type ContentRepository struct {
	db *sqlx.DB
}

func NewContentRepository(db *sqlx.DB) *ContentRepository {
	return &ContentRepository{
		db: db,
	}
}

func (r *ContentRepository) FindContent(ctx context.Context, condition repository.Condition) ([]repository.ContentContentMapping, error) {
	var res []repository.ContentContentMapping

	query := `
	SELECT DISTINCT
		cm.rating AS content_mapping_rating,
        c.name AS content_name,
        c.content_type_id AS content_type,
		c."content" AS content_data
    FROM
        users u
		JOIN scenario_user su 		ON u.id = su.user_id
		JOIN scenarios s 			ON s.id = su.scenario_id
		JOIN scenario_mapping sm 	ON sm.scenario_id = s.id
		JOIN blocks b 				ON b.key = sm.key
		JOIN content_mapping cm 	ON cm.content_block_id = b.id
		JOIN contents c 			ON c.id = cm.content_id
	`

	query += condition.String()

	if err := r.db.SelectContext(ctx, &res, query, condition.Args()...); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *ContentRepository) FindContentID(ctx context.Context, condition repository.Condition) error {
	var res any

	query := `
	SELECT
		c.id
	FROM
		contents c
	WHERE
		c.id = $1;
	`
	return r.db.GetContext(ctx, &res, query, condition.Args()...)

}

func (r *ContentRepository) FindContentName(ctx context.Context, condition repository.Condition) error {
	var res any

	query := `
	SELECT
		c.name
	FROM
		contents c
	WHERE
		c.name = $1;
	`
	return r.db.GetContext(ctx, &res, query, condition.Args()...)
}
