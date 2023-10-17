package person

import (
	"bmstu-dips-lab1/models"
	"context"
)

type Repo interface {
	Create(ctx context.Context, modelBL *models.Person) (*models.Person, error)
	GetById(ctx context.Context, id string) (*models.Person, error)
	GetAll(ctx context.Context) ([]*models.Person, error)
	Update(ctx context.Context, modelBL *models.Person, toUpdate *models.Person) (*models.Person, error)
	Delete(ctx context.Context, id string) error
}
