package models

// REF: https://stackoverflow.com/questions/60092840/string-to-jsonb-with-gorm-and-postgres
import (
	"ewage-api/utils"
	"fmt"
	"github.com/Jeffail/gabs"
	general_goutils "github.com/danielcomboni/general-go-utils"
	"github.com/danielcomboni/general-repo-service-controller-utils/repo"
	"github.com/danielcomboni/general-repo-service-controller-utils/responses"
	"github.com/gin-gonic/gin"
	"github.com/ohler55/ojg/pretty"
	"reflect"

	//gormjsonb "github.com/dariubs/gorm-jsonb"
	"github.com/jackc/pgtype"
)

type AuditHistoryUpdate struct {
	BaseModelAuditHistory
	RecordId     int64        `json:"recordId" validate:"required"`
	WorkflowId   int64        `json:"workflowId" validate:"required"`
	WorkflowType string       `json:"workflowType" validate:"required"`
	TableName    string       `json:"tableName"`
	ModelName    string       `json:"modelName"`
	CustomerId   int64        `json:"customerId"`
	BranchId     int64        `json:"branchId"`
	CaseData     pgtype.JSONB `json:"caseData" gorm:"type:jsonb;default:'{}';not null"`
}

func SaveUpdateAuditHistoryFromGenericResponse[T any](res responses.GenericResponse, c *gin.Context, shouldCommit func(flag bool)) {

	if res.Status == responses.Created || res.Status == responses.OK {

		model := general_goutils.SafeGetFromInterface(res.Data, "$.result")

		var caseData pgtype.JSONB
		var customerId int64
		var branchId int64

		customerIdStr := c.Query("customerId")
		branchIdStr := c.Query("branchId")
		username := c.Query("username")

		user := User{
			Username: "",
		}

		if !general_goutils.IsNullOrEmpty(username) {
			user.Username = username
		}

		obj := gabs.New()
		obj.Set(model, "record")
		obj.Set(user, "user")

		if !general_goutils.IsNullOrEmpty(customerIdStr) {
			customerId = general_goutils.ConvertStrToInt64(customerIdStr)
		}

		if !general_goutils.IsNullOrEmpty(branchIdStr) {
			branchId = general_goutils.ConvertStrToInt64(branchIdStr)
		}

		err := caseData.Set(obj.Data())
		if err != nil {
			if general_goutils.ConvertStrToInt64(pretty.JSON(general_goutils.SafeGet(pretty.JSON(model), "$.id"))) == 0 {
				general_goutils.Logger.Error(fmt.Sprintf("failed to set caseData: %v", err.Error()))
			}
		}

		modelName := reflect.TypeOf(*new(T)).Name()
		objToSave := AuditHistoryUpdate{
			RecordId:     general_goutils.ConvertStrToInt64(pretty.JSON(general_goutils.SafeGet(pretty.JSON(model), "$.id"))),
			WorkflowId:   general_goutils.ConvertStrToInt64(pretty.JSON(general_goutils.SafeGet(pretty.JSON(model), "$.id"))),
			WorkflowType: "UPDATE",
			TableName:    utils.ToPlural(utils.ToSnakeCase(modelName)),
			ModelName:    modelName,
			CustomerId:   customerId,
			BranchId:     branchId,
			CaseData:     caseData,
		}
		general_goutils.Logger.Info(fmt.Sprintf("saving an update audit history for entity: %v and id: %v for branchId: %v and customerId: %v and user: %v", modelName, objToSave.RecordId, branchId, customerId, username))
		_, tx, err := repo.CreateWithPropertyCheck[AuditHistoryUpdate](&objToSave)
		if err != nil {
			tx.Rollback()
			return
		}
		//_, err = repo.CreateWithPropertyCheck[AuditHistoryCreate](&objToSave)

		if err != nil {
			if shouldCommit != nil {
				shouldCommit(false)
			}
			return
		}
		if shouldCommit != nil {
			shouldCommit(true)
		}

	}
}
