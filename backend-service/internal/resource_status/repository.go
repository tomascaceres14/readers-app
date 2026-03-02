package resourcestatus

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tomascaceres14/readers-app/backend-service/internal/errs"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByName(name string) (*ResourceStatus, error) {
	var status ResourceStatus
	if err := r.db.Where("name = ?", name).First(&status).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return &status, nil
}

func (r *Repository) FindById(id uuid.UUID) (*ResourceStatus, error) {
	var status ResourceStatus
	if err := r.db.First(&status, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return &status, nil
}
