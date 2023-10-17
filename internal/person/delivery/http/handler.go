package http

import (
	"bmstu-dips-lab1/internal/person"
	"bmstu-dips-lab1/models"
	"bmstu-dips-lab1/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PersonCreatRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	Work    string `json:"work" binding:"required"`
	Age     int    `json:"age" binding:"required"`
}

type PersonUpdRequest struct {
	Name    *string `json:"name"`
	Address *string `json:"address"`
	Work    *string `json:"work"`
	Age     *int    `json:"age"`
}

type PersonResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Work    string `json:"work"`
	Age     int    `json:"age"`
}

type PersonGetAllResponse struct {
	Persons []*PersonResponse `json:"persons"`
}

type PersonHandlers struct {
	personUC person.UseCase
}

func NewPersonHandlers(personUC person.UseCase) person.Handlers {
	return &PersonHandlers{
		personUC: personUC,
	}
}

func (p *PersonHandlers) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(PersonCreatRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		modelBL := PersonCreatRequestToBL(request)

		createdperson, err := p.personUC.Create(c, modelBL)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.Header("Location", "/api/v1/persons/"+createdperson.Id)

		c.Status(http.StatusCreated)
	}
}

func (p *PersonHandlers) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(PersonUpdRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		const updStr string = "y"
		const updInt int = 1

		modelBL := &models.Person{Id: c.Param("personid")}
		toUpdate := &models.Person{}

		if request.Name != nil {
			modelBL.Name = *request.Name
			toUpdate.Name = updStr
		}

		if request.Address != nil {
			modelBL.Address = *request.Address
			toUpdate.Address = updStr
		}

		if request.Work != nil {
			modelBL.Work = *request.Work
			toUpdate.Work = updStr
		}

		if request.Age != nil {
			modelBL.Age = *request.Age
			toUpdate.Age = updInt
		}

		updatedperson, err := p.personUC.Update(c, modelBL, toUpdate)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, PersonBLToResponse(updatedperson))
	}
}

func (p *PersonHandlers) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		foundperson, err := p.personUC.GetById(c, c.Param("personid"))
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, PersonBLToResponse(foundperson))
	}
}

func (p *PersonHandlers) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := p.personUC.Delete(c, c.Param("personid"))
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.Status(http.StatusOK)
	}
}

func (p *PersonHandlers) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		foundpersons, err := p.personUC.GetAll(c)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, &PersonGetAllResponse{
			Persons: PersonsBLToResponse(foundpersons),
		})
	}
}

func PersonCreatRequestToBL(dto *PersonCreatRequest) *models.Person {
	return &models.Person{
		Name:    dto.Name,
		Address: dto.Address,
		Work:    dto.Work,
		Age:     dto.Age,
	}
}

func PersonBLToResponse(modelBL *models.Person) *PersonResponse {
	return &PersonResponse{
		Id:      modelBL.Id,
		Name:    modelBL.Name,
		Address: modelBL.Address,
		Work:    modelBL.Work,
		Age:     modelBL.Age,
	}
}

func PersonsBLToResponse(persons []*models.Person) []*PersonResponse {
	if persons == nil {
		return nil
	}

	res := make([]*PersonResponse, len(persons))

	for i, p := range persons {
		res[i] = PersonBLToResponse(p)
	}

	return res
}
