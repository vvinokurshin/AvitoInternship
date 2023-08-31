package postgres

import (
	pkgErrors "github.com/pkg/errors"
	"github.com/vvinokurshin/AvitoInternship/internal/config"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
	"github.com/vvinokurshin/AvitoInternship/internal/user/repository"
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
	"gorm.io/gorm"
)

type userRepo struct {
	cfg *config.Config
	db  *gorm.DB
}

func New(cfg *config.Config, db *gorm.DB) repository.RepositoryI {
	return &userRepo{
		cfg: cfg,
		db:  db,
	}
}

func (repo *userRepo) InsertUser(user *models.User) (uint64, error) {
	var dbUser User
	dbUser.FromUserModel(user)

	tx := repo.db.Table(User{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBUserTableName)).Create(&dbUser)
	if err := tx.Error; err != nil {
		return 0, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return dbUser.UserID, nil
}

func (repo *userRepo) UpdateUser(user *models.User) error {
	var dbUser User
	dbUser.FromUserModel(user)

	tx := repo.db.Table(User{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBUserTableName)).
		Omit("user_id").Updates(&dbUser)
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (repo *userRepo) DeleteUser(userID uint64) error {
	tx := repo.db.Table(User{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBUserTableName)).
		Where("user_id = ?", userID).Delete(User{})
	if err := tx.Error; err != nil {
		return pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return nil
}

func (repo *userRepo) SelectUserByID(userID uint64) (*models.User, error) {
	var dbUser User

	tx := repo.db.Table(User{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBUserTableName)).
		Where("user_id = ?", userID).Take(&dbUser)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrUserNotFound
		}

		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return dbUser.ToUserModel(), nil
}

func (repo *userRepo) SelectUserByUsername(username string) (*models.User, error) {
	var dbUser User

	tx := repo.db.Table(User{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBUserTableName)).
		Where("username = ?", username).Take(&dbUser)
	if err := tx.Error; err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrUserNotFound
		}

		return nil, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return dbUser.ToUserModel(), nil
}

func (repo *userRepo) SelectUserIDs() ([]uint64, error) {
	var IDs []uint64

	tx := repo.db.Table(User{}.TableName(repo.cfg.DB.DBSchemaName, repo.cfg.DB.DBUserTableName)).Select("user_id").Find(&IDs)
	if err := tx.Error; err != nil {
		return IDs, pkgErrors.WithMessage(errors.ErrInternal, err.Error())
	}

	return IDs, nil
}
