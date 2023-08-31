package usecase

import (
	pkgErr "github.com/pkg/errors"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	"github.com/vvinokurshin/AvitoInternship/internal/user/repository"
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
)

//go:generate mockgen -destination=./mocks/usecase.go -source=./usecase.go -package=mocks

type UseCaseI interface {
	CreateUser(form models.FormUser) (*models.User, error)
	EditUser(userID uint64, form models.FormUser) (*models.User, error)
	DeleteUser(userID uint64) error
	GetUserByID(userID uint64) (*models.User, error)
}

type UseCase struct {
	cfg  *config.Config
	repo repository.RepositoryI
}

func New(cfg *config.Config, repo repository.RepositoryI) UseCaseI {
	return &UseCase{
		cfg:  cfg,
		repo: repo,
	}
}

func (uc *UseCase) CreateUser(form models.FormUser) (*models.User, error) {
	_, err := uc.repo.SelectUserByUsername(form.Username)
	if err != errors.ErrUserNotFound {
		return nil, errors.ErrUserExists
	}

	user := &models.User{
		Username:  form.Username,
		FirstName: form.FirstName,
		LastName:  form.LastName,
	}

	userID, err := uc.repo.InsertUser(user)
	if err != nil {
		return nil, pkgErr.Wrap(err, "insert user")
	}

	user.UserID = userID
	return user, nil
}

func (uc *UseCase) EditUser(userID uint64, form models.FormUser) (*models.User, error) {
	user, err := uc.repo.SelectUserByID(userID)
	if err != nil {
		return nil, pkgErr.Wrap(err, "get user by ID")
	}

	if user.Username != form.Username {
		_, err := uc.repo.SelectUserByUsername(form.Username)
		if err != errors.ErrUserNotFound {
			return nil, errors.ErrUserExists
		}

		user.Username = form.Username
	}

	user.FirstName = form.FirstName
	user.LastName = form.LastName

	err = uc.repo.UpdateUser(user)
	if err != nil {
		return nil, pkgErr.Wrap(err, "update user info")
	}

	return user, nil
}

func (uc *UseCase) DeleteUser(userID uint64) error {
	_, err := uc.repo.SelectUserByID(userID)
	if err != nil {
		return pkgErr.Wrap(err, "select user by ID")
	}

	err = uc.repo.DeleteUser(userID)
	if err != nil {
		return pkgErr.Wrap(err, "delete user")
	}

	return nil
}

func (uc *UseCase) GetUserByID(userID uint64) (*models.User, error) {
	user, err := uc.repo.SelectUserByID(userID)
	if err != nil {
		return nil, pkgErr.Wrap(err, "select user by ID")
	}

	return user, nil
}

//func (uc *UseCase) GetUserByUsername(username string) (*models.User, error) {
//	user, err := uc.repo.SelectUserByUsername(username)
//	if err != nil {
//		return nil, pkgErr.Wrap(err, "select user by username")
//	}
//
//	return user, nil
//}
