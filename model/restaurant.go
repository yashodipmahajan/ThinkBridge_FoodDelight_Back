package model

import "FoodDelight/model/basemodel"

type Restaurant struct {
	basemodel.BaseModel
	Name        string `json:"name" gorm:"unique;type:varchar(100)"`
	Description string `json:"description" gorm:"type:varchar(200)"`
	Location    string `json:"location" gorm:"type:varchar(200)"`
}
