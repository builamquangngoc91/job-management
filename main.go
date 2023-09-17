package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jobmanagement/enums"
	"jobmanagement/jobs"
	"jobmanagement/models"
	"jobmanagement/services"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

	go func() {
		for {
			id := uuid.NewString()
			now := time.Now()
			job := models.Job{
				ID:        id,
				Name:      fmt.Sprintf("job %s", id),
				Data:      "{}",
				RunAt:     time.Now(),
				Times:     3,
				TTL:       100,
				Status:    string(enums.JobStatusReady),
				Level:     1,
				Type:      string(enums.JobTypePrinterJob),
				CreatedAt: &now,
				UpdatedAt: &now,
			}
			result := db.Table("jobs").Create(&job)
			if result.Error != nil {
				fmt.Printf("error: %s", result.Error.Error())
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	<-signChan
}
