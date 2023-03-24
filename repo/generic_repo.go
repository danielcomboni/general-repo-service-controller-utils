package repo

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	gen_utils "github.com/danielcomboni/general-go-utils"
	"github.com/danielcomboni/general-repo-service-controller-utils/responses"
	"github.com/gobeam/stringy"
	"github.com/mitchellh/mapstructure"
	"github.com/ohler55/ojg/pretty"
	"gorm.io/gorm"
)

var Instance *gorm.DB

func RepoInitializer(instance *gorm.DB) {
	Instance = instance
}

func RunMigrations(entityModels ...interface{}) {
	gen_utils.Logger.Info("Database Migration Started...")
	err := Instance.AutoMigrate(entityModels...)
	if err != nil {
		msg := "failed to run migrations: " + err.Error()
		gen_utils.Logger.Error(msg)
	}
	gen_utils.Logger.Info("Database Migration Completed...")
}

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

var TablePagination *Pagination

func SetPagination(limit, page int, sort string) {
	TablePagination = &Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}

func paginationParams() (offset int, limit int) {
	page := TablePagination.Page
	if page == 0 {
		page = 1
	}

	pageSize := TablePagination.Limit

	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset = (page - 1) * pageSize
	limit = pageSize
	return offset, pageSize
}

func Create[T any](model *T) (T, *gorm.DB, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\ncreating a new record: %v", reflect.TypeOf(*new(T)).Name()))

	tx := Instance.Begin()
	var t T
	result := Instance.Create(&model).Scan(&t)
	_, err := result.DB()
	if err != nil {
		msg := fmt.Sprintf("failed to create: %v", err.Error())
		gen_utils.Logger.Error(msg)
		return *model, tx, err
	}

	if gen_utils.IsNullOrEmpty(gen_utils.SafeGetFromInterface(&model, "$.id")) {
		msg := fmt.Sprintf("not saved: %v", result.Error.Error())
		log.Println(msg)
		return *model, tx, errors.New(msg)
	}

	if result.RowsAffected > 0 {
		gen_utils.Logger.Info(fmt.Sprintf("saved to database: id: %v", gen_utils.SafeGetFromInterface(t, "$.id")))
	}
	return t, tx, nil
}

func CreateWithPropertyCheckHttpResponse[T any](model *T, property ...string) (T, responses.GenericResponse, *gorm.DB, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\ncreating a new record: %v", reflect.TypeOf(*new(T)).Name()))
	var t T
	tx := Instance.Begin()
	result := tx.Create(&model).Scan(&t)

	_, err := result.DB()
	if err != nil {
		msg := fmt.Sprintf("failed to create: %v", err.Error())
		gen_utils.Logger.Error(msg)
		return *model, responses.SetResponse(responses.InternalServerError, "error occurred when saving record", err), tx, err
	}

	// if property to check is provided and the entity actually contains the particular property
	if gen_utils.IsGreaterThan(len(property), 0) && gen_utils.HasField[T](property[0]) {
		if gen_utils.IsNullOrEmpty(gen_utils.SafeGetFromInterface(&model, "$."+gen_utils.ToCamelCaseLower(property[0]))) {
			msg := fmt.Sprintf("not saved: %v", result.Error.Error())
			gen_utils.Logger.Error(msg)
			return *model, responses.SetResponse(responses.InternalServerError, "error occurred when saving record", result.Error), tx, errors.New(msg)
		}
	}

	if gen_utils.HasField[T]("Id") {
		if result.RowsAffected > 0 {
			gen_utils.Logger.Info(fmt.Sprintf("saved to database: id: %v", gen_utils.SafeGetFromInterface(t, "$.id")))
			return t, responses.SetResponse(responses.Created, "successful", t), tx, nil
		}
	}

	return t, responses.SetResponse(responses.InternalServerError, "not created. something went wrong", t), tx, nil
}

