package service

import "encoding/json"

type BlockResp struct {
	ID              int           `json:"id,omitempty"`
	Title           string        `json:"title,omitempty"`
	BackgroundImage string        `json:"background_image,omitempty"`
	Contents        []ContentResp `json:"contents,omitempty"`
}

type ContentResp struct {
	ID          int             `json:"id,omitempty"`
	Rating      int             `json:"rating,omitempty"`
	Name        string          `json:"name"`
	ContentType int             `json:"content_type"`
	Content     json.RawMessage `json:"content,omitempty"`
}

type TemplateResp struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	ContentType int             `json:"content_type"`
	Content     json.RawMessage `json:"content,omitempty"`
}

type ApiOpts struct {
	UserID                int
	BlockID               int
	BlockTitle            string
	ContentID             int
	ContentName           string
	ContentTypeID         int
	TemplateID            int
	TemplateName          string
	TemplateContentTypeID int
}
