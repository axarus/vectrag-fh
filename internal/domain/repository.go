package domain

type Repository interface {
	CreateModel(model Model) error
	UpdateModel(model Model) error
	DeleteModel(slug string) error
	GetModel(slug string) (Model, error)
	GetModels() ([]Model, error)
}
