package interfaces

import "path/filepath"

type OperationFileAnalysis struct {
	FileName          string   `json:"file_name" bson:"file_name"`
	TotalWords        int      `json:"total_words" bson:"total_words"`
	TotalLines        int      `json:"total_lines" bson:"total_lines"`
	MostFrequentWords []string `json:"most_frequent_words,omitempty" bson:"most_frequent_words"`
	TotalCharacters   int      `json:"total_characters" bson:"total_characters"`
	Summary           string   `json:"summary,omitempty" bson:"summary"`
}

type OperationAnalysisResult struct {
	ID       string                  `json:"process_id" bson:"process_id"`       // Unique identifier for the process
	Status   OperationStatusEnum     `json:"status" bson:"status"`               // Current status of the process (e.g., "pending", "in_progress", "completed", "failed")
	Analysis []OperationFileAnalysis `json:"analysis,omitempty" bson:"analysis"` // Result of the process (e.g., summary, statistics)
}

func (o *OperationAnalysisResult) Initialize(id string) {
	o.ID = id
	o.Status = Pending
	o.Analysis = []OperationFileAnalysis{}
}

func (o *OperationAnalysisResult) AppendBatchAnalysis(file_analysis []map[string]interface{}) {
	// o.Analysis = append(o.Analysis, file_analysis)
	for _, fa := range file_analysis {
		analysis := OperationFileAnalysis{}
		if fileName, ok := fa["file_name"].(string); ok {
			analysis.FileName = filepath.Base(fileName)
		}
		if totalWords, ok := fa["total_words"].(int); ok {
			analysis.TotalWords = totalWords
		}
		if totalLines, ok := fa["total_lines"].(int); ok {
			analysis.TotalLines = totalLines
		}
		if frequentWords, ok := fa["most_frequent_words"].([]string); ok {
			analysis.MostFrequentWords = frequentWords
		}
		if totalCharacters, ok := fa["total_characters"].(int); ok {
			analysis.TotalCharacters = totalCharacters
		}
		if summary, ok := fa["summary"].(string); ok {
			analysis.Summary = summary
		}

		o.Analysis = append(o.Analysis, analysis)
	}
}

func (o *OperationAnalysisResult) MarkAsCompleted() {
	o.Status = Completed
}

func (o *OperationAnalysisResult) MarkAsFailed(errorMsg string) {
	o.Status = Failed
}

func (o *OperationAnalysisResult) MarkAsStopped() {
	o.Status = Stopped
}
