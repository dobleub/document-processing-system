package interfaces

import "sort"

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
	ID                  string   `json:"id"`                   // Unique identifier for the process
	Status              string   `json:"status"`               // Current status of the process (e.g., "pending", "in_progress", "completed", "failed")
	Error               string   `json:"error"`                // Error message if the process failed
	StartedAt           string   `json:"started_at"`           // Timestamp when the process started
	EstimatedCompletion string   `json:"estimated_completion"` // Estimated completion time
	FilesProcessed      []string `json:"files_processed"`      // List of files processed, this is added to give more insights about the process, it can be used to show the progress of the process in the frontend
	FilesToProcess      []string `json:"files_to_process"`     // List files to process
	CompletedAt         string   `json:"completed_at"`         // Timestamp when the process completed
}

func (o *OperationReview) Initialize(id string) {
	o.ID = id
	o.Status = "pending"
	o.Error = ""
	o.StartedAt = ""
	o.EstimatedCompletion = ""
	o.FilesProcessed = []string{}
	o.FilesToProcess = []string{}
	o.CompletedAt = ""
}

type OperationListResponse struct {
	Processes []OperationReview `json:"processes"` // List of processes with their current status
}

func (o *OperationListResponse) Initialize() {
	o.Processes = []OperationReview{}
}

func (o *OperationListResponse) AddProcess(process OperationReview) {
	o.Processes = append(o.Processes, process)
}

func (o *OperationListResponse) OrderProcesses() {
	// order by startedAt and then by ID DESC to ensure consistent ordering
	// this is important to prevent unnecessary updates due to ordering changes
	// when the status and files processed count are the same but the order changes
	sort.SliceStable(o.Processes, func(i, j int) bool {
		if o.Processes[i].StartedAt == o.Processes[j].StartedAt {
			return o.Processes[i].ID > o.Processes[j].ID
		}
		return o.Processes[i].StartedAt > o.Processes[j].StartedAt
	})
}
