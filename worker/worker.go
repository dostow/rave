package worker

import (
	// "context"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/dostow/rave/api/rave"
	"github.com/dostow/rave/queues/machinery"
	"github.com/osiloke/dostow-contrib/api"
	"github.com/tidwall/gjson"
)

// Keys flutterwave keys
type Keys struct {
	Secret string `json:"secret"`
	Public string `json:"public"`
}

// Config addon config
type Config struct {
	APIKey string `json:"apiKey"`
	Keys   Keys   `json:"keys"`
}

// Params linked store params
type Params struct {
	Action   string `json:"action"`
	Callback string `json:"callback"`
	Options  struct {
		AccountNumber string `json:"accountNumber"`
		BankCode      string `json:"bankCode"`
		Amount        string `json:"amount"`
		Recipient     string `json:"recipient"`
		Currency      string `json:"currency"`
		Reference     string `json:"reference"`
		Narration     string `json:"narration"`
		BankLocation  string `json:"bankLocation"`
		Meta          string `json:"meta"`
	} `json:"options"`
}

// Data data from linked store
type Data struct {
	Data       map[string]interface{}
	Method     string
	GroupName  string
	Owner      string
	StoreTitle string
	StoreID    string `json:"StoreId`
	StoreName  string `json:"StoreName`
}

func doRave(apiURL, addonConfig, addonParams, data, traceID string, dry bool) error {
	var err error
	logger := log.WithField("trace", traceID)
	defer logger.Trace("doRave").Stop(&err)
	config := Config{}
	ctx := context.Background()
	err = json.Unmarshal([]byte(addonConfig), &config)
	if err != nil {
		return err
	}
	params := Params{}
	err = json.Unmarshal([]byte(addonParams), &params)
	if err != nil {
		return err
	}
	options := params.Options
	switch params.Action {
	case "createTransferRecipient":
		if len(options.AccountNumber) == 0 {
			return errors.New("missing account number template")
		}
		if len(options.BankCode) == 0 {
			return errors.New("missing account bank template")
		}
		accountNumber := gjson.Get(data, options.AccountNumber)
		accountBank := gjson.Get(data, options.BankCode)
		logger.WithFields(log.Fields{"account": accountNumber.String(), "bank": accountBank.String()}).Debug("CreateTransferRecipient")
		if !dry {
			resp, err := rave.CreateTransferRecipient(ctx,
				config.Keys.Secret,
				accountNumber.String(), accountBank.String())
			if err != nil {
				return err
			}
			if strings.Contains(resp.Status, "success") || strings.Contains(resp.Status, "ok") {
				c := api.NewClient(apiURL, config.APIKey)
				_, err = c.Store.Update(
					gjson.Get(data, "StoreName").String(),
					gjson.Get(data, "Data.id").String(),
					map[string]interface{}{
						"status": "done",
						"rave":   resp.Data,
					},
				)
				return err
			}
			return errors.New("failed creating transfer recipient - " + resp.Message)
		}
		log.Debugf(`rave.CreateTransferRecipient(ctx, "%s", "%s", "%s")`, config.Keys.Secret, accountNumber.String(), accountBank.String())
		log.Debugf(`rave.UpdateStore("%s", "%s")`, gjson.Get(data, "StoreName").String(), gjson.Get(data, "Data.id").String())
		// update store with
		return err
	case "createTransfer":
		if len(options.Amount) == 0 {
			return errors.New("missing amount template")
		}
		if len(options.Recipient) == 0 {
			return errors.New("missing recipient template")
		}
		if len(options.Reference) == 0 {
			return errors.New("missing reference template")
		}
		if len(options.Currency) == 0 {
			return errors.New("missing currency template")
		}
		if len(options.Narration) == 0 {
			return errors.New("missing narration template")
		}
		if len(options.BankLocation) == 0 {
			return errors.New("missing bank location template")
		}
		amount := gjson.Get(data, options.Amount)
		recipient := gjson.Get(data, options.Recipient)
		reference := gjson.Get(data, options.Reference)
		currency := gjson.Get(data, options.Currency)
		narration := gjson.Get(data, options.Narration)
		bankLocation := gjson.Get(data, options.BankLocation)
		meta := gjson.Get(data, options.Meta)

		metaMap := map[string]interface{}{}
		if meta.Exists() {
			for k, v := range meta.Map() {
				metaMap[k] = v.Value()
			}
		}
		resp, err := rave.CreateTransfer(ctx,
			config.Keys.Secret,
			reference.String(),
			fmt.Sprintf("%v", amount.Int()),
			recipient.String(),
			currency.String(),
			narration.String(),
			bankLocation.String(),
			params.Callback,
			metaMap,
		)
		if err != nil {
			return err
		}
		if strings.Contains(resp.Status, "success") || strings.Contains(resp.Status, "ok") {
			c := api.NewClient(apiURL, config.APIKey)
			status := "pending"
			raveStatus := strings.ToLower(resp.Data.Status)
			if raveStatus == "successful" {
				status = "completed"
			}
			if raveStatus == "failed" {
				status = "error"
			}
			_, err = c.Store.Update(
				gjson.Get(data, "StoreName").String(),
				gjson.Get(data, "Data.id").String(),
				map[string]interface{}{
					"status": status,
					"rave":   resp.Data,
				},
			)
			return err
		}
		return errors.New("failed creating transfer recipient - " + resp.Message)
	}
	return errors.New("not implemented")
}

// Worker a rave worker that sends messages to centrifuge
type Worker struct {
	Addr      string        `help:"centrifuge web address"`
	Key       string        `help:"centrifuge key"`
	Timeout   time.Duration `help:"gocent timeout"`
	ID        string        `help:"worker id"`
	Build     string        `help:"build"`
	Dry       bool          `help:"dry run"`
	DostowAPI string        `help:"dostow api url"`
}

// Run run the worker
func (w *Worker) Run() error {
	return machinery.Worker(w.ID, map[string]interface{}{
		"rave": func(args ...string) error {
			return doRave(w.DostowAPI, args[0], args[1], args[2], args[3], w.Dry)
		},
	})
}

// Send a job to another worker
func (w *Worker) Send() error {
	return errors.New("not implemented")
}

// NewWorker new worker
func NewWorker(build string) *Worker {
	return &Worker{Timeout: 5 * time.Second, Build: build, Dry: false}
}
