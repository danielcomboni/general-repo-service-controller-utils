package routes_utils

import (
	"fmt"

	"github.com/danielcomboni/general-repo-service-controller-utils/controller"
	"github.com/danielcomboni/general-repo-service-controller-utils/models"
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
	RelativePath         string           `json:"relativePath,omitempty"`
	DomainResource       string           `json:"domainResource,omitempty"`
	CustomerHandlers     *gin.RouterGroup `json:"handlers,omitempty"`
	DefaultEndpointGroup *gin.RouterGroup `json:"defaultEndpointGroup,omitempty"`
	AuthDefaultGetAll    models.Auth      `json:"authDefaultGetAll"`
	AuthDefaultGetById   models.Auth      `json:"authDefaultGetById"`
	AuthDefaultPost      models.Auth      `json:"authDefaultPost"`
	AuthDefaultPut       models.Auth      `json:"authDefaultPut"`
	AuthDefaultDelete    models.Auth      `json:"authDefaultDelete"`
}

func (s *SingleEntityGroupedRouteDefinition[T]) SetAuthDefaultDelete(funcAuth func(*gin.Context) (*gin.Context, bool, string)) *SingleEntityGroupedRouteDefinition[T] {
	s.AuthDefaultDelete = funcAuth
	return s
}

func (s *SingleEntityGroupedRouteDefinition[T]) SetAuthDefaultPost(funcAuth func(*gin.Context) (*gin.Context, bool, string)) *SingleEntityGroupedRouteDefinition[T] {
	s.AuthDefaultPost = funcAuth
	return s
}

func (s *SingleEntityGroupedRouteDefinition[T]) SetAuthDefaultGetAll(funcAuth func(*gin.Context) (*gin.Context, bool, string)) *SingleEntityGroupedRouteDefinition[T] {
	s.AuthDefaultGetAll = funcAuth
	return s
}

func (s *SingleEntityGroupedRouteDefinition[T]) SetAuthDefaultGetById(funcAuth func(*gin.Context) (*gin.Context, bool, string)) *SingleEntityGroupedRouteDefinition[T] {
	s.AuthDefaultGetById = funcAuth
	return s
}

func (s *SingleEntityGroupedRouteDefinition[T]) SetAuthDefaultPutfuncAuth(funcAuth func(*gin.Context) (*gin.Context, bool, string)) *SingleEntityGroupedRouteDefinition[T] {
	s.AuthDefaultPut = funcAuth
	return s
}

func (s *SingleEntityGroupedRouteDefinition[T]) SetSingleEntityGroupedRouteDefinitionPaginationQueryParams(router *gin.Engine, relativePath, domainResource string,
	customerHandlers []*gin.RouterGroup, paramNames []models.QueryStructure, preloads []string) *SingleEntityGroupedRouteDefinition[T] {

	s.DomainResource = domainResource
	s.RelativePath = relativePath
	if len(customerHandlers) == 0 || customerHandlers == nil {
		endpointGroupV1 := router.Group(relativePath)
		{
			Post(endpointGroupV1, domainResource, relativePath, controller.CreateWithoutServiceFuncSpecified_AndCheckPropertyPresence[T]([]string{}, s.AuthDefaultPost))
			Get(endpointGroupV1, domainResource, relativePath, controller.GetAllWithoutServiceFuncSpecifiedWithDefaultPagination[T](paramNames, preloads, s.AuthDefaultGetAll))
			Get(endpointGroupV1, domainResource, fmt.Sprintf("%v/:id", relativePath), controller.GetOneByIdWithoutServiceFuncSpecifiedWith[T](preloads, s.AuthDefaultGetById))
			Put(endpointGroupV1, domainResource, fmt.Sprintf("%v/:id", relativePath), controller.UpdateByIdWithoutServiceFuncSpecified_AndCheckPropertyPresence[T](s.AuthDefaultPut))
			Del(endpointGroupV1, domainResource, fmt.Sprintf("%v/:id", relativePath), controller.DeletePermanentlyById_WithoutServiceFuncSpecified[T](s.AuthDefaultDelete))
		}
		s.DefaultEndpointGroup = endpointGroupV1
	} else {
		s.CustomerHandlers = customerHandlers[0]
	}

	return s
}

func (s *SingleEntityGroupedRouteDefinition[T]) AddGetOneUsingPathParams(relativePath, domainResource string, queryParams []models.QueryStructure,
	funcAuth func(*gin.Context) (*gin.Context, bool, string)) *SingleEntityGroupedRouteDefinition[T] {
	path := concatEndpoint(domainResource, relativePath)
	s.DefaultEndpointGroup.GET(path, controller.GetOneByParamsWithoutServiceFuncSpecifiedWith[T](queryParams, funcAuth))
	return s
}

func (s *SingleEntityGroupedRouteDefinition[T]) AddGetAllUsingPathParams(relativePath, domainResource string,
	queryParams []models.QueryStructure, funcAuth func(*gin.Context) (*gin.Context, bool, string), preloads ...string) *SingleEntityGroupedRouteDefinition[T] {
	path := concatEndpoint(domainResource, relativePath)
	s.DefaultEndpointGroup.GET(path, controller.GetAllByParamsWithoutServiceFuncSpecifiedWith[T](queryParams, preloads, funcAuth))
	return s
}

func (s *SingleEntityGroupedRouteDefinition[T]) AddPostWithDuplicateCheckUsingProperties(relativePath, domainResource string,
	duplicateCheckerParams []string, rowAffectedCheckProperties []string, funcAuth func(*gin.Context) (*gin.Context, bool, string)) *SingleEntityGroupedRouteDefinition[T] {
	path := concatEndpoint(domainResource, relativePath)
	s.DefaultEndpointGroup.POST(path, controller.CreateWithoutServiceFuncSpecified_CheckDuplicatesFirst_AndCheckPropertyPresence[T](duplicateCheckerParams, rowAffectedCheckProperties, funcAuth))
	return s
}
