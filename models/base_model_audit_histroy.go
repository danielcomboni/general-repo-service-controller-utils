package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModelAuditHistory struct {
	Id        int64          `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (base *BaseModelAuditHistory) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("CreatedAt", time.Now().UTC())
	tx.Statement.SetColumn("UpdatedAt", time.Now().UTC())
	return nil
}

func (base *BaseModelAuditHistory) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("CreatedAt", base.CreatedAt)
	tx.Statement.SetColumn("UpdatedAt", time.Now().UTC())
	return nil
}
