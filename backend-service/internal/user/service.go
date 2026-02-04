package user

import "github.com/google/uuid"

type Service struct {
	repository *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) Create(u User) (uuid.UUID, error) {
	return s.repository.Save(u)
}
