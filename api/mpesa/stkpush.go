package mpesa

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cstockton/go-conv"
	"github.com/dostow/rave/api/models"
	mpesa_go "github.com/ndunyu/mpesa-go"
)

func (p *MPESA) ValidateTransaction(ctx context.Context, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	return nil, errors.New("not implemented")
}

// InitializePayment initialize a payment
func (p *MPESA) InitializePayment(ctx context.Context, req *models.PaymentRequest) (resp *models.PaymentResponse, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%v", p)
		}
	}()
	mpesa := mpesa_go.New(p.Config.Public, p.Config.Secret, false)
	mpesa.SetDefaultTimeOut(10 * time.Second)
	mpesa.SetDefaultPassKey("bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919")
	mpesa.SetDefaultB2CShortCode("174379")
	amount, _ := conv.Int(req.Amount)
	response, err := mpesa.StkPushRequest(mpesa_go.StKPushRequestBody{
		BusinessShortCode: "",
		Amount:            fmt.Sprintf("%d", amount/100),
		PhoneNumber:       req.Customer.Phonenumber,
		CallBackURL:       req.Customizations.CallbackURL,
		AccountReference:  req.TxRef,
		TransactionDesc:   req.Narration,
	})
	if err == nil {
		if response.ResponseCode == "0" {
			resp = &models.PaymentResponse{
				Message: response.ResponseDescription,
			}
		} else {
			err = fmt.Errorf("transaction failed")
		}
	}
	return
}
