package controller

import (
	"encoding/json"
	"errors"
	"fmt"

	// "reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/ohler55/ojg/pretty"

	general_goutils "github.com/danielcomboni/general-go-utils"
	"github.com/danielcomboni/general-repo-service-controller-utils/responses"
	"github.com/danielcomboni/general-repo-service-controller-utils/service"
)

func CreateWithServiceFuncSpecified[T any](model *T, c *gin.Context, fnServiceCreate func(t T) (T, responses.GenericResponse, error)) {

	//fnLogger func(useDefaultLogMessage bool,i ...interface{}) (interface, error), fnLoggeArgs ...interface{}
	// )

	data, err := c.GetRawData()
	if err != nil {
		msg := "get raw failed: " + err.Error()
		// if general_goutils.IsGreaterThan(len(fnLoggeArgs), 0) {
		// 	fnLogger(fnLoggeArgs)
		// }
		general_goutils.Logger.Error(msg)
		return
	}

	//var t T
	array, err := general_goutils.GetFromByteArray(data)
	if err != nil {
		msg := "failed to unmarshal " + err.Error()
		general_goutils.Logger.Error(msg)
		return
	}

	err = mapstructure.Decode(general_goutils.SafeGetFromInterface(array.Data(), "$.data"), &model)
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

	// save (insert) to database
	created, res, err := fnServiceCreate(*model)
	if err != nil {

		msg := fmt.Sprintf("failed to save record: %v", err)
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

	// g, _ := gabs.Consume(map[string]interface{}{
	// 	"message": "successful",
	// 	"body": created,
	// })

	c.JSON(responses.Created, responses.SetResponse(responses.Created, "successful", created))

}

// CreateWithoutServiceFuncSpecified_AndCheckPropertyPresence inserts into the database and checks for the isnerted row based on the
// property provided
func CreateWithoutServiceFuncSpecified_AndCheckPropertyPresence[T any](property []string, funcAuth func(*gin.Context) (*gin.Context, bool, string)) gin.HandlerFunc {
	return func(c *gin.Context) {

		if funcAuth != nil {
			_, flag, msg := funcAuth(c)
			if !flag {
				c.JSON(responses.UnAuthorized, responses.SetResponse(responses.UnAuthorized, msg, errors.New("failed to authenticate")))
				return
			}
		}

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

		// save (insert) to database
		created, res, err := service.CreateHttpWithPropertyCheck[T](&model, property...)
		if err != nil {

			msg := fmt.Sprintf("failed to save record: %v", err)
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

		c.JSON(responses.Created, responses.SetResponse(responses.Created, "successful", created))

	}

}

func CreateWithoutServiceFuncSpecified_CheckDuplicatesFirst_AndCheckPropertyPresence[T any](duplicateCheckProperties []string, property []string,
	funcAuth func(*gin.Context) (*gin.Context, bool, string)) gin.HandlerFunc {
	return func(c *gin.Context) {

		if funcAuth != nil {
			_, flag, msg := funcAuth(c)
			if !flag {
				c.JSON(responses.UnAuthorized, responses.SetResponse(responses.UnAuthorized, msg, errors.New("failed to authenticate")))
				return
			}
		}

		data, err := c.GetRawData()
		general_goutils.Logger.Info(fmt.Sprintf("incoming request data %v", string(data)))
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

		dupChecks := make(map[string]interface{})

		for _, p := range duplicateCheckProperties {

			dupChecks[strings.ToLower(general_goutils.ToSnakeCase(p))] = general_goutils.SafeGetFromInterfaceGeneric[T](model, fmt.Sprintf("$.%v", general_goutils.ToCamelCaseLower(p)))

		}

		created, res, err := service.CreateWithPriorCheckForDuplicateOfAssociatedEntity[T, T](model, dupChecks, property...)

		if err != nil {
			msg := fmt.Sprintf("failed to save record: %v", err)
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

		c.JSON(responses.Created, responses.SetResponse(responses.Created, "successful", created))

	}

}
