package worker

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

var secret = "FLWSECK_TEST-a923c443c50b874aa8c5c1e039560c02-X"

func Test_doRave(t *testing.T) {
	type args struct {
		addonConfig string
		addonParams string
		data        string
		traceID     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"",
			args{
				addonConfig: fmt.Sprintf(`{"keys": {"secret":"%s"}}`, secret),
				addonParams: `{
						"action": "createTransfer", 
						"options": {
							"account": "account", 
							"bank": "bank"
						}
					}`,
				data:    `{"account": "0690000031", "bank": "044"}`,
				traceID: "",
			},
			false,
		},
		// {
		// 	"",
		// 	args{
		// 		addonConfig: fmt.Sprintf(`{"keys": {"secret":"%s"}}`, secret),
		// 		addonParams: `{
		// 				"action": "createTransferRecipient",
		// 				"options": {
		// 					"account": "account",
		// 					"bank": "bank"
		// 				}
		// 			}`,
		// 		data:    `{"account": "0690000031", "bank": "044"}`,
		// 		traceID: "",
		// 	},
		// 	false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := doRave("http://localhost:4445/v1/", tt.args.addonConfig, tt.args.addonParams, tt.args.data, tt.args.traceID, true); (err != nil) != tt.wantErr {
				t.Errorf("doRave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createTransactionLink(t *testing.T) {
	type args struct {
		addonConfig string
		addonParams string
		data        string
		traceID     string
	}
	u, _ := uuid.NewUUID()
	req := map[string]interface{}{
		"ref":            u.String(),
		"amount":         "100",
		"currency":       "NGN",
		"redirectURL":    "/",
		"paymentOptions": "card",
		"customer": map[string]interface{}{
			"email": "hovaitis@gmail.com",
			"name":  "Osiloke Emoekpere",
		},
		"customizations": map[string]interface{}{
			"title":       "Dostow Top-up",
			"description": "Top up",
		},
	}
	data, _ := json.Marshal(map[string]interface{}{
		"data": req,
	})
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"",
			args{
				addonConfig: fmt.Sprintf(`{"keys": {"secret":"%s"}}`, secret),
				addonParams: `{
						"action": "createTransactionLink", 
						"options": {
							"tx_ref": "data.ref",
							"amount": "data.amount", 
							"currency": "data.currency",
							"redirect_url": "data.redirectURL",
							"payment_options": "data.paymentOptions",
							"customer": "data.customer",
							"customizations": "data.customizations"
						}
					}`,
				data:    string(data),
				traceID: "",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := doRave("http://localhost:4445/v1/", tt.args.addonConfig, tt.args.addonParams, tt.args.data, tt.args.traceID, true); (err != nil) != tt.wantErr {
				t.Errorf("doRave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
