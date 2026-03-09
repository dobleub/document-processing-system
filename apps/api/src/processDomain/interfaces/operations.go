package interfaces

import (
	"time"
)

type OperationStatusEnum string

// ● PENDING - Process created but not started
// ● RUNNING - Processing in progress
// ● PAUSED - Process temporarily paused
// ● COMPLETED - Process finished successfully
// ● FAILED - Process terminated with errors
// ● STOPPED - Process manually stopped
const (
	Pending   OperationStatusEnum = "PENDING"
	Running   OperationStatusEnum = "RUNNING"
	Paused    OperationStatusEnum = "PAUSED"
	Completed OperationStatusEnum = "COMPLETED"
	Failed    OperationStatusEnum = "FAILED"
	Stopped   OperationStatusEnum = "STOPPED"
)

type OperationProgress struct {
	TotalFiles     int `json:"total_files"`
	ProcessedFiles int `json:"processed_files"`
	Percentage     int `json:"percentage"`
}

type OperationResults struct {
	TotalWords        int      `json:"total_words"`
	TotalLines        int      `json:"total_lines"`
	MostFrequentWords []string `json:"most_frequent_words,omitempty"`
	FilesProcessed    []string `json:"files_processed,omitempty"`
}

type OperationStatus struct {
	ID                  string              `json:"process_id"`                     // Unique identifier for the process
	Status              OperationStatusEnum `json:"status,omitempty"`               // Current status of the process
	Progress            OperationProgress   `json:"progress,omitempty"`             // Progress details of the process
	Result              OperationResults    `json:"results,omitempty"`              // Result of the process (e.g., summary, statistics)
	Error               string              `json:"error,omitempty"`                // Error message if the process failed
	StartedAt           string              `json:"started_at,omitempty"`           // Timestamp when the process started
	EstimatedCompletion string              `json:"estimated_completion,omitempty"` // Time taken for the process to complete
	CompletedAt         string              `json:"completed_at,omitempty"`         // Timestamp when the process completed
}

func (o *OperationStatus) Initialize(id string) {
	o.ID = id
	o.Status = Pending
	o.Progress = OperationProgress{
		TotalFiles:     0,
		ProcessedFiles: 0,
		Percentage:     0,
	}
	o.Result = OperationResults{
		TotalWords:        0,
		TotalLines:        0,
		MostFrequentWords: []string{},
		FilesProcessed:    []string{},
	}
	o.Error = ""
	o.StartedAt = time.Now().Format(time.RFC3339)
	o.EstimatedCompletion = ""
}

func (o *OperationStatus) UpdateOperationStatus(updateData map[string]interface{}) {
	if status, ok := updateData["status"].(OperationStatusEnum); ok {
		o.Status = status
	}

	if progress, ok := updateData["progress"].(OperationProgress); ok {
		o.Progress = progress
	} else if progressMap, ok := updateData["progress"].(map[string]interface{}); ok {
		oldProgress := o.Progress
		newProgress := OperationProgress{
			TotalFiles:     int(progressMap["total_files"].(int)),
			ProcessedFiles: int(progressMap["processed_files"].(int)),
			Percentage:     int(progressMap["percentage"].(int)),
		}

		if oldProgress.TotalFiles != newProgress.TotalFiles {
			o.Progress.TotalFiles = newProgress.TotalFiles
		}
		if oldProgress.ProcessedFiles != newProgress.ProcessedFiles {
			o.Progress.ProcessedFiles = newProgress.ProcessedFiles
		}
		if oldProgress.Percentage != newProgress.Percentage {
			o.Progress.Percentage = newProgress.Percentage
		}

	}

	if result, ok := updateData["results"].(OperationResults); ok {
		o.Result = result
	} else if resultMap, ok := updateData["results"].(map[string]interface{}); ok {
		oldResult := o.Result
		newResult := OperationResults{
			TotalWords: int(resultMap["total_words"].(int)),
			TotalLines: int(resultMap["total_lines"].(int)),
		}

		if oldResult.TotalWords != newResult.TotalWords {
			o.Result.TotalWords = newResult.TotalWords
		}
		if oldResult.TotalLines != newResult.TotalLines {
			o.Result.TotalLines = newResult.TotalLines
		}

		// Handle MostFrequentWords and FilesProcessed if they are present in the update
		mfWordsTmp := resultMap["most_frequent_words"]
		filesProcessedTmp := resultMap["files_processed"]

		// cast to []string if they are present
		if mfWordsTmp != nil {
			if mfWords, ok := mfWordsTmp.([]string); ok {
				o.Result.MostFrequentWords = []string{}
				for _, word := range mfWords {
					o.Result.MostFrequentWords = append(o.Result.MostFrequentWords, word)
				}
			}
		}
		if filesProcessedTmp != nil {
			if filesProcessed, ok := filesProcessedTmp.([]string); ok {
				o.Result.FilesProcessed = []string{}
				for _, file := range filesProcessed {
					o.Result.FilesProcessed = append(o.Result.FilesProcessed, file)
				}
			}
		}
	}

	if errorMsg, ok := updateData["error"].(string); ok {
		o.Error = errorMsg
	}

	if elapsedTime, ok := updateData["estimated_completion"].(string); ok {
		if elapsedTime != "" {
			o.EstimatedCompletion = elapsedTime
		}
	}
}

func (o *OperationStatus) MarkAsCompleted() {
	o.Status = Completed
	o.Progress.Percentage = 100
	o.CompletedAt = time.Now().Format(time.RFC3339)
}

func (o *OperationStatus) MarkAsFailed(errorMsg string) {
	o.Status = Failed
	o.Error = errorMsg
	o.CompletedAt = time.Now().Format(time.RFC3339)
}

func (o *OperationStatus) MarkAsStopped() {
	o.Status = Stopped
	o.CompletedAt = time.Now().Format(time.RFC3339)
}
