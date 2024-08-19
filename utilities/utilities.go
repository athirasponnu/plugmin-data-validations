package utilities

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"schema_validations/entities"

	"github.com/tidwall/gjson"
)

type ValidationFunc func(fieldName string, value interface{}, ruleValue interface{}) string

var validationFuncs = map[string]ValidationFunc{
	"required": RequiredValidation,
	"max":      MaxValidation,
	"min":      MinValidation,
	"pattern":  PatternValidation,
}

func RequiredValidation(fieldName string, value interface{}, ruleValue interface{}) string {
	if ruleValue.(bool) && (value == nil || reflect.ValueOf(value).IsZero()) {
		return fmt.Sprintf("field '%s' is required", fieldName)
	}
	return ""
}

// Max validation
func MaxValidation(fieldName string, value interface{}, ruleValue interface{}) string {
	val := reflect.ValueOf(value)
	maxVal := reflect.ValueOf(ruleValue)

	switch val.Kind() {
	case reflect.String:
		if val.Len() > int(maxVal.Float()) {
			return fmt.Sprintf("field '%s' exceeds maximum length of %v", fieldName, ruleValue)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() > int64(maxVal.Float()) {
			return fmt.Sprintf("field '%s' exceeds maximum value of %v", fieldName, ruleValue)
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() > maxVal.Float() {
			return fmt.Sprintf("field '%s' exceeds maximum value of %v", fieldName, ruleValue)
		}
	case reflect.Slice, reflect.Array:
		if val.Len() > int(maxVal.Float()) {
			return fmt.Sprintf("field '%s' exceeds maximum array length of %v", fieldName, ruleValue)
		}
	default:
		return fmt.Sprintf("field '%s' has unsupported type for 'max' validation", fieldName)
	}

	return ""
}

// Min validation
func MinValidation(fieldName string, value interface{}, ruleValue interface{}) string {
	val := reflect.ValueOf(value)
	minVal := reflect.ValueOf(ruleValue)

	switch val.Kind() {
	case reflect.String:
		if val.Len() < int(minVal.Float()) {
			return fmt.Sprintf("field '%s' is less than minimum length of %v", fieldName, ruleValue)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() < int64(minVal.Float()) {
			return fmt.Sprintf("field '%s' is less than minimum value of %v", fieldName, ruleValue)
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() < minVal.Float() {
			return fmt.Sprintf("field '%s' is less than minimum value of %v", fieldName, ruleValue)
		}
	case reflect.Slice, reflect.Array:
		if val.Len() < int(minVal.Float()) {
			return fmt.Sprintf("field '%s' is less than minimum array length of %v", fieldName, ruleValue)
		}
	default:
		return fmt.Sprintf("field '%s' has unsupported type for 'min' validation", fieldName)
	}

	return ""
}

// Pattern validation
func PatternValidation(fieldName string, value interface{}, ruleValue interface{}) string {
	if str, ok := value.(string); ok {
		re := regexp.MustCompile(ruleValue.(string))
		if !re.MatchString(str) {
			return fmt.Sprintf("field '%s' does not match the required pattern", fieldName)
		}
	}
	return ""
}

// Method to apply validations dynamically based on the schema
func ValidateData(data map[string]interface{}, schema *entities.JSONSchema) map[string]string {
	errors := make(map[string]string)
	for fieldName, fieldSchema := range schema.Properties {

		value, exists := data[fieldName]

		if exists {
			// Apply validations for the current field
			errMsg := ApplyValidations(fieldName, value, fieldSchema.Validations)
			if errMsg != "" {
				errors[fieldName] = errMsg
			}
		}
	}

	return errors
}

// Function to dynamically apply validation rules based on the validation functions map
func ApplyValidations(fieldName string, value interface{}, validations map[string]interface{}) string {
	var errors []string

	for validationType, ruleValue := range validations {
		validationFunc, exists := validationFuncs[validationType]
		if !exists {
			errors = append(errors, fmt.Sprintf("validation '%s' not supported", validationType))
			continue
		}

		if errMsg := validationFunc(fieldName, value, ruleValue); errMsg != "" {
			errors = append(errors, errMsg)
		}
	}

	if len(errors) > 0 {
		return fmt.Sprintf("%s", errors)
	}

	return ""
}

func ExtractData(payload string, payloadConfigurations map[string]string) map[string]interface{} {
	data := make(map[string]interface{})

	switch {
	case len(payloadConfigurations) > 0:
		for column, path := range payloadConfigurations {
			data[column] = gjson.Get(payload, path).Value()
		}
	default:
		if err := json.Unmarshal([]byte(payload), &data); err != nil {
			log.Fatalf("Error unmarshaling JSON: %v", err)
		}
	}

	return data
}

func GetColumn(schema entities.JSONSchema) map[string]interface{} {
	data := make(map[string]interface{})
	for fieldName := range schema.Properties {
		data[fieldName] = ""
	}

	return data
}
func TableDetails(tableID string) (string, entities.JSONSchema, error) {
	var schema entities.JSONSchema
	// This is a mock implementation.
	tableName := "users"
	// Unmarshal the schema JSON into a Go struct
	err := json.Unmarshal([]byte(entities.SchemaJSON), &schema)
	if err != nil {
		return "", entities.JSONSchema{}, err
	}
	return tableName, schema, nil
}
