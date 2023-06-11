package controller

import "github.com/go-playground/validator/v10"

// import (
//
//	"encoding/json"
//	"errors"
//	"fmt"
//	general_goutils "github.com/danielcomboni/general-go-utils"
//	"github.com/danielcomboni/general-repo-service-controller-utils/models"
//	"github.com/danielcomboni/general-repo-service-controller-utils/repo"
//	"github.com/danielcomboni/general-repo-service-controller-utils/responses"
//	"github.com/gin-gonic/gin"
//	"github.com/go-playground/validator/v10"
//	"github.com/ohler55/ojg/pretty"
//	"gorm.io/gorm"
//	"reflect"
//	"strconv"
//	"strings"
//	"time"
//
// )
var validate = validator.New()

//
//type DeleteHandler[T any] struct {
//	EntityName                string `json:"entityName"`
//	ShouldAuthenticate        bool   `json:"shouldAuthenticate"`
//	ShouldRecordAuditHistory  bool   `json:"ShouldRecordAuditHistory"`
//	AuthenticateUserPerEntity func(c *gin.Context, entityName string, action string) (context *gin.Context, flag bool, message string, JWTGenericUserClaim interface{}, response responses.GenericResponse)
//	SetLogTag                 func(c *gin.Context) models.LogTag
//}
//
//func DeleteSoftlyById[T any](handler DeleteHandler[T]) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		var logTag models.LogTag
//		if handler.SetLogTag != nil {
//			handler.SetLogTag(c)
//		}
//		// authenticate
//		if handler.ShouldAuthenticate {
//			_, flag, _, _, res := handler.AuthenticateUserPerEntity(c, handler.EntityName, "DELETE")
//
//			if !flag {
//				c.JSON(res.Status, res.Data)
//				return
//			}
//		}
//
//		id := c.Param("id")
//		name := reflect.TypeOf(*new(T)).Name()
//		general_goutils.Logger.Info(fmt.Sprintf("%v --> deleting %v by ID: %v", models.AddLogTag(logTag), getEntityName[T](), id))
//
//		if general_goutils.IsNullOrEmpty(id) || general_goutils.ConvertStrToInt64(id) <= 0 {
//			msg := "record id is required"
//			general_goutils.Logger.Error(fmt.Sprintf("%v --> %v", models.AddLogTag(logTag), msg))
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, errors.New(msg).Error()))
//			return
//		}
//
//		rowsAffected, tx, err := repo.DeleteSoftById[T](id)
//
//		if err != nil {
//			msg := "something went wrong"
//			general_goutils.Logger.Error(fmt.Sprintf("%v --> %v", models.AddLogTag(logTag), msg))
//			tx.Rollback()
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, err.Error()))
//			return
//		}
//
//		if rowsAffected <= 0 {
//			msg := "record not deleted"
//			reflect.TypeOf(*new(T)).Name()
//			general_goutils.Logger.Info(fmt.Sprintf("%v --> deleted %v of ID %v", models.AddLogTag(logTag), name, id))
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, errors.New(msg)))
//			return
//		}
//
//		tx.Commit()
//		c.JSON(200, responses.SetResponse(200, "deleted", general_goutils.ConvertStrToInt64(id)))
//
//		if handler.ShouldRecordAuditHistory {
//			models.SaveDeleteAuditHistoryFromGenericResponse[T](rowsAffected, c)
//		}
//
//	}
//}
//
//type CreateHandler[T any] struct {
//	EntityName string `json:"entityName"`
//
//	ShouldAuthenticate bool `json:"shouldAuthenticate"`
//	AutoCommit         bool
//
//	ShouldRecordAuditHistory bool `json:"ShouldRecordAuditHistory"`
//
//	UseCustomValidator bool `json:"useCustomValidator"`
//
//	CustomValidatedModel func(t T, tag ...models.LogTag) interface{} `json:"customerValidatorModel"`
//
//	ShouldCheckDuplicates bool `json:"shouldCheckDuplicates"`
//
//	DuplicateParams func(t T, tag ...models.LogTag) map[string]interface{} `json:"duplicateParams"`
//
//	ShouldCheckDuplicatesManually bool
//
//	DuplicateParamsManually func(t T, tag ...models.LogTag) (bool, string, responses.GenericResponse)
//
//	ShouldValidateModelValues bool `json:"ShouldValidateModelValues"`
//
//	ValidateModelValues func(t T, tag ...models.LogTag) (flag bool, msg string, err error)
//
//	ShouldModifyModelValue bool `json:"ShouldModifyModelValue"`
//
//	ModifyModelValue func(t T, tag ...models.LogTag) T
//	CheckIdCreated   bool
//
//	AfterSuccessfulCreate     func(resp responses.GenericResponse, priorActionResults interface{}, t T, tx *gorm.DB, tag models.LogTag)
//	PriorActions              func(t T, tag ...models.LogTag) (priorActions PriorActionsModels[T])
//	CheckId                   bool
//	ResponseAfter             bool
//	ReturnValue               responses.GenericResponse
//	AuthenticateUserPerEntity func(c *gin.Context, entityName string, action string) (context *gin.Context, flag bool, message string, JWTGenericUserClaim interface{}, response responses.GenericResponse)
//	SetLogTag                 func(c *gin.Context) models.LogTag
//}
//
//type PriorActionsModels[T any] struct {
//	ClientMessage  string
//	LoggingMessage string
//	ModelValue     T
//	ResponseValue  responses.GenericResponse
//	ShouldHalt     bool
//}
//
//func Create[T any](handler *CreateHandler[T]) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		var logTag models.LogTag
//		if handler.SetLogTag != nil {
//			handler.SetLogTag(c)
//		}
//		// authenticate
//		if handler.ShouldAuthenticate {
//			_, flag, _, _, res := handler.AuthenticateUserPerEntity(c, handler.EntityName, "CREATE")
//
//			if !flag {
//				c.JSON(res.Status, res.Data)
//				return
//			}
//		}
//
//		general_goutils.Logger.Info(fmt.Sprintf("%v --> creating %v", models.AddLogTag(logTag), getEntityName[T]()))
//
//		// map to extract the json object
//		model, flag, err, msg := GetRawData[T](c)
//		if err != nil {
//			general_goutils.Logger.Error(fmt.Sprintf("%v --> %v", models.AddLogTag(logTag), msg))
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, err.Error()))
//			return
//		}
//
//		if !flag {
//			msg = fmt.Sprintf("%v: something went wrong", msg)
//			general_goutils.Logger.Error(fmt.Sprintf("%v --> %v", models.AddLogTag(logTag), msg))
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, err.Error()))
//			return
//		}
//		if handler.UseCustomValidator {
//			flag, err = ValidateModel(handler.CustomValidatedModel(model, logTag))
//			if err != nil {
//				general_goutils.Logger.Error(fmt.Sprintf("%v --> %v", models.AddLogTag(logTag), err.Error()))
//				c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, "failed to validate incoming object",
//					fmt.Sprintf("failed to validate incoming object: %v", err.Error())))
//				return
//			}
//		} else {
//
//			flag, err = ValidateModel[T](model)
//
//			if err != nil {
//				general_goutils.Logger.Error(fmt.Sprintf("%v --> failed to validate incoming object%v", models.AddLogTag(logTag), err.Error()))
//				c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, "failed to validate incoming object",
//					fmt.Sprintf("failed to validate incoming object: %v", err.Error())))
//				return
//			}
//
//			if !flag {
//				general_goutils.Logger.Error(fmt.Sprintf("%v --> false flag for incoming object map", models.AddLogTag(logTag)))
//				c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, "failed to validate request object"))
//				return
//			}
//		}
//
//		if !flag {
//			general_goutils.Logger.Error(fmt.Sprintf("%v --> failed to validate request object", models.AddLogTag(logTag)))
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, "failed to validate request object"))
//			return
//		}
//
//		if handler.ShouldValidateModelValues {
//			flag, msg, _ := handler.ValidateModelValues(model, logTag)
//			if !flag {
//				general_goutils.Logger.Error(msg)
//				c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, msg))
//				return
//			}
//		}
//
//		if handler.ShouldCheckDuplicates && !handler.ShouldCheckDuplicatesManually {
//			if CheckDuplicatesIntId[T](handler.DuplicateParams(model, logTag), logTag) {
//				res := responses.SetResponse(responses.ConflictOrDuplicateOrAlreadyExists, "already exists", nil)
//				c.JSON(res.Status, res)
//				return
//			}
//		}
//
//		if handler.ShouldCheckDuplicatesManually {
//			isDup, _, resp := handler.DuplicateParamsManually(model, logTag)
//			if isDup {
//				res := responses.SetResponse(resp.Status, resp.Message, nil)
//				c.JSON(res.Status, res)
//				return
//			}
//		}
//
//		if handler.ShouldModifyModelValue {
//			model = handler.ModifyModelValue(model, logTag)
//		}
//
//		if handler.PriorActions != nil {
//			priorRes := handler.PriorActions(model, logTag)
//
//			// prior checks/actions
//			if priorRes.ShouldHalt {
//
//				general_goutils.Logger.Error(fmt.Sprintf("%v --> failed prior action tests: %v", models.AddLogTag(logTag), priorRes.LoggingMessage))
//
//				c.JSON(priorRes.ResponseValue.Status, priorRes.ResponseValue)
//				return
//			}
//
//		}
//
//		_, res, tx, err := repo.CreateWithPropertyCheckHttpResponse[T](&model)
//
//		if err != nil {
//			tx.Rollback()
//			general_goutils.Logger.Error(fmt.Sprintf("%v --> failed to create and rollback done: %v", models.AddLogTag(logTag), err.Error()))
//			c.JSON(res.Status, res)
//			return
//		}
//
//		modelId := general_goutils.ConvertStrToInt64(pretty.JSON(general_goutils.SafeGet(pretty.JSON(model), "$.id")))
//		general_goutils.Logger.Info(fmt.Sprintf("%v --> created ID: %v", models.AddLogTag(logTag), modelId))
//		if general_goutils.IsNullOrEmpty(modelId) || modelId <= 0 {
//			tx.Rollback()
//			general_goutils.Logger.Error(fmt.Sprintf("%v --> failed to create and rollback done: ID returned is zero or empty", models.AddLogTag(logTag)))
//			c.JSON(responses.BadRequest, res)
//			return
//		}
//
//		if handler.AfterSuccessfulCreate != nil {
//			handler.AfterSuccessfulCreate(res, nil, model, tx, logTag)
//			if handler.ResponseAfter {
//
//				if handler.ShouldRecordAuditHistory && modelId > 0 {
//					models.SaveNewAuditHistoryFromGenericResponse[T](handler.ReturnValue, c, nil)
//				}
//
//				c.JSON(handler.ReturnValue.Status, handler.ReturnValue)
//				return
//			}
//		}
//
//		if handler.ShouldRecordAuditHistory && modelId > 0 {
//			models.SaveNewAuditHistoryFromGenericResponse[T](res, c, nil)
//		}
//
//		if handler.AutoCommit {
//			tx.Commit()
//		}
//
//		general_goutils.Logger.Warn(fmt.Sprintf("%v --> successfully saved", models.AddLogTag(logTag)))
//		c.JSON(responses.Created, res)
//
//	}
//}
//
//type UpdateHandler[T any] struct {
//	EntityName string `json:"entityName"`
//
//	ShouldAuthenticate bool `json:"shouldAuthenticate"`
//
//	ShouldRecordAuditHistory bool `json:"ShouldRecordAuditHistory"`
//
//	UseCustomValidator   bool                                        `json:"useCustomValidator"`
//	CustomValidatedModel func(t T, tag ...models.LogTag) interface{} `json:"customerValidatorModel"`
//
//	ShouldCheckDuplicates bool `json:"shouldCheckDuplicates"`
//
//	DuplicateParams func(t T, tag ...models.LogTag) map[string]interface{} `json:"duplicateParams"`
//
//	ShouldValidateModelValues bool `json:"ShouldValidateModelValues"`
//
//	ValidateModelValues func(t T, tag ...models.LogTag) (flag bool, msg string, err error)
//
//	ShouldModifyModelValue bool `json:"ShouldModifyModelValue"`
//
//	ModifyModelValue func(t T, tag ...models.LogTag) T
//
//	AuthenticateUserPerEntity func(c *gin.Context, entityName string, action string) (context *gin.Context, flag bool, message string, JWTGenericUserClaim interface{}, response responses.GenericResponse)
//	SetLogTag                 func(c *gin.Context) models.LogTag
//}
//
//func UpdateById[T any](handler *UpdateHandler[T]) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		var logTag models.LogTag
//		if handler.SetLogTag != nil {
//			handler.SetLogTag(c)
//		}
//		// authenticate
//		if handler.ShouldAuthenticate {
//			_, flag, _, _, res := handler.AuthenticateUserPerEntity(c, handler.EntityName, "CREATE")
//
//			if !flag {
//				c.JSON(res.Status, res.Data)
//				return
//			}
//		}
//
//		general_goutils.Logger.Info(fmt.Sprintf("%v --> updating %v", models.AddLogTag(logTag), getEntityName[T]()))
//
//		// map to extract the json object
//		model, flag, err, msg := GetRawData[T](c)
//		if err != nil {
//			general_goutils.Logger.Error(fmt.Sprintf("%v --> %v", models.AddLogTag(logTag), msg))
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, err.Error()))
//			return
//		}
//
//		if !flag {
//			msg = fmt.Sprintf("%v: something went wrong", msg)
//			general_goutils.Logger.Error(fmt.Sprintf("%v --> %v", models.AddLogTag(logTag), msg))
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, err.Error()))
//			return
//		}
//
//		if handler.UseCustomValidator {
//			flag, err = ValidateModel(handler.CustomValidatedModel(model, logTag))
//			if err != nil {
//				general_goutils.Logger.Error(fmt.Sprintf("%v --> %v", models.AddLogTag(logTag), err.Error()))
//				c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, "failed to validate incoming object",
//					fmt.Sprintf("failed to validate incoming object: %v", err.Error())))
//				return
//			}
//		} else {
//
//			flag, err = ValidateModel[T](model)
//
//			if err != nil {
//				general_goutils.Logger.Error(fmt.Sprintf("%v --> failed to validate incoming object%v", models.AddLogTag(logTag), err.Error()))
//				c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, "failed to validate incoming object",
//					fmt.Sprintf("failed to validate incoming object: %v", err.Error())))
//				return
//			}
//
//			if !flag {
//				general_goutils.Logger.Error(fmt.Sprintf("%v --> false flag for incoming object map", models.AddLogTag(logTag)))
//				c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, "failed to validate request object"))
//				return
//			}
//		}
//
//		if !flag {
//			general_goutils.Logger.Error(fmt.Sprintf("%v --> failed to validate request object", models.AddLogTag(logTag)))
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, "failed to validate request object"))
//			return
//		}
//
//		if handler.ShouldValidateModelValues {
//			flag, msg, _ := handler.ValidateModelValues(model, logTag)
//			if !flag {
//				general_goutils.Logger.Error(msg)
//				c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, msg))
//				return
//			}
//		}
//
//		id := c.Param("id")
//		if general_goutils.IsNullOrEmpty(id) || general_goutils.ConvertStrToInt64(id) <= 0 {
//			msg := "record id is required"
//			general_goutils.Logger.Error(msg)
//			general_goutils.Logger.Error(fmt.Sprintf("id found: %v", id))
//
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, errors.New(msg).Error()))
//			return
//		}
//
//		if handler.ShouldModifyModelValue {
//			model = handler.ModifyModelValue(model, logTag)
//		}
//
//		//_, res, err := service.UpdateHttpWithPropertyCheck[T](&model, id)
//		updated, tx, err := repo.UpdateByIdWithPropertyCheck[T](model, id)
//
//		if err != nil {
//			tx.Rollback()
//			general_goutils.Logger.Info(fmt.Sprintf("%v --> rolling back: %v", models.AddLogTag(logTag), err.Error()))
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, "not successful", nil))
//			return
//		}
//
//		res := responses.SetResponse(responses.Created, "successful", updated)
//		var commit bool
//
//		if handler.ShouldRecordAuditHistory {
//			models.SaveUpdateAuditHistoryFromGenericResponse[T](res, c, func(flag bool) {
//				commit = flag
//			})
//		}
//
//		if commit {
//			tx.Commit()
//			general_goutils.Logger.Info(fmt.Sprintf("%v --> committing", models.AddLogTag(logTag)))
//			c.JSON(res.Status, res)
//		} else {
//			tx.Rollback()
//			general_goutils.Logger.Warn(fmt.Sprintf("%v --> rolling back", models.AddLogTag(logTag)))
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, "not successful. try again or contact admin for help.", nil))
//			return
//		}
//
//	}
//}
//
//type GetByIdHandler[T any] struct {
//	EntityName                string   `json:"entityName"`
//	ShouldAuthenticate        bool     `json:"shouldAuthenticate"`
//	Preloads                  []string `json:"preloads"`
//	AuthenticateUserPerEntity func(c *gin.Context, entityName string, action string) (context *gin.Context, flag bool, message string, JWTGenericUserClaim interface{}, response responses.GenericResponse)
//	SetLogTag                 func(c *gin.Context) models.LogTag
//}
//
//func GetById[T any](handler GetByIdHandler[T]) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		var logTag models.LogTag
//		if handler.SetLogTag != nil {
//			handler.SetLogTag(c)
//		}
//		// authenticate
//		if handler.ShouldAuthenticate {
//			_, flag, _, _, res := handler.AuthenticateUserPerEntity(c, handler.EntityName, "CREATE")
//
//			if !flag {
//				c.JSON(res.Status, res.Data)
//				return
//			}
//		}
//
//		id := c.Param("id")
//
//		general_goutils.Logger.Info(fmt.Sprintf("%v --> fetching single %v by ID: %v", models.AddLogTag(logTag), getEntityName[T](), id))
//
//		if general_goutils.IsNullOrEmpty(id) || general_goutils.ConvertStrToInt64(id) <= 0 {
//			msg := "record id is required"
//			general_goutils.Logger.Error(msg)
//			c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, msg, errors.New(msg).Error()))
//			return
//		}
//
//		rows, err := repo.GetOneById[T](general_goutils.ConvertStrToInt64(id), handler.Preloads...)
//		if err != nil {
//			c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "failed to retrieve value", err.Error()))
//			return
//		}
//
//		if general_goutils.IsNullOrEmpty(general_goutils.SafeGet(pretty.JSON(rows), "$.id")) ||
//			general_goutils.IsLessThanOrEqualTo(general_goutils.ConvertStrToInt64(general_goutils.SafeGetToString(pretty.JSON(rows), "$.id")), 0) {
//			c.JSON(responses.OK, responses.SetResponse(responses.OK, "record not found", nil))
//			return
//		}
//
//		c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))
//
//	}
//}
//
//type URLParamType interface {
//	int | int64 | float64 | string | bool
//}
//
//type URLParam struct {
//	Name              string      `json:"name"`
//	Type              string      `json:"type"`
//	DefaultValue      interface{} `json:"defaultValue"`
//	AddToFilterFields bool        `json:"addToFilterFields"`
//	FilterFieldName   string
//	FailureMessage    string
//	FailureMessageLog string
//	ShouldHalt        bool
//	Validate          func(param interface{}) bool
//	Ignore            bool
//}
//
//type QueryParam struct {
//	Name              string      `json:"name"`
//	Type              string      `json:"type"`
//	DefaultValue      interface{} `json:"defaultValue"`
//	AddToFilterFields bool        `json:"addToFilterFields"`
//	FilterFieldName   string
//	FailureMessage    string
//	FailureMessageLog string
//	ShouldHalt        bool
//	Validate          func(param interface{}) bool
//	Ignore            bool
//}
//
//type GetAllHandler[T any] struct {
//	EntityName         string   `json:"entityName"`
//	ShouldAuthenticate bool     `json:"shouldAuthenticate"`
//	Preloads           []string `json:"preloads"`
//	//URLParams          map[string]interface{} `json:"URLParams"`
//	URLParamValidator         func(c *gin.Context) (flag bool, msg string, err error)
//	URLParams                 []URLParam
//	QueryParams               []QueryParam
//	PriorManipulations        func(c *gin.Context)
//	RunPriorManipulations     bool
//	AuthenticateUserPerEntity func(c *gin.Context, entityName string, action string) (context *gin.Context, flag bool, message string, JWTGenericUserClaim interface{}, response responses.GenericResponse)
//	SetLogTag                 func(c *gin.Context) models.LogTag
//}
//
//func URLParamParse(param URLParam, paramVal string) interface{} {
//	switch strings.ToLower(param.Type) {
//	case "int":
//		val, _ := strconv.Atoi(paramVal)
//		return val
//	case "int64":
//		return general_goutils.ConvertStrToInt64(paramVal)
//	case "float64":
//		return general_goutils.ConvertStrToFloat64(paramVal)
//	case "bool":
//		val, _ := strconv.ParseBool(paramVal)
//		return val
//	}
//	return paramVal
//}
//
//func QueryParamParse(param QueryParam, paramVal string) interface{} {
//	switch strings.ToLower(param.Type) {
//	case "int":
//		val, _ := strconv.Atoi(paramVal)
//		return val
//	case "int64":
//		return general_goutils.ConvertStrToInt64(paramVal)
//	case "float64":
//		return general_goutils.ConvertStrToFloat64(paramVal)
//	case "bool":
//		val, _ := strconv.ParseBool(paramVal)
//		return val
//	case "date":
//		const (
//			layout1 = "2006/01/02"
//			layout2 = "15:04:05"
//			layout3 = "2006-01-02T15:04:05"
//			layout4 = "2006-01-02"
//		)
//		val, _ := time.Parse(layout4, paramVal)
//		return val.Format(layout4)
//	}
//	return paramVal
//}
//
//func GetAllByClientId[T any](handler GetAllHandler[T]) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		var logTag models.LogTag
//		if handler.SetLogTag != nil {
//			handler.SetLogTag(c)
//		}
//		// authenticate
//		if handler.ShouldAuthenticate {
//			_, flag, _, _, res := handler.AuthenticateUserPerEntity(c, handler.EntityName, "CREATE")
//
//			if !flag {
//				c.JSON(res.Status, res.Data)
//				return
//			}
//		}
//
//		general_goutils.Logger.Info(fmt.Sprintf("%v --> fetching many %v by client ID: %v", models.AddLogTag(logTag), getEntityName[T](), logTag.CompanyId))
//
//		if handler.RunPriorManipulations {
//			handler.PriorManipulations(c)
//		}
//
//		queryParams := map[string]interface{}{}
//
//		// add path params
//		for _, param := range handler.URLParams {
//
//			val := URLParamParse(param, c.Param(param.Name))
//			if !param.Validate(val) && param.ShouldHalt {
//				general_goutils.Logger.Error(fmt.Sprintf(param.FailureMessageLog))
//				c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, param.FailureMessage, errors.New(param.FailureMessageLog).Error()))
//				return
//			}
//
//			if param.AddToFilterFields {
//				queryParams[param.FilterFieldName] = val
//			}
//
//		}
//
//		// add query params
//		for _, param := range handler.QueryParams {
//			val := QueryParamParse(param, c.Request.URL.Query().Get(param.Name))
//			isValid := param.Validate(val)
//			if !isValid && param.ShouldHalt && !param.Ignore {
//				general_goutils.Logger.Error(fmt.Sprintf(param.FailureMessageLog))
//				c.JSON(responses.BadRequest, responses.SetResponse(responses.BadRequest, param.FailureMessage, errors.New(param.FailureMessageLog).Error()))
//				return
//			}
//
//			if param.AddToFilterFields && (!param.Ignore && isValid) {
//				queryParams[param.FilterFieldName] = val
//			}
//
//		}
//
//		if IsPaginationLimitPresent(c) && IsPaginationPagePresent(c) {
//			SetPageAndLimitForPagination(c)
//
//			rows, err := repo.GetAllByFieldsWithPagination[T](queryParams, handler.Preloads...)
//
//			if err != nil {
//				c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "failed to retrieve record", err.Error()))
//				return
//			}
//
//			c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))
//		}
//
//		rows, err := repo.GetAllByFieldsWithNoPagination[T](queryParams, handler.Preloads...)
//
//		if err != nil {
//			c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "failed to retrieve records", err.Error()))
//			return
//		}
//
//		c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))
//
//	}
//}
//
//func CheckDuplicatesIntId[T any](queryMap map[string]interface{}, tag models.LogTag) bool {
//	general_goutils.Logger.Info(fmt.Sprintf("%v --> checking for duplicates...", models.AddLogTag(tag)))
//
//	record, err := repo.GetOneByModelPropertiesCheckPropertyPresence[T](queryMap)
//
//	if err != nil {
//		general_goutils.Logger.Error(fmt.Sprintf("%v --> check dup failed: %v", models.AddLogTag(tag), err.Error()))
//		return true
//	}
//
//	r := reflect.ValueOf(record)
//	f := reflect.Indirect(r).FieldByName("Id")
//	general_goutils.Logger.Info(fmt.Sprintf("%v --> actual value: %v", models.AddLogTag(tag), f))
//	if f.String() == "" || f.Int() == 0 {
//		msg := "duplicate not found"
//		general_goutils.Logger.Info(fmt.Sprintf("%v --> %v", models.AddLogTag(tag), msg))
//		return false
//	}
//
//	return true
//}
//
//func CheckDuplicatesReturnRecord[T any](queryMap map[string]interface{}, tags ...models.LogTag) T {
//	tag := tags[0]
//	general_goutils.Logger.Info(fmt.Sprintf("%v --> checking for duplicates...", models.AddLogTag(tag)))
//	record, err := repo.GetOneByModelPropertiesCheckPropertyPresence[T](queryMap)
//
//	if err != nil {
//		general_goutils.Logger.Error(fmt.Sprintf("%v --> check dup failed: %v", models.AddLogTag(tag), err.Error()))
//		return record
//	}
//
//	r := reflect.ValueOf(record)
//	f := reflect.Indirect(r).FieldByName("Id")
//	general_goutils.Logger.Info(fmt.Sprintf("%v --> actual value: %v", models.AddLogTag(tag), f))
//	if f.String() == "" || f.Int() == 0 {
//		msg := "duplicate not found"
//		general_goutils.Logger.Info(fmt.Sprintf("%v --> %v", models.AddLogTag(tag), msg))
//		return record
//	}
//
//	return record
//}
//
//func IsRawDataMappable[T any](c *gin.Context) bool {
//	data, err := c.GetRawData()
//	var model T
//
//	if err != nil {
//		msg := "get raw failed: " + err.Error()
//		general_goutils.Logger.Error(msg)
//		return false
//	}
//
//	err = json.Unmarshal(data, &model)
//
//	if err != nil {
//		msg := fmt.Sprintf("failed to map incoming object---: %v", err.Error())
//		general_goutils.Logger.Error(msg)
//		return false
//	}
//
//	return true
//}
//
//// GetRawData returns model, boolean flag for model mapping,
//func GetRawData[T any](c *gin.Context, tags ...models.LogTag) (T, bool, error, string) {
//
//	var tag models.LogTag
//	if len(tags) > 0 {
//		tag = tags[0]
//	}
//	general_goutils.Logger.Info(fmt.Sprintf("%v --> getting raw data", models.AddLogTag(tag)))
//	data, err := c.GetRawData()
//	var model T
//
//	general_goutils.Logger.Info(fmt.Sprintf("%v raw request body --> %v", models.AddLogTag(tag), string(data)))
//
//	if err != nil {
//		msg := "get raw failed: " + err.Error()
//		general_goutils.Logger.Error(fmt.Sprintf("%v --> msg", models.AddLogTag(tag)))
//		return model, false, err, msg
//	}
//
//	err = json.Unmarshal(data, &model)
//
//	if err != nil {
//		msg := fmt.Sprintf("failed to map incoming object: %v", err.Error())
//		general_goutils.Logger.Error(fmt.Sprintf("%v --> %v", models.AddLogTag(tag), msg))
//		general_goutils.Logger.Info("")
//		general_goutils.Logger.Info(fmt.Sprintf("%v --> the object: %v", models.AddLogTag(tag), string(data)))
//		return model, false, err, "failed to map incoming object"
//	}
//
//	return model, true, nil, "model mapped"
//}
//
//func IsPaginationLimitPresent(c *gin.Context) bool {
//	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
//	return limit > 0
//}
//
//func IsPaginationPagePresent(c *gin.Context) bool {
//	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
//	return page > 0
//}
//
//func SetPageAndLimitForPagination(c *gin.Context, tags ...models.LogTag) {
//	tag := tags[0]
//	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
//	sort := c.Request.URL.Query().Get("sort")
//	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
//	general_goutils.Logger.Info(fmt.Sprintf("%v --> page: %v and limit: %v", models.AddLogTag(tag), page, limit))
//
//	if general_goutils.IsGreaterThanOrEqualTo(page, 1) && general_goutils.IsGreaterThanOrEqualTo(limit, 1) {
//		repo.SetPagination(limit, page, sort)
//	}
//
//}
//
//type GenericResponse struct {
//	Status  int                    `json:"status"`
//	Message string                 `json:"message"`
//	Data    map[string]interface{} `json:"data"`
//}
//
//func getActionType(c *gin.Context) string {
//	general_goutils.Logger.Info(fmt.Sprintf("request action: %v", c.Request.Method))
//	switch strings.ToLower(c.Request.Method) {
//	case "get":
//		return "view"
//	case "post":
//		return "create"
//	case "put":
//		return "edit"
//	}
//	return "delete"
//}
//
//func getEntityName[T any]() string {
//	return reflect.TypeOf(*new(T)).Name()
//}
//
//// ValidateModel returns true if model is valid
//func ValidateModel[T any](model T, tags ...models.LogTag) (bool, error) {
//	var tag models.LogTag
//	if len(tags) > 0 {
//		tag = tags[0]
//	}
//	general_goutils.Logger.Info(fmt.Sprintf("%v -->  validate model", models.AddLogTag(tag)))
//	//use the validator library to Validate required fields
//	if validationErr := validate.Struct(model); validationErr != nil {
//		general_goutils.Logger.Info(fmt.Sprintf("%v --> incoming: %v", models.AddLogTag(tag), pretty.JSON(model)))
//		msg := fmt.Sprintf("%v --> failed to Validate incoming object: %v", models.AddLogTag(tag), validationErr)
//		general_goutils.Logger.Error(msg)
//		return false, validationErr
//	}
//
//	return true, nil
//
//}
