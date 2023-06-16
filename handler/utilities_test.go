package handler

import (
	"testing"

	"github.com/dostow/rave/api/models"
)

func Test_parseStructFields(t *testing.T) {
	type args struct {
		data string
		o    models.PaymentRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"",
			args{
				`{
					"data": {
						"ref":            "mref",
						"amount":         "100",
						"currency":       "NGN",
						"redirectURL":    "/",
						"paymentOptions": "card",
						"customer": {
							"email": "hovaitis@gmail.com",
							"name":  "Osiloke Emoekpere",
						},
						"customizations": {
							"title":       "Dostow Top-up",
							"description": "Top up",
						},
					}
				}`,
				models.PaymentRequest{
					TxRef:          "data.ref",
					Amount:         "data.amount",
					Currency:       "data.currency",
					RedirectURL:    "data.redirectURL",
					PaymentOptions: "data.paymentOptions",
					Customer: &models.Customer{
						Email: "data.customer.email",
						Name:  "data.customer.name",
					},
					Customizations: &models.Customizations{
						Title: "data.customizations.title",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseStructFields(tt.args.data, &tt.args.o)
		})
	}
}
