package store

type serviceStore struct {
	repo Repository
}

type Service interface {
	Get(id string) (Store, error)
	GetAll() ([]Store, error)
	Delete(id string) error
	Post(s Store) (Store, error)
}

func New(repository Repository) Service {
	s := serviceStore{
		repo: repository,
	}
	return s
}

func (ss serviceStore) Get(id string) (Store, error) {
	s, err := ss.repo.Get(id)
	if err != nil {
		return Store{}, err
	}
	return s, nil
}

func (ss serviceStore) GetAll() ([]Store, error) {
	s, err := ss.repo.GetAll()
	if err != nil {
		return s, err
	}
	return s, nil
}

func (ss serviceStore) Delete(id string) error {
	err := ss.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (ss serviceStore) Post(s Store) (Store, error) {
	sto, err := ss.repo.Post(s)
	if err != nil {
		return Store{}, err
	}
	return sto, nil
}
