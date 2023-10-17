package repo

import (
	"bmstu-dips-lab1/internal/person"
	"bmstu-dips-lab1/models"
	"bmstu-dips-lab1/pkg/errs"
	"bmstu-dips-lab1/pkg/postgres"
	"context"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
)

type PersonDB struct {
	id, age             int
	name, address, work string
}

type PersonRepo struct {
	*postgres.Postgres
}

func NewPersonRepo(db *postgres.Postgres) person.Repo {
	return &PersonRepo{db}
}

func (p *PersonRepo) Create(ctx context.Context, modelBL *models.Person) (*models.Person, error) {
	modelDB, err := PersonBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := p.Builder.
		Insert("persons_").
		Columns("name_, address_, work_, age_").
		Values(modelDB.name, modelDB.address, modelDB.work, modelDB.age).
		Suffix("RETURNING \"id_\"").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = p.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.id)
	if err != nil {
		return nil, err
	}

	return PersonDBToBL(modelDB)
}

func (p *PersonRepo) GetById(ctx context.Context, id string) (*models.Person, error) {
	intid, err := strconv.Atoi(id)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := p.Builder.
		Select("name_, address_, work_, age_").
		From("persons_").
		Where(squirrel.Eq{"id_": intid}).
		ToSql()
	if err != nil {
		return nil, err
	}

	modelDB := PersonDB{id: intid}
	err = p.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.name, &modelDB.address, &modelDB.work, &modelDB.age)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, errs.ErrNotFound
		}

		return nil, err
	}

	return PersonDBToBL(&modelDB)
}

func (p *PersonRepo) GetAll(ctx context.Context) ([]*models.Person, error) {
	sql, args, err := p.Builder.
		Select("id_, name_, address_, work_, age_").
		From("persons_").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*models.Person, 0)

	for rows.Next() {
		modelDB := PersonDB{}

		err = rows.Scan(&modelDB.id, &modelDB.name, &modelDB.address, &modelDB.work, &modelDB.age)
		if err != nil {
			return nil, err
		}

		formBL, err := PersonDBToBL(&modelDB)
		if err != nil {
			return nil, err
		}

		res = append(res, formBL)
	}

	return res, nil
}

func (p *PersonRepo) Update(ctx context.Context, modelBL *models.Person, toUpdate *models.Person) (*models.Person, error) {
	modelDB, err := PersonBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	builder := p.Builder.
		Update("persons_")

	if toUpdate.Name != "" {
		builder = builder.
			Set("name_", modelDB.name)
	}

	if toUpdate.Address != "" {
		builder = builder.
			Set("address_", modelDB.address)
	}

	if toUpdate.Work != "" {
		builder = builder.
			Set("work_", modelDB.work)
	}

	if toUpdate.Age != 0 {
		builder = builder.
			Set("age_", modelDB.age)
	}

	sql, args, err := builder.
		Where(squirrel.Eq{"id_": modelDB.id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	res, err := p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	if res.RowsAffected() == 0 {
		return nil, errs.ErrNoContent
	}

	return modelBL, nil
}

func (p *PersonRepo) Delete(ctx context.Context, id string) error {
	intid, err := strconv.Atoi(id)
	if err != nil {
		return errs.ErrInvalidContent
	}

	sql, args, err := p.Builder.
		Delete("persons_").
		Where(squirrel.Eq{"id_": intid}).
		ToSql()
	if err != nil {
		return err
	}

	res, err := p.Pool.Exec(ctx, sql, args...)

	if res.RowsAffected() == 0 {
		return errs.ErrNoContent
	}

	return nil
}

func PersonDBToBL(modelDB *PersonDB) (*models.Person, error) {
	return &models.Person{
		Id:      strconv.Itoa(modelDB.id),
		Name:    modelDB.name,
		Address: modelDB.address,
		Work:    modelDB.work,
		Age:     modelDB.age,
	}, nil
}

func PersonBLToDB(modelBL *models.Person) (*PersonDB, error) {
	var (
		err error
		id  int
	)

	if modelBL.Id != "" {
		id, err = strconv.Atoi(modelBL.Id)
		if err != nil {
			return nil, err
		}
	}

	return &PersonDB{
		id:      id,
		name:    modelBL.Name,
		address: modelBL.Address,
		work:    modelBL.Work,
		age:     modelBL.Age,
	}, nil
}
