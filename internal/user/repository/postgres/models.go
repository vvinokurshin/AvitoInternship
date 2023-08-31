package postgres

import (
	"fmt"
	"github.com/vvinokurshin/AvitoInternship/internal/models"
)

type User struct {
	UserID    uint64 `gorm:"primary_key"`
	Username  string
	FirstName string
	LastName  string
}

func (User) TableName(schemaName, tableName string) string {
	return fmt.Sprintf("%s.%s", schemaName, tableName)
}

func (u *User) FromUserModel(user *models.User) {
	u.UserID = user.UserID
	u.Username = user.Username
	u.FirstName = user.FirstName
	u.LastName = user.LastName
}

func (u *User) ToUserModel() (user *models.User) {
	return &models.User{
		UserID:    u.UserID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}
