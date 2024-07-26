package service

import (
	"context"
)

type ValidationService interface {
	UserID(context.Context, string) (int, error)
	BlockID(context.Context, string) (int, error)
	BlockTitle(context.Context, string) (string, error)
	ContentTypeID(context.Context, string) (int, error)
	ContentID(context.Context, string) (int, error)
	ContentName(context.Context, string) (string, error)
	TemplateID(context.Context, string) (int, error)
	TemplateName(context.Context, string) (string, error)
	LinkedScenarios(context.Context, int) error
	LinkedScenarioBlock(context.Context, ApiOpts) error
}
