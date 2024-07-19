package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/modylegi/service/internal/domain/repository"
)

type ScenarioUserRepository struct {
	db *sqlx.DB
}

func NewScenarioUserRepository(db *sqlx.DB) *ScenarioUserRepository {
	return &ScenarioUserRepository{
		db: db,
	}
}

func (r *ScenarioUserRepository) GetScenarioByUserID(ctx context.Context, condition repository.Condition) error {
	var res any
	query := `
	SELECT
		su.id
	FROM
		scenario_user su
	WHERE
		su.user_id=$1;
	`
	return r.db.GetContext(ctx, &res, query, condition.Args()...)

}
