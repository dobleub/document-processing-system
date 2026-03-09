package interfaces

type OperationResponse struct {
	Status   OperationStatus         `json:"status"`   // Current status of the process (e.g., "pending", "in_progress", "completed", "failed")
	Analysis OperationAnalysisResult `json:"analysis"` // Result of the process (e.g., summary, statistics)
	Stopper  chan struct{}           `json:"-"`        // Channel to signal stopping the process
}

func (o *OperationResponse) Initialize(id string) {
	operationStatus := &OperationStatus{}
	operationStatus.Initialize(id)
	o.Status = *operationStatus

	operationAnalysis := &OperationAnalysisResult{}
	operationAnalysis.Initialize(id)
	o.Analysis = *operationAnalysis
	o.Stopper = make(chan struct{})
}

type OperationReview struct {
	ID                  string `json:"id"`                             // Unique identifier for the process
	Status              string `json:"status"`                         // Current status of the process (e.g., "pending", "in_progress", "completed", "failed")
	Error               string `json:"error,omitempty"`                // Error message if the process failed
	StartedAt           string `json:"started_at,omitempty"`           // Timestamp when the process started
	EstimatedCompletion string `json:"estimated_completion,omitempty"` // Estimated completion time
}

func (o *OperationReview) Initialize(id string) {
	o.ID = id
	o.Status = "pending"
	o.Error = ""
	o.StartedAt = ""
	o.EstimatedCompletion = ""
}

type OperationListResponse struct {
	Processes []OperationReview `json:"processes"` // List of processes with their current status
}
