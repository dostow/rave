package plugin

var addonconfig = `{
	"name": "rave",
	"title": "Rave Addon",
	"description": "A
		"apiKey": {"type": "string"},
		"keys": {
			"type": "object",
			"properties": {
				"secret": {
					"type": "string",
					"description": "Secret key"
				},
				"public": {
					"type": "string",
					"description": "Public key"
				}
			},
			"required": ["secret", "public"],
			"additionalProperties": false
		}
	},
	"required": ["keys"],
	"type": "object",
	"additionalProperties": false
}`

// TODO: add jsonschema template format
var linkparams = `{
	"name": "params",
	"title": "Mail Params",
	"properties": {
		"callback": {
			"type": "string"
		},
		"action": {
			"type": "string",
			"enum": ["createTransferRecipient", "createTransfer"]
		},
		"options": {
			"type": "object",
			"oneOf": [
				{ 
					"type": "object", 
					"properties": {
						"accountNumber": {"type": "string"},
						"bankCode": {
							"type": "string"
						}
					},
					"required": ["accountNumber", "bankCode"]
				},
				{ 
					"type": "object", 
					"properties": {
						"amount": {"type": "string"},
						"recipient": {"type": "string"},
						"reference": {"type": "string"},
						"currency": {"type": "string"},
						"narration": {"type": "string"},
						"bankLocation": {"type": "string"}
					},
					"required": ["amount", "recipient", "reference", "currency", "narration", "bankLocation"]
				}
			]
		}
	},
	"required": ["callback", "action", "options"],
	"type": "object",
	"additionalProperties": false
}`

// AddonRegistrar an addon registrar
type AddonRegistrar interface {
	Add(name, config, params string)
}

// Register injects an addon into a registry
func Register(ar AddonRegistrar) {
	ar.Add("rave", addonconfig, linkparams)
}