func CreateWithPropertyCheck[T any](model *T, property ...string) (T, *gorm.DB, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\ncreating a new record: %v", reflect.TypeOf(*new(T)).Name()))
	var t T
	tx := Instance.Begin()
	result := tx.Create(&model).Scan(&t)

	_, err := result.DB()
	if err != nil {
		msg := fmt.Sprintf("failed to create: %v", err.Error())
		gen_utils.Logger.Error(msg)
		return *model, tx, err
	}

	// if property to check is provided and the entity actually contains the particular property
	if gen_utils.IsGreaterThan(len(property), 0) && gen_utils.HasField[T](property[0]) {
		if gen_utils.IsNullOrEmpty(gen_utils.SafeGetFromInterface(&model, "$."+gen_utils.ToCamelCaseLower(property[0]))) {
			msg := fmt.Sprintf("not saved: %v", result.Error.Error())
			log.Println(msg)
			return *model, tx, errors.New(msg)
		}
	}

	if gen_utils.HasField[T]("Id") {
		if result.RowsAffected > 0 {
			gen_utils.Logger.Info(fmt.Sprintf("saved to database: id: %v", gen_utils.SafeGetFromInterface(t, "$.id")))
		}
	}
	return t, tx, nil
}

func GetAllWithNoPagination[T any]() ([]T, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nretreiving collection: %v", reflect.TypeOf(*new(T)).Name()))
	var all []T
	result := Instance.Find(&all)
	_, err := result.DB()
	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return all, err
	}
	return all, nil
}

func GetAll[T any]() ([]T, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nretreiving collection: %v", reflect.TypeOf(*new(T)).Name()))
	var all []T
	offset, limit := paginationParams()
	gen_utils.Logger.Info(fmt.Sprintf("offset: %v, limit: %v", offset, limit))
	result := Instance.Offset(offset).Limit(limit).Find(&all)
	_, err := result.DB()
	if err != nil {
		log.Println(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return all, err
	}
	return all, nil
}

func GetAllWithNoParams[T any]() ([]T, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nretreiving collection: %v", reflect.TypeOf(*new(T)).Name()))
	var all []T
	result := Instance.Find(&all)
	_, err := result.DB()
	if err != nil {
		log.Println(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return all, err
	}
	return all, nil
}

func preloadsHandler(preloads ...string) *gorm.DB {

	instance := Instance

	for _, s := range preloads {
		instance = instance.Preload(s)
	}

	return instance
}

// todo ... fields to be omitted
func omitsHandler(omits ...string) *gorm.DB {

	instance := Instance

	for _, s := range omits {
		instance = instance.Omit(s)
	}

	return instance
}

func GetAllByFieldsWithPagination[T any](queryMap map[string]interface{}, preloads ...string) ([]T, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nretreiving collection: %v", reflect.TypeOf(*new(T)).Name()))
	var all []T
	offset, limit := paginationParams()
	gen_utils.Logger.Info(fmt.Sprintf("offset: %v, limit: %v", offset, limit))

	var instance *gorm.DB

	if len(preloads) == 0 {
		instance = Instance.Offset(offset).Limit(limit).Where(queryMap).Find(&all)
	} else {
		instance = preloadsHandler(preloads...).Offset(offset).Limit(limit).Where(queryMap).Find(&all)
	}

	result := instance
	_, err := result.DB()
	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return all, err
	}
	return all, nil
}

func GetAllByFieldsWithNoPagination[T any](queryMap map[string]interface{}, preloads ...string) ([]T, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nretreiving collection: %v", reflect.TypeOf(*new(T)).Name()))
	var all []T

	var instance *gorm.DB

	if len(preloads) == 0 {
		instance = Instance.Where(queryMap).Find(&all)
	} else {
		instance = preloadsHandler(preloads...).Where(queryMap).Find(&all)
	}

	result := instance
	_, err := result.DB()
	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return all, err
	}
	return all, nil
}

