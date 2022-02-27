package quikk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	conv "github.com/cstockton/go-conv"
	"github.com/dostow/rave/api/models"
	"github.com/go-resty/resty/v2"
)

var nairobi *time.Location

func init() {
	nairobi, _ = time.LoadLocation("Africa/Nairobi")
}

type Attributes struct {
	Amount       int       `json:"amount,omitempty"`
	CustomerType string    `json:"customer_type,omitempty"`
	CustomerNo   string    `json:"customer_no,omitempty"`
	ShortCode    string    `json:"short_code,omitempty"`
	PostedAt     time.Time `json:"posted_at,omitempty"`
	Reference    string    `json:"reference,omitempty"`
	ResourceID   string    `json:"resource_id,omitempty"`
}

type Data struct {
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

// PaymentRequest payment request
type PaymentRequest struct {
	Data Data `json:"data"`
}

// PaymentResult result of payment
type PaymentResult struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			ResourceID string `json:"resource_id"`
		} `json:"attributes"`
	} `json:"data"`
	Meta struct {
		Status string `json:"status"`
		Code   string `json:"code"`
		Detail string `json:"detail"`
	} `json:"meta,omitempty"`
}

type APIErrors struct {
	Errors []struct {
		Status string `json:"status"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

func encrypt(key, secret, noww string) string {
	to_encode := fmt.Sprintf("x-aux-date: %s", noww)
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(to_encode))
	buf := hash.Sum(nil)
	encoded := base64.StdEncoding.Strict().EncodeToString(buf)
	url_encoded := url.QueryEscape(encoded)
	return fmt.Sprintf(`keyId="%s",algorithm="hmac-sha256",headers="x-aux-date",signature="%s"`, key, url_encoded)
}

// InitializePayment initialize a payment
func (r *Quikk) Charge(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	ct := time.Now().UTC().In(nairobi)
	ts := ct.Format("Mon, 02 Jan 2006 15:04:05 MST")
	authorization := encrypt(r.Config.Public, r.Config.Secret, ts)
	amount, _ := conv.Int(req.Amount)
	reqBody := &PaymentRequest{
		Data: Data{
			Type: "charge",
			Attributes: Attributes{
				Amount:       amount / 100,
				CustomerType: "msisdn",
				CustomerNo:   strings.Replace(req.Customer.Phonenumber, "+", "", -1),
				ShortCode:    r.ShortCode,
				PostedAt:     ct,
				Reference:    req.TxRef,
			},
		},
	}
	rr, _ := json.Marshal(&reqBody)
	fmt.Println("x-aux-date: ", ts)
	fmt.Println("authorization: ", authorization)
	fmt.Println("\nBody\n", string(rr))
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetHeader("content-type", "application/json").
		SetHeader("authorization", authorization).
		SetResult(&PaymentResult{}).
		SetError(&errorResponse{}).
		SetBody(reqBody).
		Post(fmt.Sprintf("%s/mpesa/charge", stagingAPIURL))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 {
		result := resp.Result().(*PaymentResult)
		in, _ := json.Marshal(result.Data)
		raw := json.RawMessage(in)
		if result.Meta.Status != "FAIL" {
			return &models.PaymentResponse{Link: "", Original: &raw}, nil
		}
		fmt.Println(result)
		fmt.Println(string(resp.Body()))
		return nil, errors.New(result.Meta.Detail)
	}
	fmt.Println(string(resp.Body()))
	var apiErrors APIErrors
	if err := json.Unmarshal(resp.Body(), &apiErrors); err == nil {
		return nil, errors.New(apiErrors.Errors[0].Detail)
	}
	respBody := resp.Error().(*errorResponse)
	return nil, errors.New(respBody.Message)
}
