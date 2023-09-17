package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"jobmanagement/controllers"
	"jobmanagement/enums"
	"jobmanagement/jobs"
	"jobmanagement/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type config struct {
	Host     string
	Database string
	Username string
	Password string
	Port     string
}

func (c *config) Load() {
	c.Host = os.Getenv("POSTGRES_HOST")
	c.Database = os.Getenv("POSTGRES_DATABASE")
	c.Username = os.Getenv("POSTGRES_USERNAME")
	c.Password = os.Getenv("POSTGRES_PASSWORD")
	c.Port = os.Getenv("POSTGRES_PORT")
}

func main() {
	cf := config{}
	cf.Load()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", cf.Host, cf.Username, cf.Password, cf.Database, cf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}

	printerJob := jobs.NewPrinterJob()

	ctx := context.Background()
	jobManagement := services.NewJobManagement(ctx, db, 10000)
	jobManagement.Register(enums.JobTypePrinterJob, printerJob)
	jobManagement.Run(ctx)
	jobManagement.RunJobWatcher(ctx)

	r := gin.Default()
	controller := controllers.NewController(r, db)
	controller.Routes()

	r.Run("localhost:8080")

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	<-signChan
}
