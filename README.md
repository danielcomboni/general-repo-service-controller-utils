```console
go get github.com/danielcomboni/general-repo-service-controller-utils@v0.3.0
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

### main 

```go
func main() {

	router := gin.Default()

	router.Use(CORS())

	Connect()
	repo.RepoInitializer(Instance)
	repo.RunMigrations(&User{}, &Student{})

        // NOTE: if authentication is required, then start with set the authentication for the default endpints
	new(routes_utils.SingleEntityGroupedRouteDefinition[Student]).
        // authenticates GET method that fetches all records from database
		SetAuthDefaultGetAll(func(ctx *gin.Context) (*gin.Context, bool, string) {
			limit, _ := strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
			println(fmt.Sprintf("that made all the difference: %v", limit))
			return ctx, true, ""
		}).
        // authenticates GET method that fetches a single record by ID
		SetAuthDefaultGetById(func(ctx *gin.Context) (*gin.Context, bool, string) {
			limit, _ := strconv.Atoi(ctx.Request.URL.Query().Get("limit"))
			println(fmt.Sprintf("that made all the difference: %v", limit))
			return ctx, false, ""
		}).
        // generates CRUD operations endpoints
		SetSingleEntityGroupedRouteDefinitionPaginationQueryParams(router, "", "/api/students/v1", nil,
			[]models.QueryStructure{}, []string{}).
        // generates an endpoint to fetch by provided path params
		AddGetOneUsingPathParams("", "/api/students/v1/getByEmail/:email",
			[]models.QueryStructure{
				{
					ParamName:     "email", // this is the name in the struct
					DbTableColumn: "email", // this is the name of the db table column
				},
			}, func(ctx *gin.Context) (*gin.Context, bool, string) {
				return ctx, true, ""
			}).
        // fetches all/multiple records based on the provided path params    
        AddGetAllUsingPathParams("", "/api/students/v1/allByEmail/:email", []models.QueryStructure{
			{
				ParamName:     "email",
				DbTableColumn: "email",
			},
		}, func(ctx *gin.Context) (*gin.Context, bool, string) {
			return ctx, true, ""
		}).
        // inserts/creates a record but first checks for duplicate(s) using the specified properties/fields
        AddPostWithDuplicateCheckUsingProperties("", "/api/students/v1/createWithDuplicateCheck", []string{"email"}, []string{}, func(ctx *gin.Context) (*gin.Context, bool, string) {
			return ctx, false, ""
		})

        // finally run the application
    	err := router.Run("localhost:6000")
	    if err != nil {
		    general_goutils.Logger.Error(fmt.Sprintf("failed to to run application: %v", err))
	    }
}
```