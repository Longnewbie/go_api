package repository

import "gialong-dev04/internal/models"

type UserRepository interface {
	Create(user *models.User)
	FindById(id int)
}
