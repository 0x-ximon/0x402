package cmd

import (
	"fmt"
	"os"

	"github.com/0x-ximon/0x402/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Initializing project...")

		p := tea.NewProgram(models.InitialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		return nil
	},
}
