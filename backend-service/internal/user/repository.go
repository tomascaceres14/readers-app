package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll() ([]User, error) {
	var Users []User
	if err := r.db.Find(Users).Error; err != nil {
		return nil, err
	}

	return Users, nil
}

func (r *Repository) FindByID(id uuid.UUID) (User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) Save(user User) (uuid.UUID, error) {
	if err := r.db.Create(user).Error; err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}
