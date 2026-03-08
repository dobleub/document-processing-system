package interfaces

type OperationResponse struct {
	Status   OperationStatus         `json:"status"`   // Current status of the process (e.g., "pending", "in_progress", "completed", "failed")
	Analysis OperationAnalysisResult `json:"analysis"` // Result of the process (e.g., summary, statistics)
}

func (o *OperationResponse) Initialize(id string) {
	operationStatus := &OperationStatus{}
	operationStatus.Initialize(id)
	o.Status = *operationStatus

	operationAnalysis := &OperationAnalysisResult{}
	operationAnalysis.Initialize(id)
	o.Analysis = *operationAnalysis
}
