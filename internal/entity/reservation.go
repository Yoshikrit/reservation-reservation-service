package entity

import "time"

func (r Reservation) TableName() string {
	return "reservations"
}

type Reservation struct {
	ReservationID string    `gorm:"primaryKey;type:varchar(14);not null"`
	ProductID     string    `gorm:"type:varchar(14);index;not null"`
	Quantity      uint      `gorm:"type:int;not null;default:0;check:quantity >= 0"`
	Status        string    `gorm:"type:varchar(20);index;not null;default:'HELD'"`
	BasePrice     float64   `gorm:"type:decimal(12,2);not null;check:base_price >= 0"`
	DiscountRate  float64   `gorm:"type:decimal(5,4);not null"`
	Price         float64   `gorm:"type:decimal(12,2);not null;check:price >= 0"`
	ExpiresAt     time.Time `gorm:"index;not null"`
	AuditModel    `gorm:"embedded"`
}
