package testepppppRepository

import (
	"silveirinha/internal/app/domain/model"
	"silveirinha/internal/infra/database"
	"time"
)

var _ TestepppppRepository = (*TestepppppRepositoryImpl)(nil)

type TestepppppRepositoryImpl struct {
	db *database.Databases
}

func NewTestepppppRepository(db *database.Databases) *TestepppppRepositoryImpl {
	return &TestepppppRepositoryImpl{db: db}
}

func (r *TestepppppRepositoryImpl) Create(testeppppp *model.Testeppppp) error {
	return r.db.Write.Create(testeppppp).Error
}

func (r *TestepppppRepositoryImpl) Update(id uint, testeppppp *model.Testeppppp) error {
	existing := &model.Testeppppp{}
	if err := r.db.Write.First(existing, id).Error; err != nil {
		return err
	}
	return r.db.Write.Model(existing).Updates(testeppppp).Error
}

func (r *TestepppppRepositoryImpl) Delete(id uint) error {
	testeppppp := &model.Testeppppp{}
	if err := r.db.Write.First(testeppppp, id).Error; err != nil {
		return err
	}
	return r.db.Write.Model(testeppppp).Update("deleted_at", time.Now()).Error
}

func (r *TestepppppRepositoryImpl) FindAll() ([]*model.Testeppppp, error) {
	var testeppppps []*model.Testeppppp
	err := r.db.Read.
		Where("deleted_at IS NULL").
		Find(&testeppppps).Error
	return testeppppps, err
}

func (r *TestepppppRepositoryImpl) FindById(id uint) (*model.Testeppppp, error) {
	var testeppppp model.Testeppppp
	err := r.db.Read.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&testeppppp).Error
	return &testeppppp, err
}
