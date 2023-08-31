package usecase

import (
	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	mockSegmentRepo "github.com/vvinokurshin/AvitoInternship/internal/segment/repository/mocks"
	mockUserRepo "github.com/vvinokurshin/AvitoInternship/internal/user/repository/mocks"
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
	"testing"
)

func createConfig() *config.Config {
	return new(config.Config)
}

func generateFakeData(data any) {
	faker.SetRandomMapAndSliceMaxSize(10)
	faker.SetRandomMapAndSliceMinSize(1)
	faker.SetRandomStringLength(30)

	faker.FakeData(data)
}

func TestUseCase_CreateSegment(t *testing.T) {
	cfg := createConfig()

	var fakeForm models.FormSegment
	generateFakeData(&fakeForm)
	fakeForm.Percent = nil
	fakeSegment := &models.Segment{
		Slug: fakeForm.Slug,
	}
	fakeSegmentResponse := &models.Segment{
		SegmentID: 1,
		Slug:      fakeForm.Slug,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmentRepo := mockSegmentRepo.NewMockRepositoryI(ctrl)
	userRepo := mockUserRepo.NewMockRepositoryI(ctrl)
	segmentUC := New(cfg, segmentRepo, userRepo)

	segmentRepo.EXPECT().SelectSegmentBySlug(fakeForm.Slug).Return(nil, errors.ErrSegmentNotFound)
	segmentRepo.EXPECT().InsertSegment(fakeSegment).Return(uint64(1), nil)
	response, err := segmentUC.CreateSegment(fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeSegmentResponse, response)
	}
}

func TestUseCase_DeleteSegment(t *testing.T) {
	cfg := createConfig()

	fakeSegment := &models.Segment{
		SegmentID: 1,
		Slug:      "test",
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmentRepo := mockSegmentRepo.NewMockRepositoryI(ctrl)
	userRepo := mockUserRepo.NewMockRepositoryI(ctrl)
	segmentUC := New(cfg, segmentRepo, userRepo)

	segmentRepo.EXPECT().SelectSegmentBySlug(fakeSegment.Slug).Return(fakeSegment, nil)
	segmentRepo.EXPECT().DeleteSegment(fakeSegment.Slug).Return(nil)
	err := segmentUC.DeleteSegment(fakeSegment.Slug)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestUseCase_GetSegmentBySlug(t *testing.T) {
	cfg := createConfig()

	fakeSegment := &models.Segment{
		SegmentID: 1,
		Slug:      "test",
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmentRepo := mockSegmentRepo.NewMockRepositoryI(ctrl)
	userRepo := mockUserRepo.NewMockRepositoryI(ctrl)
	segmentUC := New(cfg, segmentRepo, userRepo)

	segmentRepo.EXPECT().SelectSegmentBySlug(fakeSegment.Slug).Return(fakeSegment, nil)
	response, err := segmentUC.GetSegmentBySlug(fakeSegment.Slug)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeSegment, response)
	}
}

func TestUseCase_GetUserSegments(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	var fakeUserSegments []models.Segment
	generateFakeData(&fakeUser)
	generateFakeData(&fakeUserSegments)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmentRepo := mockSegmentRepo.NewMockRepositoryI(ctrl)
	userRepo := mockUserRepo.NewMockRepositoryI(ctrl)
	segmentUC := New(cfg, segmentRepo, userRepo)

	userRepo.EXPECT().SelectUserByID(fakeUser.UserID).Return(fakeUser, nil)
	segmentRepo.EXPECT().SelectSegmentsByUser(fakeUser.UserID).Return(fakeUserSegments, nil)
	response, err := segmentUC.GetUserSegments(fakeUser.UserID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUserSegments, response)
	}
}

func TestUseCase_EditUserSegments(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)
	segmentsToAdd := []models.AddUserToSegment{
		{
			SegmentSlug: "test",
			SegmentID:   1,
			Until:       nil,
		},
	}
	segmentsToRemove := []string{"tmp"}
	fakeSegments := []models.Segment{
		{
			SegmentID: 1,
			Slug:      "test",
			Percent:   nil,
		},
		{
			SegmentID: 2,
			Slug:      "tmp",
			Percent:   nil,
		},
	}
	fakeUserSegments := []models.Segment{
		{
			SegmentID: fakeSegments[0].SegmentID,
			Slug:      fakeSegments[0].Slug,
			Percent:   fakeSegments[0].Percent,
		},
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmentRepo := mockSegmentRepo.NewMockRepositoryI(ctrl)
	userRepo := mockUserRepo.NewMockRepositoryI(ctrl)
	segmentUC := New(cfg, segmentRepo, userRepo)

	userRepo.EXPECT().SelectUserByID(fakeUser.UserID).Return(fakeUser, nil)
	segmentRepo.EXPECT().SelectSegmentBySlug(fakeSegments[0].Slug).Return(&fakeSegments[0], nil)
	segmentRepo.EXPECT().SelectSegmentBySlug(fakeSegments[1].Slug).Return(&fakeSegments[1], nil)
	segmentRepo.EXPECT().InsertSegmentsToUser(fakeUser.UserID, segmentsToAdd).Return(nil)
	segmentRepo.EXPECT().DeleteSegmentsFromUser(fakeUser.UserID, []uint64{fakeSegments[1].SegmentID}).Return(nil)
	segmentRepo.EXPECT().SelectSegmentsByUser(fakeUser.UserID).Return(fakeUserSegments, nil)

	response, err := segmentUC.EditUserSegments(fakeUser.UserID, segmentsToAdd, segmentsToRemove)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUserSegments, response)
	}
}
