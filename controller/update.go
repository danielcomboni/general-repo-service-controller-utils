package controller

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ohler55/ojg/pretty"

	general_goutils "github.com/danielcomboni/general-go-utils"
	// "github.com/danielcomboni/general-repo-service-controller-utils/models"
	"github.com/danielcomboni/general-repo-service-controller-utils/responses"
	"github.com/danielcomboni/general-repo-service-controller-utils/service"
)

func UpdateByIdWithoutServiceFuncSpecified_AndCheckPropertyPresence[T any](funcAuth func(*gin.Context) (*gin.Context, bool, string)) gin.HandlerFunc {
	return func(c *gin.Context) {

		if funcAuth != nil {
			_, flag, msg := funcAuth(c)
			if !flag {
				c.JSON(responses.UnAuthorized, responses.SetResponse(responses.UnAuthorized, msg, errors.New("failed to authenticate")))
				return
			}
		}

		id := c.Param("id")
		data, err := c.GetRawData()

		if err != nil {
			msg := "get raw failed: " + err.Error()
			general_goutils.Logger.Error(msg)
			return
		}

		var model T
		err = json.Unmarshal(data, &model)

		if err != nil {
			msg := fmt.Sprintf("failed to map incoming object: %v", err.Error())
			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, err.Error()))
			return
		}

		//use the validator library to Validate required fields
		if validationErr := validate.Struct(model); validationErr != nil {
			general_goutils.Logger.Info(fmt.Sprintf("incoming: %v", pretty.JSON(model)))
			msg := fmt.Sprintf("failed to Validate incoming object: %v", validationErr)
			general_goutils.Logger.Error(msg)
			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, "error", validationErr.Error()))
			return
		}

		// save (update) to database
		updated, res, _, err := service.UpdateHttpWithPropertyCheck[T](&model, id)
		if err != nil {
			msg := fmt.Sprintf("failed to update record: %v", err)
			general_goutils.Logger.Error(msg)
			c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
			return
		}

		if !general_goutils.IsNullOrEmpty(res.Message) {
			msg := fmt.Sprintf("%v", res.Message)
			general_goutils.Logger.Info(msg)
			c.JSON(res.Status, res)
			return
		}

		c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", updated))

	}
}
