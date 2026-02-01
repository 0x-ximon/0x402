package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aptos-labs/aptos-go-sdk"
)

var networkMap = map[aptos.NetworkConfig]string{
	aptos.MainnetConfig:  "aptos:1",
	aptos.TestnetConfig:  "aptos:2",
	aptos.DevnetConfig:   "aptos:3",
	aptos.LocalnetConfig: "aptos:4",
}

type X402Extra struct {
	Sponsored bool `json:"sponsored"`
}

type X402Requirements struct {
	Scheme            string    `json:"scheme"`
	Network           string    `json:"network"`
	Amount            string    `json:"amount"`
	Asset             string    `json:"asset"`
	PayTo             string    `json:"payTo"`
	MaxTimeoutSeconds uint      `json:"maxTimeoutSeconds"`
	Extra             X402Extra `json:"extra"`
}

type X402Resource struct {
	URL         string `json:"url"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

type X402Payload struct {
	X402Version int                `json:"x402Version"`
	Error       string             `json:"error"`
	Resource    X402Resource       `json:"resource"`
	Accepts     []X402Requirements `json:"accepts"`
}

type BuildPaymentParams struct {
	Request     *http.Request
	Network     *aptos.NetworkConfig
	Description string
	Amount      string
	Asset       string
	Receiver    string
	Timeout     uint
}

type PaymentParams struct {
	Facilitator string
	Payload     X402Payload      `json:"paymentPayload"`
	Requirement X402Requirements `json:"paymentRequirements"`
}

type PaymentResult struct {
	Success       bool    `json:"success"`
	InvalidReason *string `json:"invalidReason"`
	Error         *string `json:"error"`
	Payer         string  `json:"payer"`
	Transaction   string  `json:"transaction"`
	Network       string  `json:"network"`
}

type PaymentTransation struct {
	ID       string `json:"id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Network  string `json:"network"`
}

type PaymentTransactionKey struct{}

func BuildPayment(params BuildPaymentParams) (*X402Payload, error) {
	if params.Request == nil {
		return nil, errors.New("request is nil")
	}

	if params.Network == nil {
		return nil, errors.New("network is nil")
	}

	if params.Receiver == "" {
		return nil, errors.New("receiver is empty")
	}

	return &X402Payload{
		X402Version: 1,
		Error:       "",
		Resource: X402Resource{
			URL:         params.Request.URL.String(),
			Description: params.Description,
			MimeType:    "application/json",
		},
		Accepts: []X402Requirements{
			{
				Scheme:            "exact",
				Network:           networkMap[*params.Network],
				Amount:            params.Amount,
				Asset:             params.Asset,
				PayTo:             params.Receiver,
				MaxTimeoutSeconds: params.Timeout,

				Extra: X402Extra{
					Sponsored: true,
				},
			},
		},
	}, nil
}

func VerifyPayment(params PaymentParams) error {
	url := fmt.Sprintf("%s/%s", params.Facilitator, "verify")

	var req bytes.Buffer
	err := json.NewEncoder(&req).Encode(params)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", &req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type Data struct {
		IsValid       bool    `json:"isValid"`
		InvalidReason *string `json:"invalidReason"`
		Error         *string `json:"error"`
	}

	var data Data
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	if !data.IsValid {
		if data.Error != nil {
			return errors.New(*data.Error)
		}

		if data.InvalidReason != nil {
			return errors.New(*data.InvalidReason)
		}

		return errors.New("unknown error")
	}

	return nil
}

func SettlePayment(params PaymentParams) (*PaymentResult, error) {
	url := fmt.Sprintf("%s/%s", params.Facilitator, "settle")

	var req bytes.Buffer
	err := json.NewEncoder(&req).Encode(params)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", &req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PaymentResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		if result.Error != nil {
			return nil, errors.New(*result.Error)
		}

		if result.InvalidReason != nil {
			return nil, errors.New(*result.InvalidReason)
		}

		return nil, errors.New("unknown error")
	}

	return &result, nil
}
