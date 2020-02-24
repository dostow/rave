package worker

import (
	"fmt"
	"testing"
)

var secret = "FLWSECK_TEST-ffbb8b381df35ceba7c87a8dac738018-X"

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
