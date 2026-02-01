package main

import (
	"log"
	"net/http"

	"github.com/0x-ximon/0x402/cmd"
	"github.com/0x-ximon/0x402/middlewares"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var key = "ed25519-priv-0xc528840f3f9686831c9ef3a816319fa6846cee244b56b734490ca37c75a15f9a"

var rootCmd = &cobra.Command{
	Use:   "0x402",
	Short: "A CLI to integrate x402 into your Go application",
	Long: `This CLI simplifies the process of integrating x402 into your Go application,
and provides a variety of useful features that you'd probably want to use. For example:
   `,
}

func init() {
	godotenv.Load()

	rootCmd.AddCommand(cmd.InitCmd)
	rootCmd.AddCommand(cmd.AddCmd)
	rootCmd.AddCommand(cmd.QueryCmd)
	rootCmd.AddCommand(cmd.InspectCmd)
}

func main() {
	mux := http.NewServeMux()
	addr := "http://localhost:8080"

	cfg := middlewares.PaymentConfig{
		Amount:      10000, // 0.01 USDC
		Receiver:    "YOUR_ADDRESS",
		Description: "Aptos is Awesome",
		Asset:       "0x69091fbab5f7d635ee7ac5098cf0c1efbe31d68fec0f2cd565e8d168daf52832", // USDC on Aptos
	}

	// The Guard accepts a configuration struct and provides 2 Paywall middleware
	// functions that are compatible with the Go standard http package. You can
	// set it to nil and use the default configuration
	guard := middlewares.NewGuard(nil)

	// Helper function to create chain of middlewares
	chain := middlewares.NewChain(
		// The StandardPaywall middleware is used to protect the entire
		// application and takes the payment config
		guard.StandardPaywall(cfg),
	)

	server := http.Server{
		Addr:    addr,
		Handler: chain(mux),
	}

	log.Printf("Starting server on %s", addr)
	server.ListenAndServe()

	// ctx := context.Background()
	// cfg, err := NewConfig(ctx)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// client, err := aptos.NewClient(cfg.Network)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// _ = client

	// if err := rootCmd.Execute(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}
