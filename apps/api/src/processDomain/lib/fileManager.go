package lib

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

type FileManager struct {
	// File Manager
	// [ ] Handle file operations (e.g., reading, writing)
	// [ ] Manage file paths and directories
	// [ ] Ensure efficient memory usage when processing large files
	Path  string
	Files []map[string]interface{}
	Log   *zap.Logger
}

func (f *FileManager) logger() *zap.Logger {
	if f.Log != nil {
		return f.Log
	}
	return zap.NewNop()
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
		f.logger().Error("Error reading directory", zap.String("path", f.Path), zap.Error(err))
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
