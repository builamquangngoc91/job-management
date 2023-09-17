package models

import (
	"time"
)

type (
	Job struct {
		ID            string     `json:"id"`
		Name          string     `json:"name"`
		Data          string     `json:"data"`
		RunAt         time.Time  `json:"run_at"`
		ExecuteAt     *time.Time `json:"execute_at"`
		Status        string     `json:"status"`
		TTL           int64      `json:"ttl"`            // how long need to complete this job (duration) default 3s
		Times         int64      `json:"times"`          // how many times you want to retry this job: default = 1
		ExecutedTimes int64      `json:"executed_times"` // how many times were retried this job
		Note          *string    `json:"note"`
		Level         int64      `json:"level"`
		Type          string     `json:"type"`
		Logs          *string    `json:"logs"`
		CreatedAt     *time.Time `json:"created_at"`
		UpdatedAt     *time.Time `json:"updated_at"`
	}
)
