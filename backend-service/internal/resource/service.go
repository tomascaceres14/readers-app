package resource

import (
	"net/url"

	"github.com/google/uuid"
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

	existing, err := s.r.FindByUrl(cleaned)
	if err != nil {
		return err
	}

	if existing != nil {
		resource.ID = existing.ID
		return nil
	}

	return s.r.Create(resource)
}

func (s *Service) FindById(id string) (*Resource, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return s.r.FindById(uuid)
}

func cleanURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return u.Host + u.Path
}
