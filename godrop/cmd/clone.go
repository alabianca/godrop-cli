package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone [INSTANCE]",
	Short: "Clone a file hosted by a peer in the local network",
	Long:  `Run godrop clone download a file that is hosted by a peer.`,
	Run:   runClone,
}

func runClone(command *cobra.Command, args []string) {

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

	sesh, err := drop.Connect(peer)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buf := make([]byte, 100)

	n, _ := sesh.Read(buf)

	fmt.Println(string(buf[:n]))

}

func init() {
	RootCmd.AddCommand(cloneCmd)

}
