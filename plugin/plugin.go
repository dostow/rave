package plugin

var addonconfig = `{
	"name": "rave",
	"title": "Rave Addon",
	"description": "An addon for communicating with a payment api",
	"properties": {
		"platform": {
			"type": "string",
			"description": "Platform",
			"enum": ["paystack", "rave", "quikk", "mpesa"]
		},
		"callback": {
			"properties": { 
				"url": {
					"type": "string",
					"description": "callback url for response"
				},
				"method": {
					"type": "string",
					"enum": ["POST", "PUT"]
				},
				"headers": {
					"type": "object",
					"description": "headers for request"
				}
			},
			"type": "object"
		},
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
		},
		"paystack": {
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
		},
		"mpesa": {
			"type": "object",
			"properties": {
				"secret": {
					"type": "string",
					"description": "Secret key"
				},
				"public": {
					"type": "string",
					"description": "Public key"
				},
				"passkey": {
					"type": "string",
					"description": "Pass Key"
				},
				"shortcode": {
					"type": "string",
					"description": "Business Short Code"
				}
			},
			"required": ["secret", "public"],
			"additionalProperties": false
		},
		"quikk": {
			"type": "object",
			"properties": {
				"secret": {
					"type": "string",
					"description": "Secret key"
				},
				"public": {
					"type": "string",
					"description": "Public key"
				},
				"passkey": {
					"type": "string",
					"description": "Pass Key"
				},
				"shortcode": {
					"type": "string",
					"description": "Business Short Code"
				}
			},
			"required": ["secret", "public", "shortcode"],
			"additionalProperties": false
		}
	},
	"required": ["platform"],
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
			"enum": ["createTransactionLink", "createTransferRecipient", "createTransfer", "deleteTransferRecipient", "validateTransaction", "validateTransfer"]
		},
		"options": {
			"type": "object",
			"oneOf": [
				{
					"type": "object",
					"name": "Create Transaction Link",
					"description": "create transaction link for making payments",
					"properties": {
						"tx_ref": {
							"type": "string"
						},
						"amount": {
							"type": "string"
						},
						"currency": {
							"type": "string"
						},
						"redirect_url": {
							"type": "string"
						},
						"payment_options": {
							"type": "string"
						},
						"meta": {
							"type": "string"
						},
						"plan": {
							"type": "string"
						},
						"subaccount": {
							"type": "string"
						},
						"customer": {
							"type": "object",
							"properties": {
								"email": {
									"type": "string"
								},
								"phonenumber": {
									"type": "string"
								},
								"name": {
									"type": "string"
								}
							},
							"required": [
								"email"
							]
						},
						"customizations": {
							"type": "object",
							"properties": {
								"title": {
									"type": "string"
								},
								"description": {
									"type": "string"
								},
								"logo": {
									"type": "string"
								},
								"callback_url": {
									"type": "string"
								}
							},
							"required": [
								"title"
							]
						}
					},
					"required": [
						"tx_ref",
						"amount",
						"currency",
						"redirect_url",
						"payment_options", 
						"customer"
					],
                  	"additionalProperties": false
				},
				{ 
					"type": "object", 
					"description": "verify transaction",
					"properties": {
						"tx_ref": {"type": "string"}
					},
					"required": ["tx_ref"],
                  	"additionalProperties": false
				},
				{ 
					"type": "object", 
					"description": "validate transfer",
					"properties": {
						"store": {"type": "string"},
						"storeID": {"type": "string"},
						"reference": {"type": "string"}
					},
					"required": ["reference"],
                  	"additionalProperties": false
				},
				{ 
					"type": "object",
					"description": "delete transfer recipient", 
					"properties": {
						"recipient": {"type": "string"}
					},
					"required": ["recipient"],
                  	"additionalProperties": false
				},
				{ 
					"type": "object", 
					"description": "createTransferRecipient create transfer recipient",
					"properties": {
						"accountNumber": {"type": "string"},
						"bankCode": {"type": "string"}
					},
					"required": ["accountNumber", "bankCode"],
                  	"additionalProperties": false
				},
				{ 
					"type": "object", 
					"description": "create transfer",
					"properties": {
						"accountNumber": {"type": "string"},
						"amount": {"type": "string"},
						"bankCode": {"type": "string"},
						"bankLocation": {"type": "string"},
						"currency": {"type": "string"},
						"meta": {"type": "string"},
						"narration": {"type": "string"},
						"recipient": {"type": "string"},
						"recipientName": {"type": "string"},
						"recipientPhone": {"type": "string"},
						"reference": {"type": "string"}
					},
					"required": [
						"accountNumber",
						"amount", 
						"recipient", 
						"reference", 
						"currency", 
						"narration", 
						"bankLocation"
					],
                  	"additionalProperties": false
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
