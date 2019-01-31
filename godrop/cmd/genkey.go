package cmd

import (
	"github.com/spf13/cobra"
)

var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "generate a public/private key pair",
	Run:   genkey,
}

func genkey(command *cobra.Command, args []string) {}

func init() {
	RootCmd.AddCommand(genkeyCmd)
}
