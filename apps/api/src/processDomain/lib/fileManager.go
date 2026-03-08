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
	Path string
}

func (f *FileManager) ListFilesFromPath() []string {
	// List Files from a Path
	// [ ] Read files from a specified directory
	// [ ] Filter files based on type (e.g., .txt)
	// [ ] Return a list of file paths
	files := []string{}

	// go to the specified directory and list all files
	dirEntries, err := os.ReadDir(f.Path)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return files
	}

	for _, entry := range dirEntries {
		// filter only .txt files
		if !entry.IsDir() && entry.Name()[len(entry.Name())-4:] == ".txt" {
			files = append(files, f.Path+"/"+entry.Name())
		}
	}

	return files
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
