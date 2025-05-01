package models

import (
	"github.com/mohamidsaiid/uniclipboard/internal/ADT"
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
	UpdateSignal ADT.Sig
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
	user, exists := u.Get()
	if !exists {
		u.DB.Create(&Users{SecretKey: string(hashed)})
		return nil
	}
	u.DB.Model(&user).Update("secretkey", hashed)
	u.UpdateSignal <- struct{}{}
	return nil
}

func (u *UsersModel) Get() (Users, bool) {
	var user Users
	u.DB.First(&user, 1)
	if user.ID == 0 {
		return Users{}, false
	}
	return user, true
}
