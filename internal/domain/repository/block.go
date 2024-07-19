package repository

import (
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name BlockRepository
type BlockRepository interface {
	GetBlockID(context.Context, Condition) error
	GetBlockTitle(context.Context, Condition) error
	GetWithoutContentData(context.Context, Condition) ([]BlockContent, error)
	GetAll(context.Context, Condition) ([]BlockContentContentMapping, error)
	GetIDAndTitleList(context.Context, Condition) ([]Block, error)
	GetBlockIDByUserID(context.Context, Condition) error
}
