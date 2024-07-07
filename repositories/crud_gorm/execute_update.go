package crudgorm

import (
	sqlConf "go-gorm/configs/mySql"
	"go-gorm/models"
	"log"

	"gorm.io/gorm"
)

type UpdateGorm struct {
	db gorm.DB
}

type UpdateGormImpl interface {
	UpdateUser(user *models.User) error
	SaveAllFields(user *models.User) error
	UpdateSingleColumn(userID uint, column string, value interface{}) error
	UpdatesMultipleColumns(userID uint, updates map[string]interface{}) error
	UpdateSelectedFields(userID uint, selectedFields map[string]interface{}) error
	UpdateHooks(userID uint, updates map[string]interface{}) error
	BatchUpdates(users []models.User) error
	BlockGlobalUpdates() error
	UpdatedRecordsCount(userID uint, updates map[string]interface{}) (int64, error)
	UpdateWithSQLExpression(userID uint, sqlExpr string) error
	UpdateFromSubQuery(userID uint, subQuery *gorm.DB) error
	UpdateWithoutHooks(userID uint, updates map[string]interface{}) error
	ReturnModifiedData(userID uint, updates map[string]interface{}) (*models.User, error)
	CheckFieldHasChanged(userID uint, fieldName string, newValue interface{}) (bool, error)
	ChangeUpdatingValues(userID uint, updates map[string]interface{}) error
}

func NewInstanceUpdateGorm() UpdateGormImpl {
	db, err := sqlConf.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	return &UpdateGorm{
		db: *db,
	}
}

// UpdateUser: Update the entire user record
func (u *UpdateGorm) UpdateUser(user *models.User) error {
	result := u.db.Save(user)
	return result.Error
}

// SaveAllFields: Save all fields of a user, including zero values
func (u *UpdateGorm) SaveAllFields(user *models.User) error {
	result := u.db.Model(user).Select("*").Updates(user)
	return result.Error
}

// UpdateSingleColumn: Update a single column for a specific user
func (u *UpdateGorm) UpdateSingleColumn(userID uint, column string, value interface{}) error {
	result := u.db.Model(&models.User{}).Where("id = ?", userID).Update(column, value)
	return result.Error
}

// UpdatesMultipleColumns: Update multiple columns for a specific user
func (u *UpdateGorm) UpdatesMultipleColumns(userID uint, updates map[string]interface{}) error {
	result := u.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	return result.Error
}

// UpdateSelectedFields: Update selected fields of a specific user
func (u *UpdateGorm) UpdateSelectedFields(userID uint, selectedFields map[string]interface{}) error {
	result := u.db.Model(&models.User{}).Where("id = ?", userID).Select("name", "email").Updates(selectedFields)
	return result.Error
}

// UpdateHooks: Use GORM hooks when updating a user
func (u *UpdateGorm) UpdateHooks(userID uint, updates map[string]interface{}) error {
	result := u.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	return result.Error
}

// BatchUpdates: Batch update multiple user records
func (u *UpdateGorm) BatchUpdates(users []models.User) error {
	result := u.db.Model(&models.User{}).Updates(users)
	return result.Error
}

// BlockGlobalUpdates: Block global updates for safety
func (u *UpdateGorm) BlockGlobalUpdates() error {
	result := u.db.Session(&gorm.Session{AllowGlobalUpdate: false}).Model(&models.User{}).Updates(map[string]interface{}{"name": "new name"})
	return result.Error
}

// UpdatedRecordsCount: Get the number of records updated
func (u *UpdateGorm) UpdatedRecordsCount(userID uint, updates map[string]interface{}) (int64, error) {
	result := u.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	return result.RowsAffected, result.Error
}

// UpdateWithSQLExpression: Update a user using SQL expression
func (u *UpdateGorm) UpdateWithSQLExpression(userID uint, sqlExpr string) error {
	result := u.db.Model(&models.User{}).Where("id = ?", userID).Update("age", gorm.Expr(sqlExpr))
	return result.Error
}

// UpdateFromSubQuery: Update a user from a subquery
func (u *UpdateGorm) UpdateFromSubQuery(userID uint, subQuery *gorm.DB) error {
	result := u.db.Model(&models.User{}).Where("id = ?", userID).Update("age", subQuery)
	return result.Error
}

// UpdateWithoutHooks: Update without triggering hooks or updating time fields
func (u *UpdateGorm) UpdateWithoutHooks(userID uint, updates map[string]interface{}) error {
	result := u.db.Model(&models.User{}).Where("id = ?", userID).Omit("UpdatedAt").Updates(updates)
	return result.Error
}

// ReturnModifiedData: Return the modified user data after update
func (u *UpdateGorm) ReturnModifiedData(userID uint, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	result := u.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// CheckFieldHasChanged: Check if a specific field has changed
func (u *UpdateGorm) CheckFieldHasChanged(userID uint, fieldName string, newValue interface{}) (bool, error) {
	var user models.User
	result := u.db.First(&user, userID)
	if result.Error != nil {
		return false, result.Error
	}

	switch fieldName {
	case "name":
		return user.Name != newValue.(string), nil
	case "email":
		return user.Email != newValue.(string), nil
	// Add more cases as needed
	default:
		return false, nil
	}
}

// ChangeUpdatingValues: Change the values being updated dynamically
func (u *UpdateGorm) ChangeUpdatingValues(userID uint, updates map[string]interface{}) error {
	result := u.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	return result.Error
}
