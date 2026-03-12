package interfaces

import "testing"

func TestOperationStatusInitializeDefaults(t *testing.T) {
	status := &OperationStatus{}
	status.Initialize("proc-1")

	if status.ID != "proc-1" {
		t.Fatalf("expected ID proc-1, got %s", status.ID)
	}
	if status.Status != Pending {
		t.Fatalf("expected status %s, got %s", Pending, status.Status)
	}
	if status.StartedAt == "" {
		t.Fatal("expected StartedAt to be set")
	}
	if len(status.Result.FilesProcessed) != 0 || len(status.Result.FilesToProcess) != 0 {
		t.Fatal("expected empty files slices after initialization")
	}
}

func TestOperationStatusUpdateFromMaps(t *testing.T) {
	status := &OperationStatus{}
	status.Initialize("proc-2")

	status.UpdateOperationStatus(map[string]interface{}{
		"status": Running,
		"progress": map[string]interface{}{
			"total_files":     4,
			"processed_files": 2,
			"percentage":      50,
		},
		"results": map[string]interface{}{
			"total_words":         120,
			"total_lines":         10,
			"most_frequent_words": []string{"go", "test"},
			"files_processed":     []string{"a.txt"},
			"files_to_process":    []string{"b.txt"},
		},
		"estimated_completion": "10s",
	})

	if status.Status != Running {
		t.Fatalf("expected status %s, got %s", Running, status.Status)
	}
	if status.Progress.TotalFiles != 4 || status.Progress.ProcessedFiles != 2 || status.Progress.Percentage != 50 {
		t.Fatalf("unexpected progress: %#v", status.Progress)
	}
	if status.Result.TotalWords != 120 || status.Result.TotalLines != 10 {
		t.Fatalf("unexpected result counters: %#v", status.Result)
	}
	if len(status.Result.MostFrequentWords) != 2 || status.Result.MostFrequentWords[0] != "go" {
		t.Fatalf("unexpected frequent words: %#v", status.Result.MostFrequentWords)
	}
	if len(status.Result.FilesProcessed) != 1 || status.Result.FilesProcessed[0] != "a.txt" {
		t.Fatalf("unexpected files processed: %#v", status.Result.FilesProcessed)
	}
	if status.EstimatedCompletion != "10s" {
		t.Fatalf("expected estimated completion 10s, got %s", status.EstimatedCompletion)
	}
}

func TestOperationStatusMarkTransitions(t *testing.T) {
	status := &OperationStatus{}
	status.Initialize("proc-3")

	status.MarkAsCompleted()
	if status.Status != Completed || status.Progress.Percentage != 100 || status.CompletedAt == "" {
		t.Fatalf("unexpected completed state: %#v", status)
	}

	status.MarkAsFailed("boom")
	if status.Status != Failed || status.Error != "boom" || status.CompletedAt == "" {
		t.Fatalf("unexpected failed state: %#v", status)
	}

	status.MarkAsStopped()
	if status.Status != Stopped || status.CompletedAt == "" {
		t.Fatalf("unexpected stopped state: %#v", status)
	}
}
