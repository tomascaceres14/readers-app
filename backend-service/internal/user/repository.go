package user

import (
	"github.com/tomascaceres14/readers-app/backend-service/internal/resource"
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
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) FindById(id string) (*User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByIdWithResources(id string, desc bool) (*User, error) {
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

func (r *Repository) FindByIdWithResourcesAndRelationDate(id string, desc bool) (*User, error) {
	u, err := r.FindById(id)
	if err != nil {
		return nil, err
	}

	order := "ASC"
	if desc {
		order = "DESC"
	}
	var results []resource.Resource
	err = r.db.
		Model(&resource.Resource{}).
		Select("resources.url, user_resources.created_at").
		Joins("JOIN user_resources ON resources.id = user_resources.resource_id").
		Where("user_resources.user_id = ?", id).
		Order("user_resources.created_at " + order).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	u.Resources = results

	return u, err

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
