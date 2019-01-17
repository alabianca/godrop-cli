package cmd

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the app",
}

func init() {
	RootCmd.AddCommand(runCmd)
}
