package routes_utils

import (
	"fmt"

	"github.com/danielcomboni/general-repo-service-controller-utils/controller"
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

type SingleEntityGroupedRouteDefinition[T any] struct {
	RelativePath     string           `json:"relativePath,omitempty"`
	DomainResource   string           `json:"domainResource,omitempty"`
	CustomerHandlers *gin.RouterGroup `json:"handlers,omitempty"`
}

func (s *SingleEntityGroupedRouteDefinition[T]) SetSingleEntityGroupedRouteDefinition(router *gin.Engine, relativePath, domainResource string, customerHandlers ...*gin.RouterGroup) {

	s.DomainResource = domainResource
	s.RelativePath = relativePath

	if len(customerHandlers) == 0 {
		endpointGroupV1 := router.Group(relativePath)
		{
			Post(endpointGroupV1, domainResource, relativePath, controller.CreateWithoutServiceFuncSpecified_AndCheckPropertyPresence[T]())
			Get(endpointGroupV1, domainResource, relativePath, controller.GetAllWithoutServiceFuncSpecifiedWithNoPagination[T]())
			Get(endpointGroupV1, domainResource, fmt.Sprintf("%v/:id", relativePath), controller.GetOneByIdWithoutServiceFuncSpecifiedWith[T]())
			Put(endpointGroupV1, domainResource, fmt.Sprintf("%v/:id", relativePath), controller.UpdateByIdWithoutServiceFuncSpecified_AndCheckPropertyPresence[T]())
			Del(endpointGroupV1, domainResource, fmt.Sprintf("%v/:id", relativePath), controller.DeletePermanentlyById_WithoutServiceFuncSpecified[T]())
		}
	}else {
		s.CustomerHandlers = customerHandlers[0]
	}


}