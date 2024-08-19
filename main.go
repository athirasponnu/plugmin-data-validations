package main

import (
	"net/http"
	"reflect"

	"schema_validations/entities"
	"schema_validations/utilities"

	"github.com/gin-gonic/gin"
)

func handlePost(c *gin.Context) {
	// Define schema, errors, and payload
	var (
		errors  = make(map[string]string)
		payload string
	)

	// Get table ID from query parameters
	tableID := c.GetHeader("table_id")
	if tableID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
		return
	}

	// Derive table name from table ID
	tableName, schema, err := utilities.TableDetails(tableID)
	switch {
	case err != nil:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case tableName == "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
		return
	case reflect.DeepEqual(schema, entities.JSONSchema{}):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schema not found"})
		return
	}

	// Get raw data from the request
	byteData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Convert raw data to payload
	payload = string(byteData)

	// Define payload configurations
	payloadConfigurations := map[string]string{
		"count": "data.count",
		"email": "data.email",
		"name":  "data.name",
		"id":    "data.id",
	}

	// Extract data from payload
	data := utilities.ExtractData(payload, payloadConfigurations)

	// Validate if all fields in data are available in schema properties
	for key := range data {
		if _, exists := schema.Properties[key]; !exists {
			errors[key] = "Field not defined in schema"
		}
	}

	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
		return
	}

	// Validate the data against the schema
	errors = utilities.ValidateData(data, &schema)
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
		return
	}
	tableDetails := map[string]entities.Record{
		tableName: {
			ColumnVals: data,
		},
	}
	c.JSON(http.StatusOK, tableDetails)
}
func handlePatch(c *gin.Context) {
	// Define schema, errors, and payload
	var (
		schema         entities.JSONSchema
		errors         = make(map[string]string)
		payload        string
		identityValMap = make(map[string]interface{})
	)

	// Get table ID from query parameters
	tableID := c.GetHeader("table_id")
	if tableID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
		return
	}

	// Derive table name from table ID
	tableName, schema, err := utilities.TableDetails(tableID)
	switch {
	case err != nil:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case tableName == "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
		return
	case reflect.DeepEqual(schema, entities.JSONSchema{}):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schema not found"})
		return
	}
	userId := c.Param("id")

	// Get raw data from the request
	byteData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Convert raw data to payload
	payload = string(byteData)

	// Define payload configurations
	payloadConfigurations := map[string]string{
		"count": "data.count",
		"email": "data.email",
		"name":  "data.name",
		"id":    "data.id",
	}

	// Extract data from payload
	data := utilities.ExtractData(payload, payloadConfigurations)

	// Validate if all fields in data are available in schema properties
	for key := range data {
		if _, exists := schema.Properties[key]; !exists {
			errors[key] = "Field not defined in schema"
		}
	}

	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
		return
	}

	// Validate the data against the schema
	errors = utilities.ValidateData(data, &schema)
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
		return
	}

	identityValMap["id"] = userId
	tableDetails := map[string]entities.Record{
		tableName: {
			ColumnVals: data,
			IdenityVal: identityValMap,
		},
	}
	c.JSON(http.StatusOK, tableDetails)
}
func handleGet(c *gin.Context) {
	// Define schema, errors, and payload
	var (
		schema entities.JSONSchema
		params entities.Params
	)

	// Get table ID from query parameters
	tableID := c.GetHeader("table_id")
	if tableID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
		return
	}

	err := c.BindQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Derive table name from table ID
	tableName, schema, err := utilities.TableDetails(tableID)
	switch {
	case err != nil:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case tableName == "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
		return
	case reflect.DeepEqual(schema, entities.JSONSchema{}):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schema not found"})
		return
	}

	// Extract data from payload
	data := utilities.GetColumn(schema)
	result := &entities.RecordWithParams{
		ColumnVals: data,
		Params:     params,
	}

	// Respond with JSON data
	c.JSON(http.StatusOK, result)
}
func main() {
	r := gin.Default()

	// Define the POST endpoints for different resources
	r.POST("/users", handlePost)
	r.PATCH("/users/:id", handlePatch)
	r.GET("/users", handleGet)

	// Start the Gin server
	r.Run(":8080")
}
