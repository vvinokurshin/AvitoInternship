package postgres

import (
	pkgErrors "github.com/pkg/errors"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/history/repository"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type historyRepo struct {
	cfg *config.Config
	db  *gorm.DB
}

func New(cfg *config.Config, db *gorm.DB) repository.RepositoryI {
	return &historyRepo{
		cfg: cfg,
		db:  db,
	}
}

func (repo *historyRepo) SelectRecordsByDate(year int, month time.Month) ([]models.History, error) {
	var dbRecords []History
	datetime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	tx := repo.db.Table(History{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBHistoryTableName)).Omit("record_id").
		Where("date_trunc('month', datetime) = ?", datetime).Find(&dbRecords)
	if err := tx.Error; err != nil {
		return []models.History{}, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	result := make([]models.History, len(dbRecords))
	for idx, dbRecord := range dbRecords {
		result[idx] = *dbRecord.ToHistoryModel()
	}

	return result, nil
}
