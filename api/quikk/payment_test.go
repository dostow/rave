package quikk

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dostow/rave/api/models"
	"github.com/google/uuid"
)

func uniqueid_from_uuid() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
func mustParseTime(d string) time.Time {
	t, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", d)
	return t
}
func Test_encrypt(t *testing.T) {
	// time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 MST")
	type args struct {
		key    string
		secret string
		noww   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"",
			args{
				key:    "488da317cae042d19e0a1afbb47200c1",
				secret: "1c547740aa7daf4f412c4f81284dd3ab",
				noww:   mustParseTime("Sun, 27 Feb 2022 10:41:43 UTC").Format("Mon, 02 Jan 2006 15:04:05 MST"),
			},
			`keyId="488da317cae042d19e0a1afbb47200c1",algorithm="hmac-sha256",signature="c7kU4i%2FGMbr3MflWjBW%2F6bPu7w8gUBlMSiPnKlrsfbg%3D"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encrypt(tt.args.key, tt.args.secret, tt.args.noww); got != tt.want {
				t.Errorf("encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuikk_Charge(t *testing.T) {
	public, _ := os.LookupEnv("QUIKK_PUBLIC")
	secret, _ := os.LookupEnv("QUIKK_SECRET")
	phone, _ := os.LookupEnv("QUIKK_PHONE")
	shortCode, _ := os.LookupEnv("QUIKK_SHORT_CODE")
	type fields struct {
		ShortCode string
		Config    models.Keys
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
			"",
			fields{
				ShortCode: shortCode,
				Config: models.Keys{
					Public: public,
					Secret: secret,
				},
			},
			args{
				context.Background(),
				&models.PaymentRequest{
					Customer: &models.Customer{Phonenumber: phone},
					Currency: "566",
					Amount:   "100",
					TxRef:    "199300",
				},
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Quikk{
				ShortCode: tt.fields.ShortCode,
				Public:    tt.fields.Config.Public,
				Secret:    tt.fields.Config.Secret,
			}
			got, err := r.Charge(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Quikk.Charge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				res2B, _ := json.Marshal(got)
				fmt.Println(string(res2B))
				t.Errorf("Quikk.Charge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuikk_Refund(t *testing.T) {
	public, _ := os.LookupEnv("QUIKK_PUBLIC")
	secret, _ := os.LookupEnv("QUIKK_SECRET")
	phone, _ := os.LookupEnv("QUIKK_PHONE")
	shortCode, _ := os.LookupEnv("QUIKK_SHORT_CODE")
	type fields struct {
		ShortCode string
		Config    models.Keys
		Staging   bool
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
			"",
			fields{
				ShortCode: shortCode,
				Config: models.Keys{
					Public: public,
					Secret: secret,
				},
			},
			args{
				context.Background(),
				&models.PaymentRequest{
					Customer: &models.Customer{Phonenumber: phone},
					TxRef:    "6860-82751156-1",
				},
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Quikk{
				ShortCode: tt.fields.ShortCode,
				Public:    tt.fields.Config.Public,
				Secret:    tt.fields.Config.Secret,
				Staging:   tt.fields.Staging,
			}
			got, err := r.Refund(tt.args.ctx, tt.args.req)
			res2B, _ := json.Marshal(got)
			fmt.Println(string(res2B))
			if (err != nil) != tt.wantErr {
				t.Errorf("Quikk.Refund() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Quikk.Refund() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuikk_Payout(t *testing.T) {
	type fields struct {
		ShortCode string
		Config    models.Keys
		Staging   bool
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Quikk{
				ShortCode: tt.fields.ShortCode,
				Public:    tt.fields.Config.Public,
				Secret:    tt.fields.Config.Secret,
				Staging:   tt.fields.Staging,
			}
			got, err := r.Payout(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Quikk.Payout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Quikk.Payout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuikk_ValidateTransaction(t *testing.T) {
	public, _ := os.LookupEnv("QUIKK_PUBLIC")
	secret, _ := os.LookupEnv("QUIKK_SECRET")
	shortCode, _ := os.LookupEnv("QUIKK_SHORT_CODE")
	type fields struct {
		ShortCode string
		Public    string
		Secret    string
		PassKey   string
		Staging   bool
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
			"",
			fields{
				ShortCode: shortCode,
				Public:    public,
				Secret:    secret,
				PassKey:   "",
				Staging:   false,
			},
			args{
				context.Background(),
				&models.PaymentRequest{
					TxRef: "61137-62085768-1",
				},
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Quikk{
				ShortCode: tt.fields.ShortCode,
				Public:    tt.fields.Public,
				Secret:    tt.fields.Secret,
				PassKey:   tt.fields.PassKey,
				Staging:   tt.fields.Staging,
			}
			got, err := r.ValidateTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Quikk.ValidateTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Quikk.ValidateTransaction() = %v, want %v", fmt.Sprintf("%v", string(*got.Original)), tt.want)
			}
		})
	}
}
