package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/0x-ximon/0x402/cmd"
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

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
	ctx := context.Background()
	cfg, err := NewConfig(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := aptos.NewClient(cfg.Network)
	if err != nil {
		log.Fatalln(err)
	}
	_ = client

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
