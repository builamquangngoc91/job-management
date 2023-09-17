package jobs

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"jobmanagement/models"
)

type PrinterJob struct{}

func NewPrinterJob() *PrinterJob {
	return &PrinterJob{}
}

func (p *PrinterJob) Execute(ctx context.Context, job *models.Job) error {
	randomNumber := rand.Int31()
	fmt.Printf("Job (%s) random number (%d)\n", job.ID, randomNumber)
	if randomNumber%2 == 0 {
		return errors.New("error")
	}
	return nil
}
