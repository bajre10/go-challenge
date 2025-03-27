package models

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"shortcut-challenge/api/appErrors"
)

type User struct {
	gorm.Model `swaggerignore:"true"`
	Name       string `json:"name" gorm:"not null"`
	Email      string `json:"email" gorm:"unique" validate:"required,email"`
	Password   string `json:"password" gorm:"not null" validate:"required,min=8"`
	RoleID     uint   `json:"roleID" gorm:"not null"`
	Role       Role   `json:"-" gorm:"foreignKey:RoleID" swaggerignore:"true"`
}

func (r *User) BeforeSave(*gorm.DB) (err error) {
	// Validate inputs
	if err := ValidateUser(r); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			firstErr := validationErrors[0]

			errMessage := fmt.Sprintf("%s is not valid", firstErr.Field())

			// Return error message based on error type
			switch firstErr.Field() {
			case "Password":
				if firstErr.Tag() == "required" {
					errMessage = "Password is required."
				} else if firstErr.Tag() == "min" {
					errMessage = fmt.Sprintf("Password must be at least %s characters long.", firstErr.Param())
				}
			case "Email":
				if firstErr.Tag() == "required" {
					errMessage = "Email is required."
				} else if firstErr.Tag() == "email" {
					errMessage = "Invalid email format."
				}
			}

			// Retrun a DBError
			return &appErrors.DBError{
				Field:   firstErr.Field(),
				Type:    appErrors.Validation,
				Message: errMessage,
			}
		} else {
			return err
		}
	}

	// Hash password
	hashedPassword, err := r.HashPassword()
	if err != nil {
		return err
	}

	r.Password = string(hashedPassword)

	return nil
}

// HashPassword Hashes user's password
func (r *User) HashPassword() ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func (r *User) CheckPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(password)); err != nil {
		return err
	}

	return nil
}

var validate = validator.New()

func ValidateUser(user *User) error {
	return validate.Struct(user)
}
