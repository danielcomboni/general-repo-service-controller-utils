package auth

import "github.com/gin-gonic/gin"

type Auth func(*gin.Context) (*gin.Context, bool)

type AuthHandler[T any] struct {
	ShouldAuth  bool
	GeneralAuth Auth
}
