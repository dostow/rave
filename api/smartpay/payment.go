package smartpay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dostow/rave/api/models"
	"github.com/go-resty/resty/v2"
)

// PaymentRequest payment request
type PaymentRequest struct {
	ClientID           string `json:"clientId"`
	ClientAppID        string `json:"clientAppId"`
	MobileNumber       string `json:"mobileNumber"`
	Amount             string `json:"amount"`
	PaymentDescription string `json:"paymentDescription"`
	PaymentTypeID      int    `json:"paymentTypeId"`
	Channel            string `json:"channel"`
	TransactionRef     string `json:"transactionRef"`
	RedirectURL        string `json:"redirectURL"`
	ApprovedCurrency   string `json:"approvedCurrency"`
}

// PaymentResult result of payment
type PaymentResult struct {
	Status         bool             `json:"status"`
	Message        string           `json:"message"`
	TransactionRef string           `json:"transactionRef"`
	Data           *json.RawMessage `json:"data"`
}

func (p *SmartPay) ValidateTransaction(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	return nil, errors.New("not implemented")
}

// InitializePayment initialize a payment
func (r *SmartPay) InitializePayment(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	reqBody := &PaymentRequest{
		ClientID:           r.ClientID,
		ClientAppID:        r.ClientAppID,
		PaymentDescription: fmt.Sprintf("Payment of %s %s for %s", req.Currency, req.Amount, req.Meta.User),
		TransactionRef:     req.TxRef,
		RedirectURL:        req.RedirectURL,
		MobileNumber:       req.Meta.User,
		ApprovedCurrency:   req.Currency,
		Amount:             req.Amount,
		Channel:            "APP",
		PaymentTypeID:      3,
	}
	rr, _ := json.Marshal(&reqBody)
	fmt.Println(
		string(rr),
	)
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetHeader("x-api-key", r.Config.Secret).
		SetResult(&PaymentResult{}).
		SetError(&errorResponse{}).
		SetBody(reqBody).
		Post(fmt.Sprintf("%s/ubank", url))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 {
		result := resp.Result().(*PaymentResult)
		if result.Status {
			return &models.PaymentResponse{Link: "", Original: result.Data}, nil
		}
		fmt.Println(result)
		fmt.Println(string(resp.Body()))
		return nil, errors.New(result.Message)
	}
	respBody := resp.Error().(*errorResponse)
	return nil, errors.New(respBody.Message)
}
