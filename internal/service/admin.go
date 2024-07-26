package service

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/modylegi/service/internal/domain/repository"
	"github.com/modylegi/service/internal/domain/service"
	repositoryImpl "github.com/modylegi/service/internal/repository"
)

type AdminService struct {
	contentRepo     repository.ContentRepository
	contentTypeRepo repository.ContentTypeRepository
	templateRepo    repository.TemplateRepository
	blockRepo       repository.BlockRepository
}

func NewAdminService(db *sqlx.DB) *AdminService {
	return &AdminService{
		contentRepo:     repositoryImpl.NewContentRepository(db),
		contentTypeRepo: repositoryImpl.NewContentTypeRepository(db),
		templateRepo:    repositoryImpl.NewTemplateRepository(db),
		blockRepo:       repositoryImpl.NewBlockRepository(db),
	}
}

func (s *AdminService) FindBlockIDAndTitleList(ctx context.Context) ([]service.BlockResp, error) {
	condition := &repositoryImpl.Condition{}
	res, err := s.blockRepo.GetIDAndTitleList(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	var resp []service.BlockResp
	for _, v := range res {
		block := service.BlockResp{
			ID:    v.BlockID,
			Title: v.BlockTitle,
		}
		resp = append(resp, block)
	}
	return resp, nil
}

func (s *AdminService) FindBlockByIDAndOrTitle(ctx context.Context, opts service.ApiOpts) (*service.BlockResp, error) {
	condition := &repositoryImpl.Condition{
		BlockID:    opts.BlockID,
		BlockTitle: opts.BlockTitle,
	}
	res, err := s.blockRepo.GetAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	resp := &service.BlockResp{
		ID:              res[0].BlockID,
		Title:           res[0].BlockTitle,
		BackgroundImage: res[0].BlockBackgroundImage.String,
	}
	for _, item := range res {
		content := service.ContentResp{
			Rating:      int(item.ContentMappingRating.Int64),
			Name:        item.ContentName,
			ContentType: item.ContentType,
			Content:     json.RawMessage(item.ContentData),
		}
		resp.Contents = append(resp.Contents, content)
	}
	return resp, nil
}

func (s *AdminService) FindBlockWithoutContentData(ctx context.Context, opts service.ApiOpts) (*service.BlockResp, error) {
	condition := &repositoryImpl.Condition{
		BlockID: opts.BlockID,
	}
	res, err := s.blockRepo.GetWithoutContentData(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	resp := &service.BlockResp{
		ID:    res[0].BlockID,
		Title: res[0].BlockTitle,
	}

	for _, item := range res {
		content := service.ContentResp{
			ID:          item.ContentID,
			Name:        item.ContentName,
			ContentType: item.ContentType,
		}

		resp.Contents = append(resp.Contents, content)
	}
	return resp, nil
}

func (s *AdminService) FindBlockContentByIDAndOrTitleAndOrContentType(ctx context.Context, opts service.ApiOpts) (*service.BlockResp, error) {
	condition := &repositoryImpl.Condition{
		BlockID:       opts.BlockID,
		ContentID:     opts.ContentID,
		ContentName:   opts.ContentName,
		ContentTypeID: opts.ContentTypeID,
	}
	res, err := s.contentRepo.FindContent(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	resp := &service.BlockResp{}
	for _, item := range res {

		content := service.ContentResp{
			Rating:      int(item.ContentMappingRating.Int64),
			Name:        item.ContentName,
			ContentType: item.ContentType,
			Content:     item.ContentData,
		}

		resp.Contents = append(resp.Contents, content)
	}
	return resp, nil
}

func (s *AdminService) FindTemplateList(ctx context.Context) ([]service.TemplateResp, error) {
	condition := &repositoryImpl.Condition{}
	res, err := s.templateRepo.FindTemplates(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	var resp []service.TemplateResp

	for _, v := range res {
		template := service.TemplateResp{
			ID:          v.TemplateContentID,
			Name:        v.TemplateContentName,
			ContentType: v.TemplateContentType,
		}
		resp = append(resp, template)
	}

	return resp, nil

}

func (s *AdminService) FindTemplateByIDAndOrNameAndOrContentType(ctx context.Context, opts service.ApiOpts) (*service.TemplateResp, error) {
	condition := &repositoryImpl.Condition{
		TemplateID:            opts.TemplateID,
		TemplateName:          opts.TemplateName,
		TemplateContentTypeID: opts.TemplateContentTypeID,
	}
	res, err := s.templateRepo.FindTemplates(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	resp := &service.TemplateResp{
		ID:          res[0].TemplateContentID,
		Name:        res[0].TemplateContentName,
		ContentType: res[0].TemplateContentType,
		Content:     json.RawMessage(res[0].TemplateContent),
	}

	return resp, nil
}
