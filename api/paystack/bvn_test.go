package paystack

import (
	"context"
	"reflect"
	"testing"
)

func TestResolveBVN(t *testing.T) {
	type args struct {
		ctx context.Context
		bvn string
	}
	tests := []struct {
		name    string
		args    args
		want    *BVNResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveBVN(tt.args.ctx, tt.args.bvn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveBVN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveBVN() = %v, want %v", got, tt.want)
			}
		})
	}
}
