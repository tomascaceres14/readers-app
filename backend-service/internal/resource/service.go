package resource

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/errs"
	userresources "github.com/tomascaceres14/readers-app/app-server/backend-service/internal/user_resource"
)

type Service struct {
	r     *Repository
	urSvc *userresources.Repository
}

func NewService(repo *Repository, userResourceSvc *userresources.Repository) *Service {
	return &Service{r: repo, urSvc: userResourceSvc}
}

func (s *Service) GetAll() ([]Resource, error) {
	return s.r.FindAll()
}

func (s *Service) Create(resource *Resource, uid string) error {
	if uid == "" {
		return errs.ErrBadRequest
	}

	cleaned, err := cleanURL(resource.Url)
	if err != nil {
		return err
	}

	resource.Url = cleaned

	existing, err := s.r.FindByUrl(cleaned)
	switch err {
	case nil:
		if s.urSvc.Exists(uid, existing.ID) {
			fmt.Println("exists!")
			return errs.ErrAlreadyExists
		}
		resource.ID = existing.ID
	case errs.ErrNotFound:
		if error := s.r.Create(resource); error != nil {
			fmt.Println(err)
			return error
		}
	default:
		return err
	}

	userResource := userresources.UserResource{
		UserID:     uid,
		ResourceID: resource.ID,
	}

	if err := s.urSvc.Update(&userResource); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
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

func cleanURL(rawURL string) (string, error) {
	if !strings.Contains(rawURL, ".com") {
		return "", errs.ErrBadRequest
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := strings.TrimPrefix(u.Host, "www.")
	return host + u.Path, nil
}
