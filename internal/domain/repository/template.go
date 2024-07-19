package repository

import (
	"context"
)

type TemplateRepository interface {
	FindTemplateID(context.Context, Condition) error
	FindTemplateName(context.Context, Condition) error
	FindTemplates(context.Context, Condition) ([]TemplateContent, error)
}
