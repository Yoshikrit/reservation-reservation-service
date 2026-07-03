package entity

import "time"

func (c CfgPromotionDate) TableName() string {
	return "cfg_promotion_date"
}

type CfgPromotionDate struct {
	ID int64 `gorm:"primaryKey;type:bigint;not null;column:reservation_id"`
	Date          uint32    `gorm:"type:int;not null;check:date > 0"`
	EndDate       time.Time `gorm:"type:timestamptz;not null"`
	IsEnabled     bool      `gorm:"not null;default:false"`
	AuditModel    `gorm:"embedded"`
}
