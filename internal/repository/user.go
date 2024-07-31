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

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*repository.User, error) {
	var res repository.User
	query := `
    SELECT
        u.id AS user_id,
        u.username AS user_name,
        u.password AS user_password,
        u.created_at AS user_created_at,
        u.deleted AS user_deleted
    FROM
        users u
    WHERE
        u.username=$1;
    `
	if err := r.db.GetContext(ctx, &res, query, username); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *UserRepository) Create(ctx context.Context, user *repository.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.UserPassword)
	if err != nil {
		return err
	}
	return nil
}
