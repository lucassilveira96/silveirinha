package domain

import (
	"teste/internal/infra/database"
)

type Services struct {
}

func NewServices(dbs *database.Databases) *Services {

	services := &Services{}
	return services
}
