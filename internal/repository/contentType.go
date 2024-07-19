package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/modylegi/service/internal/domain/repository"
)

type ContentTypeRepository struct {
	db *sqlx.DB
}

func NewContentTypeRepository(db *sqlx.DB) *ContentTypeRepository {
	return &ContentTypeRepository{
		db: db,
	}
}

func (r *ContentTypeRepository) FindContentTypeID(ctx context.Context, condition repository.Condition) error {
	var res any
	query := `
	SELECT 
		id
	FROM
		content_type_dic
	WHERE
		id = $1;		
	`
	return r.db.GetContext(ctx, &res, query, condition.Args()...)
}
