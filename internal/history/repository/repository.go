package repository

import (
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	"time"
)

//go:generate mockgen -destination=./mocks/repository.go -source=./repository.go -package=mocks

type RepositoryI interface {
	SelectRecordsByDate(year int, month time.Month) ([]models.History, error)
}
