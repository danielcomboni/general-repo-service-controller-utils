package controller

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"

	general_goutils "github.com/danielcomboni/general-go-utils"
	"github.com/danielcomboni/general-repo-service-controller-utils/model"
	"github.com/danielcomboni/general-repo-service-controller-utils/repo"
	"github.com/danielcomboni/general-repo-service-controller-utils/responses"
)

func GetAllWithServiceFuncSpecified_And_WithPagination[T any](c *gin.Context, fnServiceGetAll func() ([]T, error)) {
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	sort := c.Request.URL.Query().Get("sort")
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))

	repo.SetPagination(limit, page, sort)

	rows, err := fnServiceGetAll()
	if err != nil {

		c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
		return
	}
	c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))
}

func GetAllWithServiceFuncSpecified_WithNoPagination[T any](c *gin.Context, fnServiceGetAll func() ([]T, error)) {
	rows, err := fnServiceGetAll()
	if err != nil {
		c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
		return
	}
	c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))
}

func GetAllWithoutServiceFuncSpecifiedWithNoPagination[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {

		rows, err := repo.GetAllWithNoPagination[T]()
		if err != nil {
			c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
			return
		}
		c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))
	}
}

func GetAllWithoutServiceFuncSpecifiedWithPagination[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {

		page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
		sort := c.Request.URL.Query().Get("sort")
		limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))

		repo.SetPagination(limit, page, sort)

		rows, err := repo.GetAll[T]()
		if err != nil {
			c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
			return
		}
		c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))
	}
}

func GetAllWithoutServiceFuncSpecifiedWithPaginationAndWithFieldParams[T any](queryMap map[string]interface{}, preloads ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
		sort := c.Request.URL.Query().Get("sort")
		limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))

		repo.SetPagination(limit, page, sort)

		rows, err := repo.GetAllByFieldsWithPagination[T](queryMap, preloads...)
		if err != nil {
			c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
			return
		}
		c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))
	}
}

func GetAllByParamsWithoutServiceFuncSpecifiedWith[T any](paramNames []model.QueryStructure, preloads ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		queryParams := make(map[string]interface{})

		for _, param := range paramNames {
			queryParams[param.DbTableColumn] = c.Param(param.ParamName)
		}

		page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
		sort := c.Request.URL.Query().Get("sort")
		limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
		general_goutils.Logger.Info(fmt.Sprintf("page: %v and limit: %v",page, limit))
		if general_goutils.IsLessThanOrEqualTo(page,0) && general_goutils.IsLessThanOrEqualTo(limit,0) {

			rows, err := repo.GetAllByFieldsWithNoPagination[T](queryParams, preloads...)
			if err != nil {
				c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
				return
			}
			c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))
			return
		}

		repo.SetPagination(limit, page, sort)

		rows, err := repo.GetAllByFieldsWithPagination[T](queryParams, preloads...)
		if err != nil {
			c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
			return
		}
		c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))

	}
}

func GetOneByIdWithoutServiceFuncSpecifiedWith[T any](preloads ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		rows, err := repo.GetOneById[T](id, preloads...)
		if err != nil {
			c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
			return
		}
		c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))
	}
}

func GetOneByParamsWithoutServiceFuncSpecifiedWith[T any](paramNames ...model.QueryStructure) gin.HandlerFunc {
	return func(c *gin.Context) {

		queryParams := make(map[string]interface{})

		for _, param := range paramNames {
			queryParams[param.DbTableColumn] = c.Param(param.ParamName)
		}

		rows, err := repo.GetOneByModelProperties[T](queryParams)
		if err != nil {
			c.JSON(responses.InternalServerError, responses.SetResponse(responses.InternalServerError, "error", err.Error()))
			return
		}

		firstField := reflect.ValueOf(queryParams).MapKeys()[0]

		r := reflect.ValueOf(rows)
		f := reflect.Indirect(r).FieldByName(general_goutils.ToCamelCase(firstField.String()))

		if general_goutils.IsNullOrEmpty(f) {
			c.JSON(responses.OK, responses.SetResponse(responses.OK, "record does not exist", nil))
			return
		}

		c.JSON(responses.OK, responses.SetResponse(responses.OK, "successful", rows))
	}
}
