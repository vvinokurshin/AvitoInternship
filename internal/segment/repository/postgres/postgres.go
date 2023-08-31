package postgres

import (
	pkgErrors "github.com/pkg/errors"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	"github.com/vvinokurshin/AvitoInternship/internal/segment/repository"
	"github.com/vvinokurshin/AvitoInternship/pkg"
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type segmentRepo struct {
	cfg *config.Config
	db  *gorm.DB
}

func New(cfg *config.Config, db *gorm.DB) (repository.RepositoryI, error) {
	segRepo := &segmentRepo{
		cfg: cfg,
		db:  db,
	}

	err := pkg.CronInit("@every 5m", segRepo.ClearExpiredConnections)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "cron init")
	}

	return segRepo, nil
}

func (repo *segmentRepo) InsertSegment(segment *models.Segment) (uint64, error) {
	var dbSegment Segment
	dbSegment.FromSegmentModel(segment)

	tx := repo.db.Table(Segment{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBSegmentTableName)).Create(&dbSegment)
	if err := tx.Error; err != nil {
		return 0, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return dbSegment.SegmentID, nil
}

func (repo *segmentRepo) DeleteSegment(slug string) error {
	tx := repo.db.Table(Segment{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBSegmentTableName)).
		Where("slug = ?", slug).Delete(Segment{})
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (repo *segmentRepo) SelectSegmentBySlug(slug string) (*models.Segment, error) {
	var dbSegment Segment

	tx := repo.db.Table(Segment{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBSegmentTableName)).
		Where("slug = ?", slug).Take(&dbSegment)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrSegmentNotFound
		}

		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return dbSegment.ToSegmentModel(), nil
}

func (repo *segmentRepo) SelectSegmentsByUser(userID uint64) ([]models.Segment, error) {
	var dbSegments []Segment
	SegmentsTablename := Segment{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBSegmentTableName)
	U2STableName := Users2Segments{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBU2STableName)

	tx := repo.db.Table(SegmentsTablename).Select(SegmentsTablename+".*").Joins("JOIN "+U2STableName+
		" using(segment_id)").Where("user_id = ?", userID).Find(&dbSegments)
	if err := tx.Error; err != nil {
		return []models.Segment{}, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	result := make([]models.Segment, len(dbSegments))
	for idx, dbSegment := range dbSegments {
		result[idx] = *dbSegment.ToSegmentModel()
	}

	return result, nil
}

func (repo *segmentRepo) InsertSegmentsToUser(userID uint64, segments []models.AddUserToSegment) error {
	dbU2S := make([]Users2Segments, len(segments))
	for idx, segment := range segments {
		dbU2S[idx].UserID = userID
		dbU2S[idx].SegmentID = segment.SegmentID
		dbU2S[idx].Until = segment.Until
	}

	tx := repo.db.Table(Users2Segments{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBU2STableName)).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "segment_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"until"}),
		}).Create(&dbU2S)
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (repo *segmentRepo) DeleteSegmentsFromUser(userID uint64, segmentIDs []uint64) error {
	tx := repo.db.Table(Users2Segments{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBU2STableName)).Where(
		"user_id = ? AND segment_id IN ?", userID, segmentIDs).Delete(&Users2Segments{})
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (repo *segmentRepo) InsertUsersToSegment(segmentID uint64, userIDs []uint64) error {
	dbU2S := make([]Users2Segments, len(userIDs))
	for idx, userID := range userIDs {
		dbU2S[idx].UserID = userID
		dbU2S[idx].SegmentID = segmentID
	}

	tx := repo.db.Table(Users2Segments{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBU2STableName)).
		Clauses(clause.OnConflict{DoNothing: true}).Create(&dbU2S)
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (repo *segmentRepo) ClearExpiredConnections() {
	repo.db.Raw("SELECT delete_old_accesses()").Rows()
}