func GetOneById[T any](id interface{}, preloads ...string) (T, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nretreiving single row of: %v by id: %v", reflect.TypeOf(*new(T)).Name(), id))
	var row T

	var instance *gorm.DB
	if len(preloads) == 0 {
		instance = Instance.Where("id=?", id).Find(&row)
	} else {
		instance = preloadsHandler(preloads...).Where("id=?", id).Find(&row)
	}

	result := instance

	_, err := result.DB()
	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return row, err
	}

	return row, nil
}

func GetOneByModelPropertiesCheckIdPresence[T any](queryMap map[string]interface{}) (T, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nretreiving single row of: %v by values: %#v", reflect.TypeOf(*new(T)).Name(), queryMap))
	var row T
	result := Instance.Where(queryMap).First(&row)
	_, err := result.DB()
	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return row, err
	}

	r := reflect.ValueOf(row)
	f := reflect.Indirect(r).FieldByName("Id")
	if f.String() == "" {
		msg := "record not found"
		log.Println(msg)
		return row, errors.New(msg)
	}
	return row, nil
}

func GetOneByModelPropertiesCheckPropertyPresence[T any](queryMap map[string]interface{}, propertiesCheck ...string) (T, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nretreiving single row of: %v by values: %#v", reflect.TypeOf(*new(T)).Name(), queryMap))
	var row T
	result := Instance.Where(queryMap).First(&row)
	_, err := result.DB()

	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return row, err
	}

	if gen_utils.IsGreaterThan(len(propertiesCheck), 0) {
		for _, property := range propertiesCheck {
			r := reflect.ValueOf(row)
			f := reflect.Indirect(r).FieldByName(property)
			if f.String() == "" {
				msg := "record not found"
				gen_utils.Logger.Warn(msg)
				return row, errors.New(msg)
			}
		}
	}

	// r := reflect.ValueOf(row)
	// f := reflect.Indirect(r).FieldByName("Id")
	// if f.String() == "" {
	// 	msg := "record not found"
	// 	log.Println(msg)
	// 	return row, errors.New(msg)
	// }
	return row, nil
}

func GetOneByModelProperties[T any](queryMap map[string]interface{}) (T, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nretreiving single row of: %v by values: %#v", reflect.TypeOf(*new(T)).Name(), queryMap))
	var row T
	result := Instance.Where(queryMap).First(&row)
	_, err := result.DB()

	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to retrieve: %v %v", reflect.TypeOf(*new(T)).Name(), err))
		return row, err
	}

	firstField := reflect.ValueOf(queryMap).MapKeys()[0]

	r := reflect.ValueOf(row)
	f := reflect.Indirect(r).FieldByName(gen_utils.ToCamelCase(firstField.String()))

	if gen_utils.IsNullOrEmpty(f) && gen_utils.ToCamelCase(firstField.String()) != "id" {
		gen_utils.Logger.Info("no record found")
		return *new(T), nil
	}

	return row, nil
}

func PatchById[T any](id, columnName string, value interface{}) (T, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\npatch column: %v row of: %v by id: %v", columnName, reflect.TypeOf(*new(T)).Name(), id))
	one, err := GetOneById[T](id)
	var t2 T
	if err != nil {
		return t2, err
	}

	if err != nil {
		log.Println(fmt.Sprintf("failed to map structure: %v", err))
		return t2, err
	}
	//result := Instance.Where("id=?", id).Update(stringy.New(columnName).SnakeCase("?", "").ToLower(), value).Scan(&one)
	result := Instance.Model(&one).Where("id=?", id).Update(stringy.New(columnName).SnakeCase("?", "").ToLower(), value).Scan(&one)
	rowsAffected := result.RowsAffected
	gen_utils.Logger.Info(fmt.Sprintf("rows affected: %v", gen_utils.ConvertInt64ToStr(rowsAffected)))

	if result.Error != nil {
		log.Println(fmt.Sprintf("failed to patch env: %v", result.Error))
		return t2, result.Error
	}

	if result.RowsAffected == 0 {
		log.Println(fmt.Sprintf("not updated: affected rows: %v", result.RowsAffected))
		return t2, errors.New("not patched")
	}

	return one, nil
}

