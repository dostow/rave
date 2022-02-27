package mpesa

import (
	"context"
	"reflect"
	"testing"

	"github.com/dostow/rave/api/models"
)

func TestMPESA_InitializePayment(t *testing.T) {
	type fields struct {
		Config models.Keys
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
			"test push",
			fields{
				models.Keys{
					Public: "AuZ56Nt7rJR1jTGpp9735eOFw9EB6JSs",
					Secret: "mGONreB8mQSVXEnH",
				},
			},
			args{
				context.Background(),
				&models.PaymentRequest{
					TxRef:     "tx1",
					Amount:    "100",
					Narration: "Test payment",
					Customizations: &models.Customizations{
						CallbackURL: "https://eozizwfq633zvy2.m.pipedream.net",
					},
					Customer: &models.Customer{
						Phonenumber: "254746378652",
					},
				},
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &MPESA{
				Config: tt.fields.Config,
			}
			got, err := p.InitializePayment(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MPESA.InitializePayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MPESA.InitializePayment() = %v, want %v", got, tt.want)
			}
		})
	}
}
