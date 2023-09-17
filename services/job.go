package services

import (
	"context"
	"fmt"

	"time"

	"jobmanagement/enums"
	"jobmanagement/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JobHandler interface {
	Execute(ctx context.Context, job *models.Job) error
}

type workerID int

type JobManagement struct {
	db *gorm.DB

	availabeWorker chan workerID
	workers        map[workerID]chan *models.Job
	jobHandlers    map[enums.JobType]JobHandler
}

func NewJobManagement(ctx context.Context, db *gorm.DB, numberOfWorkers int) *JobManagement {
	jobHandlers := make(map[enums.JobType]JobHandler)
	workers := make(map[workerID]chan *models.Job)
	availabelWorkers := make(chan workerID, numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		workerChan := make(chan *models.Job)
		workers[workerID(i)] = workerChan
		go func(ctx context.Context, db *gorm.DB, ch <-chan *models.Job, workerID workerID) {
			availabelWorkers <- workerID

			for job := range ch {
				fmt.Printf("job %s %s is running\n", job.ID, job.Type)
				jobHandler, ok := jobHandlers[enums.JobType(job.Type)]
				if !ok {
					fmt.Printf("jobType %s not found in jobHandlers\n", job.Type)
				}
				if err := jobHandler.Execute(ctx, job); err != nil {
					db.Table("jobs").
						Where("id = ?", job.ID).
						Update("status", "FAILED").
						Update("logs", err.Error())
				} else {
					db.Table("jobs").
						Where("id = ?", job.ID).
						Update("status", "SUCCEEDED")
				}

				availabelWorkers <- workerID
			}

		}(ctx, db, workerChan, workerID(i))
	}

	return &JobManagement{
		db:             db,
		availabeWorker: availabelWorkers,
		workers:        workers,
		jobHandlers:    jobHandlers,
	}
}

func (jm *JobManagement) Register(typ enums.JobType, jobHandler JobHandler) {
	fmt.Printf("register job type %s\n", typ)
	jm.jobHandlers[typ] = jobHandler
}

func (jm *JobManagement) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case workerID := <-jm.availabeWorker:
				fmt.Printf("run %d %s\n", workerID, time.Now())
				jm.loadAndHandeJobs(ctx, workerID)
			}
		}
	}()
}

func (jm *JobManagement) loadAndHandeJobs(ctx context.Context, workerID workerID) {
	var job models.Job
	err := jm.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Table("jobs").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("status = ? AND run_at <= NOW() AND executed_times < times", "READY").
			First(&job).
			Order("level DESC, run_at ASC").
			WithContext(ctx)
		if err := result.Error; err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return err
		}

		updateJobResult := tx.Table("jobs").
			Where("id = ?", job.ID).
			Update("status", "PICKED").
			Update("executed_times", job.ExecutedTimes+1).
			Update("execute_at", "NOW()").
			WithContext(ctx)
		if err := updateJobResult.Error; err != nil {
			fmt.Printf("updateJobResult error: %s\n", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		time.Sleep(1 * time.Second)
		jm.availabeWorker <- workerID
		return
	}
	workerCh := jm.workers[workerID]
	workerCh <- &job
}

func (jm *JobManagement) RunJobWatcher(ctx context.Context) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				jm.loadAndUpdateJobs(ctx)
				time.Sleep(1000)
			}
		}
	}()
}

func (jm *JobManagement) loadAndUpdateJobs(ctx context.Context) {
	err := jm.db.Transaction(func(tx *gorm.DB) error {
		updateJobResult := tx.Table("jobs").
			Where("status IN ('PICKED', 'FAILED') AND execute_at <= NOW() - INTERVAL '2 seconds' AND executed_times < times").
			Updates(map[string]interface{}{
				"status": "READY",
				"logs":   nil,
			}).
			WithContext(ctx)
		if err := updateJobResult.Error; err != nil {
			fmt.Printf("updateJobResult error: %s\n", err.Error())
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Printf("loadAndUpdateJobs error %s", err.Error())
	}
}
