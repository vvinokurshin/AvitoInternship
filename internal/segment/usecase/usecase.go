package usecase

import (
	pkgErr "github.com/pkg/errors"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	segmentRepository "github.com/vvinokurshin/AvitoInternship/internal/segment/repository"
	userRepository "github.com/vvinokurshin/AvitoInternship/internal/user/repository"
	"github.com/vvinokurshin/AvitoInternship/pkg"
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
)

//go:generate mockgen -destination=./mocks/usecase.go -source=./usecase.go -package=mocks

type UseCaseI interface {
	CreateSegment(form models.FormSegment) (*models.Segment, error)
	DeleteSegment(slug string) error
	GetSegmentBySlug(slug string) (*models.Segment, error)
	GetUserSegments(userID uint64) ([]models.Segment, error)
	EditUserSegments(userID uint64, segmentsToAdd []models.AddUserToSegment, segmentsToRemove []string) ([]models.Segment, error)
}

type UseCase struct {
	cfg         *config.Config
	segmentRepo segmentRepository.RepositoryI
	userRepo    userRepository.RepositoryI
}

func New(cfg *config.Config, segmentRepo segmentRepository.RepositoryI, userRepo userRepository.RepositoryI) UseCaseI {
	return &UseCase{
		cfg:         cfg,
		segmentRepo: segmentRepo,
		userRepo:    userRepo,
	}
}

func (uc *UseCase) CreateSegment(form models.FormSegment) (*models.Segment, error) {
	_, err := uc.segmentRepo.SelectSegmentBySlug(form.Slug)
	if err != errors.ErrSegmentNotFound {
		return nil, errors.ErrSegmentExists
	}

	segment := &models.Segment{
		Slug:    form.Slug,
		Percent: form.Percent,
	}

	segmentID, err := uc.segmentRepo.InsertSegment(segment)
	if err != nil {
		return nil, pkgErr.Wrap(err, "insert segment")
	}

	segment.SegmentID = segmentID

	if segment.Percent != nil {
		userIDs, err := uc.userRepo.SelectUserIDs()
		if err != nil {
			return nil, pkgErr.Wrap(err, "get user IDs")
		}

		IDsToAdd := pkg.PercentageIDs(userIDs, *segment.Percent)
		err = uc.segmentRepo.InsertUsersToSegment(segmentID, IDsToAdd)
		if err != nil {
			return nil, pkgErr.Wrap(err, "insert users to segment")
		}
	}

	return segment, nil
}

func (uc *UseCase) DeleteSegment(slug string) error {
	_, err := uc.segmentRepo.SelectSegmentBySlug(slug)
	if err != nil {
		return pkgErr.Wrap(err, "select segment by slug")
	}

	err = uc.segmentRepo.DeleteSegment(slug)
	if err != nil {
		return pkgErr.Wrap(err, "delete segment")
	}

	return nil
}

func (uc *UseCase) GetSegmentBySlug(slug string) (*models.Segment, error) {
	segment, err := uc.segmentRepo.SelectSegmentBySlug(slug)
	if err != nil {
		return nil, pkgErr.Wrap(err, "select segment by slug")
	}

	return segment, nil
}

func (uc *UseCase) GetUserSegments(userID uint64) ([]models.Segment, error) {
	_, err := uc.userRepo.SelectUserByID(userID)
	if err != nil {
		return []models.Segment{}, pkgErr.Wrap(err, "select user by ID")
	}

	segments, err := uc.segmentRepo.SelectSegmentsByUser(userID)
	if err != nil {
		return []models.Segment{}, pkgErr.Wrap(err, "select segments by userID")
	}

	return segments, nil
}

func (uc *UseCase) EditUserSegments(userID uint64, segmentsToAdd []models.AddUserToSegment, segmentsToRemove []string) ([]models.Segment, error) {
	_, err := uc.userRepo.SelectUserByID(userID)
	if err != nil {
		return []models.Segment{}, pkgErr.Wrap(err, "select user by ID")
	}

	for idx, currentSegment := range segmentsToAdd {
		segment, err := uc.segmentRepo.SelectSegmentBySlug(currentSegment.SegmentSlug)
		if err != nil {
			return []models.Segment{}, pkgErr.Wrap(err, "select segment by slug")
		}

		segmentsToAdd[idx].SegmentID = segment.SegmentID
	}

	segmentIDsToRemove := make([]uint64, len(segmentsToRemove))
	for idx, segmentSlug := range segmentsToRemove {
		segment, err := uc.segmentRepo.SelectSegmentBySlug(segmentSlug)
		if err != nil {
			return []models.Segment{}, pkgErr.Wrap(err, "select segment by slug")
		}

		segmentIDsToRemove[idx] = segment.SegmentID
	}

	if len(segmentsToAdd) != 0 {
		err = uc.segmentRepo.InsertSegmentsToUser(userID, segmentsToAdd)
		if err != nil {
			return []models.Segment{}, pkgErr.Wrap(err, "insert segments to user")
		}
	}

	if len(segmentIDsToRemove) != 0 {
		err = uc.segmentRepo.DeleteSegmentsFromUser(userID, segmentIDsToRemove)
		if err != nil {
			return []models.Segment{}, pkgErr.Wrap(err, "delete segments from user")
		}
	}

	segments, err := uc.segmentRepo.SelectSegmentsByUser(userID)
	if err != nil {
		return []models.Segment{}, pkgErr.Wrap(err, "select segments by userID")
	}

	return segments, nil
}
