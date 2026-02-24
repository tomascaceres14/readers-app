package userresources

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

func (r *Repository) Create(userResource *UserResource) error {
	return r.db.Create(userResource).Error
}

func (r *Repository) Update(userResource *UserResource) error {
	return r.db.Save(userResource).Error
}

func (r *Repository) Delete(userResource *UserResource) {
	r.db.Delete(userResource)
}

func (r *Repository) Exists(uid string, resourceID uuid.UUID) bool {
	var count int64
	err := r.db.Model(&UserResource{}).Where("user_id = ? AND resource_id = ?", uid, resourceID).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}
