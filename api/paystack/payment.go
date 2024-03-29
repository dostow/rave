package paystack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dostow/rave/api/models"
	"github.com/go-resty/resty/v2"
)

// PaymentRequest payment request
type PaymentRequest struct {
	TxRef             string   `json:"reference,omitempty"`
	Amount            string   `json:"amount,omitempty"`
	Currency          string   `json:"currency,omitempty"`
	CallbackURL       string   `json:"callbask_url,omitempty"`
	Channels          []string `json:"channels,omitempty"`
	Plan              string   `json:"plan,omitempty"`
	Email             string   `json:"email,omitempty"`
	TransactionCharge string   `json:"transaction_charge,omitempty"`
	Subaccount        string   `json:"subaccount,omitempty"`
	SplitCode         string   `json:"split_code,omitempty"`
	Metadata          string   `json:"metadata,omitempty"`
}

// PaymentLinkData payment link payload
type PaymentLinkData struct {
	Link string `json:"link"`
}

// InitializePaymentResultData   result
type InitializePaymentResultData struct {
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"acecess_code"`
	Reference        string `json:"reference"`
}

// InitializePaymentResult delete result
type InitializePaymentResult struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data"`
}

// ValidateTransaction checks status of a transaction
func (p *Paystack) ValidateTransaction(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetHeader("Authorization", "Bearer "+p.Config.Secret).
		SetResult(&InitializePaymentResult{}).
		SetError(&errorResponse{}).
		Get(fmt.Sprintf("%s/transaction/verify/%s", url, req.TxRef))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 {
		result := resp.Result().(*InitializePaymentResult)
		if result.Status {
			rsp := InitializePaymentResultData{}
			if err := json.Unmarshal(*result.Data, &rsp); err != nil {
				return nil, err
			}
			return &models.PaymentResponse{Original: result.Data}, nil
		}
		return nil, errors.New(result.Message)
	}
	respBody := resp.Error().(*errorResponse)
	return nil, errors.New(respBody.Message)
}

// InitializePayment initialize a payment
func (p *Paystack) InitializePayment(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	client := resty.New()
	metadata, _ := json.Marshal(&req.Meta)
	amount, _ := strconv.Atoi(req.Amount)
	preq := PaymentRequest{
		TxRef:       req.TxRef,
		Amount:      fmt.Sprintf("%v", amount*100),
		Currency:    req.Currency,
		CallbackURL: req.RedirectURL,
		Channels:    strings.Split(req.PaymentOptions, ","),
		Email:       req.Customer.Email,
		Metadata:    string(metadata),
		Plan:        req.Plan,
		Subaccount:  req.Subaccount,
	}
	resp, err := client.R().
		EnableTrace().
		SetHeader("Authorization", "Bearer "+p.Config.Secret).
		SetResult(&InitializePaymentResult{}).
		SetError(&errorResponse{}).
		SetBody(preq).
		Post(fmt.Sprintf("%s/transaction/initialize", url))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 {
		result := resp.Result().(*InitializePaymentResult)
		if result.Status {
			rsp := InitializePaymentResultData{}
			if err := json.Unmarshal(*result.Data, &rsp); err != nil {
				return nil, err
			}
			return &models.PaymentResponse{Link: rsp.AuthorizationURL, Original: result.Data}, nil
		}
		return nil, errors.New(result.Message)
	}
	respBody := resp.Error().(*errorResponse)
	return nil, errors.New(respBody.Message)
}
