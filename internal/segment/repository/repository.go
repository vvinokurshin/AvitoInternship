package repository

import "github.com/vvinokurshin/AvitoInternship/internal/models"

//go:generate mockgen -destination=./mocks/repository.go -source=./repository.go -package=mocks

type RepositoryI interface {
	InsertSegment(segment *models.Segment) (uint64, error)
	DeleteSegment(slug string) error
	SelectSegmentBySlug(slug string) (*models.Segment, error)
	SelectSegmentsByUser(userID uint64) ([]models.Segment, error)
	InsertSegmentsToUser(userID uint64, segments []models.AddUserToSegment) error
	DeleteSegmentsFromUser(userID uint64, segmentIDs []uint64) error
	InsertUsersToSegment(segmentID uint64, userIDs []uint64) error
}
