package repository

import (
	"context"
)

type ContentRepository interface {
	FindContentID(context.Context, Condition) error
	FindContentName(context.Context, Condition) error
	FindContent(context.Context, Condition) ([]ContentContentMapping, error)
}
