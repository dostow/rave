package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

var url = "https://api.ravepay.co/v2/gpx"

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// TransferCreatedResult holds information of a completed create transfer request
type TransferCreatedResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID               int       `json:"id"`
		AccountNumber    string    `json:"account_number"`
		BankCode         string    `json:"bank_code"`
		Fullname         string    `json:"fullname"`
		DateCreated      time.Time `json:"date_created"`
		Currency         string    `json:"currency"`
		Amount           string    `json:"amount"`
		Fee              int       `json:"fee"`
		Status           string    `json:"status"`
		Reference        string    `json:"reference"`
		Narration        string    `json:"narration"`
		CompleteMessage  string    `json:"complete_message"`
		RequiresApproval int       `json:"requires_approval"`
		IsApproved       int       `json:"is_approved"`
		BankName         string    `json:"bank_name"`
	} `json:"data"`
}

// CreateTransferRecipientResult result of a create transfer recipient
type CreateTransferRecipientResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID            int       `json:"id"`
		AccountNumber string    `json:"account_number"`
		BankCode      string    `json:"bank_code"`
		Fullname      string    `json:"fullname"`
		DateCreated   time.Time `json:"date_created"`
		BankName      string    `json:"bank_name"`
	} `json:"data"`
}

// DeleteTransferRecipientResult delete result
type DeleteTransferRecipientResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const (
	// AFRICAN an african bank
	AFRICAN string = "African" // 0
)

// ErrLocationNotSupported location not supported
var ErrLocationNotSupported = errors.New("location not supported")

// CreateTransfer create a transfer request
func CreateTransfer(ctx context.Context, seckey, reference, amount,
	recipient, currency, narration, bankLocation string) (*TransferCreatedResult, error) {
	switch bankLocation {
	case AFRICAN:
		client := resty.New()
		resp, err := client.R().
			EnableTrace().
			SetResult(&TransferCreatedResult{}).
			SetError(&errorResponse{}).
			SetBody(map[string]interface{}{
				"reference": reference,
				"recipient": recipient,
				"narration": narration,
				"amount":    amount,
				"currency":  currency,
				"seckey":    seckey,
			}).
			Post(fmt.Sprintf("%s/transfers/create", url))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode() == 200 {
			result := resp.Result().(*TransferCreatedResult)
			if result.Status == "success" || result.Status == "ok" {
				return result, nil
			}
			return nil, errors.New(result.Message)
		}
		respBody := resp.Error().(*errorResponse)
		return nil, errors.New(respBody.Message)
	}
	return nil, ErrLocationNotSupported
}

// CreateTransferRecipient create a transfer recipient
func CreateTransferRecipient(ctx context.Context,
	seckey, accountNumber, accountBank string) (*CreateTransferRecipientResult, error) {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetResult(&CreateTransferRecipientResult{}).
		SetError(&errorResponse{}).
		SetBody(map[string]interface{}{
			"account_number": accountNumber,
			"account_bank":   accountBank,
			"seckey":         seckey,
		}).
		Post(fmt.Sprintf("%s/transfers/beneficiaries/create", url))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 {
		result := resp.Result().(*CreateTransferRecipientResult)
		if result.Status == "success" || result.Status == "ok" {
			return result, nil
		}
		return nil, errors.New(result.Message)
	}
	respBody := resp.Error().(*errorResponse)
	return nil, errors.New(respBody.Message)
}

// DeleteTransferRecipient create a transfer recipient
func DeleteTransferRecipient(ctx context.Context,
	seckey, rid string) (*DeleteTransferRecipientResult, error) {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetResult(&DeleteTransferRecipientResult{}).
		SetError(&errorResponse{}).
		SetBody(map[string]interface{}{
			"id":     rid,
			"seckey": seckey,
		}).
		Post(fmt.Sprintf("%s/transfers/beneficiaries/delete", url))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 {
		result := resp.Result().(*DeleteTransferRecipientResult)
		if result.Status == "success" || result.Status == "ok" {
			return result, nil
		}
		return nil, errors.New(result.Message)
	}
	respBody := resp.Error().(*errorResponse)
	return nil, errors.New(respBody.Message)
}
