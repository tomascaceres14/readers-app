package resource

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/errs"
)

type Service struct {
	r *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{r: repo}
}

func (s *Service) GetAll() ([]Resource, error) {
	return s.r.FindAll()
}

func (s *Service) Create(resource *Resource) error {
	cleaned := cleanURL(resource.Url)
	resource.Url = cleaned

	fmt.Println("received", cleaned)
	existing, err := s.r.FindByUrl(cleaned)
	fmt.Println("existing:", existing)
	fmt.Println("error:", err)
	if err == nil {
		resource.ID = existing.ID
		return nil
	}

	if errors.Is(err, errs.ErrNotFound) {
		return s.r.Create(resource)
	}

	return err
}

func (s *Service) FindById(id string) (*Resource, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return s.r.FindById(uuid)
}

func (s *Service) DeleteAll() {
	s.r.DeleteAll()
}

func cleanURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return u.Host + u.Path
}
