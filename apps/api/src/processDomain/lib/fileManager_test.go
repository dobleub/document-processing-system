package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestListFilesFromPathFiltersTxtAndEmptyFiles(t *testing.T) {
	dir := t.TempDir()

	txtPath := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(txtPath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to write txt file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "empty.txt"), []byte{}, 0o644); err != nil {
		t.Fatalf("failed to write empty file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("markdown"), 0o644); err != nil {
		t.Fatalf("failed to write md file: %v", err)
	}

	manager := &FileManager{Path: dir}
	files := manager.ListFilesFromPath()

	if len(files) != 1 {
		t.Fatalf("expected exactly one valid txt file, got %d", len(files))
	}
	if files[0]["name"].(string) != txtPath {
		t.Fatalf("expected path %s, got %s", txtPath, files[0]["name"].(string))
	}
}

func TestEstimateProcessFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "big.txt")
	content := make([]byte, 200000)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	manager := &FileManager{
		Files: []map[string]interface{}{{"name": path, "info": info}},
	}

	base := float64(info.Size()) * 10 / (1024 * 1024)
	withOneActive := base * 1.05
	expected := fmt.Sprintf("%.2f seconds", withOneActive)

	got := manager.EstimateProcessFiles(1)
	if got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestEstimateProcessFilesWithoutFiles(t *testing.T) {
	manager := &FileManager{}
	got := manager.EstimateProcessFiles(0)
	if got != "0.00 seconds" {
		t.Fatalf("expected 0.00 seconds, got %s", got)
	}
}
