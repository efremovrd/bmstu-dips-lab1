package usecase

import (
	"bmstu-dips-lab1/internal/person"
	"bmstu-dips-lab1/models"
	"context"
)

type PersonUseCase struct {
	personRepo person.Repo
}

func NewPersonUseCase(personRepo person.Repo) person.UseCase {
	return &PersonUseCase{
		personRepo: personRepo,
	}
}

func (p *PersonUseCase) Create(ctx context.Context, model *models.Person) (*models.Person, error) {
	return p.personRepo.Create(ctx, model)
}

func (p *PersonUseCase) GetAll(ctx context.Context) ([]*models.Person, error) {
	return p.personRepo.GetAll(ctx)
}

func (p *PersonUseCase) GetById(ctx context.Context, id int) (*models.Person, error) {
	return p.personRepo.GetById(ctx, id)
}

func (p *PersonUseCase) Update(ctx context.Context, model *models.Person, toUpdate *models.Person) (*models.Person, error) {
	_, err := p.personRepo.Update(ctx, model, toUpdate)
	if err != nil {
		return nil, err
	}

	// crutch
	return p.personRepo.GetById(ctx, model.Id)
}

func (p *PersonUseCase) Delete(ctx context.Context, id int) error {
	return p.personRepo.Delete(ctx, id)
}
