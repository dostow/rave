package smartpay

import "github.com/dostow/rave/api/models"

const url = "https://dashboard.smartpay.ng/api/v1"

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type SmartPay struct {
	Config      models.Keys
	ClientID    string `json:"clientId"`
	ClientAppID string `json:"clientAppId"`
}
