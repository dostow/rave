package quikk

import "github.com/dostow/rave/api/models"

type Quikk struct {
	ShortCode string `json:"short_code"`
	Config    models.Keys
	Staging   bool
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var stagingAPIURL = "https://tryapi.quikk.dev/v1/mpesa/charge"
var productionAPIURL = "https://tryapi.quikk.dev/v1/mpesa/charge"
