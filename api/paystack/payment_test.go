package paystack

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/dostow/rave/api/models"
	"github.com/google/uuid"
)

func TestInitializePayment(t *testing.T) {

	u, _ := uuid.NewUUID()
	type args struct {
		ctx    context.Context 
		req    *models.PaymentRequest
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"", args{context.Background(), , &models.PaymentRequest{
			TxRef:       u.String(),
			Amount:      "100",
			Currency:    "NGN",
			RedirectURL: "/",
			Meta:        &models.Meta{User: "osiloke"},
		}}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Paystack{Config: models.Keys{Secret: "FLWSECK_TEST-a923c443c50b874aa8c5c1e039560c02-X"}}
			got, err := r.InitializePayment(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitializePayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Println(string(*got))
				t.Errorf("InitializePayment() = %v, want %v", got, tt.want)
			}
		})
	}
}
