package models

import (
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
	"time"
)

type History struct {
	UserID      uint64 `json:"userID"`
	SegmentSlug string `json:"segmentSlug"`
	Operation   string `json:"operation"`
	Datetime    string `json:"datetime"`
}

type FormHistory struct {
	Year  int        `json:"year"`
	Month time.Month `json:"month"`
}

func (form *FormHistory) Validate() error {
	if form.Year < errors.MinYear || form.Year > errors.MaxYear {
		return errors.ErrYearIsInvalid
	}

	if form.Month < errors.MinMonth || form.Month > errors.MaxMonth {
		return errors.ErrMonthIsRequired
	}

	return nil
}
