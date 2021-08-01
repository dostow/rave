package paystack

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/apex/log"

	"github.com/go-resty/resty/v2"
)

type errorResponse struct {
	Status  bool   `json:"status"`
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
		Fee              float64   `json:"fee"`
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

// GetTransferResult result of a get transfer request
type GetTransferResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		PageInfo struct {
			Total       int `json:"total"`
			CurrentPage int `json:"current_page"`
			TotalPages  int `json:"total_pages"`
		} `json:"page_info"`
		Transfers []struct {
			ID            int         `json:"id"`
			AccountNumber string      `json:"account_number"`
			BankCode      string      `json:"bank_code"`
			Fullname      string      `json:"fullname"`
			DateCreated   time.Time   `json:"date_created"`
			Currency      string      `json:"currency"`
			DebitCurrency interface{} `json:"debit_currency"`
			Amount        int         `json:"amount"`
			Fee           float64     `json:"fee"`
			Status        string      `json:"status"`
			Reference     string      `json:"reference"`
			Meta          struct {
				User string `json:"user"`
			} `json:"meta"`
			Narration        string      `json:"narration"`
			Approver         interface{} `json:"approver"`
			CompleteMessage  string      `json:"complete_message"`
			RequiresApproval int         `json:"requires_approval"`
			IsApproved       int         `json:"is_approved"`
			BankName         string      `json:"bank_name"`
		} `json:"transfers"`
	} `json:"data"`
}

// TransferBalance result of a transfer balance request
type TransferBalance struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID               int     `json:"Id"`
		ShortName        string  `json:"ShortName"`
		WalletNumber     string  `json:"WalletNumber"`
		AvailableBalance float64 `json:"AvailableBalance"`
		LedgerBalance    float64 `json:"LedgerBalance"`
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
	recipient, currency, narration, bankLocation, accountNumber, bankCode, callback string, meta map[string]interface{}) (*TransferCreatedResult, error) {
	switch bankLocation {
	case AFRICAN:
		client := resty.New()
		data := map[string]interface{}{
			"account_number": accountNumber,
			"amount":         amount,
			"bank_code":      bankCode,
			"callback_url":   callback,
			"currency":       currency,
			"meta":           meta,
			"narration":      narration,
			"recipient":      recipient,
			"reference":      reference,
			"seckey":         seckey,
		}
		log.Debugf("CreateTransfer - %v", data)
		resp, err := client.R().
			EnableTrace().
			SetResult(&TransferCreatedResult{}).
			SetError(&errorResponse{}).
			SetBody(data).
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

// GetTransfer create a transfer recipient
func GetTransfer(ctx context.Context,
	seckey, reference string) (*GetTransferResult, error) {
	data := map[string]string{
		"reference": reference,
		"seckey":    seckey,
	}
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetResult(&GetTransferResult{}).
		SetError(&errorResponse{}).
		SetQueryParams(data).
		Get(fmt.Sprintf("%s/transfers", url))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 {
		result := resp.Result().(*GetTransferResult)
		if result.Status == "success" || result.Status == "ok" {
			return result, nil
		}
		return nil, errors.New(result.Message)
	}
	respBody := resp.Error().(*errorResponse)
	return nil, errors.New(respBody.Message)
}

// GetTransferBalance get balance for transfer
func GetTransferBalance(ctx context.Context,
	seckey string) (*TransferBalance, error) {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetResult(&TransferBalance{}).
		SetError(&errorResponse{}).
		SetBody(map[string]interface{}{
			"seckey": seckey,
		}).
		Get(fmt.Sprintf("%s/transfers/balance", url))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 {
		result := resp.Result().(*TransferBalance)
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
