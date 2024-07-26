package service

import (
	"context"
)

type UserService interface {
	FindBlockList(context.Context, ApiOpts) ([]BlockResp, error)
	FindBlockListWithCache(context.Context, ApiOpts) ([]BlockResp, error)
	FindBlockIDAndTitleList(context.Context, ApiOpts) ([]BlockResp, error)
	FindBlockByIDAndOrTitle(context.Context, ApiOpts) (*BlockResp, error)
	FindBlockBWithoutContentData(context.Context, ApiOpts) (*BlockResp, error)
	FindBlockContentByIDAndOrTitleAndOrContentType(context.Context, ApiOpts) (*BlockResp, error)
}
