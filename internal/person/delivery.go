package person

import "github.com/gin-gonic/gin"

type Handlers interface {
	Create() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Update() gin.HandlerFunc
	GetById() gin.HandlerFunc
	GetAll() gin.HandlerFunc
}
