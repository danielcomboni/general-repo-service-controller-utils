package models

import (
	"fmt"
	general_goutils "github.com/danielcomboni/general-go-utils"
	"strings"
	"time"
)

type LogTag struct {
	On          time.Time `json:"on" gorm:"-:migration"`
	CompanyName string    `json:"companyName" gorm:"-:migration"`
	CompanyId   int64     `json:"companyId" gorm:"-:migration"`
	Username    string    `json:"username" gorm:"-:migration"`
	UserId      int64     `json:"userId" gorm:"-:migration"`
	Action      string    `json:"action" gorm:"-:migration"`
	Entity      string    `json:"entity"`
	ReferenceId string    `json:"referenceId"`
}

func AddLogTag(l LogTag) string {
	var sb strings.Builder
	// add on date
	t := time.Now()
	sb.WriteString(fmt.Sprintf(" >> On: %v", t))

	if !general_goutils.IsNullOrEmpty(l.ReferenceId) {
		sb.WriteString(fmt.Sprintf(" >> ReferenceId: %v", l.ReferenceId))
	}

	if !general_goutils.IsNullOrEmpty(l.CompanyName) {
		sb.WriteString(fmt.Sprintf(" >> Customer: %v", l.CompanyName))
	}

	if l.CompanyId > 0 {
		sb.WriteString(fmt.Sprintf(" >> Customer ID: %v", l.CompanyId))
	}

	if !general_goutils.IsNullOrEmpty(l.Username) {
		sb.WriteString(fmt.Sprintf(" >> Username: %v", l.Username))
	}

	if l.UserId > 0 {
		sb.WriteString(fmt.Sprintf(" >> User ID: %v", l.UserId))
	}
	if !general_goutils.IsNullOrEmpty(l.Action) {
		sb.WriteString(fmt.Sprintf(" >> Action: %v", l.Action))
	}

	if !general_goutils.IsNullOrEmpty(l.Entity) {
		sb.WriteString(fmt.Sprintf(" >> Entity: %v", l.Entity))
	}

	return sb.String()
}
