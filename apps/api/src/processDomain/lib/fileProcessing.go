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
	// [ ] Integrate with File Manager to read and analyze files
	// [ ] Handle batch processing of documents
	// [ ] Ensure efficient processing for large files
	Path  string
	State *sync.Map
}

func (f *FileProcessing) UpdateState(process_id string, status pd_interfaces.OperationStatusEnum, progress int, result string, errorMsg string, elapsedTime string) {
	// Helper to update status safely.
	if val, ok := f.State.Load(process_id); ok {
		op := val.(*pd_interfaces.OperationStatus)
		op.UpdateOperationStatus(status, progress, result, errorMsg, elapsedTime)
		
    if status == pd_interfaces.Completed || status == pd_interfaces.Failed {
			op.CompletedAt = time.Now().Format(time.RFC3339)
		}
		
    f.State.Store(process_id, op)
	}
}

func (f *FileProcessing) ProcessFilesFromDirectory(process_id string) map[string]interface{} {
	// Processing System
	// [ ] Load text files from a specific folder
	// [ ] Process documents in batches (batch processing)
	// [ ] Extract basic statistics: word count, lines, characters
	// [ ] Identify most frequent words
	// [●] Generate a content summary
	f.UpdateState(process_id, pd_interfaces.Running, 0, "", "", "")

	results := map[string]interface{}{}
  progress := 0

	FileManager := &FileManager{
		Path: f.Path,
	}
	files := FileManager.ListFilesFromPath()

	nfiles := len(files)
	if nfiles == 0 {
		f.UpdateState(process_id, pd_interfaces.Completed, 100, "No files to process", "", "")
		return results
	}

	for idx, file := range files {
		file_name := strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
		fmt.Printf("Processing file: %s\n", file_name)
		time.Sleep(5 * time.Second)

		content := FileManager.ReadFileContent(file)
		summary := FileManager.GenerateSummary(content)
		frequentWords := FileManager.IdentifyFrequentWords(content)
		statistics := FileManager.ExtractStatistics(content)

		// Update state after processing each file (for demonstration purposes)
		progress = int(float64(idx+1) / float64(nfiles) * 100)
		f.UpdateState(process_id, pd_interfaces.Running, progress, "", "", "")
    
		// Store the results in a database or return them
		// For demonstration, we are just logging the results in memory
		// log the results (for demonstration purposes)
		results[file_name] = map[string]interface{}{
      "file":          file,
			"summary":       summary,
			"frequentWords": frequentWords,
			"statistics":    statistics,
		}
	}

  if (progress >= 100) {
    f.UpdateState(process_id, pd_interfaces.Completed, progress, "", "", "")
  }

	return results
}

func (f *FileProcessing) ProcessScannerBytes([]byte) {
	// Process Scanner Bytes
	// [ ] Read bytes from a scanner (e.g., file upload)
	// [ ] Convert bytes to a readable format (e.g., string)
	// [ ] Process the content using the same logic as file processing
}
