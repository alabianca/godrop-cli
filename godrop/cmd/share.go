package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share [INSTANCE]",
	Short: "Share a file via mdns or tcp holepunch",
	Long: `Share a file via mdns or tcp holepunch. 
	Run godrop share mdns to share a file via mdns.
	Run godrop share hp to share a file with a remote peer via tcp holepunch.`,
	Run: runShare,
}

func runShare(command *cobra.Command, args []string) {

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
	RootCmd.AddCommand(shareCmd)

}
