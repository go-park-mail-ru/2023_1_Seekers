package connectors

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewGormDb(connStr, tablePrefix string, maxOpenCons int) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: connStr}), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix:   tablePrefix,
		SingularTable: false,
	},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	pgDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	pgDb.SetMaxOpenConns(maxOpenCons)
	return db, nil
}
