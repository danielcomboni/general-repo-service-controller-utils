package controller

// //

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"reflect"
// 	"strconv"

// 	"github.com/Jeffail/gabs"
// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/mitchellh/mapstructure"
// 	"github.com/ohler55/ojg/pretty"

// 	general_goutils "github.com/danielcomboni/general-go-utils"
// 	"github.com/danielcomboni/general-repo-service-controller-utils/responses"
// )

// const BadRequest = http.StatusBadRequest
// const InternalServerError = http.StatusInternalServerError
// const Created = http.StatusCreated
// const OK = http.StatusOK
// const NotFound = http.StatusNotFound
// const UnAuthorized = http.StatusUnauthorized

// func getRequestReference[T any](row T) string {
// 	r := reflect.ValueOf(row)
// 	f := reflect.Indirect(r).FieldByName("RequestReference")
// 	if f.String() == "" {
// 		value := uuid.New().String()
// 		reflect.ValueOf(&row).Elem().FieldByName("RequestReference").SetString(value)
// 		// reflect.Indirect(r).FieldByName("RequestReference").Set(reflect.ValueOf(value))
// 		return value
// 	}
// 	return f.String()
// }

// func create[T any](model *T, c *gin.Context, fnServiceCreate func(t T) (T, responses.GenericResponse, error)) {

// 	data, err := c.GetRawData()

// 	if err != nil {
// 		msg := "get raw failed: " + err.Error()
// 		general_goutils.Logger.Error(msg)
// 		return
// 	}

// 	var model T
// 	err = json.Unmarshal(data, &model)

// 	if err != nil {
// 		msg := fmt.Sprintf("failed to map incoming object: %v", err.Error())
// 		c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, err.Error()))
// 		return
// 	}

// 	//use the validator library to Validate required fields
// 	if validationErr := validate.Struct(model); validationErr != nil {
// 		general_goutils.Logger.Info(fmt.Sprintf("incoming: %v", pretty.JSON(model)))
// 		msg := fmt.Sprintf("failed to Validate incoming object: %v", validationErr)
// 		general_goutils.Logger.Error(msg)
// 		c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, "error", validationErr.Error()))
// 		return
// 	}

// 	// save (insert) to database
// 	created, res, err := service.CreateHttpWithPropertyCheck[T](&model, property...)
// 	if err != nil {

// 		msg := fmt.Sprintf("failed to save record: %v", err)
// 		general_goutils.Logger.Error(msg)

// 		c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
// 		return
// 	}

// 	if !general_goutils.IsNullOrEmpty(res.Message) {
// 		msg := fmt.Sprintf("%v", res.Message)
// 		general_goutils.Logger.Info(msg)

// 		c.JSON(res.Status, res)
// 		return
// 	}

// 	c.JSON(responses.Created, responses.SetResponse(responses.Created, "successful", created))

// }

// func getSingleByFields[T any](model *T, c *gin.Context, fnServiceGetSingleByFields func(t T) (T, error)) {

// 	var action models.Action

// 	actionMessage := "GET BY SINGLE FIELDS"

// 	var requestReference string
// 	data, err := c.GetRawData()
// 	if err != nil {
// 		msg := "get raw failed: " + err.Error()
// 		utils.Logger.Error(msg)
// 		if general_goutils.HasField[T]("RequestReference") {
// 			requestReference = getRequestReference[T](*model)
// 		}

// 		action.Action = actionMessage

// 		action.RequestReference = requestReference
// 		action.Data = make(map[string]interface{})
// 		action.Data["message"] = msg
// 		action.Message = msg
// 		sendToCollector(action)

// 		return
// 	}

// 	//var t T
// 	array, err := utils.GetFromByteArray(data)
// 	if err != nil {
// 		msg := "failed to unmarshal " + err.Error()
// 		utils.Logger.Error(msg)

// 		action.Action = actionMessage
// 		action.RequestReference = requestReference
// 		action.Data = make(map[string]interface{})
// 		action.Data["message"] = msg
// 		action.Message = msg
// 		sendToCollector(action)

// 		return
// 	}

// 	err = mapstructure.Decode(utils.SafeGetFromInterface(array.Data(), "$.data"), &model)
// 	if err != nil {
// 		msg := fmt.Sprintf("failed to map incoming object: %v", err.Error())
// 		c.JSON(BadRequest, responses.SetResponse(BadRequest, msg, err.Error()))
// 		if general_goutils.HasField[T]("RequestReference") {
// 			requestReference = getRequestReference[T](*model)
// 		}

