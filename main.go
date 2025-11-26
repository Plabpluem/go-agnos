package main

import (
	"agnos/internal/entities"
	"agnos/internal/routes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func main() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5433")
	user := getEnv("DB_USER", "myuser")
	password := getEnv("DB_PASSWORD", "mypassword")
	dbname := getEnv("DB_NAME", "mydatabase")

	dbs := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dbs), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic("Can't connect database")
	}

	db.AutoMigrate(&entities.Patient{}, &entities.Staff{})
	router := gin.Default()

	routes.StaffRoutes(&router.RouterGroup, db)

	routes.PatientRoutes(&router.RouterGroup, db)

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	router.Run(":8080")
}
