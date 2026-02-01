package cmd

import (
	"fmt"
	"os"

	"github.com/0x-ximon/0x402/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var InspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a x402 transaction on the Aptos blockchain",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Inspecting Transaction...")

		p := tea.NewProgram(models.InitialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		return nil
	},
}
