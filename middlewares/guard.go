package middlewares

import (
	"net/http"

	"github.com/aptos-labs/aptos-go-sdk"
)

type GuardConfig struct {
	Facilitator string
	Network     aptos.NetworkConfig
}

type Guard struct {
	GuardConfig
}

func NewGuard(config *GuardConfig) *Guard {
	if config == nil {
		config = &GuardConfig{
			Facilitator: "https://faucet.devnet.aptoslabs.com",
			Network:     aptos.TestnetConfig,
		}
	}

	return &Guard{
		GuardConfig: *config,
	}
}

func NewChain(xs ...middleware) middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			next = xs[i](next)
		}

		return next
	}
}
