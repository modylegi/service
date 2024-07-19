package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/modylegi/service/internal/domain/repository"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindUserID(ctx context.Context, condition repository.Condition) error {
	var res any
	query := `
	SELECT
		u.id
	FROM
		users u
	WHERE
		u.id = $1;
	`

	return r.db.GetContext(ctx, &res, query, condition.Args()...)
}

func (r *UserRepository) FindScenario(ctx context.Context, condition repository.Condition) error {
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
