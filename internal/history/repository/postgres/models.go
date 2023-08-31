package postgres

import (
	"fmt"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
)

type History struct {
	RecordID    uint64
	UserID      uint64
	SegmentSlug string
	Operation   string
	Datetime    string
}

func (History) TableName(schemaName, tableName string) string {
	return fmt.Sprintf("%s.%s", schemaName, tableName)
}

func (s *History) ToHistoryModel() *models.History {
	return &models.History{
		UserID:      s.UserID,
		SegmentSlug: s.SegmentSlug,
		Operation:   s.Operation,
		Datetime:    s.Datetime,
	}
}
