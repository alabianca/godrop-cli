package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the app",
	Run: func(command *cobra.Command, args []string) {
		fmt.Println("In run")
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
