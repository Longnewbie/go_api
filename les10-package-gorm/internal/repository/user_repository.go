package repository

import (
	"fmt"
	"gialong-dev04/internal/models"

	"gorm.io/gorm"
)

type SQLUserRepository struct {
	db *gorm.DB
}

func NewSQLUserRepository(DB *gorm.DB) UserRepository {
	return &SQLUserRepository{
		db: DB,
	}
}

func (ur *SQLUserRepository) Create(user *models.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}

	return nil
}

func (ur *SQLUserRepository) FindById(user *models.User, id int) error {
	if err := ur.db.First(user, id).Error; err != nil {
		return err
	}

	return nil
}
