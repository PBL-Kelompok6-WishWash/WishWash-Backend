package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

// 1. Interface
type UserRepository interface {
	CreateUser(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	UpdateUser(user *model.User) error
	FindByID(userID uint) (*model.User, error)
	DeleteUser(id uint) error
}

// 2. Struct
type userRepository struct {
	db *gorm.DB
}

// 3. Constructor
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// 4. Implementasi CreateUser
func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

// 5. Implementasi FindByUsername
func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Model(user).Updates(user).Error
}

func (r *userRepository) FindByID(userID uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, userID).Error
	return &user, err
}

func (r *userRepository) DeleteUser(id uint) error {
	return r.db.Where("id_user = ?", id).Delete(&model.User{}).Error
}