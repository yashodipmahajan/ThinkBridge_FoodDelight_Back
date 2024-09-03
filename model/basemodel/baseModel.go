package basemodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type BaseModel struct {
	ID        uuid.UUID  `json:"id" gorm:"type:char(36);primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" sql:"index"`
}

// Hook to Autocreate uuid Entry
func (base *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New().String())
	return nil
}
