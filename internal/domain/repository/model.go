package repository

import (
	"database/sql"
	"time"
)

type User struct {
	UserID        int       `db:"user_id"`
	Username      string    `db:"user_name"`
	UserPassword  string    `db:"user_password"`
	UserCreatedAt time.Time `db:"user_created_at"`
	UserDeleted   bool      `db:"user_deleted"`
}

type Block struct {
	BlockID              int            `db:"block_id"`
	BlockKey             string         `db:"block_key"`
	BlockTitle           string         `db:"block_title"`
	BlockBackgroundImage sql.NullString `db:"block_background_image"`
	BlockCreatedAt       time.Time      `db:"block_created_at"`
	BlockDeleted         bool           `db:"block_deleted"`
}

type Content struct {
	ContentID        int       `db:"content_id"`
	ContentName      string    `db:"content_name"`
	ContentType      int       `db:"content_type"`
	ContentData      []byte    `db:"content_data"`
	ContentCreatedAt time.Time `db:"content_created_at"`
	ContentDeleted   bool      `db:"content_deleted"`
	ContentUpdatedAt time.Time `db:"content_updated_at"`
}

type BlockContent struct {
	Block
	Content
}

type ContentMapping struct {
	ContentMappingID             int           `db:"content_mapping_id"`
	ContentMappingContentBlockID int           `db:"content_block_id"`
	ContentMappingContentID      int           `db:"content_id"`
	ContentMappingRating         sql.NullInt64 `db:"content_mapping_rating"`
	ContentMappingDeleted        bool          `db:"content_mapping_deleted"`
}

type BlockContentContentMapping struct {
	Block
	Content
	ContentMapping
}

type ContentContentMapping struct {
	Content
	ContentMapping
}

type ContentTypeDic struct {
	ContentTypeDicID          int    `db:"content_type_dic_id"`
	ContentTypeDicName        string `db:"content_type_dic_name"`
	ContentTypeDicDescription string `db:"content_type_dic_description"`
}

type TemplateContent struct {
	TemplateContentID        int       `db:"template_content_id"`
	TemplateContentName      string    `db:"template_content_name"`
	TemplateContent          []byte    `db:"template_content_data"`
	TemplateContentType      int       `db:"template_content_type"`
	TemplateContentCreatedAt time.Time `db:"template_content_created_at"`
	TemplateContentDeleted   bool      `db:"template_content_deleted"`
}
