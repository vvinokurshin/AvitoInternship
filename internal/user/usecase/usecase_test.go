package usecase

import (
	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	pkgErr "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
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

func TestUseCase_CreateUser(t *testing.T) {
	cfg := createConfig()

	var fakeForm models.FormUser
	generateFakeData(&fakeForm)
	fakeUser := &models.User{
		Username:  fakeForm.Username,
		FirstName: fakeForm.FirstName,
		LastName:  fakeForm.LastName,
	}
	fakeUserResponse := &models.User{
		UserID:    1,
		Username:  fakeForm.Username,
		FirstName: fakeForm.FirstName,
		LastName:  fakeForm.LastName,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepositoryI(ctrl)
	userUC := New(cfg, userRepo)

	userRepo.EXPECT().SelectUserByUsername(fakeForm.Username).Return(nil, errors.ErrUserNotFound)
	userRepo.EXPECT().InsertUser(fakeUser).Return(uint64(1), nil)
	response, err := userUC.CreateUser(fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUserResponse, response)
	}
}

func TestUseCase_EditUser(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var fakeForm models.FormUser
	generateFakeData(&fakeForm)
	fakeUser := &models.User{
		UserID:    userID,
		Username:  fakeForm.Username,
		FirstName: fakeForm.FirstName,
		LastName:  fakeForm.LastName,
	}

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepositoryI(ctrl)
	userUC := New(cfg, userRepo)

	userRepo.EXPECT().SelectUserByID(userID).Return(fakeUser, nil)
	userRepo.EXPECT().UpdateUser(fakeUser).Return(nil)
	response, err := userUC.EditUser(userID, fakeForm)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}

func TestUseCase_DeleteUser(t *testing.T) {
	cfg := createConfig()

	userID := uint64(1)
	var fakeUser *models.User
	generateFakeData(&fakeUser)
	fakeUser.UserID = userID

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepositoryI(ctrl)
	userUC := New(cfg, userRepo)

	userRepo.EXPECT().SelectUserByID(userID).Return(fakeUser, nil)
	userRepo.EXPECT().DeleteUser(userID).Return(nil)
	err := userUC.DeleteUser(userID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	}
}

func TestUseCase_GetUser(t *testing.T) {
	cfg := createConfig()

	var fakeUser *models.User
	generateFakeData(&fakeUser)

	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mockUserRepo.NewMockRepositoryI(ctrl)
	userUC := New(cfg, userRepo)

	userRepo.EXPECT().SelectUserByID(fakeUser.UserID).Return(fakeUser, nil)
	response, err := userUC.GetUserByID(fakeUser.UserID)
	causeErr := pkgErr.Cause(err)

	if causeErr != nil {
		t.Errorf("[TEST] simple: expected err \"%v\", got \"%v\"", nil, causeErr)
	} else {
		require.Equal(t, fakeUser, response)
	}
}
