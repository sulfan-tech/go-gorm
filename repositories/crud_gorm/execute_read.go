package crudgorm

import (
	sqlConf "go-gorm/configs/mySql"
	"go-gorm/models"
	"log"
	"time"

	"gorm.io/gorm"
)

type ReadGorm struct {
	db gorm.DB
}

type ReadGormImpl interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUsersByStatus(status string) ([]models.User, error)
	GetUsersByCondition(condition models.User) ([]models.User, error)
	GetUsersByMapCondition(condition map[string]interface{}) ([]models.User, error)
	GetUsersBySpecificFields(name string, age int) ([]models.User, error)
	GetActiveUsers() ([]models.User, error)
	GetNonActiveUsers() ([]models.User, error)
	GetUsersByMultipleStatus(status1, status2 string) ([]models.User, error)
	GetUserNamesAndEmails() ([]map[string]interface{}, error)
	GetUsersOrderedByAge() ([]models.User, error)
	GetLimitedUsers(limit, offset int) ([]models.User, error)
	GetUsersGroupedByStatus() ([]map[string]interface{}, error)
	GetDistinctStatuses() ([]string, error)
	GetUsersWithProfiles() ([]models.User, error)
	GetUsersWithOrders() ([]models.User, error)
	GetUsersWithRecentOrders() ([]models.User, error)
	ScanUserEmails() ([]string, error)
}

func NewInstanceReadGorm() ReadGormImpl {
	db, err := sqlConf.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	return &ReadGorm{
		db: *db,
	}
}

// Retrieving a single object by condition
func (r *ReadGorm) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Retrieving objects with primary key
func (r *ReadGorm) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Retrieving all objects
func (r *ReadGorm) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// String Conditions
func (r *ReadGorm) GetUsersByStatus(status string) ([]models.User, error) {
	var users []models.User
	result := r.db.Where("status = ?", status).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Struct & Map Conditions
func (r *ReadGorm) GetUsersByCondition(condition models.User) ([]models.User, error) {
	var users []models.User
	result := r.db.Where(&condition).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *ReadGorm) GetUsersByMapCondition(condition map[string]interface{}) ([]models.User, error) {
	var users []models.User
	result := r.db.Where(condition).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Specify Struct search fields
// Specify Struct search fields
func (r *ReadGorm) GetUsersBySpecificFields(name string, age int) ([]models.User, error) {
	var users []models.User
	result := r.db.Where(&models.User{Name: name, Age: age}, "Name", "Age").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Inline Condition
func (r *ReadGorm) GetActiveUsers() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users, "status = ?", "active")
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Not Conditions
func (r *ReadGorm) GetNonActiveUsers() ([]models.User, error) {
	var users []models.User
	result := r.db.Not("status = ?", "active").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Or Conditions
func (r *ReadGorm) GetUsersByMultipleStatus(status1, status2 string) ([]models.User, error) {
	var users []models.User
	result := r.db.Where("status = ?", status1).Or("status = ?", status2).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Selecting Specific Fields
func (r *ReadGorm) GetUserNamesAndEmails() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	rows, err := r.db.Table("users").Select("name, email").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name, email string
		rows.Scan(&name, &email)
		result = append(result, map[string]interface{}{"name": name, "email": email})
	}
	return result, nil
}

// Order
func (r *ReadGorm) GetUsersOrderedByAge() ([]models.User, error) {
	var users []models.User
	result := r.db.Order("age asc").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Limit & Offset
func (r *ReadGorm) GetLimitedUsers(limit, offset int) ([]models.User, error) {
	var users []models.User
	result := r.db.Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Group By & Having
func (r *ReadGorm) GetUsersGroupedByStatus() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	rows, err := r.db.Table("users").Select("status, COUNT(*) as count").Group("status").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		result = append(result, map[string]interface{}{"status": status, "count": count})
	}
	return result, nil
}

// Distinct
func (r *ReadGorm) GetDistinctStatuses() ([]string, error) {
	var statuses []string
	result := r.db.Model(&models.User{}).Distinct("status").Pluck("status", &statuses)
	if result.Error != nil {
		return nil, result.Error
	}
	return statuses, nil
}

// Joins
func (r *ReadGorm) GetUsersWithProfiles() ([]models.User, error) {
	var users []models.User
	result := r.db.Preload("Profile").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Joins Preloading
func (r *ReadGorm) GetUsersWithOrders() ([]models.User, error) {
	var users []models.User
	result := r.db.Preload("Orders").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Joins a Derived Table
func (r *ReadGorm) GetUsersWithRecentOrders() ([]models.User, error) {
	subQuery := r.db.Model(&models.Order{}).Select("user_id").Where("created_at > ?", time.Now().AddDate(0, -1, 0))
	var users []models.User
	result := r.db.Joins("JOIN (?) AS recent_orders ON recent_orders.user_id = users.id", subQuery).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Scan
func (r *ReadGorm) ScanUserEmails() ([]string, error) {
	var emails []string
	result := r.db.Model(&models.User{}).Pluck("email", &emails)
	if result.Error != nil {
		return nil, result.Error
	}
	return emails, nil
}
