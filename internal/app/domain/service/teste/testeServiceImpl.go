package testeService

import (
	"silveirinha/internal/app/domain/model"
	testeRepository "silveirinha/internal/app/domain/repository/teste"
)

var _ TesteService = (*TesteServiceImpl)(nil)

type TesteServiceImpl struct {
	repository testeRepository.TesteRepository
}

func NewTesteService(repository testeRepository.TesteRepository) *TesteServiceImpl {
	return &TesteServiceImpl{repository: repository}
}

func (s *TesteServiceImpl) Create(teste *model.Teste) error {
	return s.repository.Create(teste)
}

func (s *TesteServiceImpl) Update(id uint, teste *model.Teste) error {
	return s.repository.Update(id, teste)
}

func (s *TesteServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}

func (s *TesteServiceImpl) FindAll() ([]*model.Teste, error) {
	return s.repository.FindAll()
}

func (s *TesteServiceImpl) FindById(id uint) (*model.Teste, error) {
	return s.repository.FindById(id)
}
