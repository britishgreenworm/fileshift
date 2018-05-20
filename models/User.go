package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//Recruit endpoint
type User struct {
	gorm.Model
	Name     string
	Role     string
	Password string
	Files    []File
}

func CreateUser(db *gorm.DB, user User) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	if err != nil {
		return User{}, fmt.Errorf("Model: Generating Password Hash: %v", err)
	}

	result := db.Create(&user)

	if result.Error != nil {
		return User{}, fmt.Errorf("Model: Creating User: %v", result.Error)
	}

	return user, nil
}

func GetUser(db *gorm.DB, id int64) (User, error) {
	user := User{}
	result := db.First(&user, id).Select("Name, Role, ID")
	if result.Error != nil {
		return User{}, fmt.Errorf("Model: Getting User: %v", result.Error)
	}

	return user, nil
}

func GetUsers(db *gorm.DB, id int64) ([]User, error) {
	users := []User{}
	result := db.Find(&users).Select("Name, Role, ID")

	if result.Error != nil {
		return []User{}, fmt.Errorf("model: Getting Users: %v", result.Error)
	}

	return users, nil
}

func UpdateUser(db *gorm.DB, user User) (int64, error) {
	result := db.Save(&user)
	if result.Error != nil {
		return result.RowsAffected, fmt.Errorf("Model: Updating User: %v", result.Error)
	}

	return result.RowsAffected, nil
}

func DeleteUser(db *gorm.DB, id int64) (int64, error) {
	user := User{}
	result := db.Unscoped().Delete(&user, id)
	if result.Error != nil {
		return result.RowsAffected, fmt.Errorf("model: Deleting User: %v", result.Error)
	}

	return result.RowsAffected, nil
}
