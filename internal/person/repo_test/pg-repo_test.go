package repo_test

import (
	"bmstu-dips-lab1/internal/person/repo"
	"bmstu-dips-lab1/models"
	"bmstu-dips-lab1/pkg/errs"
	"bmstu-dips-lab1/pkg/postgres"
	"context"
	"errors"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

var (
	_builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func TestPersonRepo_Create(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewPersonRepo(&db)

	type mockBehavior func(ctx context.Context, form *models.Person)

	testTable := []struct {
		nameTest       string
		ctx            context.Context
		person         models.Person
		mockBehavior   mockBehavior
		expectedPerson models.Person
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			person: models.Person{
				Name:    "qwerty",
				Work:    "sdcsd",
				Address: "ecefvc",
				Age:     12,
			},
			mockBehavior: func(ctx context.Context, person *models.Person) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_"}).AddRow(345).ToPgxRows()
				pgxRows.Next()
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO persons_ (name_, address_, work_, age_) VALUES ($1,$2,$3,$4) RETURNING \"id_\"", person.Name, person.Address, person.Work, person.Age).Return(pgxRows)
			},
			expectedPerson: models.Person{
				Id:      345,
				Name:    "qwerty",
				Work:    "sdcsd",
				Address: "ecefvc",
				Age:     12,
			},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			person: models.Person{
				Name:    "qwerty",
				Work:    "sdcsd",
				Address: "ecefvc",
				Age:     12,
			},
			mockBehavior: func(ctx context.Context, person *models.Person) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO persons_ (name_, address_, work_, age_) VALUES ($1,$2,$3,$4) RETURNING \"id_\"", person.Name, person.Address, person.Work, person.Age).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.person)

			got, err := r.Create(testCase.ctx, &testCase.person)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedPerson, *got)
			case "no_rows":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestPersonRepo_GetById(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewPersonRepo(&db)

	type mockBehavior func(ctx context.Context, id int)

	testTable := []struct {
		nameTest       string
		ctx            context.Context
		id             int
		person         models.Person
		mockBehavior   mockBehavior
		expectedPerson models.Person
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			id:       345,
			mockBehavior: func(ctx context.Context, id int) {
				pgxRows := pgxpoolmock.NewRows([]string{"name_", "address_", "work_", "age_"}).AddRow("qwerty", "ecefvc", "sdcsd", 12).ToPgxRows()
				pgxRows.Next()
				mockPool.EXPECT().QueryRow(ctx, "SELECT name_, address_, work_, age_ FROM persons_ WHERE id_ = $1", id).Return(pgxRows)
			},
			expectedPerson: models.Person{
				Id:      345,
				Name:    "qwerty",
				Work:    "sdcsd",
				Address: "ecefvc",
				Age:     12,
			},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			id:       345,
			mockBehavior: func(ctx context.Context, id int) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				mockPool.EXPECT().QueryRow(ctx, "SELECT name_, address_, work_, age_ FROM persons_ WHERE id_ = $1", id).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.id)

			got, err := r.GetById(testCase.ctx, testCase.id)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedPerson, *got)
			case "invalid_inputs":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "no_rows":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestPersonRepo_GetAll(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewPersonRepo(&db)

	type mockBehavior func(ctx context.Context)

	testTable := []struct {
		nameTest        string
		ctx             context.Context
		mockBehavior    mockBehavior
		expectedPersons []*models.Person
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			mockBehavior: func(ctx context.Context) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_", "name_", "address_", "work_", "age_"}).AddRow(345, "qwerty1", "address1", "work1", 11).AddRow(346, "qwerty2", "address2", "work2", 12).ToPgxRows()
				mockPool.EXPECT().Query(ctx, "SELECT id_, name_, address_, work_, age_ FROM persons_").Return(pgxRows, nil)
			},
			expectedPersons: []*models.Person{
				{
					Id:      345,
					Address: "address1",
					Work:    "work1",
					Name:    "qwerty1",
					Age:     11,
				},
				{
					Id:      346,
					Address: "address2",
					Work:    "work2",
					Name:    "qwerty2",
					Age:     12,
				},
			},
		},
		{
			nameTest: "query_error",
			ctx:      context.Background(),
			mockBehavior: func(ctx context.Context) {
				mockPool.EXPECT().Query(ctx, "SELECT id_, name_, address_, work_, age_ FROM persons_").Return(nil, errors.New("query_error"))
			},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			mockBehavior: func(ctx context.Context) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				pgxRows.Next()
				mockPool.EXPECT().Query(ctx, "SELECT id_, name_, address_, work_, age_ FROM persons_").Return(pgxRows, nil)
			},
			expectedPersons: []*models.Person{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx)

			got, err := r.GetAll(testCase.ctx)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedPersons, got)
			case "query_error":
				assert.NotEqual(t, nil, err)
			case "no_rows":
				assert.Equal(t, nil, err)
				assert.Equal(t, []*models.Person{}, got)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestFormRepo_Update(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewPersonRepo(&db)

	type mockBehavior func(ctx context.Context, person *models.Person, toUpdate *models.Person)

	testTable := []struct {
		nameTest       string
		ctx            context.Context
		person         models.Person
		toUpdate       models.Person
		mockBehavior   mockBehavior
		expectedPerson models.Person
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			person: models.Person{
				Id:      345,
				Name:    "qwerty",
				Work:    "work",
				Address: "address",
				Age:     12,
			},
			toUpdate: models.Person{
				Id:      0,
				Name:    "",
				Work:    "newwork",
				Address: "",
				Age:     0,
			},
			mockBehavior: func(ctx context.Context, person *models.Person, toUpdate *models.Person) {
				mockPool.EXPECT().Exec(ctx, "UPDATE persons_ SET work_ = $1 WHERE id_ = $2", person.Work, person.Id).Return(pgxmock.NewResult("UPDATE", 1), nil)
			},
			expectedPerson: models.Person{
				Id:      345,
				Name:    "qwerty",
				Work:    "work",
				Address: "address",
				Age:     12,
			},
		},
		{
			nameTest: "no_person_to_update",
			ctx:      context.Background(),
			person: models.Person{
				Id:      345,
				Name:    "qwerty",
				Work:    "work",
				Address: "address",
				Age:     12,
			},
			toUpdate: models.Person{
				Id:      0,
				Name:    "",
				Work:    "newwork",
				Address: "",
				Age:     0,
			},
			mockBehavior: func(ctx context.Context, person *models.Person, toUpdate *models.Person) {
				mockPool.EXPECT().Exec(ctx, "UPDATE persons_ SET work_ = $1 WHERE id_ = $2", person.Work, person.Id).Return(pgxmock.NewResult("UPDATE", 0), nil)
			},
		},
		{
			nameTest: "exec_error",
			ctx:      context.Background(),
			person: models.Person{
				Id:      345,
				Name:    "qwerty",
				Work:    "work",
				Address: "address",
				Age:     12,
			},
			toUpdate: models.Person{
				Id:      0,
				Name:    "",
				Work:    "newwork",
				Address: "",
				Age:     0,
			},
			mockBehavior: func(ctx context.Context, person *models.Person, toUpdate *models.Person) {
				mockPool.EXPECT().Exec(ctx, "UPDATE persons_ SET work_ = $1 WHERE id_ = $2", person.Work, person.Id).Return(nil, errors.New("exec_error"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.person, &testCase.toUpdate)

			got, err := r.Update(testCase.ctx, &testCase.person, &testCase.toUpdate)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedPerson, *got)
			case "invalid_inputs_id":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "no_person_to_update":
				assert.Equal(t, errs.ErrNoContent, err)
			case "exec_error":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestFormRepo_Delete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewPersonRepo(&db)

	type mockBehavior func(ctx context.Context, id int)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		id           int
		mockBehavior mockBehavior
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			id:       345,
			mockBehavior: func(ctx context.Context, id int) {
				mockPool.EXPECT().Exec(ctx, "DELETE FROM persons_ WHERE id_ = $1", id).Return(pgxmock.NewResult("DELETE", 1), nil)
			},
		},
		{
			nameTest: "no_person_to_delete",
			ctx:      context.Background(),
			id:       345,
			mockBehavior: func(ctx context.Context, id int) {
				mockPool.EXPECT().Exec(ctx, "DELETE FROM persons_ WHERE id_ = $1", id).Return(pgxmock.NewResult("DELETE", 0), nil)
			},
		},
		{
			nameTest: "exec_error",
			ctx:      context.Background(),
			id:       345,
			mockBehavior: func(ctx context.Context, id int) {
				mockPool.EXPECT().Exec(ctx, "DELETE FROM persons_ WHERE id_ = $1", id).Return(nil, errors.New("exec_error"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.id)

			err := r.Delete(testCase.ctx, testCase.id)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
			case "invalid_inputs":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "no_person_to_delete":
				assert.Equal(t, errs.ErrNoContent, err)
			case "exec_error":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}
