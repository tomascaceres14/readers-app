package user

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll() ([]User, error) {
	var users []User
	if err := r.db.Preload("Resources").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) FindById(id string, desc bool) (*User, error) {
	var user User
	order := "asc"
	if desc {
		order = "desc"
	}
	if err := r.db.Preload("Resources", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at " + order)
	}).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Register(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *Repository) Delete(user *User) {
	r.db.Delete(user)
}
