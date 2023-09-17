package enums

type JobStatus string

const (
	JobStatusReady     JobStatus = "READY"
	JobStatusPicked    JobStatus = "PICKED"
	JobStatusFailed    JobStatus = "FAILED"
	JobStatusSucceeded JobStatus = "SUCCEEDED"
)
