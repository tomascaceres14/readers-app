package resource

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tomascaceres14/readers-app/backend-service/internal/errs"
	"github.com/tomascaceres14/readers-app/backend-service/internal/messaging"
	resourcestatus "github.com/tomascaceres14/readers-app/backend-service/internal/resource_status"
	userresources "github.com/tomascaceres14/readers-app/backend-service/internal/user_resource"
)

type Service struct {
	r          *Repository
	urSvc      *userresources.Repository
	statusRepo *resourcestatus.Repository
	publisher  *messaging.Publisher
}

func NewService(repo *Repository, userResourceSvc *userresources.Repository, statusRepo *resourcestatus.Repository, publisher *messaging.Publisher) *Service {
	return &Service{r: repo, urSvc: userResourceSvc, statusRepo: statusRepo, publisher: publisher}
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

	pendingStatus, err := s.statusRepo.FindByName(resourcestatus.PENDING)
	if err != nil {
		return err
	}

	resource.Status = *pendingStatus

	existing, err := s.r.FindByUrl(cleaned)
	switch err {
	case nil:
		if s.urSvc.Exists(uid, existing.ID) {
			return errs.ErrAlreadyExists
		}
		*resource = *existing
	case errs.ErrNotFound:
		if error := s.r.Create(resource); error != nil {
			return error
		}
	default:
		return err
	}

	userResource := userresources.UserResource{
		UserID:     uid,
		ResourceID: resource.ID,
		CreatedAt:  time.Now(),
	}

	if err := s.urSvc.Create(&userResource); err != nil {
		return err
	}

	msg := messaging.Message{
		UserID:     uid,
		ResourceID: resource.ID.String(),
		URL:        resource.Url,
	}

	fmt.Println("msg:", msg)
	if err := s.publisher.PublishScrapingTask(msg); err != nil {
		return err
	}

	resource.CreatedAt = userResource.CreatedAt

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

func (s *Service) UpdateAfterScrape(id uuid.UUID, statusName, title, excerpt, language string) error {
	status, err := s.statusRepo.FindByName(statusName)
	if err != nil {
		return fmt.Errorf("failed to find status %s: %w", statusName, err)
	}

	return s.r.UpdateAfterScrape(id, status.ID, title, excerpt, language)
}

func (s *Service) UpdateStatusFailed(id uuid.UUID, statusName string) error {
	status, err := s.statusRepo.FindByName(statusName)
	if err != nil {
		return fmt.Errorf("failed to find status %s: %w", statusName, err)
	}

	return s.r.UpdateStatusFailed(id, status.ID)
}

func cleanURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if u.Scheme != "https" {
		return "", errs.ErrBadURL
	}

	if u.Host == "" || !strings.Contains(u.Host, ".") {
		return "", errs.ErrBadURL
	}

	u.RawQuery = ""
	u.Fragment = ""

	return u.String(), nil
}
