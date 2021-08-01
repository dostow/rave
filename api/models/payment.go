package models

import (
	"context"
	"encoding/json"
)

type Meta struct {
	User string `json:"user"`
	Type string `json:"type,omitempty"`
}
type Customer struct {
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
	Name        string `json:"name"`
}

type Customizations struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
}

type PaymentRequest struct {
	TxRef          string          `json:"tx_ref"`
	Amount         string          `json:"amount"`
	Currency       string          `json:"currency"`
	RedirectURL    string          `json:"redirect_url"`
	Narration      string          `json:"narration"`
	PaymentOptions string          `json:"payment_options"`
	Meta           *Meta           `json:"meta"`
	Customer       *Customer       `json:"customer"`
	Customizations *Customizations `json:"customizations,omitempty"`
}

type PaymentResponse struct {
	Link     string `json:"link"`
	Original *json.RawMessage
}

type Payment interface {
	InitializePayment(ctx context.Context, seckey string, req *PaymentRequest) (*PaymentResponse, error)
}