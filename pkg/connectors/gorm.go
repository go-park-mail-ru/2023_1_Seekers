package connectors

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewGormDb(connStr, tablePrefix string) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{DSN: connStr}), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix:   tablePrefix,
		SingularTable: false,
	}})
}
