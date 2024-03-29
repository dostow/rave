package worker

import (
	// "context"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/dostow/rave/api/models"
	"github.com/dostow/rave/api/paystack"
	"github.com/dostow/rave/api/quikk"
	"github.com/dostow/rave/api/rave"
	"github.com/osiloke/dostow-contrib/api"
	"github.com/sethgrid/pester"
	"github.com/tidwall/gjson"
)

type Callback struct {
	URL     string                 `json:"url"`
	Headers map[string]interface{} `json:"headers"`
	Method  string                 `json:"method"`
}

// Config addon config
type Config struct {
	Callback *Callback    `json:"callback"`
	Platform string       `json:"platform"`
	APIKey   string       `json:"apiKey"`
	Keys     models.Keys  `json:"keys"`
	Paystack models.Keys  `json:"paystack"`
	Quikk    *quikk.Quikk `json:"quikk"`
}

// Params linked store params
type Params struct {
	Action   string `json:"action"`
	Callback string `json:"callback"`
	Options  struct {
		AccountNumber  string `json:"accountNumber"`
		Amount         string `json:"amount"`
		BankCode       string `json:"bankCode"`
		BankLocation   string `json:"bankLocation"`
		Currency       string `json:"currency"`
		Meta           string `json:"meta"`
		Narration      string `json:"narration"`
		Recipient      string `json:"recipient"`
		RecipientName  string `json:"recipientName"`
		RecipientPhone string `json:"recipientPhone"`
		Reference      string `json:"reference"`
		Store          string `json:"store"`
		StoreID        string `json:"storeID"`
	} `json:"options"`
}

// CreateTransactionParams create transaction params
type CreateTransactionParams struct {
	Action   string                 `json:"action"`
	Callback string                 `json:"callback"`
	Options  *models.PaymentRequest `json:"options"`
}

// Data data from linked store
type Data struct {
	Data       map[string]interface{}
	GroupName  string
	Method     string
	Owner      string
	StoreID    string `json:"StoreId"`
	StoreName  string `json:"StoreName"`
	StoreTitle string
}

