package entities

// Define a struct for field properties
type FieldSchema struct {
	Type        string                 `json:"type"`
	Label       string                 `json:"label"`
	Validations map[string]interface{} `json:"validations"`
}

type Params struct {
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	Order string `form:"order"`
	Sort  string `form:"sort"`
}

// Define the JSONSchema struct
type JSONSchema struct {
	Properties map[string]FieldSchema `json:"properties"`
}

// Define validation function signatures

type ColumnVals map[string]interface{}

type Record struct {
	ColumnVals map[string]interface{} `json:"columnVals"`
	IdenityVal map[string]interface{} `json:"identitiyVal,omitempty"`
}
type RecordWithParams struct {
	ColumnVals map[string]interface{} `json:"columnVals"`
	IdenityVal map[string]interface{} `json:"identitiyVal,omitempty"`
	Params     Params                 `json:"params,omitempty"`
}
