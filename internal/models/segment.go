package models

import (
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
	"time"
)

type Segment struct {
	SegmentID uint64 `json:"segmentID"`
	Slug      string `json:"slug"`
	Percent   *int   `json:"percent"`
}

type FormSegment struct {
	Slug    string `json:"slug" validate:"required"`
	Percent *int   `json:"percent"`
}

func (form *FormSegment) Validate() error {
	if form.Percent != nil {
		if *form.Percent < errors.MinPercent || *form.Percent > errors.MaxPercent {
			return errors.ErrPercentIsInvalid
		}
	}

	return nil
}

type SegmentResponse struct {
	Segment Segment `json:"segment"`
}

type SegmentsResponse struct {
	Segments []Segment `json:"segments"`
	Count    int       `json:"count"`
}

type AddUserToSegment struct {
	SegmentSlug string  `json:"segmentSlug"`
	SegmentID   uint64  `json:"-"`
	Until       *string `json:"until"`
}

type FormEditSegments struct {
	SegmentsToAdd    []AddUserToSegment `json:"segmentsToAdd"`
	SegmentsToRemove []string           `json:"segmentsToRemove"`
}

func (form *FormEditSegments) Validate() error {
	for idx, segment := range form.SegmentsToAdd {
		if segment.Until != nil {
			timeValue, err := time.Parse("2006-01-02 15:04", *segment.Until)
			if err != nil {
				return errors.ErrUntilIsInvalid
			}
			subtractedTime := timeValue.Add(-3 * time.Hour).Format("2006-01-02 15:04")
			form.SegmentsToAdd[idx].Until = &subtractedTime
		}
	}

	return nil
}
