package lib

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	pd_interfaces "nx-recipes/dps/lambda/src/processDomain/interfaces"
	summarizerDomainLib "nx-recipes/dps/lambda/src/summarizerDomain/lib"

	"go.uber.org/zap"
	"google.golang.org/genai"
)

type FileProcessing struct {
	// File Processing
	// [x] Implement the core logic for processing text files
	// [x] Integrate with File Manager to read and analyze files
	// [x] Handle batch processing of documents
	// [●] Ensure efficient processing for large files
	Path      string
	State     *sync.Map
	Log       *zap.Logger
	McpClient *genai.Client
}

func (f *FileProcessing) logger() *zap.Logger {
	if f.Log != nil {
		return f.Log
	}
	return zap.NewNop()
}

func (f *FileProcessing) UpdateState(process_id string, status pd_interfaces.OperationStatusEnum, progress map[string]interface{}, results map[string]interface{}, errorMsg string, estimated_completion string, batchAnalysis []map[string]interface{}) {
	// Helper to update status safely.
	if val, ok := f.State.Load(process_id); ok {
		opResponse := val.(*pd_interfaces.OperationResponse)

		opStatus := opResponse.Status
		opAnalysis := opResponse.Analysis

		if opStatus.Status == pd_interfaces.Stopped {
			f.logger().Info("Process is stopped, skipping state update", zap.String("process_id", process_id))
			return
		}

		opStatus.UpdateOperationStatus(map[string]interface{}{
			"status":               status,
			"progress":             progress,
			"results":              results,
			"error":                errorMsg,
			"estimated_completion": estimated_completion,
		})

		if status != "" {
			opAnalysis.Status = status
		}

		if progress["percentage"].(int) >= 100 {
			opStatus.MarkAsCompleted()
			opAnalysis.Status = pd_interfaces.Completed
		}

		if errorMsg != "" {
			opStatus.MarkAsFailed(errorMsg)
			opAnalysis.Status = pd_interfaces.Failed
		}

		if len(batchAnalysis) > 0 {
			opAnalysis.AppendBatchAnalysis(batchAnalysis)
		}

		opResponse.Status = opStatus
		opResponse.Analysis = opAnalysis
		f.State.Store(process_id, opResponse)
	}
}

