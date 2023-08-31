package postgres

import (
	"fmt"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
)

type Segment struct {
	SegmentID uint64 `gorm:"primary_key"`
	Slug      string
	Percent   *int `gorm:"null"`
}

func (Segment) TableName(schemaName, tableName string) string {
	return fmt.Sprintf("%s.%s", schemaName, tableName)
}

func (s *Segment) FromSegmentModel(segment *models.Segment) {
	s.SegmentID = segment.SegmentID
	s.Slug = segment.Slug
	s.Percent = segment.Percent
}

func (s *Segment) ToSegmentModel() *models.Segment {
	return &models.Segment{
		SegmentID: s.SegmentID,
		Slug:      s.Slug,
		Percent:   s.Percent,
	}
}

type Users2Segments struct {
	UserID    uint64
	SegmentID uint64
	Until     *string
}

func (Users2Segments) TableName(schemaName, tableName string) string {
	return fmt.Sprintf("%s.%s", schemaName, tableName)
}