func UpdateById[T any](t T, id string) (T, *gorm.DB, error) {

	gen_utils.Logger.Info(fmt.Sprintf("\n\nupdating row of: %v by id: %v", reflect.TypeOf(*new(T)).Name(), id))
	one, err := GetOneById[T](id)
	var t2 T
	tx := Instance.Begin()
	if err != nil {
		return t2, tx, err
	}

	err = mapstructure.Decode(t, &one)

	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to map structure: %v", err.Error()))
		return t2, tx, err
	}

	// set the createdAt date and updatedAt

	result := tx.Where("id=?", id).Updates(&one).Scan(&one)
	rowsAffected := result.RowsAffected
	gen_utils.Logger.Info(fmt.Sprintf("rows affected: %v", gen_utils.ConvertInt64ToStr(rowsAffected)))

	if result.Error != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to update env: %v", result.Error))
		return t2, tx, result.Error
	}

	if result.RowsAffected == 0 {
		gen_utils.Logger.Warn(fmt.Sprintf("not updated: affected rows: %v", result.RowsAffected))
		return t2, tx, errors.New("not updated")
	}

	return one, tx, nil
}

func UpdateByIdWithPropertyCheck[T any](t T, id interface{}, property ...string) (T, *gorm.DB, error) {

	gen_utils.Logger.Info(fmt.Sprintf("\n\nupdating row of: %v by id: %v", reflect.TypeOf(*new(T)).Name(), id))
	one, err := GetOneById[T](id)
	var t2 T
	tx := Instance.Begin()
	if err != nil {
		return t2, tx, err
	}

	err = mapstructure.Decode(t, &one)

	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to map structure: %v", err.Error()))
		return t2, tx, err
	}

	// result := Instance.Where("id=?", id).Updates(&one).Scan(&one)
	// result := Instance.Updates(&one).Scan(&one)
	result := Instance.Updates(&one).Scan(&one)
	rowsAffected := result.RowsAffected
	gen_utils.Logger.Info(fmt.Sprintf("rows affected: %v", gen_utils.ConvertInt64ToStr(rowsAffected)))

	if result.Error != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to update env: %v", result.Error))
		return t2, tx, result.Error
	}

	if result.RowsAffected == 0 {
		gen_utils.Logger.Warn(fmt.Sprintf("not updated: affected rows: %v", result.RowsAffected))
		return t2, tx, errors.New("not updated")
	}

	// if property to check is provided and the entity actually contains the particular property
	if gen_utils.IsGreaterThan(len(property), 0) && gen_utils.HasField[T](property[0]) {
		if gen_utils.IsNullOrEmpty((gen_utils.SafeGetFromInterface(one, "$."+gen_utils.ToCamelCaseLower(property[0])))) {
			msg := fmt.Sprintf("not updated: %v", result.Error.Error())
			gen_utils.Logger.Error(msg)
			return one, tx, errors.New(msg)
		}
		gen_utils.Logger.Info(fmt.Sprintf("updated row with %v: %v", property[0], gen_utils.SafeGetFromInterface(one, "$."+gen_utils.ToCamelCaseLower(property[0]))))
	} else {
		gen_utils.Logger.Info(fmt.Sprintf("updated row with id: %v", gen_utils.SafeGetFromInterface(one, "$.id")))
	}
	return one, tx, nil
}

