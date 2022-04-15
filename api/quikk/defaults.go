package quikk

type Quikk struct {
	ShortCode string `json:"shortcode"`
	Public    string `json:"public"`
	Secret    string `json:"secret"`
	PassKey   string `json:"passkey"`
	Staging   bool
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var stagingAPIURL = "https://tryapi.quikk.dev/v1/mpesa"
var productionAPIURL = "https://tryapi.quikk.dev/v1/mpesa"
