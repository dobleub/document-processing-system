package lib

import (
	"fmt"
	"strings"
	"sync"
	"time"

	pd_interfaces "nx-recipes/dps/lambda/src/processDomain/interfaces"
)

type FileProcessing struct {
	// File Processing
	// [ ] Implement the core logic for processing text files
	// [x] Integrate with File Manager to read and analyze files
	// [ ] Handle batch processing of documents
	// [ ] Ensure efficient processing for large files
	Path  string
	State *sync.Map
}

func (f *FileProcessing) UpdateState(process_id string, status pd_interfaces.OperationStatusEnum, progress map[string]interface{}, results map[string]interface{}, errorMsg string, estimated_completion string) {
	// Helper to update status safely.
	if val, ok := f.State.Load(process_id); ok {
		op := val.(*pd_interfaces.OperationStatus)

		op.UpdateOperationStatus(map[string]interface{}{
			"status":               status,
			"progress":             progress,
			"results":              results,
			"error":                errorMsg,
			"estimated_completion": estimated_completion,
		})

		if progress["percentage"].(int) >= 100 {
			op.MarkAsCompleted()
		}

		if errorMsg != "" {
			op.MarkAsFailed(errorMsg)
		}

		f.State.Store(process_id, op)
	}
}

func (f *FileProcessing) ProcessFilesFromDirectory(process_id string) map[string]interface{} {
	// Processing System
	// [x] Load text files from a specific folder
	// [ ] Process documents in batches (batch processing)
	// [ ] Extract basic statistics: word count, lines, characters
	// [ ] Identify most frequent words
	// [●] Generate a content summary
	progress := map[string]interface{}{"total_files": 0, "processed_files": 0, "percentage": 0}
	results := map[string]interface{}{"total_words": 0, "total_lines": 0, "most_frequent_words": []string{}, "files_processed": []string{}}
	filesProcessed := []string{}

	f.UpdateState(process_id, pd_interfaces.Running, progress, results, "", "")

	FileManager := &FileManager{
		Path: f.Path,
	}
	files := FileManager.ListFilesFromPath()

	nfiles := len(files)
	if nfiles == 0 {
		f.UpdateState(process_id, pd_interfaces.Completed, progress, results, "No files to process", "")
		return results
	}

	estimatedTimePerFile := 5 * time.Second // Placeholder for estimated time per file, in a real application this could be dynamic based on file size or type
	totalEstimatedTime := time.Duration(nfiles) * estimatedTimePerFile
	progress["total_files"] = nfiles
	f.UpdateState(process_id, pd_interfaces.Running, progress, results, "", totalEstimatedTime.String())

	for idx, file := range files {
		start_time := time.Now()
		file_name := strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
		fmt.Printf("Processing file: %s\n", file_name)
		time.Sleep(5 * time.Second)

		// TODO: process the file and extract statistics

		filesProcessed = append(filesProcessed, file_name)
		progress := map[string]interface{}{
			"total_files":     nfiles,
			"processed_files": idx + 1,
			"percentage":      int(float64(idx+1) / float64(nfiles) * 100),
		}
		results := map[string]interface{}{
			"total_words":         1000 + idx,                   // Placeholder for actual word count
			"total_lines":         100 + idx,                    // Placeholder for actual line count
			"most_frequent_words": []string{"the", "and", "to"}, // Placeholder for actual frequent words
			"files_processed":     filesProcessed,
		}
		f.UpdateState(process_id, pd_interfaces.Running, progress, results, "", "")
		elapsedTime := time.Since(start_time).String()
		fmt.Printf("Completed processing file: %s in %s\n", file_name, elapsedTime)
	}

	return results
}

func (f *FileProcessing) ProcessScannerBytes([]byte) {
	// Process Scanner Bytes
	// [ ] Read bytes from a scanner (e.g., file upload)
	// [ ] Convert bytes to a readable format (e.g., string)
	// [ ] Process the content using the same logic as file processing
}