func DeleteHardByFields[T any](queryMap map[string]interface{}) (int64, *gorm.DB, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nhard deleting a row of: %v by fields: %v", reflect.TypeOf(*new(T)).Name(), queryMap))
	one, err := GetOneByModelProperties[T](queryMap)
	var t2 T
	tx := Instance.Begin()
	if err != nil {
		return 0, tx, err
	}

	err = mapstructure.Decode(one, &t2)

	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to map structure: %v", err))
		return 0, tx, err
	}

	r := Instance.Delete(&one).Where(queryMap)

	if r.Error != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to delete row by fields: %v", queryMap))
		gen_utils.Logger.Error(fmt.Sprintf("err: %v", r.Error))
		return 0, tx, r.Error
	}

	if r.RowsAffected <= 0 {
		gen_utils.Logger.Warn(fmt.Sprintf("failed to delete row by fields: %v", queryMap))
		gen_utils.Logger.Warn(fmt.Sprintf("number of rows deleted: %v", r.RowsAffected))
		return 0, tx, r.Error
	}

	return r.RowsAffected, tx, nil
}

func DeleteHardById[T any](id interface{}) (int64, *gorm.DB, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nhard deleting a row of: %v by id: %v", reflect.TypeOf(*new(T)).Name(), id))

	one, err := GetOneById[T](id)

	if gen_utils.IsNullOrEmpty(gen_utils.SafeGet(pretty.JSON(one), "$.id")) ||
		gen_utils.IsLessThanOrEqualTo(gen_utils.ConvertStrToInt64(gen_utils.SafeGetToString(pretty.JSON(one), "$.id")), 0) {
		gen_utils.Logger.Error(fmt.Sprintf("record of id: %v not found", id))
		return 0, nil, errors.New("record not found")
	}

	tx := Instance.Begin()
	var t2 T
	if err != nil {
		return 0, tx, err
	}

	err = mapstructure.Decode(one, &t2)

	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to map structure: %v", err.Error()))
		return 0, tx, err
	}

	r := Instance.Unscoped().Delete(&one).Where("id=?", id)

	if r.Error != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to delete row by id: %v", id))
		gen_utils.Logger.Error(fmt.Sprintf("err: %v", r.Error))
		return 0, tx, r.Error
	}

	if r.RowsAffected <= 0 {
		gen_utils.Logger.Warn(fmt.Sprintf("failed to delete row by id: %v", id))
		gen_utils.Logger.Warn(fmt.Sprintf("number of rows deleted: %v", r.RowsAffected))
		if r.Error != nil {
			return 0, tx, r.Error
		}
		return 0, tx, errors.New("not successful")
	}

	return r.RowsAffected, tx, nil
}

func DeleteSoftById[T any](id string) (int64, *gorm.DB, error) {
	gen_utils.Logger.Info(fmt.Sprintf("\n\nsoft deleting a row of: %v by id: %v", reflect.TypeOf(*new(T)).Name(), id))
	one, err := GetOneById[T](id)
	if err != nil {
		gen_utils.Logger.Error(fmt.Sprintf("failed to get record by id:%v err:%v", id, err))
		return 0, nil, err
	}

	if gen_utils.IsNullOrEmpty(gen_utils.SafeGetFromInterface(one, "$.id")) {
		msg := fmt.Sprintf("no record found with id: %v", id)
		gen_utils.Logger.Warn(msg)
		return 0, nil, errors.New(msg)
	}

	tx := Instance.Begin()
	var t2 T

	err = mapstructure.Decode(one, &t2)

	if err != nil {
		log.Println(fmt.Sprintf("failed to map structure: %v", err))
		return 0, tx, err
	}

	r := Instance.Delete(&one).Where("id=?", id)

	if r.Error != nil {
		log.Println(fmt.Sprintf("failed to delete row by id: %v", id))
		log.Println(fmt.Sprintf("err: %v", r.Error))
		return 0, tx, r.Error
	}

	if r.RowsAffected <= 0 {
		gen_utils.Logger.Warn(fmt.Sprintf("failed to delete row by id: %v", id))
		gen_utils.Logger.Warn(fmt.Sprintf("number of rows deleted: %v", r.RowsAffected))
		return 0, tx, r.Error
	}

	return r.RowsAffected, tx, nil
}
