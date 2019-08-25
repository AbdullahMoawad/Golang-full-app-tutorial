package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Session struct {
	gorm.Model
	UserId uint
	SessionId uuid.UUID
}
