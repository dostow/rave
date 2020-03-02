package plugin

var addonconfig = `{
	"name": "rave",
	"title": "Rave Addon",
	"description": "An addon for communicating with flutterwave rave api",
	"properties": {
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
	"title": "Rave params",
	"properties": {
		"callback": {
			"type": "string"
		},
		"action": {
			"type": "string",
			"enum": ["createTransferRecipient", "createTransfer", "validateTransfer"]
		},
		"options": {
			"type": "object",
			"oneOf": [
				{ 
					"type": "object", 
					"properties": {
						"reference": {"type": "string"}
					},
					"required": ["reference"]
				},
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
						"bankLocation": {"type": "string"},
						"meta": {"type": "object"}
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
