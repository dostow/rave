package api

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
)

var seckey string
var accountNumber = "0690000031"
var accountBank = "044"

func init() {
	seckey, _ = os.LookupEnv("SEC_KEY")
}
func TestCreateTransfer(t *testing.T) {

	ctx := context.Background()
	resp, err := CreateTransferRecipient(ctx, seckey, accountNumber, accountBank)
	if err != nil {
		t.Errorf("CreateTransfer() error = %v", err)
		return
	}
	type args struct {
		ctx          context.Context
		seckey       string
		reference    string
		amount       string
		recipient    string
		currency     string
		narration    string
		bankLocation string
	}
	tests := []struct {
		name    string
		args    args
		want    *TransferCreatedResult
		wantErr bool
	}{
		{
			"",
			args{
				ctx:          ctx,
				seckey:       seckey,
				reference:    "ref_PMCK",
				amount:       "5000",
				recipient:    fmt.Sprintf("%v", resp.Data.ID),
				currency:     "NGN",
				narration:    "Transfer out",
				bankLocation: AFRICAN,
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTransfer(tt.args.ctx, tt.args.seckey, tt.args.reference,
				tt.args.amount, tt.args.recipient, tt.args.currency, tt.args.narration, tt.args.bankLocation)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTransfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTransfer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTransferRecipient(t *testing.T) {
	type args struct {
		ctx           context.Context
		seckey        string
		accountNumber string
		accountBank   string
	}
	tests := []struct {
		name    string
		args    args
		want    *CreateTransferRecipientResult
		wantErr bool
	}{
		{
			"",
			args{
				ctx:           context.Background(),
				seckey:        seckey,
				accountNumber: accountNumber,
				accountBank:   "044",
			},
			nil,
			false,
		},
		{
			"",
			args{
				ctx:           context.Background(),
				seckey:        "",
				accountNumber: "",
				accountBank:   "",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTransferRecipient(tt.args.ctx, tt.args.seckey, tt.args.accountNumber, tt.args.accountBank)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTransferRecipient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTransferRecipient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteTransferRecipient(t *testing.T) {
	type args struct {
		ctx    context.Context
		seckey string
		rid    string
	}
	tests := []struct {
		name    string
		args    args
		want    *DeleteTransferRecipientResult
		wantErr bool
	}{
		{
			"",
			args{
				ctx:    context.Background(),
				seckey: seckey,
				rid:    "1733",
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteTransferRecipient(tt.args.ctx, tt.args.seckey, tt.args.rid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTransferRecipient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteTransferRecipient() = %v, want %v", got, tt.want)
			}
		})
	}
}
