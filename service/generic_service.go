package service

import (
	"fmt"
	"strings"

	general_goutils "github.com/danielcomboni/general-go-utils"
	"github.com/danielcomboni/general-repo-service-controller-utils/repo"
	"github.com/danielcomboni/general-repo-service-controller-utils/responses"
	"github.com/ohler55/ojg/pretty"
)

// Create where T is the entity of interest and model is the instance of the entity
func CreateHttp[T any](model *T) (T, responses.GenericResponse, error) {
	created, err := repo.Create(model)
	if err != nil {
		msg := fmt.Sprintf("error: %v", err.Error())
		general_goutils.Logger.Error(msg)
		return created, responses.SetResponse(responses.InternalServerError, msg, created), err
	}
	return created, responses.SetResponse(responses.Created, "successful", created), err
}

func CreateHttpWithPropertyCheck[T any](model *T, property ...string) (T, responses.GenericResponse, error) {
	created, err := repo.CreateWithPropertyCheck(model, property...)
	if err != nil {
		msg := fmt.Sprintf("error: %v", err.Error())
		general_goutils.Logger.Error(msg)
		return created, responses.SetResponse(responses.InternalServerError, msg, created), err
	}
	return created, responses.SetResponse(responses.Created, "successful", created), err
}

// Create where T is the entity of interest and model is the instance of the entity
func Create[T any](model *T) (T, error) {
	created, err := repo.Create(model)
	if err != nil {
		msg := fmt.Sprintf("error: %v", err.Error())
		general_goutils.Logger.Error(msg)
		return created, err
	}
	return created, err
}

// CreateWithPriorCheckForDuplicateOfAssociatedEntity where T is the entity of interest and A in the associated entity whose duplicate check is to be performed
func CreateWithPriorCheckForDuplicateOfAssociatedEntity[T any, A any](model T, queryMap map[string]interface{}, properties ...string) (T, responses.GenericResponse, error) {

	presence, err := repo.GetOneByModelProperties[A](queryMap)

	if err != nil {
		if strings.ToLower(strings.TrimSpace(err.Error())) != "record not found" {
			// msg := fmt.Sprintf("error when checking for duplicate: %v",err.Error())
			return model, responses.SetResponse(responses.BadRequest, "error occurred when cehcking for duplicate", err), err
		}
	}

	if  !general_goutils.IsNullOrEmpty(general_goutils.SafeGet(pretty.JSON(presence), "$.id")) && pretty.JSON(general_goutils.SafeGet(pretty.JSON(presence), "$.id")) != "0" {
		general_goutils.Logger.Warn("duplicated found")
		return model, responses.SetResponse(responses.ConflictOrDuplicateOrAlreadyExists, "already exists. (duplicate found)", presence), err
	}

	return repo.CreateWithPropertyCheckHttpResponse[T](&model, properties...)
	// return Create[T](&model)

}

// CreateWithPriorCheckForDuplicateOfAssociatedEntity where T is the entity of interest and A in the associated entity whose duplicate check is to be performed
func CreateWithPriorCheckForDuplicateOfAssociatedEntityHttp[T any, A any](model T, queryMap map[string]interface{}) (T, responses.GenericResponse, error) {

	presence, err := repo.GetOneByModelPropertiesCheckIdPresence[A](queryMap)
	if err != nil {
		if strings.ToLower(strings.TrimSpace(err.Error())) != "record not found" {
			return model, responses.SetResponse(responses.BadRequest, "error occurred during duplicate check: "+err.Error(), presence), err
		}
	}

	if !general_goutils.IsNullOrEmpty(   general_goutils.SafeGet(pretty.JSON(presence),"$.id")   ) {
		return model, responses.SetResponse(responses.ConflictOrDuplicateOrAlreadyExists, "already exists", nil), err
	}

	return CreateHttp[T](&model)

}

func GetAllWithNoPagination[T any]() ([]T, error) {
	return repo.GetAllWithNoPagination[T]()
}

func GetAll[T any]() ([]T, error) {
	return repo.GetAll[T]()
}

func GetOneByModelPropertiesCheckIdPresence[T any](queryMap map[string]interface{}) (T, error) {
	return repo.GetOneByModelPropertiesCheckIdPresence[T](queryMap)
}

func GetOneByModelPropertiesCheckPropertyPresence[T any](queryMap map[string]interface{}, propertiesCheck ...string) (T, error) {
	return repo.GetOneByModelPropertiesCheckPropertyPresence[T](queryMap, propertiesCheck...)
}

func GetOneByModelProperties[T any](queryMap map[string]interface{}) (T, error) {
	return repo.GetOneByModelProperties[T](queryMap)
}

func GetOneById[T any](id interface{}) (T, error) {
	return repo.GetOneById[T](id)
}

func updateById[T any](model T, id string) (T, error) {
	return repo.UpdateById[T](model, id)
}

func UpdateHttpWithPropertyCheck[T any](model *T, id interface{}, property ...string) (T, responses.GenericResponse, error) {

	// created, err := repo.CreateWithPropertyCheck(model, property...)
	updated, err := repo.UpdateByIdWithPropertyCheck(model, id, property...)
	if err != nil {
		msg := fmt.Sprintf("error: %v", err.Error())
		general_goutils.Logger.Error(msg)
		return *updated, responses.SetResponse(responses.InternalServerError, msg, updated), err
	}
	return *updated, responses.SetResponse(responses.Created, "successful", updated), err
}

func deleteSoftlyById[T any](id string) (int64, error) {
	return repo.DeleteSoftById[T](id)
}

func DeletePermanentlyById_WithoutService[T any](id interface{}) (int64, error) {
	return repo.DeleteHardById[T](id)
}
