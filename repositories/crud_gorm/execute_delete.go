package crudgorm

import (
	sqlConf "go-gorm/configs/mySql"
	"log"

	"gorm.io/gorm"
)

type DeleteGorm struct {
	db gorm.DB
}

type DeleteGormImpl interface {
}

func NewInstanceDeleteGorm() DeleteGormImpl {
	db, err := sqlConf.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	return &DeleteGorm{
		db: *db,
	}
}

