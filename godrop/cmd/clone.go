package cmd

import (
	"fmt"
	"io"
	"log"
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

	header, err := sesh.ReadHeader()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Content-Length: %d\nFile Name: %s\n", header.Size, header.Name)

	file, err := os.Create(header.Name)

	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	var receivedByts int64

	for {
		if (header.Size - receivedByts) < BUF_SIZE {
			io.CopyN(file, sesh, (header.Size - receivedByts))
			break
		}

		io.CopyN(file, sesh, BUF_SIZE)
		receivedByts += BUF_SIZE
	}

	log.Println("Done")

}

func init() {
	RootCmd.AddCommand(cloneCmd)

}
