```console
go get github.com/danielcomboni/general-repo-service-controller-utils@v0.3.8
```

-------------------------------------- example --------------------------------------

### connect to database

```go
func Connect() {

	dsn := "host=localhost user=postgres password=password dbname=generic_repo_service_controller_dev_db port=5432 sslmode=disable TimeZone=Africa/Nairobi"

	Instance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		general_goutils.Logger.Error("failed to connect to database...")
		general_goutils.Logger.Fatal(err.Error())
	}
	general_goutils.Logger.Info("Connected to Database...")
	general_goutils.Logger.Info("db details: " + dsn)
}
```

### models
```go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Student struct {
	Id int64
	User
	School string `json:"school"`
	Email  string `json:"email"`
}
```

### create endpoints
```go
	new(routes_utils.SingleEntityGroupedRouteDefinition[Student]).
		SetAuthDefaultGetAll(func(ctx *gin.Context) (*gin.Context, bool, string) {
			limit, _ := strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
			println(fmt.Sprintf("that made all the difference: %v", limit))
			return ctx, true, ""
		}).
		SetAuthDefaultGetById(func(ctx *gin.Context) (*gin.Context, bool, string) {
			limit, _ := strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
			println(fmt.Sprintf("that made all the difference: %v", limit))
			return ctx, false, ""
		}).
		SetSingleEntityGroupedRouteDefinitionPaginationQueryParams(router, "", "/api/students/v1", nil,
			[]models.QueryStructure{}, []string{}).
		AddGetOneUsingPathParams("", "/api/students/v1/getByEmail/:email",
			[]models.QueryStructure{
				{
					ParamName:     "email",
					DbTableColumn: "email",
				},
			}, func(ctx *gin.Context) (*gin.Context, bool, string) {
				return ctx, true, ""
			}).
		AddGetAllUsingPathParams("", "/api/students/v1/allByEmail/:email", []models.QueryStructure{
			{
				ParamName:     "email",
				DbTableColumn: "email",
			},
		}, func(ctx *gin.Context) (*gin.Context, bool, string) {
			return ctx, true, ""
		}).
		AddPostWithDuplicateCheckUsingProperties("", "/api/students/v1/createWithDuplicateCheck", []string{"email"}, []string{}, func(ctx *gin.Context) (*gin.Context, bool, string) {
			return ctx, false, ""
		})
```

### main (conatins complete sample)

```go
package main

import (
	"fmt"
	"strconv"

	general_goutils "github.com/danielcomboni/general-go-utils"
	"github.com/danielcomboni/general-repo-service-controller-utils/controller"
	"github.com/danielcomboni/general-repo-service-controller-utils/models"

	"github.com/danielcomboni/general-repo-service-controller-utils/repo"
	routes_utils "github.com/danielcomboni/general-repo-service-controller-utils/route"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // <------------ here

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Instance *gorm.DB
var err error

func Connect() {

	dsn := "host=localhost user=postgres password=password dbname=generic_repo_service_controller_dev_db port=5432 sslmode=disable TimeZone=Africa/Nairobi"

	Instance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		general_goutils.Logger.Error("failed to connect to database...")
		general_goutils.Logger.Fatal(err.Error())
	}
	general_goutils.Logger.Info("Connected to Database...")
	general_goutils.Logger.Info("db details: " + dsn)

	// run migrations
	general_goutils.Logger.Info("Database Migration Started...")
	err := Instance.AutoMigrate(&User{}, &Student{})
	if err != nil {
		msg := "failed to run migrations: " + err.Error()
		general_goutils.Logger.Error(msg)
	}
	general_goutils.Logger.Info("Database Migration Completed...")

}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Student struct {
	Id int64
	User
	School string `json:"school"`
	Email  string `json:"email"`
}

func main() {

	router := gin.Default()

	router.Use(CORS())

	Connect()
	repo.RepoInitializer(Instance)

	new(routes_utils.SingleEntityGroupedRouteDefinition[Student]).
		SetAuthDefaultGetAll(func(ctx *gin.Context) (*gin.Context, bool, string) {
			limit, _ := strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
			println(fmt.Sprintf("that made all the difference: %v", limit))
			return ctx, true, ""
		}).
		SetAuthDefaultGetById(func(ctx *gin.Context) (*gin.Context, bool, string) {
			limit, _ := strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
			println(fmt.Sprintf("that made all the difference: %v", limit))
			return ctx, false, ""
		}).
		SetSingleEntityGroupedRouteDefinitionPaginationQueryParams(router, "", "/api/students/v1", nil,
			[]models.QueryStructure{}, []string{}).
		AddGetOneUsingPathParams("", "/api/students/v1/getByEmail/:email",
			[]models.QueryStructure{
				{
					ParamName:     "email",
					DbTableColumn: "email",
				},
			}, func(ctx *gin.Context) (*gin.Context, bool, string) {
				return ctx, true, ""
			}).
		AddGetAllUsingPathParams("", "/api/students/v1/allByEmail/:email", []models.QueryStructure{
			{
				ParamName:     "email",
				DbTableColumn: "email",
			},
		}, func(ctx *gin.Context) (*gin.Context, bool, string) {
			return ctx, true, ""
		}).
		AddPostWithDuplicateCheckUsingProperties("", "/api/students/v1/createWithDuplicateCheck", []string{"email"}, []string{}, func(ctx *gin.Context) (*gin.Context, bool, string) {
			return ctx, false, ""
		})

	err := router.Run("localhost:6000")
	if err != nil {
		general_goutils.Logger.Error(fmt.Sprintf("failed to to run application: %v", err))
	}

}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-authorization-key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

```