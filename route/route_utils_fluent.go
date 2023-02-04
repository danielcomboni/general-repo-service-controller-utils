package routes_utils

import "github.com/gin-gonic/gin"

type FluentRouteUtils[T any] struct {
	
}

func (f *FluentRouteUtils[T]) AddRequest(ednpointGroup *gin.RouterGroup) *FluentRouteUtils[T] {

	return f
}

// import (
// 	"github.com/danielcomboni/general-repo-service-controller-utils/controller"
// 	"github.com/danielcomboni/general-repo-service-controller-utils/models"
// 	"github.com/gin-gonic/gin"
// )

// type CRUD[T any] struct {
// 	GetMany      *models.GetManyRequest[T]
// 	IgnoreGetAll bool
// }

// type FluentRouter[T any] struct {
// 	Router *gin.RouterGroup
// }

// func (f *FluentRouter[T]) InitCRUD(router *gin.Engine, relativePath, domainResource string, crud *CRUD[T]) *FluentRouter[T] {
// 	endpointGroupV1 := router.Group(relativePath)
// 	{

// 		// if(crud.GetMany != nil){
// 			Get(endpointGroupV1, domainResource, relativePath, controller.GetMany[T](crud.GetMany))
// 		// }else {

// 		// }

// 		// Post(endpointGroupV1, domainResource, relativePath, controller.CreateWithoutServiceFuncSpecified_AndCheckPropertyPresence[T]([]string{}, s.AuthDefaultPost))
// 		// Get(endpointGroupV1, domainResource, relativePath, controller.GetAllWithoutServiceFuncSpecifiedWithDefaultPagination[T](paramNames, preloads, s.AuthDefaultGetAll))
// 		// Get(endpointGroupV1, domainResource, fmt.Sprintf("%v/:id", relativePath), controller.GetOneByIdWithoutServiceFuncSpecifiedWith[T](preloads, s.AuthDefaultGetById))
// 		// Put(endpointGroupV1, domainResource, fmt.Sprintf("%v/:id", relativePath), controller.UpdateByIdWithoutServiceFuncSpecified_AndCheckPropertyPresence[T](s.AuthDefaultPut))
// 		// Del(endpointGroupV1, domainResource, fmt.Sprintf("%v/:id", relativePath), controller.DeletePermanentlyById_WithoutServiceFuncSpecified[T](s.AuthDefaultDelete))
// 	}
// 	return f
// }

// func (fr *FluentRouter[T]) AddRouteGetMany(relativePath string, requestHandler *models.GetManyRequest[T]) *FluentRouter[T] {
// 	fr.Router.GET(relativePath, controller.GetMany[T](requestHandler))
// 	return fr
// }
