// Package helpers provides some utils to handle especific escenarois
package helpers

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueryField struct {
	Operation string `json:"operation" bson:"operation"`
	Field     string `json:"field" bson:"field"`
	Value     string `json:"value" bson:"value"`
	Type      string `json:"type" bson:"type"`
}

type Query struct {
	Operation string       `json:"operation" bson:"operation"`
	Fields    []QueryField `json:"fields" bson:"fields"`
}

func NewQuery(st map[string]interface{}) *Query {
	q := &Query{}
	q.Operation = st["operation"].(string)
	for _, v := range st["fields"].([]interface{}) {
		tmp := v.(map[string]interface{})
		q.Fields = append(q.Fields, QueryField{
			Operation: tmp["operation"].(string),
			Field:     tmp["field"].(string),
			Value:     tmp["value"].(string),
			Type:      tmp["type"].(string),
		})
	}
	return q
}

func (q *Query) GetFilterOperation(operation string) string {
	tmpOp := "$eq"
	if operation == "ne" {
		tmpOp = "$ne"
	}
	if operation == "gt" {
		tmpOp = "$gt"
	}
	if operation == "gte" {
		tmpOp = "$gte"
	}
	if operation == "lt" {
		tmpOp = "$lt"
	}
	if operation == "lte" {
		tmpOp = "$lte"
	}
	if operation == "in" {
		tmpOp = "$in"
	}
	if operation == "nin" {
		tmpOp = "$nin"
	}
	if operation == "exists" {
		tmpOp = "$exists"
	}
	if operation == "regex" {
		tmpOp = "$regex"
	}

	return tmpOp
}

func (q *Query) GetBoolOperation() string {
	if q.Operation == "or" {
		return "$or"
	}
	return "$and"
}

func (q *Query) GetFilterPipeline() bson.M {
	pipeline := bson.M{}
	filedsPipeline := []bson.M{}
	boolOperation := q.GetBoolOperation()

	for _, v := range q.Fields {
		tmpValue := v.Value
		tmpOperation := q.GetFilterOperation(v.Operation)
		tmpPipe := bson.M{v.Field: bson.M{tmpOperation: tmpValue}}

		if tmpOperation == "$regex" {
			tmpPipe = bson.M{v.Field: bson.M{tmpOperation: tmpValue, "$options": "i"}}
		}
		if tmpOperation == "$exists" {
			tmpPipe = bson.M{v.Field: bson.M{tmpOperation: true}}
		}
		if tmpOperation == "$nin" || tmpOperation == "$in" {
			tmpValue := StringToArray(tmpValue)
			tmpPipe = bson.M{v.Field: bson.M{tmpOperation: tmpValue}}
		}

		if v.Type == "integer" {
			intValue, _ := strconv.Atoi(v.Value)
			tmpPipe = bson.M{v.Field: bson.M{tmpOperation: intValue}}
		}
		if v.Type == "float" {
			floatValue, _ := strconv.ParseFloat(v.Value, 64)
			tmpPipe = bson.M{v.Field: bson.M{tmpOperation: floatValue}}
		}
		if v.Type == "boolean" {
			boolValue, _ := strconv.ParseBool(v.Value)
			tmpPipe = bson.M{v.Field: bson.M{tmpOperation: boolValue}}
		}
		if v.Type == "date" {
			dateValue := v.Value
			tmpPipe = bson.M{v.Field: bson.M{tmpOperation: dateValue}}
		}
		if v.Type == "objectid" {
			tmpID, _ := primitive.ObjectIDFromHex(v.Value)
			tmpPipe = bson.M{v.Field: tmpID}
		}

		filedsPipeline = append(filedsPipeline, tmpPipe)
	}
	pipeline[boolOperation] = filedsPipeline

	return pipeline
}
