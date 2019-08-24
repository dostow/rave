package worker

import (
	// "context"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/dostow/rave/api"
	"github.com/dostow/rave/queues/machinery"
	"github.com/tidwall/gjson"
)

// Keys flutterwave keys
type Keys struct {
	Secret string `json:"secret"`
	Public string `json:"public"`
}

// Config addon config
type Config struct {
	Keys Keys `json:"keys"`
}

// Params linked store params
type Params struct {
	Action   string `json:"action"`
	Callback string `json:"callback"`
	Options  struct {
		AccountNumber string `json:"account"`
		BankCode      string `json:"bank"`
		Amount        string `json:"amount"`
		Recipient     string `json:"recipient"`
		Currency      string `json:"currency"`
		Reference     string `json:"reference"`
		Narration     string `json:"narration"`
		BankLocation  string `json:"bankLocation"`
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
}

func doRave(addonConfig, addonParams, data, traceID string) error {
	config := Config{}
	ctx := context.Background()
	err := json.Unmarshal([]byte(addonConfig), &config)
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
		_, err := api.CreateTransferRecipient(ctx,
			config.Keys.Secret,
			accountNumber.String(), accountBank.String())
		// TODO: after creating a transfer recipient, it should be linked to a store entry
		// a group access key should be provided to this addon
		// the key would have access rules that allows this addon to modify a store entry
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
		_, err := api.CreateTransfer(ctx,
			config.Keys.Secret,
			reference.String(),
			amount.String(),
			recipient.String(),
			currency.String(),
			narration.String(),
			bankLocation.String(),
		)
		return err
	}
	return errors.New("not implemented")
}

// Worker a rave worker that sends messages to centrifuge
type Worker struct {
	Addr    string        `help:"centrifuge web address"`
	Key     string        `help:"centrifuge key"`
	Timeout time.Duration `help:"gocent timeout"`
	ID      string        `help:"worker id"`
	Build   string        `help:"build"`
}

// Run run the worker
func (w *Worker) Run() error {
	return machinery.Worker(w.ID, map[string]interface{}{
		"rave": doRave,
	})
}

// Send a job to another worker
func (w *Worker) Send() error {
	return errors.New("not implemented")
}

// NewWorker new worker
func NewWorker(build string) *Worker {
	return &Worker{Timeout: 5 * time.Second, Build: build}
}
