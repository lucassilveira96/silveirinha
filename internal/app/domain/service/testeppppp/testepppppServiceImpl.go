package testepppppService

import (
	"silveirinha/internal/app/domain/model"
	testepppppRepository "silveirinha/internal/app/domain/repository/testeppppp"
)

var _ TestepppppService = (*TestepppppServiceImpl)(nil)

type TestepppppServiceImpl struct {
	repository testepppppRepository.TestepppppRepository
}

func NewTestepppppService(repository testepppppRepository.TestepppppRepository) *TestepppppServiceImpl {
	return &TestepppppServiceImpl{repository: repository}
}

func (s *TestepppppServiceImpl) Create(testeppppp *model.Testeppppp) error {
	return s.repository.Create(testeppppp)
}

func (s *TestepppppServiceImpl) Update(id uint, testeppppp *model.Testeppppp) error {
	return s.repository.Update(id, testeppppp)
}

func (s *TestepppppServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}

func (s *TestepppppServiceImpl) FindAll() ([]*model.Testeppppp, error) {
	return s.repository.FindAll()
}

func (s *TestepppppServiceImpl) FindById(id uint) (*model.Testeppppp, error) {
	return s.repository.FindById(id)
}
