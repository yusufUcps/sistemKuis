package model

type Login struct {
    Email    string `json:"email" gorm:"type:varchar(20);unique;not null"`
    Password string `json:"password" gorm:"not null"`
}