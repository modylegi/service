package handlers

import (
	"github.com/modylegi/service/internal/api"
	"github.com/modylegi/service/internal/domain/service"
	"github.com/modylegi/service/pkg/encoding"
	"github.com/rs/zerolog"
	"net/http"
)

type UserHandler struct {
	cache         bool
	log           *zerolog.Logger
	userSvc       service.UserService
	validationSvc service.ValidationService
}

func NewUserHandler(
	cache bool,
	log *zerolog.Logger,
	userSvc service.UserService,
	validationSvc service.ValidationService,
) *UserHandler {
	return &UserHandler{
		cache:         cache,
		log:           log,
		userSvc:       userSvc,
		validationSvc: validationSvc,
	}
}

// AllBlocksHandler godoc
//
//	@Summary		Получение всего контента по сценарию пользователя.
//	@Description	Метод получение всего контента, который доступен пользователю по его сценарию.
//	@Tags			User
//	@Produce		json
//	@Param			user_id	path	integer	true	"id пользователя"
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/scenario/{user_id} [get]
func (h *UserHandler) AllBlocksHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userIDString := r.PathValue("user_id")

	// Параметр user_id не указан
	if userIDString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoUserID)
	}

	// Указан несуществующий/невалидный user_id
	userIDInt, err := h.validationSvc.UserID(ctx, userIDString)
	if err != nil {
		return err
	}

	opts := api.Opts{UserID: userIDInt}

	// Пользователь не привязан ни к какому сценарию
	if err := h.validationSvc.LinkedScenarios(ctx, userIDInt); err != nil {
		return err
	}

	var res []service.BlockResp
	if h.cache {
		res, err = h.userSvc.FindBlockListWithCache(ctx, opts)
	} else {
		res, err = h.userSvc.FindBlockList(ctx, opts)
	}
	if err != nil {
		return err
	}

	if len(res) > 0 {
		return encoding.Encode(w, http.StatusOK, res)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// AllBlocksHandlerIDAndTitle godoc
//
//	@Summary		Получение списка id и названий блоков, доступных пользователю по сценарию.
//	@Description	Метод получения списка id и названий блоков, доступных пользователю по сценарию.
//	@Tags			User
//	@Produce		json
//	@Param			user_id	path	integer	true	"id пользователя"
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/scenario/list/{user_id} [get]
func (h *UserHandler) AllBlocksHandlerIDAndTitle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userIDString := r.PathValue("user_id")

	// Параметр user_id не указан
	if userIDString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoUserID)
	}

	// Указан несуществующий/невалидный user_id
	userIDInt, err := h.validationSvc.UserID(ctx, userIDString)
	if err != nil {
		return err
	}

	opts := api.Opts{UserID: userIDInt}

	// Пользователь не привязан ни к какому сценарию
	if err := h.validationSvc.LinkedScenarios(ctx, userIDInt); err != nil {
		return err
	}

	res, err := h.userSvc.FindBlockIDAndTitleList(ctx, opts)
	if err != nil {
		return err
	}

	if len(res) > 0 {
		return encoding.Encode(w, http.StatusOK, res)
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// BlockByIDAndOrTitleParam godoc
//
//	@Summary		Получение блока контента по id и/или названию для пользователя.
//	@Description	Метод получение блока контента определенного типа по id и/или названию.
//	@Tags			User
//	@Produce		json
//	@Param			user_id		path	integer	true	"id пользователя"
//	@Param			block_id	query	integer	false	"id блока"
//	@Param			name		query	string	false	"название блока"
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/scenario/blocks/{user_id} [get]
func (h *UserHandler) BlockByIDAndOrTitleParam(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userIDString := r.PathValue("user_id")

	// Параметр user_id не указан
	if userIDString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoUserID)
	}

	params := r.URL.Query()
	blockIDString := params.Get("block_id")
	blockTitleString := params.Get("name")

	// Не указан ни параметр block_id, ни параметр name
	if blockIDString == "" && blockTitleString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoBlockIDAndTitle)
	}

	// Указан несуществующий/невалидный user_id
	userIDInt, err := h.validationSvc.UserID(ctx, userIDString)
	if err != nil {
		return err
	}

	opts := api.Opts{
		UserID: userIDInt,
	}

	var blockIDInt int

	// Указан несуществующий/невалидный block_id
	if blockIDString != "" {
		blockIDInt, err = h.validationSvc.BlockID(ctx, blockIDString)
		if err != nil {
			return err
		}

	}

	// Указан несуществующий/невалидный name
	if blockTitleString != "" {
		blockTitleString, err = h.validationSvc.BlockTitle(ctx, blockTitleString)
		if err != nil {
			return err
		}

	}

	if blockIDInt != 0 {
		opts.BlockID = blockIDInt
		// Указан параметр  block_id, который не соответствует сценарию пользователя
		if err := h.validationSvc.LinkedScenarioBlock(ctx, opts); err != nil {
			return err
		}
	} else {
		opts.BlockTitle = blockTitleString
		// Указан параметр  name, который не соответствует сценарию пользователя
		if err := h.validationSvc.LinkedScenarioBlock(ctx, opts); err != nil {
			return err
		}
	}

	res, err := h.userSvc.FindBlockByIDAndOrTitle(ctx, opts)
	if err != nil {
		return err
	}
	// Указан параметр  block_id или параметр name, который не соответствует сценарию пользователя
	if res != nil {
		if blockIDInt != 0 && blockTitleString != "" && res.Title != blockTitleString {
			return api.NewError(http.StatusBadRequest, api.ErrParamsDoNotMatchEachOther)
		}
		return encoding.Encode(w, http.StatusOK, res)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// BlockByID godoc
//
//	@Summary		Получение списка id и названий элементов контента для блока.
//	@Description	Метод получения списка id и названий элементов контента для блока.
//	@Tags			User
//	@Produce		json
//	@Param			user_id		path	integer	true	"id пользователя"
//	@Param			block_id	path	integer	true	"id блока"
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/scenario/block/{user_id}/list/{block_id} [get]
func (h *UserHandler) BlockByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userIDString := r.PathValue("user_id")
	blockIDString := r.PathValue("block_id")

	// Параметр user_id не указан
	if userIDString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoUserID)
	}

	// Параметр block_id не указан
	if blockIDString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoBlockID)
	}

	// Указан несуществующий/невалидный user_id
	userIDInt, err := h.validationSvc.UserID(ctx, userIDString)
	if err != nil {
		return err
	}

	// Указан несуществующий/невалидный block_id
	blockIDInt, err := h.validationSvc.BlockID(ctx, blockIDString)
	if err != nil {
		return err
	}

	opts := api.Opts{
		UserID: userIDInt,
	}

	// Пользователь не привязан ни к какому сценарию
	if err := h.validationSvc.LinkedScenarios(ctx, userIDInt); err != nil {
		return err
	}

	opts.BlockID = blockIDInt

	// Указан параметр  block_id, который не соответствует сценарию пользователя
	if err := h.validationSvc.LinkedScenarioBlock(ctx, opts); err != nil {
		return err
	}

	res, err := h.userSvc.FindBlockBWithoutContentData(ctx, opts)
	if err != nil {
		return err
	}

	if res != nil {
		return encoding.Encode(w, http.StatusOK, res)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Content godoc
//
//	@Summary		Получение элемента контента из блока по id, названию или типу контента.
//	@Description	Метод получения элемента контента из блока по id, названию или типу контента.
//	@Tags			User
//	@Produce		json
//	@Param			user_id			path	integer	true	"id пользователя"
//	@Param			block_id		path	integer	true	"id блока"
//	@Param			content_id		query	integer	false	"id контента"
//	@Param			name			query	string	false	"название элемента контента"
//	@Param			content_type	query	integer	false	"тип контента (1 - баннер, 2 - истории, 3 - тесты)"
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/scenario/block/{user_id}/{block_id} [get]
func (h *UserHandler) Content(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userIDString := r.PathValue("user_id")
	blockIDString := r.PathValue("block_id")

	// Параметр user_id не указан
	if userIDString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoUserID)
	}

	// Параметр block_id не указан
	if blockIDString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoBlockID)
	}

	// Указан несуществующий/невалидный user_id
	userIDInt, err := h.validationSvc.UserID(ctx, userIDString)
	if err != nil {
		return err
	}

	// Указан несуществующий/невалидный block_id
	blockIDInt, err := h.validationSvc.BlockID(ctx, blockIDString)
	if err != nil {
		return err
	}

	params := r.URL.Query()
	contentIDString := params.Get("content_id")
	contentNameString := params.Get("name")
	contentTypeString := params.Get("content_type")
	// Не указан ни один из параметров content_id, name, content_type
	if contentIDString == "" && contentNameString == "" && contentTypeString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoContentIDAndNameAndType)
	}

	var contentIDInt int
	var contentTypeInt int

	if contentIDString != "" {
		// Указан несуществующий/невалидный content_id
		contentIDInt, err = h.validationSvc.ContentID(ctx, contentIDString)
		if err != nil {
			return err
		}
	}

	if contentNameString != "" {
		// Указан несуществующий/невалидный name
		contentNameString, err = h.validationSvc.ContentName(ctx, contentNameString)
		if err != nil {
			return err
		}
	}

	if contentTypeString != "" {
		// Указан несуществующий/невалидный content_type
		contentTypeInt, err = h.validationSvc.ContentTypeID(ctx, contentTypeString)
		if err != nil {
			return err
		}
	}

	opts := api.Opts{
		UserID:  userIDInt,
		BlockID: blockIDInt,
	}

	// Указан параметр  block_id, который не соответствует сценарию пользователя
	if err := h.validationSvc.LinkedScenarioBlock(ctx, opts); err != nil {
		return err
	}

	if contentIDInt != 0 {
		opts.ContentID = contentIDInt
	} else {
		opts.ContentName = contentNameString
		opts.ContentTypeID = contentTypeInt
	}

	res, err := h.userSvc.FindBlockContentByIDAndOrTitleAndOrContentType(ctx, opts)
	if err != nil {
		return err
	}

	if res != nil {
		// Параметры content_id, name, content_type указаны
		// одновременно (возможны комбинации),
		// и при этом эти параметры друг другу не соответствуют (для баннера с таким названием другой id, и наоборот)
		if contentIDInt != 0 && contentNameString != "" && contentTypeInt != 0 && (res.Contents[0].Name != contentNameString || res.Contents[0].ContentType != contentTypeInt) {
			return api.NewError(http.StatusBadRequest, api.ErrParamsDoNotMatchEachOther)
		}
		return encoding.Encode(w, http.StatusOK, res)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
