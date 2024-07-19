package service

import (
	"context"

	"github.com/modylegi/service/internal/api"
)

type AdminService interface {
	FindBlockIDAndTitleList(context.Context) ([]BlockResp, error)
	FindBlockByIDAndOrTitle(context.Context, api.Opts) (*BlockResp, error)
	FindBlockWithoutContentData(context.Context, api.Opts) (*BlockResp, error)
	FindBlockContentByIDAndOrTitleAndOrContentType(context.Context, api.Opts) (*BlockResp, error)
	FindTemplateList(context.Context) ([]TemplateResp, error)
	FindTemplateByIDAndOrNameAndOrContentType(context.Context, api.Opts) (*TemplateResp, error)
}
