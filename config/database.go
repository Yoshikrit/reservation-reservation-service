package config

import (
	"time"

	"reservation/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	DatabaseUrl string `env:"DATABASE_URL,required"`
}

func InitDatabase(databaseUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MigrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Reservation{},
		&entity.Outbox{},
		&entity.CfgPromotionDate{},
	)
}

func SeedDatabase(db *gorm.DB) {
	endDate := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	promotionDates := []entity.CfgPromotionDate{
		{ID: 1, Date: 11, EndDate: endDate, IsEnabled: true},
		{ID: 2, Date: 22, EndDate: endDate, IsEnabled: true},
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&promotionDates)
}
