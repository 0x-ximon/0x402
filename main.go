package main

import (
	"fmt"
	"os"

	"github.com/0x-ximon/0x402/cmd"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "x402",
	Short: "A CLI to integrate x402 into your Go application",
	Long: `This CLI simplifies the process of integrating x402 into your Go application,
	and provides a variety of useful features that you'd probably want to use. For example:
   `,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

func init() {
	RootCmd.AddCommand(cmd.InitCmd)
	RootCmd.AddCommand(cmd.AddCmd)
}
func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
