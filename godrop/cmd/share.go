package cmd

import (
	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share [COMMAND]",
	Short: "Share a file via mdns or tcp holepunch",
	Long: `Share a file via mdns or tcp holepunch. 
	Run godrop share mdns to share a file via mdns.
	Run godrop share hp to share a file with a remote peer via tcp holepunch.`,
}

func init() {
	RootCmd.AddCommand(shareCmd)

}
