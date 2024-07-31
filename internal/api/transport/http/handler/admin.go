package handler

import (
	"net/http"

	"github.com/modylegi/service/internal/api"
	"github.com/modylegi/service/internal/domain/service"
	"github.com/modylegi/service/pkg/encoding"
	"github.com/rs/zerolog"
)

type AdminHandler struct {
	log           *zerolog.Logger
	adminSvc      service.AdminService
	validationSvc service.ValidationService
}

func NewAdminHandler(
	log *zerolog.Logger,
	adminSvc service.AdminService,
	validationSvc service.ValidationService,
) *AdminHandler {
	return &AdminHandler{
		log:           log,
		adminSvc:      adminSvc,
		validationSvc: validationSvc,
	}
}

// BlockIDAndTitleList godoc
//	@Summary		Получение списка id и названий всех блоков.
//	@Description	Метод получения списка id и названий всех блоков.
//	@Tags			Admin
//	@Produce		json
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/block/list [get]
//	@Security		BearerAuth
func (h *AdminHandler) BlockIDAndTitleList(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	res, err := h.adminSvc.FindBlockIDAndTitleList(ctx)
	if err != nil {
		return err
	}
	if len(res) > 0 {
		return encoding.Encode(w, http.StatusOK, res)
	}
	// 204 - пустое тело ответа, блоков нет в базе
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// BlockByIDAndOrTitle godoc
//	@Summary		Получение блока контента по id и/или названию.
//	@Description	Метод получение блока контента по id и/или названию.
//	@Tags			Admin
//	@Produce		json
//	@Param			block_id	query	integer	false	"id блока"
//	@Param			name		query	string	false	"название блока"
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/block [get]
//	@Security		BearerAuth
func (h *AdminHandler) BlockByIDAndOrTitle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := r.URL.Query()
	blockIDString := params.Get("block_id")
	blockTitleString := params.Get("name")

	if blockIDString == "" && blockTitleString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoBlockIDAndTitle)
	}

	var blockIDInt int
	var err error

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

	var opts service.ApiOpts
	if blockIDInt != 0 {
		opts.BlockID = blockIDInt
	} else {
		opts.BlockTitle = blockTitleString
	}

	res, err := h.adminSvc.FindBlockByIDAndOrTitle(ctx, opts)
	if err != nil {
		return err
	}

	if res != nil {
		// Указан и параметр  block_id, и параметр name, и при этом эти параметры друг другу не соответствуют (для баннера с таким названием другой id, и наоборот)
		if blockIDInt != 0 && blockTitleString != "" && res.Title != blockTitleString {
			return api.NewError(http.StatusBadRequest, api.ErrParamsDoNotMatchEachOther)
		}
		return encoding.Encode(w, http.StatusOK, res)
	}
	// По заданному запросу ничего не найдено
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// BlockWithoutContentData godoc
//	@Summary		Получение списка id и названий элементов контента для блока.
//	@Description	Метод получения списка id и названий элементов контента для блока.
//	@Tags			Admin
//	@Produce		json
//	@Param			block_id	path	integer	true	"id блока"
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/block/{block_id}/list [get]
//	@Security		BearerAuth
func (h *AdminHandler) BlockWithoutContentData(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	blockIDString := r.PathValue("block_id")

	// Параметр block_id не указан
	if blockIDString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoBlockID)
	}

	// Указан несуществующий/невалидный block_id
	blockIDInt, err := h.validationSvc.BlockID(ctx, blockIDString)
	if err != nil {
		return err
	}

	opts := service.ApiOpts{
		BlockID: blockIDInt,
	}

	res, err := h.adminSvc.FindBlockWithoutContentData(ctx, opts)
	if err != nil {
		return err
	}

	if res != nil {
		return encoding.Encode(w, http.StatusOK, res)
	}
	// Для данного block_id нет контента - 204 - пустое тело ответа
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Content godoc
//	@Summary		Получение элементов контента из блока по id, названию или типу контента.
//	@Description	Метод получения элемента контента из блока по id, названию или типу контента.
//	@Tags			Admin
//	@Produce		json
//	@Param			block_id		path	integer	true	"id блока"
//	@Param			content_id		query	integer	false	"id контента"
//	@Param			name			query	string	false	"название элемента контента"
//	@Param			content_type	query	integer	false	"тип контента (1 - баннер, 2 - истории, 3 - тесты)"
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/block/{block_id} [get]
//	@Security		BearerAuth
func (h *AdminHandler) Content(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	blockIDString := r.PathValue("block_id")

	// Параметр block_id не указан
	if blockIDString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoBlockID)
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

	// Указан несуществующий/невалидный content_id
	if contentIDString != "" {
		contentIDInt, err = h.validationSvc.ContentID(ctx, contentIDString)
		if err != nil {
			return err
		}
	}

	// Указан несуществующий/невалидный name
	if contentNameString != "" {
		contentNameString, err = h.validationSvc.ContentName(ctx, contentNameString)
		if err != nil {
			return err
		}
	}

	// Указан несуществующий/невалидный content_type
	if contentTypeString != "" {
		contentTypeInt, err = h.validationSvc.ContentTypeID(ctx, contentTypeString)
		if err != nil {
			return err
		}
	}

	var opts service.ApiOpts
	opts.BlockID = blockIDInt
	if contentIDInt != 0 {
		opts.ContentID = contentIDInt
	} else {
		opts.ContentName = contentNameString
		opts.ContentTypeID = contentTypeInt
	}

	res, err := h.adminSvc.FindBlockContentByIDAndOrTitleAndOrContentType(ctx, opts)
	if err != nil {
		return err
	}

	if res != nil {
		// Параметры content_id, name, content_type указаны одновременно (возможны комбинации), и при этом эти параметры друг другу не соответствуют (для баннера с таким названием другой id, и наоборот)
		if contentIDInt != 0 && contentNameString != "" && contentTypeInt != 0 && (res.Contents[0].Name != contentNameString || res.Contents[0].ContentType != contentTypeInt) {
			return api.NewError(http.StatusBadRequest, api.ErrParamsDoNotMatchEachOther)
		}
		return encoding.Encode(w, http.StatusOK, res)
	}
	// По заданным параметрам ничего не найдено - возвращаем пустое тело ответа и 204 No Content
	w.WriteHeader(http.StatusNoContent)
	return nil

}

