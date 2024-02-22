package product

type ServiceProduct struct {
	repo Repository
}

type Service interface {
	Get(id string) (Product, error)
	Delete(id string) error
	Post(p Product) (Product, error)
	GetAll() ([]Product, error)
	GetAllByName(name string) ([]Product, error)
}

func New(repository Repository) Service {
	s := ServiceProduct{
		repo: repository,
	}
	return s
}

func (sp ServiceProduct) Get(id string) (Product, error) {
	p, err := sp.repo.Get(id)
	if err != nil {
		return Product{}, err
	}
	return p, nil
}
func (sp ServiceProduct) Delete(id string) error {
	err := sp.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
func (sp ServiceProduct) Post(p Product) (Product, error) {
	AddTva(&p)
	pro, err := sp.repo.Post(p)
	if err != nil {
		return Product{}, err
	}
	return pro, nil
}
func (sp ServiceProduct) GetAll() ([]Product, error) {
	p, err := sp.repo.GetAll()
	if err != nil {
		return p, err
	}
	return p, nil
}
func (sp ServiceProduct) GetAllByName(name string) ([]Product, error) {
	p, err := sp.repo.GetAllByName(name)
	if err != nil {
		return p, err
	}
	return p, nil
}
