package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	SecretKey string
}

type UsersModel struct {
	DB *gorm.DB
}

func InitateDatabase(dbName string) (*UsersModel, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Users{})
	return &UsersModel{DB: db}, nil
}

func (u *UsersModel) Update(secretKey string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(secretKey), 10)
	if err != nil {
		return err
	}
	var user = u.Get()
	u.DB.Model(&user).Update("secretkey", hashed)
	return nil
}

func (u *UsersModel) Get() (Users) {
	var user Users
	u.DB.First(&user, 1)
	return user
}
