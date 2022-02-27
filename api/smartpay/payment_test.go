package smartpay

import (
	"context"
	"reflect"
	"testing"

	"github.com/dostow/rave/api/models"
)

func TestSmartPay_InitializePayment(t *testing.T) {
	type fields struct {
		Config      models.Keys
		ClientId    string
		ClientAppId string
	}
	type args struct {
		ctx context.Context
		req *models.PaymentRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.PaymentResponse
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				Config: models.Keys{
					Secret: "4c78b6653255e710edcd31621ae35eddade6524b11f5af74acacd32950225ed7",
				},
				ClientId:    "520270",
				ClientAppId: "781547",
			},

			args: args{
				context.Background(),
				&models.PaymentRequest{
					Meta:     &models.Meta{User: "09090006712"},
					Currency: "566",
					Amount:   "10",
					TxRef:    "100005",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SmartPay{
				Config:      tt.fields.Config,
				ClientID:    tt.fields.ClientId,
				ClientAppID: tt.fields.ClientAppId,
			}
			got, err := r.InitializePayment(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SmartPay.InitializePayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SmartPay.InitializePayment() = %v, want %v", got, tt.want)
			}
		})
	}
}
