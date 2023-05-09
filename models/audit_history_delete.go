package models

// REF: https://stackoverflow.com/questions/60092840/string-to-jsonb-with-gorm-and-postgres
import (
	"ewage-api/utils"
	"fmt"
	"github.com/Jeffail/gabs"
	general_goutils "github.com/danielcomboni/general-go-utils"
	"github.com/danielcomboni/general-repo-service-controller-utils/repo"
	"github.com/gin-gonic/gin"
	"reflect"

	//gormjsonb "github.com/dariubs/gorm-jsonb"
	"github.com/jackc/pgtype"
)

type AuditHistoryDelete struct {
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

func SaveDeleteAuditHistoryFromGenericResponse[T any](rowsAffected int64, c *gin.Context) {

	if rowsAffected > 0 {

		id := c.Param("id")
		var caseData pgtype.JSONB
		err := caseData.Set(map[string]int64{
			"id": general_goutils.ConvertStrToInt64(id),
		})
		if err != nil {
			general_goutils.Logger.Error(fmt.Sprintf("failed to set caseData: %v", err.Error()))
		}

		var customerId int64
		var branchId int64

		customerIdStr := c.Query("customerId")
		branchIdStr := c.Query("branchId")
		username := c.Query("username")

		if !general_goutils.IsNullOrEmpty(customerIdStr) {
			customerId = general_goutils.ConvertStrToInt64(customerIdStr)
		}

		if !general_goutils.IsNullOrEmpty(branchIdStr) {
			branchId = general_goutils.ConvertStrToInt64(branchIdStr)
		}

		user := User{
			Username: "",
		}

		if !general_goutils.IsNullOrEmpty(username) {
			user.Username = username
		}

		obj := gabs.New()
		obj.Set(rowsAffected, "rowsAffected")
		obj.Set(user, "user")

		modelName := reflect.TypeOf(*new(T)).Name()
		objToSave := AuditHistoryDelete{
			BaseModelAuditHistory: BaseModelAuditHistory{},
			RecordId:              general_goutils.ConvertStrToInt64(id),
			WorkflowId:            general_goutils.ConvertStrToInt64(id),
			WorkflowType:          "DELETE",
			TableName:             utils.ToPlural(utils.ToSnakeCase(modelName)),
			ModelName:             modelName,
			CustomerId:            customerId,
			BranchId:              branchId,
			CaseData:              caseData,
		}
		general_goutils.Logger.Info(fmt.Sprintf("saving a delete audit history for entity: %v and id: %v for branchId: %v and customerId: %v and user: %v", modelName, objToSave.RecordId, branchId, customerId, username))
		_, _, err = repo.CreateWithPropertyCheck[AuditHistoryDelete](&objToSave)
		if err != nil {
			return
		}

	}
}
