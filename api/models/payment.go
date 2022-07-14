package models

import (
	"context"
	"encoding/json"
)

type Meta struct {
	User string `json:"user,omitempty"`
	Type string `json:"type,omitempty"`
}
type Customer struct {
	Email       string `json:"email,omitempty"`
	Phonenumber string `json:"phonenumber,omitempty"`
	Name        string `json:"name,omitempty"`
}

type Customizations struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Logo        string `json:"logo,omitempty"`
	CallbackURL string `json:"callback_url,omitempty"`
}

type PaymentRequest struct {
	TxRef          string          `json:"tx_ref,omitempty"`
	Amount         string          `json:"amount,omitempty"`
	Plan           string          `json:"plan,omitempty"`
	Currency       string          `json:"currency,omitempty"`
	RedirectURL    string          `json:"redirect_url,omitempty"`
	Narration      string          `json:"narration,omitempty"`
	Subaccount     string          `json:"subaccount,omitempty"`
	PaymentOptions string          `json:"payment_options,omitempty"`
	Meta           *Meta           `json:"meta,omitempty"`
	Customer       *Customer       `json:"customer,omitempty"`
	Customizations *Customizations `json:"customizations,omitempty"`
}

type PaymentResponse struct {
	Link     string `json:"link"`
	Original *json.RawMessage
	Message  string `json:"message"`
}

type Payment interface {
	InitializePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)
	ValidateTransaction(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)
}