func (f *FileProcessing) ProcessFilesFromDirectory(process_id string) map[string]interface{} {
	// Processing System
	// [x] Load text files from a specific folder
	// [x] Process documents in batches (batch processing)
	// [x] Extract basic statistics: word count, lines, characters
	// [x] Identify most frequent words
	// [x] Generate a content summary
	start_time := time.Now()
	progress := map[string]interface{}{"total_files": 0, "processed_files": 0, "percentage": 0}
	results := map[string]interface{}{"total_words": 0, "total_lines": 0, "most_frequent_words": []string{}, "files_processed": []string{}}
	filesProcessed := []string{}

	f.UpdateState(process_id, pd_interfaces.Running, progress, results, "", "", nil)

	FileManager := &FileManager{
		Path: f.Path,
		Log:  f.logger(),
	}
	files := FileManager.ListFilesFromPath()
	filesToProcess := []string{}
	for _, file := range files {
		filesToProcess = append(filesToProcess, file["name"].(string))
	}
	results["files_to_process"] = filesToProcess

	nfiles := len(files)
	if nfiles == 0 {
		f.UpdateState(process_id, pd_interfaces.Completed, progress, results, "No files to process", "", nil)
		return results
	}

	// Add nActive process to estimate time remaining
	nActiveProcesses := 0
	f.State.Range(func(key, value interface{}) bool {
		opResponse := value.(*pd_interfaces.OperationResponse)
		if opResponse.Status.Status == pd_interfaces.Running {
			nActiveProcesses++
		}
		return true
	})
	totalEstimatedTime := FileManager.EstimateProcessFiles(nActiveProcesses)
	progress["total_files"] = nfiles
	f.UpdateState(process_id, pd_interfaces.Running, progress, results, "", totalEstimatedTime, nil)

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
		batchAnalysis := resultsBatch["file_analysis"].([]map[string]interface{})
		batchResults := resultsBatch["results"].(map[string]interface{})

		// Update results and progress after processing each batch
		results["total_words"] = results["total_words"].(int) + batchResults["total_words"].(int)
		results["total_lines"] = results["total_lines"].(int) + batchResults["total_lines"].(int)
		results["most_frequent_words"] = append(results["most_frequent_words"].([]string), batchResults["most_frequent_words"].([]string)...)
		// Update files processed/to process and current progress
		filesProcessed = append(filesProcessed, batchResults["files_processed"].([]string)...)
		results["files_processed"] = filesProcessed
		progress["processed_files"] = len(filesProcessed)
		progress["percentage"] = (progress["processed_files"].(int) * 100) / progress["total_files"].(int)
		// Update filesToProcess to remove processed files
		remainingFiles := []string{}
		for _, file := range results["files_to_process"].([]string) {
			if !contains(filesProcessed, file) {
				remainingFiles = append(remainingFiles, file)
			}
		}
		results["files_to_process"] = remainingFiles

		// Update estimated completion time based on remaining files and processing speed
		estimatedCompletion := ""
		processedFiles := progress["processed_files"].(int)
		totalFiles := progress["total_files"].(int)
		if processedFiles > 0 && processedFiles < totalFiles {
			elapsed := time.Since(start_time)
			remainingFilesCount := totalFiles - processedFiles
			// Estimate remaining time based on elapsed time and remaining files.
			// We assume that the processing speed is consistent, which may not be the case in real scenarios, but this is a simple estimation.
			remainingDuration := time.Duration(float64(elapsed) * (float64(remainingFilesCount) / float64(processedFiles)))
			estimatedCompletion = formatRemainingDuration(remainingDuration)
		}

		f.UpdateState(process_id, pd_interfaces.Running, progress, results, "", estimatedCompletion, batchAnalysis)
	}

	progress["percentage"] = 100
	f.UpdateState(process_id, pd_interfaces.Completed, progress, results, "", "0s", nil)

	enlapsed_time := time.Since(start_time)
	f.logger().Info("Total processing time", zap.String("duration", enlapsed_time.String()), zap.String("process_id", process_id))
	return results
}

