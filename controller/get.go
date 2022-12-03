package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

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

