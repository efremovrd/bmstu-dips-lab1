package person

import (
	"bmstu-dips-lab1/models"
	"context"
)

type UseCase interface {
	Create(ctx context.Context, model *models.Person) (*models.Person, error)
	GetAll(ctx context.Context) ([]*models.Person, error)
	GetById(ctx context.Context, id string) (*models.Person, error)
	Update(ctx context.Context, model *models.Person, toUpdate *models.Person) (*models.Person, error)
	Delete(ctx context.Context, id string) error
}
