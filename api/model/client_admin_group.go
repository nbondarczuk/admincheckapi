package model

import (
	"gorm.io/gorm"
)

//
// ClitnAdminGroup is the base entity of the system. It links
// Client with it's Admin Group identified by its' ID as stored in AD.
// The same ID is provided in the JWT token so that the match can be done.
//
type ClientAdminGroup struct {
	Client       string `gorm:"index"`
	AdminGroupId string
	gorm.Model
}
