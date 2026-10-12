package repository

import "gialong-dev04/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	FindById(user *models.User, id int) error
}
