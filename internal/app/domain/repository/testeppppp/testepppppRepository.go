package testepppppRepository

import "silveirinha/internal/app/domain/model"

type TestepppppRepository interface {
	Create(testeppppp *model.Testeppppp) error
	Update(id uint, testeppppp *model.Testeppppp) error
	Delete(id uint) error
	FindAll() ([]*model.Testeppppp, error)
	FindById(id uint) (*model.Testeppppp, error)
}
