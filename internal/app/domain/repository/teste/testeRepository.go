package testeRepository

import "silveirinha/internal/app/domain/model"

type TesteRepository interface {
	Create(teste *model.Teste) error
	Update(id uint, teste *model.Teste) error
	Delete(id uint) error
	FindAll() ([]*model.Teste, error)
	FindById(id uint) (*model.Teste, error)
}
