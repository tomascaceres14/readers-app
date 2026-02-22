package userresources

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(userResource *UserResource) error {
	return r.db.Create(userResource).Error
}

func (r *Repository) Update(userResource *UserResource) error {
	return r.db.Save(userResource).Error
}

func (r *Repository) Delete(userResource *UserResource) {
	r.db.Delete(userResource)
}
