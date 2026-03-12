package interfaces

import "testing"

func TestOperationResponseInitialize(t *testing.T) {
	op := &OperationResponse{}
	op.Initialize("proc-1")

	if op.ID != "proc-1" {
		t.Fatalf("expected ID proc-1, got %s", op.ID)
	}
	if op.Status.ID != "proc-1" {
		t.Fatalf("expected status ID proc-1, got %s", op.Status.ID)
	}
	if op.Analysis.ID != "proc-1" {
		t.Fatalf("expected analysis ID proc-1, got %s", op.Analysis.ID)
	}
	if op.Stopper == nil {
		t.Fatal("expected stopper channel to be initialized")
	}
}

func TestOperationListOrderProcesses(t *testing.T) {
	list := &OperationListResponse{}
	list.Initialize()

	list.AddProcess(OperationReview{ID: "1", StartedAt: "2024-01-01T10:00:00Z"})
	list.AddProcess(OperationReview{ID: "3", StartedAt: "2024-01-02T10:00:00Z"})
	list.AddProcess(OperationReview{ID: "2", StartedAt: "2024-01-02T10:00:00Z"})

	list.OrderProcesses()

	if list.Processes[0].ID != "3" || list.Processes[1].ID != "2" || list.Processes[2].ID != "1" {
		t.Fatalf("unexpected order: %#v", list.Processes)
	}
}

func TestOperationListGetFirstTenProcesses(t *testing.T) {
	list := &OperationListResponse{}
	list.Initialize()

	for i := 0; i < 12; i++ {
		list.AddProcess(OperationReview{ID: "p"})
	}

	firstTen := list.GetFirstTenProcesses()
	if len(firstTen) != 10 {
		t.Fatalf("expected 10 processes, got %d", len(firstTen))
	}
}
