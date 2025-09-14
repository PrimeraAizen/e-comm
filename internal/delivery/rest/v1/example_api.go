package v1

import (
	"e-comm/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExampleApi struct {
	exampleService service.Example
}

func NewExampleApi(srv service.Example) *ExampleApi {
	return &ExampleApi{
		exampleService: srv,
	}
}

func (api *ExampleApi) InitRoutes(router *gin.RouterGroup) {
	exampleRoutes := router.Group("/example")
	{
		exampleRoutes.GET("/", api.ExampleEndpoint)
	}
}

func (api *ExampleApi) ExampleEndpoint(c *gin.Context) {
	err := api.exampleService.ExampleMethod()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "success")
}
