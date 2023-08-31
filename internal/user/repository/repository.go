package repository

import "github.com/vvinokurshin/AvitoInternship/internal/models"

//go:generate mockgen -destination=./mocks/repository.go -source=./repository.go -package=mocks

type RepositoryI interface {
	InsertUser(user *models.User) (uint64, error)
	UpdateUser(user *models.User) error
	DeleteUser(userID uint64) error
	SelectUserByID(userID uint64) (*models.User, error)
	SelectUserByUsername(username string) (*models.User, error)
	SelectUserIDs() ([]uint64, error)
}
