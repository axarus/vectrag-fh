package application

import "github.com/axarus/vectrag/internal/domain"

type ModelService struct {
	repo domain.Repository
}

func NewModelService(repo domain.Repository) *ModelService {
	return &ModelService{
		repo: repo,
	}
}

func (ms *ModelService) Create(model domain.Model) error {
	return ms.repo.CreateModel(model)
}

func (ms *ModelService) Update(model domain.Model) error {
	return ms.repo.UpdateModel(model)
}

func (ms *ModelService) Delete(slug string) error {
	return ms.repo.DeleteModel(slug)
}

func (ms *ModelService) Get(slug string) (domain.Model, error) {
	return ms.repo.GetModel(slug)
}

func (ms *ModelService) List() ([]domain.Model, error) {
	return ms.repo.GetModels()
}
