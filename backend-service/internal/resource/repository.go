package resource

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/errs"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(resources *Resource) error {
	return r.db.Create(resources).Error
}

func (r *Repository) FindAll() ([]Resource, error) {
	var resources []Resource
	if err := r.db.Find(resources).Error; err != nil {
		return nil, err
	}

	return resources, nil
}

func (r *Repository) FindById(id uuid.UUID) (*Resource, error) {
	var resource Resource
	if err := r.db.First(&resource, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &resource, nil
}

func (r *Repository) FindByUrl(url string) (*Resource, error) {
	var resource Resource
	if err := r.db.Where("url = ?", url).First(&resource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound("resource not found", err)
		}
		return nil, err
	}
	return &resource, nil
}
