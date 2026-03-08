package lib

import (
	"path/filepath"
	"sync"
	"time"

	pd_interfaces "nx-recipes/dps/lambda/src/processDomain/interfaces"

	"go.uber.org/zap"
)

type FileProcessing struct {
	// File Processing
	// [ ] Implement the core logic for processing text files
	// [x] Integrate with File Manager to read and analyze files
	// [ ] Handle batch processing of documents
	// [ ] Ensure efficient processing for large files
	Path  string
	State *sync.Map
	Log   *zap.Logger
}

func (f *FileProcessing) logger() *zap.Logger {
	if f.Log != nil {
		return f.Log
	}
	return zap.NewNop()
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
	start_time := time.Now()
	progress := map[string]interface{}{"total_files": 0, "processed_files": 0, "percentage": 0}
	results := map[string]interface{}{"total_words": 0, "total_lines": 0, "most_frequent_words": []string{}, "files_processed": []string{}}
	filesProcessed := []string{}

	f.UpdateState(process_id, pd_interfaces.Running, progress, results, "", "")

	FileManager := &FileManager{
		Path: f.Path,
		Log:  f.logger(),
	}
	files := FileManager.ListFilesFromPath()

	nfiles := len(files)
	if nfiles == 0 {
		f.UpdateState(process_id, pd_interfaces.Completed, progress, results, "No files to process", "")
		return results
	}

	totalEstimatedTime := FileManager.EstimateProcessFiles()
	progress["total_files"] = nfiles
	f.UpdateState(process_id, pd_interfaces.Running, progress, results, "", totalEstimatedTime)

	files_bach := [][]map[string]interface{}{}
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

	enlapsed_time := time.Since(start_time)
	f.logger().Info("Total processing time", zap.String("duration", enlapsed_time.String()), zap.String("process_id", process_id))
	return results
}

func (f *FileProcessing) ProcessBatchDocuments(process_id string, files []map[string]interface{}, nfiles int) map[string]interface{} {
	// Batch Document Processing
	// [x] Accept a batch of documents (e.g., multiple file paths or byte arrays)
	// [x] Process each document in the batch
	// [●] Update progress and results for the entire batch
	results := map[string]interface{}{"total_words": 0, "total_lines": 0, "most_frequent_words": []string{}, "files_processed": []string{}}

	type fileResult struct {
		idx           int
		fileName      string
		totalWords    int
		totalLines    int
		frequentWords []string
	}

	resultsByFile := make([]fileResult, len(files))
	var wg sync.WaitGroup

	for idx, file := range files {
		idx := idx
		file := file

		wg.Add(1)
		go func() {
			defer wg.Done()

			startTime := time.Now()
			fileName := filepath.Base(file["name"].(string))
			f.logger().Info("Processing file", zap.String("file", fileName), zap.String("process_id", process_id))
			time.Sleep(5 * time.Second)

			// TODO: process the file and extract statistics.
			resultsByFile[idx] = fileResult{
				idx:           idx,
				fileName:      fileName,
				totalWords:    100,
				totalLines:    10,
				frequentWords: []string{"example", "test"},
			}

			elapsedTime := time.Since(startTime).String()
			f.logger().Info("Completed processing file", zap.String("file", fileName), zap.String("duration", elapsedTime), zap.String("process_id", process_id))
		}()
	}

	wg.Wait() // Wait for all files in the batch to be processed

	// Aggregate results from all files in the batch
	for _, fileResult := range resultsByFile {
		results["total_words"] = results["total_words"].(int) + fileResult.totalWords
		results["total_lines"] = results["total_lines"].(int) + fileResult.totalLines
		results["most_frequent_words"] = append(results["most_frequent_words"].([]string), fileResult.frequentWords...)
		results["files_processed"] = append(results["files_processed"].([]string), fileResult.fileName)
	}

	return results
}

func (f *FileProcessing) ProcessScannerBytes([]byte) {
	// Process Scanner Bytes
	// [ ] Read bytes from a scanner (e.g., file upload)
	// [ ] Convert bytes to a readable format (e.g., string)
	// [ ] Process the content using the same logic as file processing
}
