package repository

import (
	"quiz/helper"
	"quiz/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UsersModel struct {
	db *gorm.DB
}

// Init digunakan untuk menginisialisasi objek UsersModel dengan koneksi database.
func (um *UsersModel) Init(db *gorm.DB) {
	um.db = db
}

// Register digunakan untuk mendaftarkan pengguna baru.
func (um *UsersModel) Register(newUser model.Users) (*model.Users, int) {
	existingUser := &model.Users{}

	// Mengecek apakah alamat email sudah terdaftar.
	if err := um.db.Where("email = ?", newUser.Email).First(existingUser).Error; err == nil {
		logrus.Error("Model: Email already registered")
		return nil, 1
	}

	// Menyimpan pengguna baru ke dalam database.
	if err := um.db.Create(&newUser).Error; err != nil {
		logrus.Error("Model: Insert data error, ", err.Error())
		return nil, 2
	}

	return &newUser, 0
}

func (um *UsersModel) Login(email string, password string) (*model.Users, int) {
	var data = model.Users{}
	if err := um.db.Where("email = ?", email).First(&data).Error; err != nil {
		logrus.Error("Model : Login data error,", err.Error())
		if data.ID == 0 {
			logrus.Error("Model : not found")
			return nil, 2
		}
		return nil, 1
	}

	if err := helper.ComparePassword(data.Password, password); err != nil {
		logrus.Error("Model : Login data error,", err.Error())
		return nil, 2
	}

	return &data,0
}

