package helpers

import "github.com/graphql-go/graphql"

func GetFieldType(v interface{}) graphql.Output {
	switch v.(type) {
	case int:
		return graphql.Int
	case float64:
		return graphql.Float
	case string:
		return graphql.String
	case bool:
		return graphql.Boolean
	default:
		return graphql.String
	}
}

func CreateFieldTypeFromStruct(name string, s map[string]any) graphql.Fields {
	fields := graphql.Fields{}
	// Convert s to a map type
	for k, v := range s {
		// if k is _id, convert it to a graphql type
		if k == "_id" {
			fields[k] = &graphql.Field{
				Type: graphql.ID,
			}
			continue
		}
		// if v is a map, convert it to a graphql type
		if vMap, ok := v.(map[string]interface{}); ok {
			fields[k] = &graphql.Field{
				Type: CreateTypeFromStruct(k+name+"Type", vMap),
			}
			continue
		}
		// Convert v to a graphql type
		fields[k] = &graphql.Field{
			Type: GetFieldType(v),
		}
	}
	return fields
}

func CreateTypeFromStruct(name string, s map[string]interface{}) *graphql.Object {
	name = capitalize(name)
	fields := CreateFieldTypeFromStruct(name, s)

	return graphql.NewObject(graphql.ObjectConfig{
		Name:        name,
		Fields:      fields,
		Description: name + " type",
	})
}

func CreateInputFieldTypeFromStruct(name string, s map[string]interface{}) graphql.InputObjectConfigFieldMap {
	fields := graphql.InputObjectConfigFieldMap{}
	// Convert s to a map type
	for k, v := range s {
		// if k is _id, convert it to a graphql type
		if k == "_id" {
			fields[k] = &graphql.InputObjectFieldConfig{
				Type: graphql.ID,
			}
			continue
		}
		// if v is a map, convert it to a graphql type
		if vMap, ok := v.(map[string]interface{}); ok {
			fields[k] = &graphql.InputObjectFieldConfig{
				Type: CreateInputTypeFromStruct(k+name+"InputType", vMap),
			}
			continue
		}
		// Convert v to a graphql type
		fields[k] = &graphql.InputObjectFieldConfig{
			Type: GetFieldType(v),
		}
	}
	return fields
}

func CreateInputTypeFromStruct(name string, s map[string]interface{}) *graphql.InputObject {
	name = capitalize(name)
	fields := CreateInputFieldTypeFromStruct(name, s)

	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        name,
		Fields:      fields,
		Description: name + " input type",
	})
}
