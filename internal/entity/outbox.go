package entity

import "time"

func (Outbox) TableName() string {
	return "outbox"
}

type Outbox struct {
	EventID    string     `gorm:"primaryKey;type:varchar(36);not null"`
	Topic      string     `gorm:"type:varchar(100);not null"`
	EventType  string     `gorm:"type:varchar(100);not null"`
	Payload    string     `gorm:"type:text;not null"`
	Status     string     `gorm:"type:varchar(20);index;not null;default:'PENDING'"`
	RetryCount int        `gorm:"not null;default:0"`
	PublishAt  *time.Time `gorm:"type:timestamptz"`
	AuditModel `gorm:"embedded"`
}
