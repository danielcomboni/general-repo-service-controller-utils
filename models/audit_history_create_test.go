package models

import (
	"ewage-api/configs"
	"ewage-api/envs"
	"ewage-api/utils"
	"fmt"
	"github.com/Jeffail/gabs"
	general_goutils "github.com/danielcomboni/general-go-utils"
	"github.com/danielcomboni/general-repo-service-controller-utils/repo"
	"github.com/ohler55/ojg/pretty"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

/*
- https://onsi.github.io/gomega/#codeghttpcode-testing-http-clients
- https://semaphoreci.com/community/tutorials/getting-started-with-bdd-in-go-using-ginkgo
*/

var Instance *gorm.DB
var err error

func Connect(configURL ...string) {

	var connStr string
	var config configs.DBConfig
	if len(configURL) == 0 {
		connStr, config = configs.GetConnectionStringPg("./config.json")
	} else {
		connStr, config = configs.GetConnectionStringPg(configURL[0])
	}

	if envs.GetRemoteEnv() == "PRODUCTION" {
		connStr, config = configs.GetConnectionStringPg("./config.production.json")
	}

	utils.Logger.Info("connecting to database: " + connStr)

	//c := "postgresql://postgres:password@localhost:5432/configs_db_dev?sslmode=disable"
	//d := "postgresql://postgresql:password@localhost:5432/configs_db_dev?sslmode=disable"
	if config.Host == "localhost" {
		Instance, err = gorm.Open(postgres.Open(connStr), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	} else {
		Instance, err = gorm.Open(postgres.Open(connStr))
	}

	if err != nil {
		utils.Logger.Error("failed to connect to database...")
		// utils.Logger.Fatal(err.Error())
	}
	utils.Logger.Info("Connected to Database...")
	utils.Logger.Info("db details: " + connStr)
}

type Story struct {
	Name string `json:"name"`
}

func Migrate() {
	utils.Logger.Info("Database Migration Started...")
	err := Instance.AutoMigrate(

		// forex
		//&forexmodels.Currency{},
		//&forexmodels.Rate{},
		//&ClientRightOnEntityAction{},
		//&UserAndRole{},
		&AuditHistoryCreate{},
		&Story{},
	)
	if err != nil {
		msg := "failed to run migrations: " + err.Error()
		utils.Logger.Error(msg)
	}

	utils.Logger.Info("Database Migration Completed...")

}

func SetupTestParams() {
	envs.GetConnectionString()
	utils.Initialize()
	general_goutils.ReInitializeLoggingWithFileSyncEnabled()
	Connect("../config.test.json")
	Migrate()
	repo.RepoInitializer(Instance)

}

func TestSaveNewAuditHistoryFromGenericResponse(t *testing.T) {
	SetupTestParams()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Audit History Create Suite")
}

//func insertSample() {
//	// save AuditHistoryCreate
//
//	var caseData pgtype.JSONB
//
//	b := Branch{
//		BaseModel: BaseModel{
//			Id:               1,
//			CreatedAt:        time.Now(),
//			UpdatedAt:        time.Now(),
//			ReferenceId:      "",
//			RequestReference: "",
//		},
//		Name:       "Kampala Branch",
//		CustomerId: 1,
//		Customer:   Customer{},
//	}
//
//	type User struct {
//		Username string `json:"username"`
//	}
//
//	user := User{
//		Username: "comboni93@gmail.com",
//	}
//
//	obj := gabs.New()
//	obj.Set(b, "record")
//	obj.Set(user, "user")
//	println()
//	println()
//	println(fmt.Sprintf("adrenaline: %v", pretty.JSON(obj.Data())))
//	println()
//	println()
//
//	_ = caseData.Set(obj.Data())
//
//	modelName := reflect.TypeOf(*new(ClientUser)).Name()
//	objToSave := AuditHistoryCreate{
//		RecordId:     1,
//		WorkflowId:   1,
//		WorkflowType: "CREATE",
//		TableName:    utils.ToPlural(utils.ToSnakeCase(modelName)),
//		ModelName:    modelName,
//		CaseData:     caseData,
//	}
//	println()
//	println()
//	println(fmt.Sprintf("effectively: %v", general_goutils.IsNullOrEmpty(obj.Data())))
//	println()
//	println()
//	_, _ = repo.CreateWithPropertyCheck[AuditHistoryCreate](&objToSave)
//
//}

var _ = Describe("Audit history create", func() {
	Context("testing...", func() {

		It("test...", func() {

			//insertSample()

			all, err := repo.GetAllWithNoParams[AuditHistoryCreate]()
			if err != nil {
				general_goutils.Logger.Error(fmt.Sprintf("failed to fetch AuditHistoryCreate: %v", err.Error()))
				return
			}
			general_goutils.Logger.Info(fmt.Sprintf("not impressed: %v", len(all)))

			//var tables []string
			//if err := repo.Instance.Table("information_schema.tables").Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
			//	//panic(err)
			//}
			//
			//repo.Instance.Unscoped()
			//lo.ForEach(tables, func(item string, index int) {
			//	general_goutils.Logger.Info(fmt.Sprintf("tens: %v %v", index, item))
			//})

		})
	})
})

func TestRandom(t *testing.T) {

	raw := []byte(`{"phoneNumber":"0782888000"}`)
	consume, err := gabs.Consume(raw)
	if err != nil {
		println("error occured: ")
	}

	println(pretty.JSON(consume.Data()))

}

func TestGabs(t *testing.T) {
	obj := gabs.New()
	obj.SetP("the token", "token")
	obj.SetP("the token2", "token2")

	println(obj.String())

}