// TemplateList godoc
//	@Summary		Получение списка id, названий и типов всех шаблонов.
//	@Description	Метод получения списка id, названий и типов всех шаблонов.
//	@Tags			Admin
//	@Produce		json
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/template/list [get]
//	@Security		BearerAuth
func (h *AdminHandler) TemplateList(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	res, err := h.adminSvc.FindTemplateList(ctx)
	if err != nil {
		return err
	}
	if res != nil {
		return encoding.Encode(w, http.StatusOK, res)
	}

	// 204 - пустое тело ответа, шаблонов нет в базе
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// TemplateByIDNameType godoc
//	@Summary		Получение шаблона по id, названию или типу.
//	@Description	Метод получения шаблона по id, названию или типу контента.
//	@Tags			Admin
//	@Produce		json
//	@Param			template_id		query	integer	false	"id шаблона"
//	@Param			name			query	string	false	"название шаблона"
//	@Param			content_type	query	integer	false	"тип шаблона (1 - баннер, 2 - истории, 3 - тесты)"
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/template [get]
//	@Security		BearerAuth
func (h *AdminHandler) TemplateByIDNameType(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := r.URL.Query()
	templateIDString := params.Get("template_id")
	templateNameString := params.Get("name")
	contentTypeString := params.Get("content_type")

	// Не указан ни один из параметров template_id, name, content_type
	if templateIDString == "" && templateNameString == "" && contentTypeString == "" {
		return api.NewError(http.StatusBadRequest, api.ErrNoTemplateIDAndNameAndType)
	}

	var templateIDInt int
	var contentTypeInt int
	var err error

	// Указан несуществующий/невалидный template_id
	if templateIDString != "" {
		templateIDInt, err = h.validationSvc.TemplateID(ctx, templateIDString)
		if err != nil {
			return err
		}
	}

	// Указан несуществующий/невалидный name
	if templateNameString != "" {
		templateNameString, err = h.validationSvc.TemplateName(ctx, templateNameString)
		if err != nil {
			return err
		}

	}

	// Указан несуществующий/невалидный content_type
	if contentTypeString != "" {
		contentTypeInt, err = h.validationSvc.ContentTypeID(ctx, contentTypeString)
		if err != nil {
			return err
		}

	}

	var opts service.ApiOpts
	if templateIDInt != 0 {
		opts.TemplateID = templateIDInt
	} else {
		opts.TemplateName = templateNameString
		opts.TemplateContentTypeID = contentTypeInt
	}

	res, err := h.adminSvc.FindTemplateByIDAndOrNameAndOrContentType(ctx, opts)
	if err != nil {
		return err
	}

	if res != nil {
		// Параметры template_id, name, content_type указаны одновременно (возможны комбинации), и при этом эти параметры друг другу не соответствуют (для шаблона с таким названием другой id, и наоборот)
		if templateIDInt != 0 && templateNameString != "" && contentTypeInt != 0 && (res.Name != templateNameString || res.ContentType != contentTypeInt) {
			return api.NewError(http.StatusBadRequest, api.ErrParamsDoNotMatchEachOther)
		}
		return encoding.Encode(w, http.StatusOK, res)
	}

	// По заданным параметрам ничего не найдено - возвращаем пустое тело ответа и 204 No Content
	w.WriteHeader(http.StatusNoContent)
	return nil
}
