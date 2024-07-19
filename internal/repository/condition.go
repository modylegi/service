package repository

import (
	"fmt"
	"strings"
)

type Condition struct {
	UserID                int
	ScenarioUserID        int
	BlockID               int
	BlockTitle            string
	ContentID             int
	ContentName           string
	ContentTypeID         int
	TemplateID            int
	TemplateName          string
	TemplateContentTypeID int
}

func (c *Condition) Args() []any {
	var list []any
	if c.UserID != 0 {
		list = append(list, c.UserID)
	}
	if c.ScenarioUserID != 0 {
		list = append(list, c.ScenarioUserID)
	}
	if c.BlockID != 0 {
		list = append(list, c.BlockID)
	}
	if c.BlockTitle != "" {
		list = append(list, c.BlockTitle)
	}
	if c.ContentID != 0 {
		list = append(list, c.ContentID)
	}
	if c.ContentName != "" {
		list = append(list, c.ContentName)
	}
	if c.ContentTypeID != 0 {
		list = append(list, c.ContentTypeID)
	}
	if c.TemplateID != 0 {
		list = append(list, c.TemplateID)
	}
	if c.TemplateName != "" {
		list = append(list, c.TemplateName)
	}
	if c.TemplateContentTypeID != 0 {
		list = append(list, c.TemplateContentTypeID)
	}
	return list
}

func (c *Condition) String() string {
	var conditions []string
	if c.UserID != 0 {
		conditions = append(conditions, fmt.Sprintf("u.id = $%d", len(conditions)+1))
	}
	if c.ScenarioUserID != 0 {
		conditions = append(conditions, fmt.Sprintf("su.user_id = $%d", len(conditions)+1))
	}
	if c.BlockID != 0 {
		conditions = append(conditions, fmt.Sprintf("b.id = $%d", len(conditions)+1))
	}
	if c.BlockTitle != "" {
		conditions = append(conditions, fmt.Sprintf("b.title = $%d", len(conditions)+1))
	}
	if c.ContentID != 0 {
		conditions = append(conditions, fmt.Sprintf("c.id = $%d", len(conditions)+1))
	}
	if c.ContentName != "" {
		conditions = append(conditions, fmt.Sprintf("c.name = $%d", len(conditions)+1))
	}
	if c.ContentTypeID != 0 {
		conditions = append(conditions, fmt.Sprintf("c.content_type_id = $%d", len(conditions)+1))
	}
	if c.TemplateID != 0 {
		conditions = append(conditions, fmt.Sprintf("tc.id = $%d", len(conditions)+1))
	}
	if c.TemplateName != "" {
		conditions = append(conditions, fmt.Sprintf("tc.name = $%d", len(conditions)+1))
	}
	if c.TemplateContentTypeID != 0 {
		conditions = append(conditions, fmt.Sprintf("tc.content_type_id = $%d", len(conditions)+1))
	}
	if len(conditions) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(conditions, " AND ")
}

func (c *Condition) GetScenarioUserID() int {
	return c.ScenarioUserID
}
