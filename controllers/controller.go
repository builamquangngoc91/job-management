package controllers

import (
	"jobmanagement/enums"
	"jobmanagement/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Controller struct {
	r  *gin.Engine
	db *gorm.DB
}

func NewController(router *gin.Engine, db *gorm.DB) *Controller {
	return &Controller{
		r:  router,
		db: db,
	}
}

func (ctl *Controller) Routes() {
	ctl.r.POST("/jobs", ctl.CreateJobs)
}

func (ctl *Controller) CreateJobs(c *gin.Context) {
	var job models.Job
	if err := c.BindJSON(&job); err != nil {
		return
	}

	var validateRequestMessage string
	if job.Data == "" {
		validateRequestMessage = "missing data"
	}
	if job.Type == "" {
		validateRequestMessage = "missing type"
	}
	if job.Name == "" {
		validateRequestMessage = "missing name"
	}
	if validateRequestMessage != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": validateRequestMessage,
		})
		return
	}

	now := time.Now()
	job.ID = uuid.NewString()
	job.Status = string(enums.JobStatusReady)
	if job.Level == 0 {
		job.Level = 1
	}
	if job.RunAt == nil {
		job.RunAt = &now
	}
	if job.Times == 0 {
		job.Times = 1
	}
	if job.TTL == 0 {
		job.TTL = 1000 // 1000 milliseconds
	}

	result := ctl.db.Table("jobs").Create(&job)
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
