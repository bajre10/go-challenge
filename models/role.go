package models

import "gorm.io/gorm"

type RoleTypes string

const (
	ADMIN RoleTypes = "Administrator"
	USER  RoleTypes = "User"
)

type Role struct {
	gorm.Model `swaggerignore:"true"`
	Name       string `gorm:"not null unique"`
}
