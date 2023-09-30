package quikk

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var stagingAPIURL = "https://tryapi.quikk.dev/v1/mpesa"
var productionAPIURL = "https://api.quikk.dev/v1/mpesa"
