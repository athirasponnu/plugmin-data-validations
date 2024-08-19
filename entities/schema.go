package entities

var SchemaJSON = `{
	"$id": "https://example.com/product.schema.json",
  	"$schema": "https://json-schema.org/draft/2020-12/schema",
  	"description": "user details",
	"properties": {
			"count": {
				"type": "integer",
				"label":"count",
				"validations": {
					"required": true,
					"max":10				
				}
			},
			"email": {
				"type": "string",
				"label":"email",
				"validations": {
					"pattern": "^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$",
					"required": true
				}
			},
			"name": {
				"type": "string",
				"label":"name",
				"validations": {
					"max": 7,
					"min":3	
				}
			},
			"id": {
				"type": "string",
				"label":"id",
				"validations": {
					"required": true
				}
			}
		},
		 "title": "member",
  		 "type": "object",
		"required": [
			"email",
			"count"
		]
	}`
