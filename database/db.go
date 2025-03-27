package database

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"shortcut-challenge/api"
	"shortcut-challenge/models"
	"shortcut-challenge/utils"
)

var dbInstance *gorm.DB

func SetDBMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return error if dbInstance is not present
		if dbInstance == nil {
			api.ThrowInternalError(w)
		}

		ctx := context.WithValue(r.Context(), "DB", dbInstance)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetDBInstance(r *http.Request) *gorm.DB {
	// Get DB instance from request context
	if db, ok := r.Context().Value("DB").(*gorm.DB); ok {
		return db
	}
	return nil
}

func InitDB() *gorm.DB {
	var err error
	dbInstance, err = connectionDatabase()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	err = performMigration()
	if err != nil {
		log.Fatalf("Could not auto migrate: %v", err)
	}

	return dbInstance
}

func connectionDatabase() (*gorm.DB, error) {
	dbUsername := utils.GetEnv("DB_USERNAME", "root")
	dbPassword := utils.GetEnv("DB_PASSWORD", "")
	dbName := utils.GetEnv("DB_NAME", "gorm")
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "3306")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUsername, dbPassword, dbHost, dbPort)

	// Open a connection to MySQL server
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Create database if it doesn't already exist
	_ = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName + ";")

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)

	// Open a connection to the DB
	databaseConnection, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info), // Enable logging for debugging
		TranslateError: true,
	})

	if err != nil {
		return nil, err
	}

	_, err = databaseConnection.DB()
	if err != nil {
		return nil, err
	}

	return databaseConnection, nil
}

func performMigration() error {
	err := dbInstance.AutoMigrate(&models.InventoryItem{}, &models.User{}, &models.Role{}, &models.Restock{})
	if err != nil {
		return err
	}

	// Insert seed data
	err = SetRoles(dbInstance)
	if err != nil {
		return err
	}

	err = SetUsers(dbInstance)
	if err != nil {
		return err
	}

	return nil
}
