package service

import (
	"context"
)

type AdminService interface {
	FindBlockIDAndTitleList(context.Context) ([]BlockResp, error)
	FindBlockByIDAndOrTitle(context.Context, ApiOpts) (*BlockResp, error)
	FindBlockWithoutContentData(context.Context, ApiOpts) (*BlockResp, error)
	FindBlockContentByIDAndOrTitleAndOrContentType(context.Context, ApiOpts) (*BlockResp, error)
	FindTemplateList(context.Context) ([]TemplateResp, error)
	FindTemplateByIDAndOrNameAndOrContentType(context.Context, ApiOpts) (*TemplateResp, error)
}
