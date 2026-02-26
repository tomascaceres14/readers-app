package user

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

func (s *Service) FindById(id string, desc bool) (*User, error) {
	return s.repo.FindById(id)
}

func (s *Service) FindByIdWithResources(id string, desc bool) (*User, error) {
	return s.repo.FindByIdWithResources(id, desc)
}

func (s *Service) FindByIdWithResourcesAndRelationDate(id string, desc bool) (*User, error) {
	return s.repo.FindByIdWithResourcesAndRelationDate(id, desc)
}

func (s *Service) UpdateById(id, username string) error {

	user, err := s.FindById(id, false)
	if err != nil {
		return err
	}

	user.Username = username

	return s.repo.Update(user)
}

func (s *Service) DeleteById(id string) {
	user, _ := s.FindById(id, false)
	s.repo.Delete(user)
}
