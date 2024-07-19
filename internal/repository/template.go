package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/modylegi/service/internal/domain/repository"
)

type TemplateRepository struct {
	db *sqlx.DB
}

func NewTemplateRepository(db *sqlx.DB) *TemplateRepository {
	return &TemplateRepository{
		db: db,
	}
}

func (r *TemplateRepository) FindTemplateID(ctx context.Context, condition repository.Condition) error {
	var res any

	query := `
	SELECT 
		id
	FROM
		template_contents
	WHERE
		id = $1;
	`
	return r.db.GetContext(ctx, &res, query, condition.Args()...)

}

func (r *TemplateRepository) FindTemplateName(ctx context.Context, condition repository.Condition) error {
	var res any

	query := `
	SELECT 
		id
	FROM
		template_contents
	WHERE
		name = $1;
	`
	return r.db.GetContext(ctx, &res, query, condition.Args()...)
}

func (r *TemplateRepository) FindTemplates(ctx context.Context, condition repository.Condition) ([]repository.TemplateContent, error) {
	var res []repository.TemplateContent

	query := `
	SELECT
		tc.id AS template_content_id,
        tc.name AS template_content_name,
		tc.content_type_id AS template_content_type,
		tc.template_content AS template_content_data
    FROM
        template_contents tc
	`
	query += condition.String() + " ORDER BY tc.id;"

	if err := r.db.SelectContext(ctx, &res, query, condition.Args()...); err != nil {
		return nil, err
	}
	return res, nil

}
