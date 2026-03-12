package helpers

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestGetFilterOperation(t *testing.T) {
	q := &Query{}

	cases := map[string]string{
		"":       "$eq",
		"ne":     "$ne",
		"gt":     "$gt",
		"gte":    "$gte",
		"lt":     "$lt",
		"lte":    "$lte",
		"in":     "$in",
		"nin":    "$nin",
		"exists": "$exists",
		"regex":  "$regex",
	}

	for input, expected := range cases {
		got := q.GetFilterOperation(input)
		if got != expected {
			t.Fatalf("operation %q expected %q, got %q", input, expected, got)
		}
	}
}

func TestGetBoolOperation(t *testing.T) {
	if (&Query{Operation: "or"}).GetBoolOperation() != "$or" {
		t.Fatal("expected 'or' operation to map to $or")
	}

	if (&Query{Operation: "and"}).GetBoolOperation() != "$and" {
		t.Fatal("expected non-'or' operation to map to $and")
	}
}

func TestGetFilterPipeline(t *testing.T) {
	q := &Query{
		Operation: "and",
		Fields: []QueryField{
			{Operation: "gt", Field: "pages", Value: "10", Type: "integer"},
			{Operation: "regex", Field: "title", Value: "report", Type: "string"},
			{Operation: "in", Field: "status", Value: "[DONE RUNNING]", Type: "string"},
			{Operation: "exists", Field: "summary", Value: "ignored", Type: "string"},
		},
	}

	got := q.GetFilterPipeline()
	expected := bson.M{
		"$and": []bson.M{
			{"pages": bson.M{"$gt": 10}},
			{"title": bson.M{"$regex": "report", "$options": "i"}},
			{"status": bson.M{"$in": []string{"DONE", "RUNNING"}}},
			{"summary": bson.M{"$exists": true}},
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("pipeline mismatch\nexpected: %#v\ngot: %#v", expected, got)
	}
}