// 		action.Action = actionMessage
// 		action.RequestReference = requestReference
// 		action.Data = make(map[string]interface{})
// 		action.Data["message"] = msg
// 		action.Message = msg
// 		sendToCollector(action)

// 		return
// 	}

// 	//use the validator library to Validate required fields
// 	if validationErr := validate.Struct(model); validationErr != nil {
// 		msg := fmt.Sprintf("failed to Validate incoming object: %v", validationErr)
// 		utils.Logger.Info(fmt.Sprintf("incoming: %v", pretty.JSON(model)))
// 		logging.LogError(msg)
// 		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", validationErr.Error()))

// 		if general_goutils.HasField[T]("RequestReference") {
// 			requestReference = getRequestReference[T](*model)
// 		}

// 		action.Action = actionMessage
// 		action.RequestReference = requestReference
// 		action.Data = make(map[string]interface{})
// 		action.Data["message"] = msg
// 		action.Message = msg
// 		sendToCollector(action)

// 		return
// 	}

// 	rows, err := fnServiceGetSingleByFields(*model)
// 	if err != nil {
// 		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))

// 		if general_goutils.HasField[T]("RequestReference") {
// 			requestReference = getRequestReference[T](rows)
// 		}

// 		action.Action = actionMessage
// 		action.RequestReference = requestReference
// 		action.Data = make(map[string]interface{})
// 		action.Data["message"] = responses.SetResponse(InternalServerError, "error", err.Error())
// 		action.Message = fmt.Sprintf("failed fnServiceGetSingleByFields: %v ", err.Error())
// 		sendToCollector(action)

// 		return
// 	}
// 	c.JSON(OK, responses.SetResponse(OK, "successful", rows))
// }

// func createBatch[T any](model []T, c *gin.Context, fnServiceCreate func(t []T) ([]T, responses.GenericResponse, error)) {

// 	//Validate the request body
// 	if err := c.BindJSON(&model); err != nil {
// 		logging.LogError(fmt.Sprintf("failed to bind incoming object: %v", err))
// 		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", err.Error()))
// 		return
// 	}

// 	for _, t := range model {
// 		//use the validator library to Validate required fields
// 		if validationErr := validate.Struct(t); validationErr != nil {
// 			logging.LogError(fmt.Sprintf("failed to Validate incoming object: %v", validationErr))
// 			c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", validationErr.Error()))
// 			return
// 		}
// 	}

// 	// save (insert) to database
// 	created, res, err := fnServiceCreate(model)
// 	if err != nil {
// 		logging.LogError(fmt.Sprintf("failed to save record: %v", err))
// 		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
// 		return
// 	}

// 	if !utils.IsNullOrEmpty(res.Message) {
// 		logging.LogIncoming(fmt.Sprintf("%v", res.Message))
// 		c.JSON(res.Status, res)
// 		return
// 	}

// 	c.JSON(Created, responses.SetResponse(Created, "successful", created))

// }

// func updateById[T any](model *T, c *gin.Context, fnServiceUpdate func(t T, id string) (T, error)) {
// 	id := c.Param("id")
// 	data, err := c.GetRawData()
// 	if err != nil {
// 		utils.Logger.Error("get raw failed: " + err.Error())
// 		return
// 	}

// 	//var t T
// 	array, err := utils.GetFromByteArray(data)
// 	if err != nil {
// 		utils.Logger.Error("failed to unmarshal " + err.Error())
// 		return
// 	}

// 	err = mapstructure.Decode(utils.SafeGetFromInterface(array.Data(), "$.data"), &model)
// 	if err != nil {
// 		msg := fmt.Sprintf("failed to map incoming object: %v", err.Error())
// 		c.JSON(BadRequest, responses.SetResponse(BadRequest, msg, err.Error()))
// 		return
// 	}

// 	//use the validator library to Validate required fields
// 	if validationErr := validate.Struct(model); validationErr != nil {
// 		utils.Logger.Info(fmt.Sprintf("incoming: %v", pretty.JSON(model)))
// 		logging.LogError(fmt.Sprintf("failed to Validate incoming object: %v", validationErr))
// 		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", validationErr.Error()))
// 		return
// 	}

// 	// save (insert) to database
// 	created, err := fnServiceUpdate(*model, id)
// 	if err != nil {
// 		logging.LogError(fmt.Sprintf("failed to update record: %v", err))
// 		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
// 		return
// 	}

