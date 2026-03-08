package lib

import (
	"fmt"
	"os"
)

type FileManager struct {
	// File Manager
	// [ ] Handle file operations (e.g., reading, writing)
	// [ ] Manage file paths and directories
	// [ ] Ensure efficient memory usage when processing large files
	Path  string
	Files []map[string]interface{}
}

func (f *FileManager) ListFilesFromPath() []map[string]interface{} {
	// List Files from a Path
	// [ ] Read files from a specified directory
	// [ ] Filter files based on type (e.g., .txt)
	// [ ] Return a list of file paths
	files := []map[string]interface{}{}

	// go to the specified directory and list all files
	dirEntries, err := os.ReadDir(f.Path)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return files
	}

	for _, entry := range dirEntries {
		// filter only .txt files
		file_info, _ := entry.Info()
		if !entry.IsDir() && file_info.Size() > 0 && entry.Name()[len(entry.Name())-4:] == ".txt" {
			// files = append(files, f.Path+"/"+entry.Name())
			files = append(files, map[string]interface{}{
				"name": f.Path + "/" + entry.Name(),
				"info": file_info,
			})
		}
	}

	f.Files = files
	return files
}

func (f *FileManager) EstimateProcessFiles() string {
	// Estimate Process File
	// [ ] Estimate the time required to process a file based on its size
	totalEstimateTime := int64(0)

	if len(f.Files) > 0 {
		for _, file := range f.Files {
			fileInfo := file["info"].(os.FileInfo)
			size := fileInfo.Size()
			estimatedTime := size / (1024 * 1024) // Size in MB
			totalEstimateTime += estimatedTime
		}
	}

	return fmt.Sprintf("%d seconds", totalEstimateTime)
}

func (f *FileManager) ReadFileContent(filePath string) string {
	// Read File Content
	// [ ] Open and read the content of a file
	// [ ] Handle potential errors (e.g., file not found)
	// [ ] Return the content as a string
	return ""
}

func (f *FileManager) GenerateSummary(content string) string {
	// Generate Summary
	// [ ] Analyze the content of the file
	// [ ] Extract key points and main ideas
	// [ ] Create a concise summary of the content
	return ""
}

func (f *FileManager) IdentifyFrequentWords(content string) []string {
	// Identify Frequent Words
	// [ ] Count the frequency of each word in the content
	// [ ] Sort words by frequency
	// [ ] Return a list of the most frequent words
	return []string{}
}

func (f *FileManager) ExtractStatistics(content string) map[string]int {
	// Extract Statistics
	// [ ] Count the number of words, lines, and characters in the content
	// [ ] Return a map with the statistics
	return map[string]int{}
}
