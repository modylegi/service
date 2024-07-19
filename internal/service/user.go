package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/modylegi/service/internal/api"
	"github.com/modylegi/service/internal/domain/repository"
	"github.com/modylegi/service/internal/domain/service"
	repositoryImpl "github.com/modylegi/service/internal/repository"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserService struct {
	userRepo         repository.UserRepository
	contentRepo      repository.ContentRepository
	contentTypeRepo  repository.ContentTypeRepository
	blockRepo        repository.BlockRepository
	scenarioUserRepo repository.ScenarioUserRepository
	rd               *redis.Client
}

func NewUserService(db *sqlx.DB, rd *redis.Client) *UserService {
	return &UserService{
		userRepo:         repositoryImpl.NewUserRepository(db),
		contentRepo:      repositoryImpl.NewContentRepository(db),
		blockRepo:        repositoryImpl.NewBlockRepository(db),
		contentTypeRepo:  repositoryImpl.NewContentTypeRepository(db),
		scenarioUserRepo: repositoryImpl.NewScenarioUserRepository(db),
		rd:               rd,
	}
}

func (s *UserService) FindBlockList(ctx context.Context, opts api.Opts) ([]service.BlockResp, error) {
	condition := &repositoryImpl.Condition{
		ScenarioUserID: opts.UserID,
	}
	res, err := s.blockRepo.GetAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	blockMap := make(map[int]*service.BlockResp)
	for _, item := range res {
		block, exists := blockMap[item.BlockID]
		if !exists {
			block = &service.BlockResp{
				ID:              item.BlockID,
				Title:           item.BlockTitle,
				BackgroundImage: item.BlockBackgroundImage.String,
			}
			blockMap[item.BlockID] = block
		}

		content := service.ContentResp{
			Rating:      int(item.ContentMappingRating.Int64),
			Name:        item.ContentName,
			ContentType: item.ContentType,
			Content:     json.RawMessage(item.ContentData),
		}

		block.Contents = append(block.Contents, content)
	}

	var resp []service.BlockResp
	for _, block := range blockMap {
		resp = append(resp, *block)
	}

	return resp, nil
}

func (s *UserService) FindBlockListWithCache(ctx context.Context, opts api.Opts) ([]service.BlockResp, error) {
	var resp []service.BlockResp
	key := fmt.Sprintf("blockListUser:%d", opts.UserID)

	val, err := s.rd.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			resp, err = s.FindBlockList(ctx, opts)
			if err != nil {
				return nil, err
			}
			cacheData, err := json.Marshal(resp)
			if err != nil {
				return nil, err
			}
			if err := s.rd.SetNX(ctx, key, cacheData, time.Hour).Err(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if err := json.NewDecoder(bytes.NewReader(val)).Decode(&resp); err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func (s *UserService) FindBlockIDAndTitleList(ctx context.Context, opts api.Opts) ([]service.BlockResp, error) {
	condition := &repositoryImpl.Condition{
		ScenarioUserID: opts.UserID,
	}
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

func (s *UserService) FindBlockByIDAndOrTitle(ctx context.Context, opts api.Opts) (*service.BlockResp, error) {
	condition := &repositoryImpl.Condition{
		ScenarioUserID: opts.UserID,
		BlockID:        opts.BlockID,
		BlockTitle:     opts.BlockTitle,
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

func (s *UserService) FindBlockBWithoutContentData(ctx context.Context, opts api.Opts) (*service.BlockResp, error) {
	condition := &repositoryImpl.Condition{
		ScenarioUserID: opts.UserID,
		BlockID:        opts.BlockID,
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

func (s *UserService) FindBlockContentByIDAndOrTitleAndOrContentType(ctx context.Context, opts api.Opts) (*service.BlockResp, error) {
	condition := &repositoryImpl.Condition{
		ScenarioUserID: opts.UserID,
		BlockID:        opts.BlockID,
		ContentID:      opts.ContentID,
		ContentName:    opts.ContentName,
		ContentTypeID:  opts.ContentTypeID,
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
