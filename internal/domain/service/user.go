package service

import (
	"context"

	"github.com/modylegi/service/internal/api"
)

type UserService interface {
	FindBlockList(context.Context, api.Opts) ([]BlockResp, error)
	FindBlockListWithCache(context.Context, api.Opts) ([]BlockResp, error)
	FindBlockIDAndTitleList(context.Context, api.Opts) ([]BlockResp, error)
	FindBlockByIDAndOrTitle(context.Context, api.Opts) (*BlockResp, error)
	FindBlockBWithoutContentData(context.Context, api.Opts) (*BlockResp, error)
	FindBlockContentByIDAndOrTitleAndOrContentType(context.Context, api.Opts) (*BlockResp, error)
}
