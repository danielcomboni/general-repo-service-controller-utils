package models

import (
	"github.com/gin-gonic/gin"
)

type Auth func(*gin.Context) (*gin.Context, bool, string)

type FnServiceGetAll[T any] func() ([]T, error)

type FnBeforeExecution[T any] func(c *gin.Context) (*gin.Context, T, error)

type GetManyRequest[T any] struct {
	Limit int    `json:"limit,omitempty"`
	Page  int    `json:"page,omitempty"`
	Sort  string `json:"sort,omitempty"`

	Ctx               *gin.Context         `json:"ctx"`
	FnServiceGetMany  FnServiceGetAll[T]   `json:"fnServiceGetAll,omitempty"`
	FnBeforeExecution FnBeforeExecution[T] `json:"fnBeforeExecution,omitempty"`

	GenericRequest GenericRequest[T] `json:"genericRequest"`
}

func (get *GetManyRequest[T]) AddAuth(fnAuth Auth) *GetManyRequest[T] {
	get.GenericRequest.GeneralAuth = fnAuth
	return get
}

func (get *GetManyRequest[T]) BeforeExecution(fn FnBeforeExecution[T]) *GetManyRequest[T] {
	get.FnBeforeExecution = fn
	return get
}

func (get *GetManyRequest[T]) GetManyCustomHandler(fn FnServiceGetAll[T]) *GetManyRequest[T] {
	get.FnServiceGetMany = fn
	return get
}
