package http

import (
	"bmstu-dips-lab1/internal/person"

	"github.com/gin-gonic/gin"
)

func MapPersonRoutes(personGroup *gin.RouterGroup, h person.Handlers) {
	personGroup.POST("", h.Create())
	personGroup.DELETE("/:personid", h.Delete())
	personGroup.PATCH("/:personid", h.Update())
	personGroup.GET("", h.GetAll())
	personGroup.GET("/:personid", h.GetById())
}
