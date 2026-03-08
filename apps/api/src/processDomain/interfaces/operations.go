package interfaces

import "time"

type OperationStatusEnum string

// ● PENDING - Process created but not started
// ● RUNNING - Processing in progress
// ● PAUSED - Process temporarily paused
// ● COMPLETED - Process finished successfully
// ● FAILED - Process terminated with errors
// ● STOPPED - Process manually stopped
const (
	Pending   OperationStatusEnum = "pending"
	Running   OperationStatusEnum = "running"
	Paused    OperationStatusEnum = "paused"
	Completed OperationStatusEnum = "completed"
	Failed    OperationStatusEnum = "failed"
	Stopped   OperationStatusEnum = "stopped"
)

type OperationStatus struct {
	ID           string              `json:"id"`                     // Unique identifier for the process
	Status       OperationStatusEnum `json:"status"`                 // Current status of the process
	Progress     int                 `json:"progress"`               // Progress percentage (0-100)
	Result       string              `json:"result,omitempty"`       // Result of the process (e.g., summary, statistics)
	Error        string              `json:"error,omitempty"`        // Error message if the process failed
	EnlapsedTime string              `json:"elapsed_time,omitempty"` // Time taken for the process to complete
	StartedAt    string              `json:"started_at,omitempty"`   // Timestamp when the process started
	CompletedAt  string              `json:"completed_at,omitempty"` // Timestamp when the process completed
}

func (o *OperationStatus) Initialize(id string) {
	o.ID = id
	o.Status = Pending
	o.Progress = 0
	o.Result = ""
	o.Error = ""
	o.EnlapsedTime = ""
	o.StartedAt = time.Now().Format(time.RFC3339)
}

func (o *OperationStatus) UpdateOperationStatus(status OperationStatusEnum, progress int, result string, errorMsg string, elapsedTime string) {
	o.Status = status
	o.Progress = progress
	o.Result = result
	o.Error = errorMsg
	o.EnlapsedTime = elapsedTime
}

func (o *OperationStatus) MarkAsCompleted(result string, elapsedTime string) {
	o.Status = Completed
	o.Progress = 100
	o.Result = result
	o.Error = ""
	o.EnlapsedTime = elapsedTime
}

func (o *OperationStatus) MarkAsFailed(errorMsg string, elapsedTime string) {
	o.Status = Failed
	o.Progress = 0
	o.Result = ""
	o.Error = errorMsg
	o.EnlapsedTime = elapsedTime
}
