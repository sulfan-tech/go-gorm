package crudgorm

import (
	"context"
	sqlConf "go-gorm/configs/mySql"
	"go-gorm/models"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CreateGorm struct {
	db gorm.DB
}

type CreateGormImpl interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	CreateUserWithSelectedFields(ctx context.Context, user *models.User, fields ...string) (*models.User, error)
	BatchCreateUsers(ctx context.Context, users []models.User) ([]models.User, error)
	CreateUserWithHooks(ctx context.Context, user *models.User) (*models.User, error)
	CreateUserFromMap(ctx context.Context, userMap map[string]interface{}) (*models.User, error)
	CreateUserFromSQLExpression(ctx context.Context, name string, age int) (*models.User, error)
	CreateUserWithAssociations(ctx context.Context, user *models.User) (*models.User, error)
	UpsertUser(ctx context.Context, user *models.User) (*models.User, error)
}

func NewInstanceCRUDGorm() CreateGormImpl {
	db, err := sqlConf.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	return &CreateGorm{
		db: *db,
	}
}

// Create Record
func (c *CreateGorm) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	res := c.db.Create(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

// Create Record With Selected Fields
func (c *CreateGorm) CreateUserWithSelectedFields(ctx context.Context, user *models.User, fields ...string) (*models.User, error) {
	res := c.db.Select(fields).Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

// Batch Insert
func (c *CreateGorm) BatchCreateUsers(ctx context.Context, users []models.User) ([]models.User, error) {
	res := c.db.Create(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

// Create Hooks
//
//	func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
//	    // Logika sebelum pembuatan catatan
//	    user.UUID = uuid.New().String() // Contoh: mengatur UUID sebelum pembuatan catatan
//	    return
//	}
func (c *CreateGorm) CreateUserWithHooks(ctx context.Context, user *models.User) (*models.User, error) {
	res := c.db.Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

// Create From Map
func (c *CreateGorm) CreateUserFromMap(ctx context.Context, userMap map[string]interface{}) (*models.User, error) {
	res := c.db.Model(&models.User{}).Create(userMap)
	if res.Error != nil {
		return nil, res.Error
	}
	user := &models.User{}
	c.db.Where("id = ?", userMap["id"]).First(user)
	return user, nil
}

// Create From SQL Expression/Context Valuer
func (c *CreateGorm) CreateUserFromSQLExpression(ctx context.Context, name string, age int) (*models.User, error) {
	res := c.db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", name, age)
	if res.Error != nil {
		return nil, res.Error
	}
	user := &models.User{}
	c.db.Where("name = ? AND age = ?", name, age).First(user)
	return user, nil
}

// Create With Associations
func (c *CreateGorm) CreateUserWithAssociations(ctx context.Context, user *models.User) (*models.User, error) {
	res := c.db.Omit("CreditCards").Create(&user) // Jangan buat asosiasi CreditCards dulu
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

// Upsert / On Conflict
func (c *CreateGorm) UpsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	res := c.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "age"}), // Kolom yang diupdate jika konflik
	}).Create(&user)

	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
