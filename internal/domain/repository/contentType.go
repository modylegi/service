package repository

import (
	"context"
)

type ContentTypeRepository interface {
	FindContentTypeID(context.Context, Condition) error
}
