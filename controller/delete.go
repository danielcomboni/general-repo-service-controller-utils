package controller

import (
	"errors"

	"github.com/danielcomboni/general-repo-service-controller-utils/responses"
	"github.com/danielcomboni/general-repo-service-controller-utils/service"
	"github.com/gin-gonic/gin"
)

func DeletePermanentlyById_WithoutServiceFuncSpecified[T any](funcAuth func(*gin.Context) (*gin.Context, bool, string)) gin.HandlerFunc {
	return func(c *gin.Context) {

		if funcAuth != nil {
			_, flag, msg := funcAuth(c)
			if !flag {
				c.JSON(responses.UnAuthorized, responses.SetResponse(responses.UnAuthorized, msg, errors.New("failed to authenticate")))
				return
			}
		}

		id := c.Param("id")
		rowsAffected, _, err := service.DeletePermanentlyById_WithoutService[T](id)
		if err != nil {
			c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
			return
		}
		c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rowsAffected))
	}
}
