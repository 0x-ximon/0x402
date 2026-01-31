package cmd

import (
	"fmt"
	"os"

	"github.com/0x-ximon/0x402/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new resource or feature to the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Adding new resource or feature...")

		p := tea.NewProgram(models.InitialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		return nil
	},
}
