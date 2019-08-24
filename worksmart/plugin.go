package worksmart

import (
	"git.progwebtech.com/code/worksmart/addon"
	"git.progwebtech.com/code/worksmart/addon/addons"
)

var addonconfig = `{
	"name": "rave",
	"title": "Rave Addon",
	"description": "An addon for communicating with flutterwave rave api",
	"properties": {
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
	"type": "object",
	"additionalProperties": false
}`

// TODO: add jsonschema template format
var linkparams = `{
	"name": "params",
	"title": "Mail Params",
	"properties": {
		"callback": {
			"type": "object",
			"properties": {
				"store": {"type": "string"}
			}
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
						"bankCode": {"type": "string"}
					}
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
	"required": [],
	"type": "object",
	"additionalProperties": false
}`

func init() {
	addon.Register(addons.NewType("rave", linkparams, addonconfig))
}