func (f *FileProcessing) ProcessBatchDocuments(process_id string, files []map[string]interface{}, nfiles int) map[string]interface{} {
	// Batch Document Processing
	// [x] Accept a batch of documents (e.g., multiple file paths or byte arrays)
	// [x] Process each document in the batch
	// [x] Update progress and results for the entire batch
	results := map[string]interface{}{"total_words": 0, "total_lines": 0, "most_frequent_words": []string{}, "files_processed": []string{}}

	type fileResult struct {
		idx            int
		fileName       string
		totalWords     int
		totalLines     int
		frequentWords  []string
		characterCount int
		summary        string
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
			fileName := file["name"].(string)
			f.logger().Info("Processing file", zap.String("file", fileName), zap.String("process_id", process_id))

			if val, ok := f.State.Load(process_id); ok {
				opResponse := val.(*pd_interfaces.OperationResponse)
				select {
				case <-opResponse.Stopper:
					f.logger().Info("Process stopped", zap.String("process_id", process_id))
					return
				default:
				}
			}

			// Open the file, read its content, count words, lines, and identify frequent words.
			var totalWords int = 0                                  // total word count for the file
			var totalLines int = 0                                  // total line count for the file
			var wordFrequency map[string]int = make(map[string]int) // map to track word frequencies
			var mostFrequentWords []string                          // slice to store the most frequent words
			var characterCount int = 0                              // total character count for the file

			file, err := os.Open(fileName)
			if err != nil {
				f.logger().Error("Error opening file", zap.String("file", fileName), zap.Error(err))
			}
			defer file.Close()

			// Use a scanner to read the file line by line, count words and lines, and track frequent words.
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				totalLines++
				characterCount += len(line)

				words := bufio.NewScanner(strings.NewReader(line))
				words.Split(bufio.ScanWords)
				for words.Scan() {
					word := words.Text()
					totalWords++
					wordFrequency[word]++
				}
			}

			if err := scanner.Err(); err != nil {
				f.logger().Error("Error reading file", zap.String("file", fileName), zap.Error(err))
			}

			// Identify the most frequent words (for simplicity, we take the top 5)
			type wordCount struct {
				word  string
				count int
			}
			var wordCounts []wordCount
			for word, count := range wordFrequency {
				wordCounts = append(wordCounts, wordCount{word: word, count: count})
			}
			sort.Slice(wordCounts, func(i, j int) bool {
				return wordCounts[i].count > wordCounts[j].count
			})
			for i := 0; i < len(wordCounts) && i < 5; i++ {
				// if the word is not in the mostFrequentWords slice, add it
				if !contains(mostFrequentWords, wordCounts[i].word) {
					mostFrequentWords = append(mostFrequentWords, wordCounts[i].word)
				}
			}

			// Generate a summary for the file content
			// Get the content of the file as a string (for simplicity, we read the entire file into memory, but this could be optimized for large files)
			// contentBytes, err := os.ReadFile(fileName)
			// if err != nil {
			// 	f.logger().Error("Error reading file content", zap.String("file", fileName), zap.Error(err))
			// }
			// content := string(contentBytes)
			// summary := f.GenerateSummary(content)
			summary := f.GenerateSummary("")

			resultsByFile[idx] = fileResult{
				idx:            idx,
				fileName:       fileName,
				totalWords:     totalWords,
				totalLines:     totalLines,
				frequentWords:  mostFrequentWords,
				characterCount: characterCount,
				summary:        summary,
			}

			elapsedTime := time.Since(startTime).String()
			f.logger().Info("Completed processing file", zap.String("file", fileName), zap.String("duration", elapsedTime), zap.String("process_id", process_id))
		}()
	}

	wg.Wait() // Wait for all files in the batch to be processed

	// Aggregate results from all files in the batch
	fileAnalysis := []map[string]interface{}{}
	for _, fileResult := range resultsByFile {
		results["total_words"] = results["total_words"].(int) + fileResult.totalWords
		results["total_lines"] = results["total_lines"].(int) + fileResult.totalLines
		results["most_frequent_words"] = append(results["most_frequent_words"].([]string), fileResult.frequentWords...)
		results["files_processed"] = append(results["files_processed"].([]string), fileResult.fileName)

		fileAnalysis = append(fileAnalysis, map[string]interface{}{
			"file_name":           fileResult.fileName,
			"total_words":         fileResult.totalWords,
			"total_lines":         fileResult.totalLines,
			"most_frequent_words": fileResult.frequentWords,
			"total_characters":    fileResult.characterCount,
			"summary":             fileResult.summary,
		})
	}

	return map[string]interface{}{
		"results":       results,
		"file_analysis": fileAnalysis,
	}
}

// TODO: process the file and generate a summary using an NLP model or similar. For now, this is a placeholder that returns a static summary.
func (f *FileProcessing) GenerateSummary(content string) string {
	// Placeholder for summary generation logic.
	// In a real implementation, this could use an NLP model or similar to generate a summary of the content.
	context := context.Background()

	if content == "" {
		time.Sleep(5 * time.Second) // Simulate processing time for the file
		return "No content to summarize"
	}

	summaryContent, err := summarizerDomainLib.SummarizeContent(context, f.McpClient, string(content), 100)
	if err != nil {
		f.Log.Error("Failed to summarize content", zap.Error(err))
		return "Failed to generate summary"
	}

	return summaryContent
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func formatRemainingDuration(d time.Duration) string {
	if d <= 0 {
		return "0s"
	}

	seconds := int(d.Seconds())
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, secs)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, secs)
	}
	return fmt.Sprintf("%ds", secs)
}
