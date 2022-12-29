package models

import (
	"github.com/gin-gonic/gin"
)

type Auth func(*gin.Context) (*gin.Context, bool, string)

type Request struct {
	ParamNames  []QueryStructure       `json:"paramNames,omitempty"`
	Preloads    []string               `json:"preloads,omitempty"`
	QueryMap    map[string]interface{} `json:"queryMap,omitempty"`
	Property    []string               `json:"property,omitempty"`
	ShouldAuth  bool                   `json:"shouldAuth,omitempty"`
	RequestType string                 `json:"requestType"`
	GeneralAuth Auth                   `json:"generalAuth,omitempty"`
}

// HandleAuth
// returns true if authenticated
// returns the custom auth fail message
func HandleAuth(request Request, c *gin.Context) (bool, string) {
	if request.ShouldAuth {
		_, flag, msg := request.GeneralAuth(c)
		return flag, msg
	}
	return true, "authenticated"
}
