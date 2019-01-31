package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [peer instance]",
	Short: "connect to a peer at the given instance",
	Run:   connect,
}

func connect(command *cobra.Command, args []string) {
	if len(args) <= 0 {
		command.Usage()
		os.Exit(1)
	}

	peer := args[0]

	drop, err := configGodropMdns()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := drop.Connect(peer)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		_, err := io.Copy(conn, os.Stdin)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}
}

func init() {
	RootCmd.AddCommand(connectCmd)
}
