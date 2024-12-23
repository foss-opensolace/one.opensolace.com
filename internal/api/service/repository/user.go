package repository

import (
	"github.com/foss-opensolace/api.opensolace.com/internal/api/model"
	"github.com/foss-opensolace/api.opensolace.com/internal/api/model/dto"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	Create(obj *dto.UserRegister) (*model.User, error)
}

func NewUserRepository(postgres *gorm.DB) UserRepository {
	return userRepository{db: postgres}
}

func (u userRepository) Create(obj *dto.UserRegister) (*model.User, error) {
	var displayName *string

	if obj.DisplayName != "" {
		displayName = &obj.DisplayName
	}

	newUser := model.User{
		DisplayName: displayName,
		Username:    obj.Username,
		Email:       obj.Email,
	}

	var user *model.User
	if err := u.db.First(&user, "username = ? OR email = ?", newUser.Username, newUser.Email).Error; err == nil {
		return nil, gorm.ErrDuplicatedKey
	}

	if hashed, err := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost); err != nil {
		return nil, err
	} else {
		newUser.Password = string(hashed)
	}

	transaction := u.db.Begin()

	if err := transaction.Create(&newUser).Error; err != nil {
		transaction.Rollback()
		return nil, err
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return nil, err
	}

	return &newUser, nil
}
