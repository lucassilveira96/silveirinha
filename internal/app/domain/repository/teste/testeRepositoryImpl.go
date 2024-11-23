package testeRepository

import (
	"silveirinha/internal/app/domain/model"
	"silveirinha/internal/infra/database"
	"time"
)

var _ TesteRepository = (*TesteRepositoryImpl)(nil)

type TesteRepositoryImpl struct {
	db *database.Databases
}

func NewTesteRepository(db *database.Databases) *TesteRepositoryImpl {
	return &TesteRepositoryImpl{db: db}
}

func (r *TesteRepositoryImpl) Create(teste *model.Teste) error {
	return r.db.Write.Create(teste).Error
}

func (r *TesteRepositoryImpl) Update(id uint, teste *model.Teste) error {
	existing := &model.Teste{}
	if err := r.db.Write.First(existing, id).Error; err != nil {
		return err
	}
	return r.db.Write.Model(existing).Updates(teste).Error
}

func (r *TesteRepositoryImpl) Delete(id uint) error {
	teste := &model.Teste{}
	if err := r.db.Write.First(teste, id).Error; err != nil {
		return err
	}
	return r.db.Write.Model(teste).Update("deleted_at", time.Now()).Error
}

func (r *TesteRepositoryImpl) FindAll() ([]*model.Teste, error) {
	var testes []*model.Teste
	err := r.db.Read.
		Where("deleted_at IS NULL").
		Find(&testes).Error
	return testes, err
}

func (r *TesteRepositoryImpl) FindById(id uint) (*model.Teste, error) {
	var teste model.Teste
	err := r.db.Read.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&teste).Error
	return &teste, err
}
