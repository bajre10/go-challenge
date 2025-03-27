package database

import (
	"errors"
	"gorm.io/gorm"
	"shortcut-challenge/models"
)

var roles = []models.Role{
	{
		Name: string(models.ADMIN),
	},
	{
		Name: string(models.USER),
	},
}

func SetRoles(db *gorm.DB) error {
	// Check if roles already exist
	for _, role := range roles {
		var existingRole models.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			// Create role if it doesn't exist
			if dbErr := db.Create(&role).Error; dbErr != nil {
				return dbErr
			}
		}
	}

	return nil
}

func SetUsers(db *gorm.DB) error {
	var adminRole models.Role
	if err := db.Where("name = ?", string(models.ADMIN)).First(&adminRole).Error; err != nil {
		return err
	}

	defaultUser := &models.User{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: "password",
	}

	// Create user of role "Administrator" if it doesn't exist
	var users []models.User
	if err := db.Where("role_id = ?", adminRole.ID).First(&users).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create([]models.User{{
			Name:     defaultUser.Name,
			Email:    defaultUser.Email,
			Password: defaultUser.Password,
			RoleID:   adminRole.ID},
		}).Error
	}

	return nil
}
