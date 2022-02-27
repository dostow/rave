package quikk

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/dostow/rave/api/models"
)

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
				key:    "0f14b9acfedf4993ba73921e32de169b",
				secret: "b79ffe917285678a9539320d3ab7b4a2",
				noww:   mustParseTime("Sun, 27 Feb 2022 10:41:43 UTC").Format("Mon, 02 Jan 2006 15:04:05 MST"),
			},
			`keyId="0f14b9acfedf4993ba73921e32de169b",algorithm="hmac-sha256",headers="x-aux-date",signature="DY%2By63zA9e8Ip3KgML3tAGcemKn7VnyMwbS3M%2FSvxeA%3D"`,
		},
		{
			"",
			args{
				key:    "0f14b9acfedf4993ba73921e32de169b",
				secret: "b79ffe917285678a9539320d3ab7b4a2",
				noww:   mustParseTime("Sun, 27 Feb 2022 11:59:29 UTC").Format("Mon, 02 Jan 2006 15:04:05 MST"),
			},
			`keyId="0f14b9acfedf4993ba73921e32de169b",algorithm="hmac-sha256",headers="x-aux-date",signature="WjfGt89li3hpadb1uVbR4wYSP0qKAvdsT%2FEIipJN3lo%3D"`,
		},
		{
			"",
			args{
				key:    "0f14b9acfedf4993ba73921e32de169b",
				secret: "b79ffe917285678a9539320d3ab7b4a2",
				noww:   mustParseTime("Sun, 27 Feb 2022 12:00:39 UTC").Format("Mon, 02 Jan 2006 15:04:05 MST"),
			},
			`keyId="0f14b9acfedf4993ba73921e32de169b",algorithm="hmac-sha256",headers="x-aux-date",signature="Sm74g%2FCrcomK5hDX27bXK4wdiHH0Ha7U6mpRGh6Urco%3D"`,
		},
		{
			"",
			args{
				key:    "0f14b9acfedf4993ba73921e32de169b",
				secret: "b79ffe917285678a9539320d3ab7b4a2",
				noww:   mustParseTime("Fri, 22 Nov 2019 08:47:01 EAT").Format("Mon, 02 Jan 2006 15:04:05 MST"),
			},
			`keyId="0f14b9acfedf4993ba73921e32de169b",algorithm="hmac-sha256",headers="date x-custom",signature="FhJ%2FH0sMU3gcOAE%2FcUBjFaTYeUw3WQj9d8x2xGFVrNE%3D"`,
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
				ShortCode: "174379",
				Config: models.Keys{
					Public: "488da317cae042d19e0a1afbb47200c1",
					Secret: "1c547740aa7daf4f412c4f81284dd3ab",
					// Public: "0f14b9acfedf4993ba73921e32de169b",
					// Secret: "b79ffe917285678a9539320d3ab7b4a2",
				},
			},
			args{
				context.Background(),
				&models.PaymentRequest{
					Customer: &models.Customer{Phonenumber: "254746378652"},
					Currency: "566",
					Amount:   "100",
					TxRef:    fmt.Sprintf("100005%d", 1),
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
				Config:    tt.fields.Config,
			}
			got, err := r.Charge(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Quikk.Charge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Quikk.Charge() = %v, want %v", got, tt.want)
			}
		})
	}
}