// 	c.JSON(Created, responses.SetResponse(Created, "successful", created))

// }

// func patchById[T any](model *models.PatchByIdModel, c *gin.Context, fnServicePatch func(object models.PatchByIdModel) (T, error)) {
// 	//Validate the request body
// 	if err := c.BindJSON(&model); err != nil {
// 		logging.LogError(fmt.Sprintf("failed to bind incoming object: %v", err))
// 		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", err.Error()))
// 		return
// 	}

// 	//use the validator library to Validate required fields
// 	if validationErr := validate.Struct(model); validationErr != nil {
// 		logging.LogError(fmt.Sprintf("failed to Validate incoming object: %v", validationErr))
// 		c.JSON(BadRequest, responses.SetResponse(BadRequest, "error", validationErr.Error()))
// 		return
// 	}

// 	// save (insert) to database
// 	created, err := fnServicePatch(*model)
// 	if err != nil {
// 		logging.LogError(fmt.Sprintf("failed to patch record: %v", err))
// 		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
// 		return
// 	}

// 	c.JSON(Created, responses.SetResponse(Created, "successful", created))

// }

// func getAll[T any](c *gin.Context, fnServiceGetAll func() ([]T, error)) {
// 	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
// 	sort := c.Request.URL.Query().Get("sort")
// 	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))

// 	repositories.SetPagination(limit, page, sort)

// 	rows, err := fnServiceGetAll()
// 	if err != nil {

// 		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
// 		return
// 	}
// 	c.JSON(OK, responses.SetResponse(OK, "successful", rows))
// }

// func getAllByClientId[T any](c *gin.Context, fnServiceGetAll func(id string) ([]T, error)) {
// 	id := c.Param("clientId")

// 	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
// 	sort := c.Request.URL.Query().Get("sort")
// 	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))

// 	repositories.SetPagination(limit, page, sort)

// 	rows, err := fnServiceGetAll(id)
// 	if err != nil {
// 		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
// 		return
// 	}
// 	c.JSON(OK, responses.SetResponse(OK, "successful", rows))
// }

// func getAllByOtherPathParamsId[T any](c *gin.Context, fnServiceGetAll func(pathParams ...repositories.PathParams) ([]T, error), pathParams ...string) {
// 	//id := c.Param("clientId")
// 	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
// 	sort := c.Request.URL.Query().Get("sort")
// 	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
// 	repositories.SetPagination(limit, page, sort)
// 	var params []repositories.PathParams
// 	for _, param := range c.Params {
// 		params = append(params, repositories.PathParams{
// 			Key:   param.Key,
// 			Value: param.Value,
// 		})
// 	}
// 	rows, err := fnServiceGetAll(params...)
// 	if err != nil {
// 		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
// 		return
// 	}
// 	c.JSON(OK, responses.SetResponse(OK, "successful", rows))
// }

// func getOneById[T any](c *gin.Context, fnServiceGetOneById func(id string) (T, error)) {
// 	id := c.Param("id")
// 	row, err := fnServiceGetOneById(id)
// 	if err != nil {
// 		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
// 		return
// 	}
// 	if utils.IsNullOrEmpty(utils.SafeGetFromInterface(row, "$.id")) {
// 		c.JSON(OK, responses.SetResponse(NotFound, "not found", nil))
// 		return
// 	}
// 	c.JSON(OK, responses.SetResponse(OK, "successful", row))
// }

// func deleteSoftlyById[T any](c *gin.Context, fnServiceDeleteSoftlyById func(id string) (int64, error)) {
// 	id := c.Param("id")
// 	rowsAffected, err := fnServiceDeleteSoftlyById(id)
// 	if err != nil {
// 		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
// 		return
// 	}
// 	c.JSON(OK, responses.SetResponse(OK, "successful", rowsAffected))
// }

// func deletePermanentlyById[T any](c *gin.Context, fnServiceDeletePermanentlyById func(id string) (int64, error)) {
// 	id := c.Param("id")
// 	rowsAffected, err := fnServiceDeletePermanentlyById(id)
// 	if err != nil {
// 		c.JSON(InternalServerError, responses.SetResponse(InternalServerError, "error", err.Error()))
// 		return
// 	}
// 	c.JSON(OK, responses.SetResponse(OK, "successful", rowsAffected))
// }
