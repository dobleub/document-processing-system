package interfaces

import "testing"

func TestAppendBatchAnalysisMapsFields(t *testing.T) {
	analysis := &OperationAnalysisResult{}
	analysis.Initialize("proc-1")

	analysis.AppendBatchAnalysis([]map[string]interface{}{
		{
			"file_name":           "/tmp/docs/book.txt",
			"total_words":         42,
			"total_lines":         3,
			"most_frequent_words": []string{"a", "b"},
			"total_characters":    128,
			"summary":             "summary",
		},
	})

	if len(analysis.Analysis) != 1 {
		t.Fatalf("expected one analysis entry, got %d", len(analysis.Analysis))
	}
	entry := analysis.Analysis[0]
	if entry.FileName != "book.txt" {
		t.Fatalf("expected base filename book.txt, got %s", entry.FileName)
	}
	if entry.TotalWords != 42 || entry.TotalLines != 3 || entry.TotalCharacters != 128 {
		t.Fatalf("unexpected numeric fields: %#v", entry)
	}
	if len(entry.MostFrequentWords) != 2 || entry.Summary != "summary" {
		t.Fatalf("unexpected analysis payload: %#v", entry)
	}
}

func TestOperationAnalysisMarkTransitions(t *testing.T) {
	analysis := &OperationAnalysisResult{}
	analysis.Initialize("proc-2")

	analysis.MarkAsCompleted()
	if analysis.Status != Completed {
		t.Fatalf("expected %s, got %s", Completed, analysis.Status)
	}

	analysis.MarkAsFailed("err")
	if analysis.Status != Failed {
		t.Fatalf("expected %s, got %s", Failed, analysis.Status)
	}

	analysis.MarkAsStopped()
	if analysis.Status != Stopped {
		t.Fatalf("expected %s, got %s", Stopped, analysis.Status)
	}
}
