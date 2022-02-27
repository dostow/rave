package models

type API interface {
	Payment
}

// Keys flutterwave keys
type Keys struct {
	Secret    string `json:"secret"`
	Public    string `json:"public"`
	PassKey   string `json:"passkey"`
	ShortCode string `json:"shortcode"`
}
