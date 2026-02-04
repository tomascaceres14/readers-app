package user

import "github.com/google/uuid"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]User, error) {
	return s.repo.FindAll()
}

func (s *Service) Register(u *User) error {
	// validate legal email, check username/email don't exist in db,
	// hash password, generate tokens, etc...
	return s.repo.Register(u)
}

func (s *Service) FindById(id string) (*User, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return s.repo.FindById(uuid)
}

func (s *Service) UpdateById(id, username string) error {

	user, err := s.FindById(id)
	if err != nil {
		return err
	}

	user.Username = username

	return s.repo.Update(user)
}

func (s *Service) DeleteById(id string) {
	user, _ := s.FindById(id)
	s.repo.Delete(user)
}
