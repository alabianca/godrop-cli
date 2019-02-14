package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
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

	dir, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Cloning contents to ", dir)

	drop, err := configGodropMdns()
	progress := NewProgressBar(os.Stdout)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Finding %s ...\n", peer)
	sesh, err := drop.Connect(peer)

	sesh.DebugWriter = progress

	if !sesh.IsEncrypted() {
		log.Warn("Transfer is not encrypted. See godrop help gencert for generating a signed TLS certificate")
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Connected to %s. \nHostname: %s\n", peer, sesh.RemoteHost)

	header, err := sesh.Authenticate()

	if err != nil {
		log.Error("An Error Occurred During Authentication")
	}

	fmt.Printf("Cloning Remote %s. Content-Length: %s (Uncompressed)\n", header.Name, bytesToReadableFormat(header.Size))
	progress.Init(header.Size)

	if err := sesh.CloneDir(dir); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	progress.Done()
}

func init() {
	RootCmd.AddCommand(cloneCmd)

}
