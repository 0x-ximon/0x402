package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aptos-labs/aptos-go-sdk"
)

// TODO: Refactor this to be chain agnostic
var networks = map[string]aptos.NetworkConfig{
	"mainnet":  aptos.MainnetConfig,
	"testnet":  aptos.TestnetConfig,
	"localnet": aptos.LocalnetConfig,
	"devnet":   aptos.DevnetConfig,
}

type Config struct {
	Network     aptos.NetworkConfig
	Facilitator string
}

func NewConfig(ctx context.Context) (*Config, error) {
	net, ok := os.LookupEnv("NETWORK")
	if !ok {
		net = "testnet"
	}

	network, found := networks[net]
	if !found {
		return nil, fmt.Errorf("network %s not found", net)
	}

	facilator, ok := os.LookupEnv("FACILITATOR")
	if !ok {
		// TODO: Set defaults for different networks
		facilator = "https://x402-navy.vercel.app/facilitator"
	}

	return &Config{
		Network:     network,
		Facilitator: facilator,
	}, nil
}
