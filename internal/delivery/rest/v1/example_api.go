package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *Handler) InitExampleRoutes(router *gin.RouterGroup) {
	exampleRoutes := router.Group("/example")
	{
		exampleRoutes.GET("/", api.ExampleEndpoint)
	}
}

func (api *Handler) ExampleEndpoint(c *gin.Context) {
	err := api.services.ExampleService.ExampleMethod()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "success")
}
