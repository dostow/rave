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
	"github.com/google/uuid"
)

type Attributes struct {
	Amount            int       `json:"amount,omitempty"`
	CustomerType      string    `json:"customer_type,omitempty"`
	CustomerNo        string    `json:"customer_no,omitempty"`
	ShortCode         string    `json:"short_code,omitempty"`
	PostedAt          time.Time `json:"posted_at,omitempty"`
	Reference         string    `json:"reference,omitempty"`
	ResourceID        string    `json:"resource_id,omitempty"`
	RecipientNo       string    `json:"recipient_no,omitempty"`
	RecipientType     string    `json:"recipient_type,omitempty"`
	RecipientIDType   string    `json:"recipient_id_type,omitempty"`
	RecipientIDNumber string    `json:"recipient_id_number,omitempty"`
	OriginTxnId       string    `json:"origin_txn_id,omitempty"`
	Q                 string    `json:"q,omitempty"`
	On                string    `json:"on,omitempty"`
}

type Data struct {
	ID         string     `json:"id,omitempty"`
	Type       string     `json:"type,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`
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
	to_encode := fmt.Sprintf("date: %s", noww)
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(to_encode))
	buf := hash.Sum(nil)
	encoded := base64.StdEncoding.Strict().EncodeToString(buf)
	url_encoded := url.QueryEscape(encoded)
	return fmt.Sprintf(`keyId="%s",algorithm="hmac-sha256",signature="%s"`, key, url_encoded)
}

func (r *Quikk) doRequest(path string, ct time.Time, reqBody interface{}) (*models.PaymentResponse, error) {
	ts := ct.UTC().Format("Mon, 02 Jan 2006 15:04:05 MST")
	authorization := encrypt(r.Public, r.Secret, ts)
	client := resty.New()
	v, _ := json.Marshal(&reqBody)
	fmt.Println(string(v))
	resp, err := client.R().
		EnableTrace().
		SetHeader("content-type", "application/json").
		SetHeader("Authorization", authorization).
		SetHeader("Date", ts).
		SetResult(&PaymentResult{}).
		SetBody(reqBody).
		Post(fmt.Sprintf("%s/%s", r.URL, path))
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
		return nil, errors.New(result.Meta.Detail)
	} else if resp.StatusCode() == 401 {
		result := resp.Result().(*PaymentResult)
		return nil, fmt.Errorf("%s - %s", resp.Status(), result.Meta.Status)
	}
	var apiErrors APIErrors
	b := strings.Replace(string(resp.Body()), "\\", "", -1)
	err = json.Unmarshal([]byte(b), &apiErrors)
	if err == nil {
		if len(apiErrors.Errors) == 0 {
			return nil, fmt.Errorf("%s - %s", "failed", "charge failed")
		}
		return nil, fmt.Errorf("%s - %s", apiErrors.Errors[0].Title, apiErrors.Errors[0].Detail)
	}
	return nil, err
}

// Charge initialize a payment and send an stk push
func (r *Quikk) Charge(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	ct := time.Now()

	amount, _ := conv.Int(req.Amount)
	reqBody := &PaymentRequest{
		Data: Data{
			Type: "charge",
			Attributes: Attributes{
				Amount:       amount / 100,
				CustomerType: "msisdn",
				CustomerNo:   strings.Replace(req.Customer.Phonenumber, "+", "", -1),
				ShortCode:    r.ShortCode,
				Reference:    req.TxRef,
				PostedAt:     time.Now(),
			},
		},
	}
	return r.doRequest("charge", ct, reqBody)
}

// Refund initialize a payment
func (r *Quikk) Refund(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	ct := time.Now()
	reqBody := &PaymentRequest{
		Data: Data{
			Type: "refund",
			Attributes: Attributes{
				ShortCode:   r.ShortCode,
				OriginTxnId: req.TxRef,
			},
		},
	}
	return r.doRequest("refund", ct, reqBody)
}

// Refund initialize a payment
func (r *Quikk) Payout(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	ct := time.Now()
	reqBody := &PaymentRequest{
		Data: Data{
			Type: "payouts",
			Attributes: Attributes{
				ShortCode:   r.ShortCode,
				OriginTxnId: req.TxRef,
			},
		},
	}
	return r.doRequest("payouts", ct, reqBody)
}

func (r *Quikk) ValidateTransaction(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	ct := time.Now()
	reqBody := &PaymentRequest{
		Data: Data{
			ID:   uuid.NewString(),
			Type: "search",
			Attributes: Attributes{
				ShortCode: r.ShortCode,
				Q:         req.TxRef,
				On:        "resource_id",
			},
		},
	}
	return r.doRequest("searches/transaction", ct, reqBody)
}

// InitializePayment initialize a payment
func (p *Quikk) InitializePayment(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	return p.Charge(ctx, req)
}
