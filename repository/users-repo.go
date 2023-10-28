package repository

import (
	"errors"
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

func convertRegisterRes(users *model.Users) {
	panic("unimplemented")
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

func (um *UsersModel) MyProfile(id uint) (*model.Users, int) {
	var data = model.Users{}

	if err := um.db.First(&data, id).Error; err != nil {
		logrus.Error("Model : Myprofil data error,", err.Error())
		return nil, 1
	}

	return &data,0
}

func (um *UsersModel) UpdateMyProfile(updateUser *model.Users) (*model.Users, int) {
	var user = model.Users{}

	if err := um.db.First(&user, updateUser.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Error("Repo: User not found for UpdateMyProfile, ", err.Error())
			return nil, 1
		}
		logrus.Error("Repo: Select method UpdateMyProfile data error, ", err.Error())
		return nil, 2
	}

	if user.Email != updateUser.Email {
		var existingUser = model.Users{}
		if err := um.db.Where("email = ?", updateUser.Email).First(&existingUser).Error; err == nil {
			logrus.Error("Repo: UpdateMyProfile, Email already registered")
			return nil, 3
		}
	}

	user.Name = updateUser.Name
	user.Email = updateUser.Email
	user.Password = updateUser.Password

	var qry = um.db.Save(&user)
	if err := qry.Error; err != nil {
		logrus.Error("Repo: Save method UpdateMyProfile data error, ", err.Error())
		return nil, 2
	}

	if dataCount := qry.RowsAffected; dataCount < 1 {
		logrus.Error("Repo: UpdateMyProfile data error, no data affected")
		return nil, 2
	}

	return &user, 0
}


