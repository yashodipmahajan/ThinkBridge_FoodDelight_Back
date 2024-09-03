package model

import (
	db "FoodDelight/DB"
	"FoodDelight/model/basemodel"
	"errors"

	"github.com/google/uuid"
)

type Admin struct {
	basemodel.BaseModel
	Fullname string `json:"fullname" gorm:"type:varchar(100)"`
	Username string `json:"username" gorm:"unique;type:varchar(100)"`
	Password string `json:"password" gorm:"type:varchar(100)"`
	Email    string `json:"email" gorm:"type:varchar(100);unique"`
}

func FindAdminByID(id uuid.UUID) (error, *Admin) {

	db := db.NewDBConnection()

	defer db.Close()

	admin := &Admin{}

	err := db.Where("id = ?", id).First(admin).Error
	if err != nil {

		return errors.New("USER NOT FOUND"), nil
	}

	return nil, admin

}
