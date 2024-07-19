package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/modylegi/service/internal/domain/repository"
)

type BlockRepository struct {
	db *sqlx.DB
}

func NewBlockRepository(db *sqlx.DB) *BlockRepository {
	return &BlockRepository{
		db: db,
	}
}

func (r *BlockRepository) GetBlockID(ctx context.Context, condition repository.Condition) error {
	var res any

	query := `
	SELECT
		b.id
	FROM
		blocks b
	WHERE
		b.id = $1;
	`
	return r.db.GetContext(ctx, &res, query, condition.Args()...)
}

func (r *BlockRepository) GetBlockIDByUserID(ctx context.Context, condition repository.Condition) error {
	var res int
	query := `
	SELECT DISTINCT
		b.id AS block_id
	FROM
		blocks b
		JOIN scenario_mapping sm 	ON sm.key = b.key
		JOIN scenarios s 			ON s.id = sm.scenario_id
		JOIN scenario_user su 		ON su.scenario_id = s.id
	`
	query += condition.String()

	return r.db.GetContext(ctx, &res, query, condition.Args()...)
}

func (r *BlockRepository) GetBlockTitle(ctx context.Context, condition repository.Condition) error {
	var res any

	query := `
	SELECT
		b.title
	FROM
		blocks b
	WHERE
		b.title = $1;
	`
	return r.db.GetContext(ctx, &res, query, condition.Args()...)
}

func (r *BlockRepository) GetWithoutContentData(ctx context.Context, condition repository.Condition) ([]repository.BlockContent, error) {
	var res []repository.BlockContent

	query := `
	SELECT DISTINCT
        b.id AS block_id,
        b.title AS block_title,
        c.id AS content_id,
        c.name AS content_name,
        c.content_type_id AS content_type
	FROM
        blocks b
		JOIN content_mapping cm 	ON cm.content_block_id = b.id
		JOIN contents c 			ON c.id = cm.content_id
	`
	if condition.GetScenarioUserID() != 0 {
		query += `
		JOIN scenario_mapping sm 	ON sm.key = b.key
		JOIN scenarios s			ON s.id = sm.scenario_id
		JOIN scenario_user su		ON su.scenario_id = s.id
		`
	}
	query += condition.String() + " ORDER BY b.id;"
	if err := r.db.SelectContext(ctx, &res, query, condition.Args()...); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *BlockRepository) GetIDAndTitleList(ctx context.Context, condition repository.Condition) ([]repository.Block, error) {
	var res []repository.Block

	query := `
	SELECT DISTINCT
        b.id AS block_id,
        b.title AS block_title
	FROM
        blocks b
	`
	if condition.GetScenarioUserID() != 0 {
		query += `
		JOIN scenario_mapping sm 	ON sm.key = b.key
		JOIN scenarios s			ON s.id = sm.scenario_id
		JOIN scenario_user su		ON su.scenario_id = s.id
	WHERE su.user_id = $1
		`
	}

	query += " ORDER BY b.id;"

	if err := r.db.SelectContext(ctx, &res, query, condition.Args()...); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *BlockRepository) GetAll(ctx context.Context, condition repository.Condition) ([]repository.BlockContentContentMapping, error) {
	var res []repository.BlockContentContentMapping
	query := `
	SELECT DISTINCT
        b.id AS block_id,
        b.title AS block_title,
		b.background_image AS block_background_image,
		cm.rating AS content_mapping_rating,
        c.name AS content_name,
        c.content_type_id AS content_type,
		c.content AS content_data
	FROM
        blocks b
		JOIN content_mapping cm 	ON cm.content_block_id = b.id
		JOIN contents c 			ON c.id = cm.content_id
	`
	if condition.GetScenarioUserID() != 0 {
		query += `
		JOIN scenario_mapping sm 	ON sm.key = b.key
		JOIN scenarios s 			ON s.id = sm.scenario_id
		JOIN scenario_user su 		ON su.scenario_id = s.id
		`
	}

	query += condition.String() + " ORDER BY b.id;"

	if err := r.db.SelectContext(ctx, &res, query, condition.Args()...); err != nil {
		return nil, err
	}

	return res, nil
}
