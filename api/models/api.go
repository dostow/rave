package models

type API interface {
	Payment
}

// Keys flutterwave keys
type Keys struct {
	Secret string `json:"secret"`
	Public string `json:"public"`
}
