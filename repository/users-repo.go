package repository

import (
	"quiz/helper"
	"quiz/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UsersInterface interface {
	Register(newUser model.Users) (*model.Users, int)
	Login(email string, password string) (*model.Users, int)
	MyProfile(id uint) (*model.Users, int)
	UpdateMyProfile(updateUser model.Users) (*model.Users, int)
}

type UsersModel struct {
	db *gorm.DB
}

func (um *UsersModel) Init(db *gorm.DB) {
	um.db = db
}

func NewUsersModel(db *gorm.DB) UsersInterface {
	return &UsersModel{
		db: db,
	}
}


func (um *UsersModel) Register(newUser model.Users) (*model.Users, int) {
	var count int64
	if err := um.db.Model(&model.Users{}).Where("email = ?", newUser.Email).Count(&count).Error; err != nil {
		logrus.Error("Repository: Unable to check existing email:", err.Error())
		return nil, 2 
	}

	if count > 0 {
		logrus.Error("Repository: Email already registered")
		return nil, 1 
	}
	hashedPassword := helper.HashPassword(newUser.Password)
	newUser.Password = hashedPassword

	if err := um.db.Create(&newUser).Error; err != nil {
		logrus.Error("Repository: Failed to insert data:", err.Error())
		return nil, 2 
	}

	return &newUser, 0 
}

func (um *UsersModel) Login(email string, password string) (*model.Users, int) {
	var user = model.Users{}

	if err := um.db.Where("email = ?", email).First(&user).Error; err != nil {
		if user.ID == 0 {
			logrus.Error("Repository: not found user")
			return nil, 2
		}

		logrus.Error("Repository: Login data error,", err.Error())
		return nil, 1
	}

	if password == ""{
		logrus.Error("Repository: password nil")
		return nil, 2
	}

	if err := helper.ComparePassword(user.Password, password); err != nil {
		logrus.Error("Repository: Login data error,", err.Error())
		return nil, 2
	}

	return &user,0
}

func (um *UsersModel) MyProfile(id uint) (*model.Users, int) {
	var user = model.Users{}

	if err := um.db.First(&user, id).Error; err != nil {
		logrus.Error("Repository: Myprofil data error,", err.Error())
		return nil, 1
	}

	return &user,0
}

func (um *UsersModel) UpdateMyProfile(updateUser model.Users) (*model.Users, int) {
	var user = model.Users{}
	var count int64

	if err := um.db.First(&user, updateUser.ID).Error; err != nil {
		logrus.Error("Repository: Select method UpdateMyProfile data error, ", err.Error())
		return nil, 1
	}

	if user.Email != updateUser.Email {
		var existingUser = model.Users{}
		if err := um.db.Where("email = ?", updateUser.Email).First(&existingUser).Count(&count).Error; err != nil {
			logrus.Error("Repository: UpdateMyProfile, Error checking existing email", err.Error())
			return nil, 1
		}
		if count > 0 {
			logrus.Error("Repository: UpdateMyProfile, Email already registered")
			return nil, 2
		}
	}

	if updateUser.Password != "" {
		hashedPassword := helper.HashPassword(updateUser.Password)
		user.Password = hashedPassword
	}

	user.Name = updateUser.Name
	user.Email = updateUser.Email

	var qry = um.db.Save(&user)
	if err := qry.Error; err != nil {
		logrus.Error("Repository: Save method UpdateMyProfile data error, ", err.Error())
		return nil, 1
	}

	return &user, 0
}



