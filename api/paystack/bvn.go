package paystack

import (
	"context"
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/go-resty/resty/v2"
)

// BVNResult successful bvn query
type BVNResult struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Dob          string `json:"dob"`
		FormattedDob string `json:"formatted_dob"`
		Mobile       string `json:"mobile"`
		Bvn          string `json:"bvn"`
	} `json:"data"`
	Meta struct {
		CallsThisMonth int `json:"calls_this_month"`
		FreeCallsLeft  int `json:"free_calls_left"`
	} `json:"meta"`
}

// ResolveBVN resolve bvn
func ResolveBVN(ctx context.Context, bvn string) (*BVNResult, error) {
	log.Debugf("ResolveBVN - %v", bvn)
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetResult(&TransferCreatedResult{}).
		SetError(&errorResponse{}).
		Get(fmt.Sprintf("%s/bank/bvn_resolve/%s", url, bvn))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 200 {
		result := resp.Result().(*BVNResult)
		if result.Status == true {
			return result, nil
		}
		return nil, errors.New(result.Message)
	}
	respBody := resp.Error().(*errorResponse)
	return nil, errors.New(respBody.Message)
}
