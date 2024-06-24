package crudgorm

import (
	"context"
	sqlConf "go-gorm/configs/mySql"
	"go-gorm/models"
	"log"

	"gorm.io/gorm"
)

type CRUDGorm struct {
	db gorm.DB
}

type CRUDGormImpl interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
}

func NewInstanceCRUDGorm() CRUDGormImpl {
	db, err := sqlConf.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	return &CRUDGorm{
		db: *db,
	}
}

func (c *CRUDGorm) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	res := c.db.Create(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}
