package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/modylegi/service/internal/domain/service"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/modylegi/service/internal/api"
	"github.com/modylegi/service/internal/domain/repository"
	repositoryImpl "github.com/modylegi/service/internal/repository"
)

type ValidationService struct {
	userRepo         repository.UserRepository
	contentRepo      repository.ContentRepository
	contentTypeRepo  repository.ContentTypeRepository
	blockRepo        repository.BlockRepository
	templateRepo     repository.TemplateRepository
	scenarioUserRepo repository.ScenarioUserRepository
}

func NewValidationService(db *sqlx.DB) *ValidationService {
	return &ValidationService{
		userRepo:         repositoryImpl.NewUserRepository(db),
		contentRepo:      repositoryImpl.NewContentRepository(db),
		contentTypeRepo:  repositoryImpl.NewContentTypeRepository(db),
		blockRepo:        repositoryImpl.NewBlockRepository(db),
		templateRepo:     repositoryImpl.NewTemplateRepository(db),
		scenarioUserRepo: repositoryImpl.NewScenarioUserRepository(db),
	}
}

func (s *ValidationService) UserID(ctx context.Context, userIDString string) (int, error) {
	userIDInt, err := strconv.Atoi(userIDString)
	if err != nil {
		return 0, api.NewError(http.StatusBadRequest, fmt.Errorf("user_id некорректный: %s", userIDString))
	}
	condition := &repositoryImpl.Condition{
		UserID: userIDInt,
	}
	if err := s.userRepo.FindUserID(ctx, condition); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, api.NewError(http.StatusBadRequest, fmt.Errorf("user_id не существует: %d", userIDInt))
		}
		return 0, err
	}
	return userIDInt, nil
}

func (s *ValidationService) BlockID(ctx context.Context, blockIDString string) (int, error) {
	blockIDInt, err := strconv.Atoi(blockIDString)
	if err != nil {
		return 0, api.NewError(http.StatusBadRequest, fmt.Errorf("block_id некорректный: %s", blockIDString))
	}
	condition := &repositoryImpl.Condition{
		BlockID: blockIDInt,
	}
	if err := s.blockRepo.GetBlockID(ctx, condition); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, api.NewError(http.StatusBadRequest, fmt.Errorf("block_id не существует: %d", blockIDInt))
		}
		return 0, err
	}
	return blockIDInt, nil
}

func (s *ValidationService) BlockTitle(ctx context.Context, blockTitleString string) (string, error) {
	condition := &repositoryImpl.Condition{
		BlockTitle: blockTitleString,
	}
	if err := s.blockRepo.GetBlockTitle(ctx, condition); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", api.NewError(http.StatusBadRequest, fmt.Errorf("name не существует: %s", blockTitleString))
		}
		return "", err
	}
	return blockTitleString, nil
}

func (s *ValidationService) ContentTypeID(ctx context.Context, contentTypeIDString string) (int, error) {
	contentTypeIDInt, err := strconv.Atoi(contentTypeIDString)
	if err != nil {
		return 0, api.NewError(http.StatusBadRequest, fmt.Errorf("content_type некорректный: %s", contentTypeIDString))
	}
	condition := &repositoryImpl.Condition{
		ContentTypeID: contentTypeIDInt,
	}
	if err := s.contentTypeRepo.FindContentTypeID(ctx, condition); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, api.NewError(http.StatusBadRequest, fmt.Errorf("content_type не существует: %d", contentTypeIDInt))
		}
		return 0, err
	}
	return contentTypeIDInt, nil
}

func (s *ValidationService) ContentID(ctx context.Context, contentIDString string) (int, error) {
	contentIDInt, err := strconv.Atoi(contentIDString)
	if err != nil {
		return 0, api.NewError(http.StatusBadRequest, fmt.Errorf("content_id некорректный: %s", contentIDString))
	}
	condition := &repositoryImpl.Condition{
		ContentID: contentIDInt,
	}
	if err := s.contentRepo.FindContentID(ctx, condition); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, api.NewError(http.StatusBadRequest, fmt.Errorf("content_id не существует: %d", contentIDInt))
		}
		return 0, err
	}
	return contentIDInt, nil
}

func (s *ValidationService) ContentName(ctx context.Context, contentNameString string) (string, error) {
	condition := &repositoryImpl.Condition{
		ContentName: contentNameString,
	}
	if err := s.contentRepo.FindContentName(ctx, condition); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", api.NewError(http.StatusBadRequest, fmt.Errorf("name не существует: %s", contentNameString))
		}
		return "", err
	}
	return contentNameString, nil
}

func (s *ValidationService) TemplateID(ctx context.Context, templateIDString string) (int, error) {
	templateIDInt, err := strconv.Atoi(templateIDString)
	if err != nil {
		return 0, api.NewError(http.StatusBadRequest, fmt.Errorf("template_id некорректный: %s", templateIDString))
	}
	condition := &repositoryImpl.Condition{
		TemplateID: templateIDInt,
	}
	if err := s.templateRepo.FindTemplateID(ctx, condition); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, api.NewError(http.StatusBadRequest, fmt.Errorf("template_id не существует: %d", templateIDInt))
		}
		return 0, err
	}
	return templateIDInt, nil
}

func (s *ValidationService) TemplateName(ctx context.Context, templateNameString string) (string, error) {
	condition := &repositoryImpl.Condition{
		TemplateName: templateNameString,
	}
	if err := s.templateRepo.FindTemplateName(ctx, condition); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", api.NewError(http.StatusBadRequest, fmt.Errorf("name не существует: %s", templateNameString))
		}
		return "", err
	}
	return templateNameString, nil
}

func (s *ValidationService) LinkedScenarios(ctx context.Context, userID int) error {
	condition := &repositoryImpl.Condition{
		ScenarioUserID: userID,
	}
	if err := s.scenarioUserRepo.GetScenarioByUserID(ctx, condition); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.NewError(http.StatusNoContent, api.ErrNoScenario)
		}
		return err
	}
	return nil
}

func (s *ValidationService) LinkedScenarioBlock(ctx context.Context, opts service.ApiOpts) error {
	condition := &repositoryImpl.Condition{
		ScenarioUserID: opts.UserID,
	}
	if opts.BlockID != 0 {
		condition.BlockID = opts.BlockID

	} else {
		condition.BlockTitle = opts.BlockTitle
	}
	if err := s.blockRepo.GetBlockIDByUserID(ctx, condition); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if condition.BlockID != 0 {
				return api.NewError(http.StatusBadRequest, api.ErrBlockIDNotMatchesUserScenario)
			}
			if condition.BlockTitle != "" {
				return api.NewError(http.StatusBadRequest, api.ErrBlockTitleNotMatchesUserScenario)
			}

		}
		return err
	}
	return nil
}
