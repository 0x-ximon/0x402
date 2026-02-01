package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0x-ximon/0x402/services"
)

type middleware func(http.Handler) http.Handler

type PaymentErrorResponse struct {
	Error       string `json:"error"`
	X402Version int    `json:"x402Version"`
}

type PaymentConfig struct {
	Asset       string `json:"asset"`
	Amount      uint   `json:"amount"`
	Receiver    string `json:"receiver"`
	Description string `json:"description"`
}

func (g *Guard) StandardPaywall(payment PaymentConfig) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := services.BuildPaymentParams{
				Request:     r,
				Network:     &g.Network,
				Asset:       payment.Asset,
				Receiver:    payment.Receiver,
				Description: "Standard Paywall",
				Amount:      strconv.FormatUint(uint64(payment.Amount), 10),
			}

			paymentPayload, err := services.BuildPayment(params)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			paymentSignature := r.Header.Get("Payment-Signature")
			if paymentSignature == "" {
				requirementJson, err := json.Marshal(paymentPayload)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				paymentSignature := base64.StdEncoding.EncodeToString(requirementJson)

				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Payment-Required", paymentSignature)
				w.WriteHeader(http.StatusPaymentRequired)

				json.NewEncoder(w).Encode(PaymentErrorResponse{
					Error:       "Payment required",
					X402Version: 2,
				})

				return
			}

			decodedBytes, err := base64.StdEncoding.DecodeString(paymentSignature)
			if err != nil {
				http.Error(w, "Invalid payment signature", http.StatusBadRequest)
				return
			}

			var paymentRequirement services.X402Requirements
			err = json.Unmarshal(decodedBytes, &paymentRequirement)
			if err != nil {
				http.Error(w, "Invalid payment signature", http.StatusBadRequest)
				return
			}

			paymentParams := services.PaymentParams{
				Facilitator: g.Facilitator,
				Payload:     *paymentPayload,
				Requirement: paymentRequirement,
			}

			err = services.VerifyPayment(paymentParams)
			if err != nil {
				http.Error(w, "Payment verification failed", http.StatusBadRequest)
				return
			}

			result, err := services.SettlePayment(paymentParams)
			if err != nil {
				http.Error(w, "Payment settlement failed", http.StatusInternalServerError)
				return
			}

			tx := services.PaymentTransation{
				ID:       result.Transaction,
				Network:  result.Network,
				Receiver: payment.Receiver,
				Sender:   result.Payer,
			}

			transactionJson, err := json.Marshal(tx)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			transactionSignature := base64.StdEncoding.EncodeToString(transactionJson)
			w.Header().Set("Payment-Transaction", transactionSignature)

			next.ServeHTTP(w, r)
		})
	}
}

func (g *Guard) ResourcePaywall() middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implement Aptos Resource Paywall
			next.ServeHTTP(w, r)
		})
	}
}
