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

	files_bach := [][]string{}
	batchSize := 2 // Placeholder for batch size, this could be configurable
	for i := 0; i < nfiles; i += batchSize {
		end := i + batchSize
		if end > nfiles {
			end = nfiles
		}
		files_bach = append(files_bach, files[i:end])
	}

	for _, batch := range files_bach {
		resultsBatch := f.ProcessBatchDocuments(process_id, batch, nfiles)
		// Update results and progress after processing each batch
		results["total_words"] = results["total_words"].(int) + resultsBatch["total_words"].(int)
		results["total_lines"] = results["total_lines"].(int) + resultsBatch["total_lines"].(int)
		results["most_frequent_words"] = append(results["most_frequent_words"].([]string), resultsBatch["most_frequent_words"].([]string)...)
		filesProcessed = append(filesProcessed, resultsBatch["files_processed"].([]string)...)

		progress["processed_files"] = len(filesProcessed)
		progress["percentage"] = (progress["processed_files"].(int) * 100) / progress["total_files"].(int)

		f.UpdateState(process_id, pd_interfaces.Running, progress, results, "", "")
	}

	return results
}

func (f *FileProcessing) ProcessBatchDocuments(process_id string, files []string, nfiles int) map[string]interface{} {
	// Batch Document Processing
	// [x] Accept a batch of documents (e.g., multiple file paths or byte arrays)
	// [x] Process each document in the batch
	// [●] Update progress and results for the entire batch
	results := map[string]interface{}{"total_words": 0, "total_lines": 0, "most_frequent_words": []string{}, "files_processed": []string{}}

	for idx, file := range files {
		start_time := time.Now()
		file_name := strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
		fmt.Printf("Processing file: %s\n", file_name)
		time.Sleep(5 * time.Second)

		// TODO: process the file and extract statistics
		results["total_words"] = results["total_words"].(int) + 100                                                    // Placeholder for word count
		results["total_lines"] = results["total_lines"].(int) + 10                                                     // Placeholder for line count
		results["most_frequent_words"] = append(results["most_frequent_words"].([]string), fmt.Sprintf("word%d", idx)) // Placeholder for most frequent words
		results["files_processed"] = append(results["files_processed"].([]string), file_name)

		// f.UpdateState(process_id, pd_interfaces.Running, progress, results, "", "")
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
