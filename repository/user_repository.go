package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

// 1. Interface: Kontrak kerja menggunakan Username
type UserRepository interface {
	CreateUser(user *model.User) error
	FindByUsername(username string) (*model.User, error)
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

// 5. Implementasi FindByUsername: Mencari berdasarkan Username
func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	
	// Cari di database: "SELECT * FROM users WHERE username = ? LIMIT 1"
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err // Kalau username nggak ketemu
	}
	
	return &user, nil // Kalau ketemu, kembalikan data
}