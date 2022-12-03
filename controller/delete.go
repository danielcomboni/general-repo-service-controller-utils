package controller

import (
	"github.com/danielcomboni/general-repo-service-controller-utils/responses"
	"github.com/danielcomboni/general-repo-service-controller-utils/service"
	"github.com/gin-gonic/gin"
)

func DeletePermanentlyById_WithoutServiceFuncSpecified[T any]() gin.HandlerFunc {

	return func(c *gin.Context) {

		id := c.Param("id")
		rowsAffected, err := service.DeletePermanentlyById_WithoutService[T](id)
		if err != nil {
			c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
			return
		}
		c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rowsAffected))
	}
}
