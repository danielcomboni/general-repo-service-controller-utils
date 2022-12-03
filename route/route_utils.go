package routes_utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func concatEndpoint(domainResource, parameter string) string {
	if domainResource[0:1] != "/" {
		domainResource = "/" + domainResource
	}
	if len(parameter) > 0 && parameter[0:1] != "/" {
		parameter = "/" + parameter
	}
	return fmt.Sprintf("%v/%v", domainResource, parameter)
}

func Post(endpointGroup *gin.RouterGroup, domainResource, parameter string, controllerMethod gin.HandlerFunc) {
	relativePath := concatEndpoint(domainResource, parameter)
	endpointGroup.POST(relativePath, controllerMethod)
}

func Get(endpointGroup *gin.RouterGroup, domainResource, parameter string, controllerMethod gin.HandlerFunc) {
	relativePath := concatEndpoint(domainResource, parameter)
	endpointGroup.GET(relativePath, controllerMethod)
}

func Put(endpointGroup *gin.RouterGroup, domainResource, parameter string, controllerMethod gin.HandlerFunc) {
	relativePath := concatEndpoint(domainResource, parameter)
	endpointGroup.PUT(relativePath, controllerMethod)
}

func Del(endpointGroup *gin.RouterGroup, domainResource, parameter string, controllerMethod gin.HandlerFunc) {
	relativePath := concatEndpoint(domainResource, parameter)
	endpointGroup.DELETE(relativePath, controllerMethod)
}