func doRave(apiURL, addonConfig, addonParams, data, traceID string, dry bool) error {
	var err error
	logger := log.WithField("trace", traceID)
	defer logger.Trace("doRave").Stop(&err)
	config := Config{}
	ctx := context.Background()
	client := pester.New()
	client.Concurrency = 3
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialBackoff
	client.KeepLog = true
	err = json.Unmarshal([]byte(addonConfig), &config)
	if err != nil {
		return err
	}
	params := Params{}
	err = json.Unmarshal([]byte(addonParams), &params)
	if err != nil {
		return err
	}
	logger.Debugf("Received %s action", params.Action)
	c := api.NewClient(apiURL, config.APIKey)
	options := params.Options
	var paymentAPI models.API
	switch config.Platform {
	case "paystack":
		paymentAPI = &paystack.Paystack{Config: config.Paystack}
	case "quikk":
		paymentAPI = config.Quikk
	default:
		paymentAPI = &rave.Rave{Config: config.Keys}
	}
	switch params.Action {
	case "validateTransaction":
		params := &CreateTransactionParams{}
		err = json.Unmarshal([]byte(addonParams), params)
		if err != nil {
			return err
		}
		parseStructFields(data, params.Options)
		resp, err := paymentAPI.ValidateTransaction(ctx,
			params.Options,
		)
		if err != nil {
			logger.WithError(err).WithField("data", data).WithField("options", options).Error("Failed validating transfer")
			log.Debugf(`rave.UpdateStore("%s", "%s")`, gjson.Get(data, "StoreName").String(), gjson.Get(data, "Data.id").String())
			response := map[string]interface{}{
				"status":  "failed",
				"options": params.Options,
				"error":   err.Error(),
			}
			if config.Callback != nil {
				return postRequest(client, config, response)
			} else {
				_, err = c.Store.Update(
					gjson.Get(data, "StoreName").String(),
					gjson.Get(data, "Data.id").String(),
					response,
				)
			}
			return err
		}
		callbackStore := gjson.Get(data, "StoreName").String()
		callbackStoreID := gjson.Get(data, "Data.id").String()
		log.Debugf(`rave.UpdateStore("%s", "%s")`, gjson.Get(data, "StoreName").String(), gjson.Get(data, "Data.id").String())
		callbackData := map[string]interface{}{
			"data":   resp.Original,
			"status": "success",
		}
		if config.Callback != nil {
			return postRequest(client, config, callbackData)
		} else {
			_, err = c.Store.Update(
				callbackStore,
				callbackStoreID,
				callbackData,
			)
		}
		return err
	case "createTransactionLink":
		params := &CreateTransactionParams{}
		err = json.Unmarshal([]byte(addonParams), params)
		if err != nil {
			return err
		}
		parseStructFields(data, params.Options)
		resp, err := paymentAPI.InitializePayment(ctx,
			params.Options,
		)
		if err != nil {
			logger.WithError(err).WithField("data", data).WithField("options", options).Error("Failed initializing payment")
			log.Debugf(`rave.UpdateStore("%s", "%s")`, gjson.Get(data, "StoreName").String(), gjson.Get(data, "Data.id").String())
			response := map[string]interface{}{
				"status":  "failed",
				"options": params.Options,
				"error":   err.Error(),
			}
			if config.Callback != nil {
				return postRequest(client, config, response)
			} else {
				_, err = c.Store.Update(
					gjson.Get(data, "StoreName").String(),
					gjson.Get(data, "Data.id").String(),
					response,
				)
			}
			return err
		}
		callbackStore := gjson.Get(data, "StoreName").String()
		callbackStoreID := gjson.Get(data, "Data.id").String()
		// result := map[string]interface{}{"status": "done", "link": resp.Link, "options": params.Options}
		log.Debugf(`rave.UpdateStore("%s", "%s")`, gjson.Get(data, "StoreName").String(), gjson.Get(data, "Data.id").String())
		callbackData := map[string]interface{}{
			"data": resp,
		}
		if config.Callback != nil {
			return postRequest(client, config, callbackData)
		} else {
			_, err = c.Store.Update(
				callbackStore,
				callbackStoreID,
				callbackData,
			)
		}
		return err

	case "validateTransfer":
		// Get Transfer and then update transfer status
		// If the transfer was successful, update the transfer status
		//
		if len(options.Reference) == 0 {
			return errors.New("missing reference template")
		}
		if len(options.Store) == 0 {
			return errors.New("missing store template")
		}
		if len(options.StoreID) == 0 {
			return errors.New("missing store ID template")
		}
		reference := gjson.Get(data, options.Reference)
		storeName := gjson.Get(data, options.Store)
		storeNameString := ""
		if storeName.Exists() {
			storeNameString = storeName.String()
		} else {
			storeNameString = options.Store
		}
		storeID := gjson.Get(data, options.StoreID)
		if !dry {
			log.Debugf(`rave.ValidateTransfer(ctx, "%s", "%s")`, config.Keys.Secret, reference.String())
			resp, err := rave.GetTransfer(ctx,
				config.Keys.Secret,
				reference.String(),
			)
			if err != nil {
				log.Errorf(`rave.ValidateTransfer(ctx, "%s", "%s") = %s`, config.Keys.Secret, reference.String(), err.Error())
				if len(reference.String()) == 0 || strings.Contains(err.Error(), "not found") {
					// update validation and transfer to reflect not found
					// set status to retry, this should trigger an update addon link if exists
					log.Debugf(`rave.UpdateStore("%s", "%s") set to retry`, storeNameString, storeID.String())
					result := map[string]interface{}{
						"status":  "failed",
						"options": params.Options,
						"error":   err.Error(),
					}
					if config.Callback != nil {
						return postRequest(client, config, result)
					} else {
						_, err = c.Store.Update(
							gjson.Get(data, "StoreName").String(),
							gjson.Get(data, "Data.id").String(),
							result,
						)
					}
					if err == nil {
						log.Debugf(`rave.UpdateStore("%s", "%s") = %s`, gjson.Get(data, "StoreName").String(), gjson.Get(data, "Data.id").String())
						result := map[string]interface{}{
							"status":  "retry",
							"options": params.Options,
						}
						if config.Callback != nil {
							return postRequest(client, config, result)
						} else {
							_, err = c.Store.Update(
								gjson.Get(data, "StoreName").String(),
								gjson.Get(data, "Data.id").String(),
								result,
							)
						}
					}
					return err
				}
				log.WithError(err).WithField("resp", resp).Error(`rave.ValidateTransfer failed`)
				return err
			}
			log.WithField("resp", resp).Debug(`rave.ValidateTransfer - transfer retrieved`)
			if strings.Contains(resp.Status, "success") || strings.Contains(resp.Status, "ok") {
				if resp.Data != nil {
					transfer := resp.Data
					log.WithField("transfer", transfer.ID).WithField("status", transfer.Status).Debugf("got transfer")
					updatedData := map[string]interface{}{
						"rave": transfer,
					}
					transferStatus := strings.ToLower(transfer.Status)
					if strings.Contains(transferStatus, "success") {
						updatedData["status"] = "done"
						updatedData["transactionStatus"] = 3
					} else if strings.Contains(transferStatus, "failed") {
						updatedData["status"] = "failed"
						updatedData["transactionStatus"] = 4
					}
					log.Debugf(`rave.UpdateStore("%s", "%s")`, storeNameString, storeID.String())
					_, err = c.Store.Update(
						storeNameString,
						storeID.String(),
						updatedData,
					)
					if err == nil {
						log.Debugf(`rave.UpdateStore("%s", "%s")`, gjson.Get(data, "StoreName").String(), gjson.Get(data, "Data.id").String())
						_, err = c.Store.Update(
							gjson.Get(data, "StoreName").String(),
							gjson.Get(data, "Data.id").String(),
							map[string]interface{}{
								"status":  "done",
								"options": params.Options,
							},
						)
					}
					return err
				}
			}

			_, err = c.Store.Update(
				storeNameString,
				storeID.String(),
				map[string]interface{}{
					"status":  "failed",
					"options": params.Options,
				},
			)
			if err == nil {
				_, err = c.Store.Update(
					gjson.Get(data, "StoreName").String(),
					gjson.Get(data, "Data.id").String(),
					map[string]interface{}{
						"status":  "failed",
						"options": params.Options,
					},
				)
			}
			return err
		}
		// update store with
		return err
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
			c := api.NewClient(apiURL, config.APIKey)
			resp, err := rave.CreateTransferRecipient(ctx,
				config.Keys.Secret,
				accountNumber.String(), accountBank.String())
			if err != nil {
				_, err = c.Store.Update(
					gjson.Get(data, "StoreName").String(),
					gjson.Get(data, "Data.id").String(),
					map[string]interface{}{
						"status":  "error",
						"error":   err.Error(),
						"options": params.Options,
					},
				)
				return err
			}
			if strings.Contains(resp.Status, "success") || strings.Contains(resp.Status, "ok") {
				_, err = c.Store.Update(
					gjson.Get(data, "StoreName").String(),
					gjson.Get(data, "Data.id").String(),
					map[string]interface{}{
						"status":  "done",
						"rave":    resp.Data,
						"options": params.Options,
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
		if len(options.BankCode) == 0 {
			return errors.New("missing bank code template")
		}
		if len(options.AccountNumber) == 0 {
			return errors.New("missing account number template")
		}
		amount := gjson.Get(data, options.Amount)
		recipient := gjson.Get(data, options.Recipient)
		reference := gjson.Get(data, options.Reference)
		currency := gjson.Get(data, options.Currency)
		narration := gjson.Get(data, options.Narration)
		bankLocation := gjson.Get(data, options.BankLocation)
		bankCode := gjson.Get(data, options.BankCode)
		accountNumber := gjson.Get(data, options.AccountNumber)
		meta := gjson.Get(data, options.Meta)
		recipientName := gjson.Get(data, options.RecipientName)
		status := gjson.Get(data, "Data.status").String()
		if !(status == "pending" || status == "retry") {
			msg := "unable to create transfer - status is not pending or being retried"
			logger.Error(msg)
			return errors.New(msg)

		}
		fullname := recipientName.String()
		if len(fullname) == 0 {
			fullname = reference.String()
		}
		metaMap := map[string]interface{}{}
		if meta.Exists() {
			for k, v := range meta.Map() {
				metaMap[k] = v.Value()
			}
		}
		c := api.NewClient(apiURL, config.APIKey)
		resp, err := rave.CreateTransfer(ctx,
			config.Keys.Secret,
			fullname,
			reference.String(),
			fmt.Sprintf("%v", amount.Int()),
			recipient.String(),
			currency.String(),
			narration.String(),
			bankLocation.String(),
			accountNumber.String(),
			bankCode.String(),
			params.Callback,
			metaMap,
		)
		if err != nil {
			logger.Errorf("unable to create transfer - %s", err.Error())
			var existing *rave.GetTransferResult
			if strings.Contains(err.Error(), "Payout with this ref already exists") {
				// get existing
				existing, err = rave.GetTransfer(ctx, config.Keys.Secret, reference.String())
				status := "failed"
				upd := map[string]interface{}{
					"status": status,
				}
				if err == nil {
					upd["rave"] = existing.Data
				}
				_, err = c.Store.Update(
					gjson.Get(data, "StoreName").String(),
					gjson.Get(data, "Data.id").String(),
					upd,
				)
			}
			return err
		}
		if strings.Contains(resp.Status, "success") || strings.Contains(resp.Status, "ok") {
			logger.WithField("resp", resp).Debug("Transfer created")
			raveStatus := strings.ToLower(resp.Data.Status)
			if raveStatus == "successful" {
				status = "completed"
			} else {
				status = "processing"
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
	case "deleteTransferRecipient":
		if len(options.Recipient) == 0 {
			return errors.New("missing recipient template")
		}
		recipient := gjson.Get(data, options.Recipient)
		resp, err := rave.DeleteTransferRecipient(ctx,
			config.Keys.Secret,
			recipient.String(),
		)
		if err != nil {
			log.WithError(err).Error("failed from rave")
			return err
		}
		if strings.Contains(resp.Status, "success") || strings.Contains(resp.Status, "ok") {
			c := api.NewClient(apiURL, config.APIKey)
			_, err = c.Store.Remove(
				gjson.Get(data, "StoreName").String(),
				gjson.Get(data, "Data.hash").String(),
			)
			return err
		}
	}
	return errors.New("not implemented")
}
